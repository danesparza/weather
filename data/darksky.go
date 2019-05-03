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
func (s DarksyService) GetWeatherReport(ctx context.Context, lat, long string) (darksky.Forecast, error) {

	//	Start the service segment
	ctx, seg := xray.BeginSubsegment(ctx, "darksky-service")

	//	The default return value
	retval := darksky.Forecast{}

	//	Get the api key:
	apikey := os.Getenv("darksky-api-key")
	if apikey == "" {
		return retval, fmt.Errorf("{darksky-api-key} key is blank but shouldn't be")
	}

	//	Create our request:
	response, err := darksky.Get(apikey, lat, long, "now", darksky.CA, darksky.English)
	if err != nil {
		return retval, fmt.Errorf("Error getting forecast: %v", err)
	}

	retval = *response

	// Close the segment
	seg.Close(nil)

	return retval, nil
}
