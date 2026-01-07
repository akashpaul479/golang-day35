package weatherapi2

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response struct for Open-Meteo API
type WeatherResponse struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
	} `json:"current"`
}

// Handler to fetch temperature
func getTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	if lat == "" || lon == "" {
		http.Error(w, "Missing lat or lon query params", http.StatusBadRequest)
		return
	}

	// Build API URL
	apiURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=temperature_2m", lat, lon)

	// Call Open-Meteo API
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode JSON
	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		http.Error(w, "Failed to parse weather data", http.StatusInternalServerError)
		return
	}

	// Respond with temperature
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"latitude":    lat,
		"longitude":   lon,
		"temperature": weather.Current.Temperature,
	})
}

func Weather3() {

	http.HandleFunc("/temperature", getTemperatureHandler)

	fmt.Printf("Server running on port:8080")
	http.ListenAndServe(":8080", nil)
}
