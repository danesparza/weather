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

	//	Get the daily points:
	dailyPoints := []WeatherDataPoint{}
	for _, item := range response.Daily.Data {

		//	Grab the current daily point
		dailyPoint := WeatherDataPoint{
			ApparentTemperature: item.ApparentTemperature,
			CloudCover:          item.CloudCover,
			Humidity:            item.Humidity,
			Icon:                item.Icon,
			Ozone:               item.Ozone,
			PrecipAccumulation:  item.PrecipAccumulation,
			PrecipIntensity:     item.PrecipIntensity,
			PrecipIntensityMax:  item.PrecipIntensityMax,
			PrecipProbability:   item.PrecipProbability,
			PrecipType:          item.PrecipType,
			Pressure:            item.Pressure,
			Summary:             item.Summary,
			Temperature:         item.Temperature,
			TemperatureMax:      item.TemperatureMax,
			Time:                item.Time,
			Visibility:          item.Visibility,
			WindBearing:         item.WindBearing,
			WindGust:            item.WindGust,
			WindSpeed:           item.WindSpeed,
		}

		//	Add our daily point:
		dailyPoints = append(dailyPoints, dailyPoint)
	}

	//	Get the alerts
	alerts := []WeatherAlert{}
	for _, item := range response.Alerts {

		//	Grab the current alert:
		alert := WeatherAlert{
			Title:       item.Title,
			Regions:     item.Regions,
			Severity:    item.Severity,
			Description: item.Description,
			Time:        float64(item.Time),
			Expires:     item.Expires,
			URI:         item.URI,
		}

		//	Add the alert:
		alerts = append(alerts, alert)
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
		Daily: WeatherDataBlock{
			Summary: response.Daily.Summary,
			Icon:    response.Daily.Icon,
			Data:    dailyPoints,
		},
		Alerts:   alerts,
		APICalls: response.APICalls,
		Code:     response.Code,
	}

	// Close the segment
	seg.Close(nil)

	return retval, nil
}
