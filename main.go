package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type ApiWeatherResponse struct {
	Lat      float32 `json:"lat"`
	Lon      float32 `json:"lon"`
	Timezone string  `json:"timezone"`
	Current  struct {
		Temp      float32 `json:"temp"`
		FeelsLike float32 `json:"feels_like"`
		Weather   []Weather
	}
}

type ApiGeoCodingResponse struct {
	Name    string  `json:"name"`
	Lat     float32 `json:"lat"`
	Lon     float32 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}

type ApiErrorResponse struct {
	Cod     int    `json:"cod"`
	Message string `json:"message"`
}

func fetchGeocoding(url string, ch chan<- []ApiGeoCodingResponse) {
	response, err := http.Get(url)

	if err != nil {
		log.Println("Error: ", err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {

		var error ApiErrorResponse

		err = json.NewDecoder(response.Body).Decode(&error)
		if err != nil {
			log.Println("Error decoding Error JSON: ", err)
			return
		}

		log.Println(error)
		return
	}

	var data []ApiGeoCodingResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Println("Error decoding JSON: ", err)
		return
	}

	ch <- data
}

func fetchWeather(url string, ch chan<- ApiWeatherResponse) {
	response, err := http.Get(url)

	if err != nil {
		log.Println("Error: ", err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {

		var error ApiErrorResponse

		err = json.NewDecoder(response.Body).Decode(&error)
		if err != nil {
			log.Println("Error decoding Error JSON: ", err)
			return
		}

		log.Println(error)
		return
	}

	var data ApiWeatherResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Println("Error decoding JSON: ", err)
		return
	}

	ch <- data
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	openWeatherGeoCodingApiKey := os.Getenv("OPEN_WEATHER_GEOCODING_API_KEY")
	if openWeatherGeoCodingApiKey == "" {
		log.Fatal("OPEN_WEATHER_GEOCODING_API_KEY environment variable not set.")
	}

	openWeatherApiKey := os.Getenv("OPEN_WEATHER_API_KEY")
	if openWeatherApiKey == "" {
		log.Fatal("OPEN_WEATHER_API_KEY environment variable not set.")
	}

	cityName := "Accra"
	geocodingUrl := fmt.Sprintf("https://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s", cityName, openWeatherGeoCodingApiKey)
	geocodingCh := make(chan []ApiGeoCodingResponse)

	go fetchGeocoding(geocodingUrl, geocodingCh)

	geocodingData := <-geocodingCh
	log.Println("GeoCodingData: ", geocodingData)

	weatherUrl := fmt.Sprintf("https://openweathermap.org/data/2.5/onecall?lat=%f&lon=%f&units=metric&exclude=minutely,hourly,daily,alerts&appid=%v", geocodingData[0].Lat, geocodingData[0].Lon, openWeatherApiKey)
	weatherCh := make(chan ApiWeatherResponse)

	go fetchWeather(weatherUrl, weatherCh)

	weatherData := <-weatherCh
	log.Println("weatherData: ", weatherData)

	// slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	// if slackBotToken == "" {
	// 	log.Fatal("SLACK_BOT_TOKEN environment variable not set.")
	// }

	// slackAppToken := os.Getenv("SLACK_APP_TOKEN")
	// if slackAppToken == "" {
	// 	log.Fatal("SLACK_APP_TOKEN environment variable not set.")
	// }

	// bot := slacker.NewClient(slackBotToken, slackAppToken)

	// bot.AddCommand("weather in <location>", &slacker.CommandDefinition{
	// 	Description: '',
	// })

}
