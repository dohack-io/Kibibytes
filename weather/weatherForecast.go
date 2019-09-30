package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type forecast struct {
	Cod     string  `json:"cod"`
	Message float64 `json:"message"`
	Cnt     int     `json:"cnt"`
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  float64 `json:"pressure"`
			SeaLevel  float64 `json:"sea_level"`
			GrndLevel float64 `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    float64 `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   float64 `json:"deg"`
		} `json:"wind"`
		Rain struct {
			ThreeH float64 `json:"3h"`
		} `json:"rain,omitempty"`
		Sys struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Timezone   int    `json:"timezone"`
		Sunrise    int    `json:"sunrise"`
		Sunset     int    `json:"sunset"`
	} `json:"city"`
}

func (w forecast) String() string {

	var inhalt string
	var t time.Time
	var timeString string
	var currentDay = 0

	for _, elem := range w.List {

		t = time.Unix(int64(elem.Dt), 0)
		if currentDay != t.Day() {
			currentDay = t.Day()
			timeString = fmt.Sprintf("🕙: %02d:%02d", t.Hour(), t.Minute())
			inhalt += fmt.Sprintf("Date: %02d.%02d\n", t.Day(), t.Month())
			inhalt += fmt.Sprintf("%s 🌡️: %.1f ☁️: %s \n", timeString, elem.Main.Temp-273.15, elem.Weather[0].Main)
		} else {
			timeString = fmt.Sprintf("🕙: %02d:%02d", t.Hour(), t.Minute())
			inhalt += fmt.Sprintf("%s 🌡️: %.1f ☁️: %s \n", timeString, elem.Main.Temp-273.15, elem.Weather[0].Main)
		}

	}
	if inhalt == "" {
		log.Println("inhalt ist leer")
		return fmt.Sprintf("No forecast possible. Please try again")
	}
	return inhalt
}

func GetWeatherForecast(location string) string {
	res, err := http.Get("http://api.openweathermap.org/data/2.5/forecast?q=" + location + "&appid=" + appid)
	if err != nil {
		log.Println(err)
		return fmt.Sprintf("API-unreachable")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return fmt.Sprintf("API-failed(body)")
	}

	data := forecast{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		return fmt.Sprintf("This location wasnt found. Please try again! :'(")
	}

	return fmt.Sprintf("%+v", data)
}
