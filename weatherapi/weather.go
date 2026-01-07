package weatherapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type OpenWeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

type WeatherResult struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
	Status      string  `json:"status"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	place := r.URL.Query().Get("place")

	if strings.TrimSpace(place) == "" {
		http.Error(w, "place should not be empty", http.StatusBadRequest)
		return
	}

	apikey := os.Getenv("WEATHER_KEY")

	if apikey == "" {
		http.Error(w, "Weather key is not set in env", http.StatusInternalServerError)
		return
	}
	Cleanplace := strings.TrimSpace(place)

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s,IN&appid=%s&units=metric", Cleanplace, apikey)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "failed to connect weather service", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "location not found", http.StatusNotFound)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "invalid response", http.StatusBadRequest)
		return
	}
	var ow OpenWeatherResponse
	if err := json.Unmarshal(body, &ow); err != nil {
		http.Error(w, "failes to parse weather json", http.StatusBadGateway)
		return
	}

	condition := "Not available"

	if len(ow.Weather) > 0 {
		condition = strings.ToLower(ow.Weather[0].Description)

	}
	result := WeatherResult{
		Location:    ow.Name,
		Temperature: ow.Main.Temp,
		Status:      condition,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// main func
func WeatherAPI() {

	http.HandleFunc("/api/weather", WeatherHandler)

	fmt.Println("Server running on port:8080")
	http.ListenAndServe(":8080", nil)
}
