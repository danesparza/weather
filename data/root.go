package data

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
	darksky "github.com/danesparza/darksky/v2"
)

// WeatherService is the interface for all services that can fetch weather data
type WeatherService interface {
	// GetWeatherReport gets the weather report
	GetWeatherReport(ctx context.Context, lat, long string) (darksky.Forecast, error)
}

// GetWeatherReport calls all services in parallel and returns the first result
func GetWeatherReport(ctx context.Context, services []WeatherService, lat, long string) darksky.Forecast {

	ch := make(chan darksky.Forecast, 1)

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
