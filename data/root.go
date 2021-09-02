package data

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
)

// WeatherReport defines a weather report
type WeatherReport struct {
	Latitude  float64            `json:"latitude"`
	Longitude float64            `json:"longitude"`
	Currently WeatherDataPoint   `json:"currently"`
	Daily     WeatherDataBlock   `json:"daily"`
	Hourly    []WeatherDataPoint `json:"hourly"`
	APICalls  int                `json:"apicalls"`
	Code      int                `json:"code"`
	Version   string             `json:"version"`
}

// WeatherDataBlock defines a group of data points
type WeatherDataBlock struct {
	Summary string             `json:"summary"`
	Icon    string             `json:"icon"`
	Data    []WeatherDataPoint `json:"data"`
}

// WeatherDataPoint defines a weather data point
type WeatherDataPoint struct {
	Time                int64   `json:"time"`
	Summary             string  `json:"summary"`
	Icon                string  `json:"icon"`
	PrecipIntensity     float64 `json:"precipIntensity"`
	PrecipIntensityMax  float64 `json:"precipIntensityMax"`
	PrecipProbability   float64 `json:"precipProbability"`
	PrecipType          string  `json:"precipType"`
	PrecipAccumulation  float64 `json:"precipAccumulation"`
	Temperature         float64 `json:"temperature"`
	TemperatureMin      float64 `json:"temperatureMin"`
	TemperatureMax      float64 `json:"temperatureMax"`
	ApparentTemperature float64 `json:"apparentTemperature"`
	WindSpeed           float64 `json:"windSpeed"`
	WindGust            float64 `json:"windGust"`
	WindBearing         float64 `json:"windBearing"`
	CloudCover          float64 `json:"cloudCover"`
	Humidity            float64 `json:"humidity"`
	Pressure            float64 `json:"pressure"`
	Visibility          float64 `json:"visibility"`
	Ozone               float64 `json:"ozone"`
	UVIndex             float64 `json:"uvindex"`
}

// WeatherService is the interface for all services that can fetch weather data
type WeatherService interface {
	// GetWeatherReport gets the weather report
	GetWeatherReport(ctx context.Context, lat, long string) (WeatherReport, error)
}

// GetWeatherReport calls all services in parallel and returns the first result
func GetWeatherReport(ctx context.Context, services []WeatherService, lat, long string) WeatherReport {

	ch := make(chan WeatherReport, 1)

	//	Start the service segment
	ctx, seg := xray.BeginSubsegment(ctx, "weather-report")
	defer seg.Close(nil)

	//	For each passed service ...
	for _, service := range services {

		//	Launch a goroutine for each service...
		go func(c context.Context, s WeatherService, la, lo string) {

			//	Get its pollen report ...
			result, err := s.GetWeatherReport(c, la, lo)

			//	As long as we don't have an error, return what we found on the result channel
			if err == nil {
				select {
				case ch <- result:
				default:
				}
			}
		}(ctx, service, lat, long)

	}

	//	Return the first result passed on the channel
	return <-ch
}
