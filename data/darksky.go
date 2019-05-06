package data

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-xray-sdk-go/xray"
	darksky "github.com/danesparza/darksky/v2"
)

// DarksyService is a weather service for Darksky formatted data
type DarksyService struct{}

// GetWeatherReport gets the weather report
func (s DarksyService) GetWeatherReport(ctx context.Context, lat, long string) (WeatherReport, error) {

	//	Start the service segment
	ctx, seg := xray.BeginSubsegment(ctx, "darksky-service")

	//	The default return value
	retval := WeatherReport{}

	//	Get the api key:
	apikey := os.Getenv("DARKSKY_API_KEY")
	if apikey == "" {
		return retval, fmt.Errorf("{DARKSKY_API_KEY} key is blank but shouldn't be")
	}

	//	Create our request:
	response, err := darksky.Get(apikey, lat, long, "now", darksky.US, darksky.English)
	if err != nil {
		return retval, fmt.Errorf("Error getting forecast: %v", err)
	}

	//	Format our weather report
	retval = WeatherReport{
		Latitude:  response.Latitude,
		Longitude: response.Longitude,
		Currently: WeatherDataPoint{
			ApparentTemperature: response.Currently.ApparentTemperature,
			CloudCover:          response.Currently.CloudCover,
			Humidity:            response.Currently.Humidity,
			Icon:                response.Currently.Icon,
			Ozone:               response.Currently.Ozone,
			PrecipAccumulation:  response.Currently.PrecipAccumulation,
			PrecipIntensity:     response.Currently.PrecipIntensity,
			PrecipIntensityMax:  response.Currently.PrecipIntensityMax,
			PrecipProbability:   response.Currently.PrecipProbability,
			PrecipType:          response.Currently.PrecipType,
			Pressure:            response.Currently.Pressure,
			Summary:             response.Currently.Summary,
			Temperature:         response.Currently.Temperature,
			TemperatureMax:      response.Currently.TemperatureMax,
			Time:                response.Currently.Time,
			Visibility:          response.Currently.Visibility,
			WindBearing:         response.Currently.WindBearing,
			WindGust:            response.Currently.WindGust,
			WindSpeed:           response.Currently.WindSpeed,
		},
		APICalls: response.APICalls,
		Code:     response.Code,
	}

	// Close the segment
	seg.Close(nil)

	return retval, nil
}
