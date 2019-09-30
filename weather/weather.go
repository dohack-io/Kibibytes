package weather

import (
	"Kibibytes/utils/secrets"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type weatherData struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func (w weatherData) String() string {
	var emoji = ""

	switch w.Weather[0].Main {
	case "clear", "Clear":
		emoji = "☀️"
	case "clouds", "Clouds":
		emoji = "☁️"
	case "rain", "Rain":
		emoji = "🌧️"
	}

	return fmt.Sprintf("The lowest temperature is about %.1f🌡️ degrees. The highest temperature is %.1f🌡️ degrees. Currently the temperature is %.1f🌡️ degrees. "+
		"The forecast for today: %s",
		w.Main.TempMin-273.15, w.Main.TempMax-273.15, w.Main.Temp-273.15, emoji)
}

var appid = secrets.Get("OPENWEATHERMAP_APPID")

func GetWeatherNow(location string) string {
	res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + location + "&appid=" + appid)
	if err != nil {
		log.Println(err)
		return fmt.Sprintf("API-unreachable")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return fmt.Sprintf("API-failed(body)")
	}

	data := weatherData{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		return fmt.Sprintf("This location wasnt found. Please try again! :'(")
	}

	return fmt.Sprintf("%+v", data)
}
