package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response struct for Open-Meteo API
type WeatherResponse2 struct {
	Daily struct {
		Time           []string
		Temperaturemax []float64 `json:"temperature_2m_max"`
		Temperaturemin []float64 `json:"temperature_2m_min"`
	} `json:"daily"`
}

// Response struct for Geocoding API
type GeoResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Country   string  `json:"country"`
	} `json:"results"`
}

// Handler to fetch temperature
func getTemperatureHandler1(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")
	place := r.URL.Query().Get("place")

	// If place is given, resolve to lat/lon
	if place != "" {
		geoURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s", place)
		resp, err := http.Get(geoURL)
		if err != nil {
			http.Error(w, "Failed to fetch geocoding data", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var geo GeoResponse
		if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil || len(geo.Results) == 0 {
			http.Error(w, "Place not found", http.StatusBadRequest)
			return
		}

		lat = fmt.Sprintf("%f", geo.Results[0].Latitude)
		lon = fmt.Sprintf("%f", geo.Results[0].Longitude)
		place = geo.Results[0].Name + ", " + geo.Results[0].Country
	}

	if lat == "" || lon == "" {
		http.Error(w, "Missing lat/lon or place query param", http.StatusBadRequest)
		return
	}

	// Build API URL
	apiURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&daily=temperature_2m_max,temperature_2m_min&timezone=auto", lat, lon)

	// Call Open-Meteo API
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode JSON
	var weather WeatherResponse2
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		http.Error(w, "Failed to parse weather data", http.StatusInternalServerError)
		return
	}

	// Respond with temperature
	response := map[string]interface{}{
		"latitude":  lat,
		"longitude": lon,
		"place":     place,
		"forecast":  weather.Daily,
	}

	// pretty print
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		http.Error(w, "Failed to format json", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func WeatherForecast() {
	http.HandleFunc("/temperature", getTemperatureHandler1)
	fmt.Println("Server running on port:8080")
	http.ListenAndServe(":8080", nil)
}
