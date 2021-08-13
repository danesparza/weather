package data

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
)

// WeatherReport defines a weather report
type WeatherReport struct {
	Latitude  float64          `json:"latitude,omitempty"`
	Longitude float64          `json:"longitude,omitempty"`
	Currently WeatherDataPoint `json:"currently,omitempty"`
	Daily     WeatherDataBlock `json:"daily,omitempty"`
	Alerts    []WeatherAlert   `json:"alerts,omitempty"`
	APICalls  int              `json:"apicalls,omitempty"`
	Code      int              `json:"code,omitempty"`
	Version   string           `json:"version"`
}

// WeatherDataBlock defines a group of data points
type WeatherDataBlock struct {
	Summary string             `json:"summary,omitempty"`
	Icon    string             `json:"icon,omitempty"`
	Data    []WeatherDataPoint `json:"data,omitempty"`
}

// WeatherDataPoint defines a weather data point
type WeatherDataPoint struct {
	Time                int64   `json:"time,omitempty"`
	Summary             string  `json:"summary,omitempty"`
	Icon                string  `json:"icon,omitempty"`
	PrecipIntensity     float64 `json:"precipIntensity,omitempty"`
	PrecipIntensityMax  float64 `json:"precipIntensityMax,omitempty"`
	PrecipProbability   float64 `json:"precipProbability,omitempty"`
	PrecipType          string  `json:"precipType,omitempty"`
	PrecipAccumulation  float64 `json:"precipAccumulation,omitempty"`
	Temperature         float64 `json:"temperature,omitempty"`
	TemperatureMax      float64 `json:"temperatureMax,omitempty"`
	ApparentTemperature float64 `json:"apparentTemperature,omitempty"`
	WindSpeed           float64 `json:"windSpeed,omitempty"`
	WindGust            float64 `json:"windGust,omitempty"`
	WindBearing         float64 `json:"windBearing,omitempty"`
	CloudCover          float64 `json:"cloudCover,omitempty"`
	Humidity            float64 `json:"humidity,omitempty"`
	Pressure            float64 `json:"pressure,omitempty"`
	Visibility          float64 `json:"visibility,omitempty"`
	Ozone               float64 `json:"ozone,omitempty"`
	UVIndex             float64 `json:"uvindex,omitempty"`
}

// WeatherAlert defines an alert for a weather event
type WeatherAlert struct {
	Title       string   `json:"title,omitempty"`
	Regions     []string `json:"regions,omitempty"`
	Severity    string   `json:"severity,omitempty"`
	Description string   `json:"description,omitempty"`
	Time        float64  `json:"time,omitempty"`
	Expires     float64  `json:"expires,omitempty"`
	URI         string   `json:"uri,omitempty"`
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
