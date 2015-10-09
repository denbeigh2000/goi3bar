package weather

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"
)

type Fetcher interface {
	Fetch() (WeatherInfo, error)
}

const urlFormat = "http://http://api.openweathermap.org/data/2.5/weather?q=%s,%s"

type SunState int

const (
	PreDawn SunState = iota
	Day
	Night
)

type OpenWeatherInfoResult struct {
	Coord struct {
		Lon float64
		Lat float64
	}
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
		Pressure float64 `json:"pressure"`
	} `json:"main"`
	Visibility uint64 `json:"visibility"`
	Wind       struct {
		Speed   float64 `json:"speed"`
		Degrees float64 `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		Cover float64 `json:"all"`
	} `json:"clouds"`
	CurrentTime float64 `json:"dt"`
	Sys         struct {
		Country string  `json:"country"`
		Sunrise float64 `json:"sunrise"`
		Sunset  float64 `json:"sunset"`
	} `json:"sys"`
	Name string `json:"name"`
}

type Time struct {
	CurrentTime time.Time
	SunriseTime time.Time
	SunsetTime  time.Time
	SunState    SunState
}

type Temp struct {
	Current float64
	Max     float64
	Min     float64
}

type Wind struct {
	Speed   uint64
	Degrees float64
}

type MiscWeather struct {
	CloudCover uint64
	Visibility uint64
	Pressure   float64
}

type Weather struct {
	MainDesc   string
	SecondDesc string

	Temp
	MiscWeather
}

type Location struct {
	Lat         float64
	Long        float64
	City        string
	CountryCode string
}

type WeatherInfo struct {
	Time
	Weather
	Location
}

type OpenWeatherFetcher struct {
	// Name of city to search for
	City string

	// ISO 3166 Country Code
	Country string
}

func (f OpenWeatherFetcher) fetch() (result OpenWeatherInfoResult, err error) {
	result = OpenWeatherInfoResult{}
	err = nil

	url := fmt.Sprintf(urlFormat, f.City, f.Country)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(result)

	return
}

func (f OpenWeatherFetcher) Fetch() (info WeatherInfo, err error) {
	result, err := f.fetch()
	if err != nil {
		return
	}

	info = WeatherInfo{
		Time: Time{
			CurrentTime: time.Unix(int64(result.CurrentTime), 0),
			SunriseTime: time.Unix(int64(result.Sys.Sunrise), 0),
			SunsetTime:  time.Unix(int64(result.Sys.Sunset), 0),
		},

		Weather: Weather{
			MainDesc:   result.Weather[0].Main,
			SecondDesc: result.Weather[0].Description,

			Temp: Temp{
				Current: result.Main.Temp,
				Min:     result.Main.TempMin,
				Max:     result.Main.TempMax,
			},

			MiscWeather: MiscWeather{
				CloudCover: uint64(result.Clouds.Cover),
				Visibility: result.Visibility,
				Pressure:   result.Main.Pressure,
			},
		},
		Location: Location{
			Lat:         result.Coord.Lat,
			Long:        result.Coord.Lon,
			City:        result.Name,
			CountryCode: result.Sys.Country,
		},
	}

	var sunState SunState
	switch {
	case info.Time.CurrentTime.After(info.Time.SunsetTime):
		sunState = Night
	case info.Time.CurrentTime.After(info.Time.SunriseTime):
		sunState = Day
	default:
		sunState = PreDawn
	}
	info.Time.SunState = sunState

	return
}

type WeatherApplet struct {
	Fetcher

	Template *template.Template
}

func (f *WeatherApplet) Generate() ([]i3.Output, error) {
	info, err := f.Fetch()
	if err != nil {
		return nil, err
	}

	var outBuf bytes.Buffer
	err = f.Template.Execute(&outBuf, info)
	if err != nil {
		return nil, err
	}

	var color string
	switch info.Time.SunState {
	case PreSunrise:
		color = "#0077CC"
	case Day:
		color = "#FFFF00"
	case Night:
		color = "#00FFFF"
	}

	output := i3.Output{
		FullText: outBuf.String(),
		Color:    color,
	}

	return []i3.Output{output}, nil
}
