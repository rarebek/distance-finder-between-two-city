package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type OpenStreetMapResponse struct {
	Lat    string `json:"lat"`
	Lon    string `json:"lon"`
	ErrMsg string `json:"error"`
}

func getLatLong(cityName string) (float64, float64, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"
	url := fmt.Sprintf("%s?q=%s&format=json", baseURL, cityName)

	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var data []OpenStreetMapResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, err
	}

	if len(data) > 0 && data[0].ErrMsg == "" {
		lat, lon := data[0].Lat, data[0].Lon
		return convertToFloat(lat), convertToFloat(lon), nil
	}

	return 0, 0, fmt.Errorf("Location not found for %s", cityName)
}

func convertToFloat(str string) float64 {
	var result float64
	_, err := fmt.Sscanf(str, "%f", &result)
	if err != nil {
		return 0
	}
	return result
}

func main() {
	var cityName, cityName2 string
	fmt.Print("Enter the city name: ")
	fmt.Scanln(&cityName)
	fmt.Print("Enter second city name: ")
	fmt.Scanln(&cityName2)

	latitude, longitude, err := getLatLong(cityName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(latitude, longitude)

	latitude1, longitude1, err := getLatLong(cityName2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(latitude1, longitude1)

	fmt.Printf("%s va %s o'rtasidagi masofa: %.2f km\n", cityName, cityName2, Haversine(latitude, longitude, latitude1, longitude1))

	googleMapsURL := generateGoogleMapsURL(latitude, longitude, latitude1, longitude1)
	fmt.Println("Google Maps uchun LINK:", googleMapsURL)
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371

	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c
	return distance
}

func generateGoogleMapsURL(lat1, lon1, lat2, lon2 float64) string {
	baseURL := "https://www.google.com/maps/dir/?api=1"
	origin := fmt.Sprintf("&origin=%f,%f", lat1, lon1)
	destination := fmt.Sprintf("&destination=%f,%f", lat2, lon2)
	return baseURL + origin + destination
}
