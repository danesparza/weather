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
	response, err := darksky.Get(apikey, lat, long, "now", darksky.CA, darksky.English)
	if err != nil {
		return retval, fmt.Errorf("Error getting forecast: %v", err)
	}

	//	Format our weather report
	retval = WeatherReport{
		APICalls: response.APICalls,
		Code:     response.Code,
	}

	// Close the segment
	seg.Close(nil)

	return retval, nil
}
