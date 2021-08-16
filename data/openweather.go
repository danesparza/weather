package data

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-xray-sdk-go/xray"
	"golang.org/x/net/context/ctxhttp"
)

type OpenWeatherRequest struct {
	Lat     string `json:"lat"`     // Latitude
	Long    string `json:"lon"`     // Longitude
	Exclude string `json:"exclude"` // Exclude data from response
	Units   string `json:"units"`   // Units to use
	AppId   string `json:"appid"`   // Application id / token
}

type OpenWeatherResponse struct {
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Timezone       string  `json:"timezone"`
	TimezoneOffset int     `json:"timezone_offset"`
	Current        struct {
		Dt        int     `json:"dt"`
		Sunrise   int     `json:"sunrise"`
		Sunset    int     `json:"sunset"`
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		DewPoint  float64 `json:"dew_point"`
		Uvi       float64 `json:"uvi"`
		Clouds    int     `json:"clouds"`
		Rain      struct {
			LastHour float64 `json:"1h"`
		}
		Visibility int     `json:"visibility"`
		WindSpeed  float64 `json:"wind_speed"`
		WindDeg    int     `json:"wind_deg"`
		Weather    []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"current"`
	Hourly []struct {
		Dt         int     `json:"dt"`
		Temp       float64 `json:"temp"`
		FeelsLike  float64 `json:"feels_like"`
		Pressure   int     `json:"pressure"`
		Humidity   int     `json:"humidity"`
		DewPoint   float64 `json:"dew_point"`
		Uvi        int     `json:"uvi"`
		Clouds     int     `json:"clouds"`
		Visibility int     `json:"visibility"`
		WindSpeed  float64 `json:"wind_speed"`
		WindDeg    int     `json:"wind_deg"`
		WindGust   float64 `json:"wind_gust"`
		Weather    []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Pop  float64 `json:"pop"`
		Rain struct {
			OneH float64 `json:"1h"`
		} `json:"rain,omitempty"`
	} `json:"hourly"`
	Daily []struct {
		Dt        int     `json:"dt"`
		Sunrise   int     `json:"sunrise"`
		Sunset    int     `json:"sunset"`
		Moonrise  int     `json:"moonrise"`
		Moonset   int     `json:"moonset"`
		MoonPhase float64 `json:"moon_phase"`
		Temp      struct {
			Day   float64 `json:"day"`
			Min   float64 `json:"min"`
			Max   float64 `json:"max"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"temp"`
		FeelsLike struct {
			Day   float64 `json:"day"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"feels_like"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		DewPoint  float64 `json:"dew_point"`
		WindSpeed float64 `json:"wind_speed"`
		WindDeg   int     `json:"wind_deg"`
		WindGust  float64 `json:"wind_gust"`
		Weather   []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds int     `json:"clouds"`
		Pop    float64 `json:"pop"`
		Rain   float64 `json:"rain"`
		Uvi    float64 `json:"uvi"`
	} `json:"daily"`
	Alerts []struct {
		SenderName  string   `json:"sender_name"`
		Event       string   `json:"event"`
		Start       float64  `json:"start"`
		End         float64  `json:"end"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	} `json:"alerts"`
}

// OpenWeatherService is a weather service for OpenWeather formatted data
type OpenWeatherService struct{}

// GetWeatherReport gets the weather report
func (s OpenWeatherService) GetWeatherReport(ctx context.Context, lat, long string) (WeatherReport, error) {

	//	Start the service segment
	ctx, seg := xray.BeginSubsegment(ctx, "openweather-service")

	//	The default return value
	retval := WeatherReport{}

	//	Get the api key:
	apikey := os.Getenv("OPENWEATHER_API_KEY")
	if apikey == "" {
		seg.AddError(fmt.Errorf("{OPENWEATHER_API_KEY} key is blank but shouldn't be"))
		return retval, fmt.Errorf("{OPENWEATHER_API_KEY} key is blank but shouldn't be")
	}

	//	Create our request:
	clientRequest, err := http.NewRequest("GET", "https://api.openweathermap.org/data/2.5/onecall", nil)
	if err != nil {
		seg.AddError(err)
		return retval, fmt.Errorf("problem creating request to openweather: %v", err)
	}

	q := clientRequest.URL.Query()
	q.Add("lat", lat)
	q.Add("lon", long)
	q.Add("exclude", "minutely,hourly")
	q.Add("units", "imperial")
	q.Add("appid", apikey)
	clientRequest.URL.RawQuery = q.Encode()

	//	Set our headers
	clientRequest.Header.Set("Content-Type", "application/json; charset=UTF-8")

	//	Execute the request
	client := &http.Client{}
	clientResponse, err := ctxhttp.Do(ctx, xray.Client(client), clientRequest)
	if err != nil {
		seg.AddError(err)
		return retval, fmt.Errorf("error when sending request to OpenWeather API server: %v", err)
	}
	defer clientResponse.Body.Close()

	//	Decode the response:
	owResponse := OpenWeatherResponse{}
	err = json.NewDecoder(clientResponse.Body).Decode(&owResponse)
	if err != nil {
		seg.AddError(err)
		return retval, fmt.Errorf("problem decoding the response from the OpenWeather API server: %v", err)
	}

	//	Get the daily points:
	dailyPoints := []WeatherDataPoint{}

	//	Set the weather data points:
	for _, item := range owResponse.Daily {

		//	Grab the current daily point
		dailyPoint := WeatherDataPoint{
			ApparentTemperature: item.FeelsLike.Day,
			CloudCover:          float64(item.Clouds / 100),
			Humidity:            float64(item.Humidity / 100),
			Icon:                item.Weather[0].Icon,
			UVIndex:             item.Uvi,
			PrecipAccumulation:  item.Rain,
			PrecipProbability:   item.Pop,
			Pressure:            float64(item.Pressure),
			Summary:             item.Weather[0].Description,
			Temperature:         item.Temp.Min,
			TemperatureMax:      item.Temp.Max,
			Time:                int64(item.Dt),
			WindBearing:         float64(item.WindDeg),
			WindGust:            item.WindGust,
			WindSpeed:           item.WindSpeed,
		}

		//	Add our daily point:
		dailyPoints = append(dailyPoints, dailyPoint)
	}

	//	Get the alerts
	alerts := []WeatherAlert{}
	for _, item := range owResponse.Alerts {

		//	Grab the current alert:
		alert := WeatherAlert{
			Title:       item.Event,
			Regions:     []string{item.SenderName},
			Description: item.Description,
			Time:        item.Start,
			Expires:     item.End,
		}

		//	Add the alert:
		alerts = append(alerts, alert)
	}

	//	Format our weather report
	retval = WeatherReport{
		Latitude:  owResponse.Lat,
		Longitude: owResponse.Lon,
		Currently: WeatherDataPoint{
			ApparentTemperature: owResponse.Current.FeelsLike,
			CloudCover:          float64(owResponse.Current.Clouds / 100),
			Humidity:            float64(owResponse.Current.Humidity / 100),
			Icon:                owResponse.Current.Weather[0].Icon,
			UVIndex:             owResponse.Current.Uvi,
			PrecipAccumulation:  owResponse.Current.Rain.LastHour,
			Pressure:            float64(owResponse.Current.Pressure),
			Summary:             owResponse.Current.Weather[0].Description,
			Temperature:         owResponse.Current.Temp,
			Time:                int64(owResponse.Current.Dt),
			WindBearing:         float64(owResponse.Current.WindDeg),
			WindSpeed:           owResponse.Current.WindSpeed,
		},
		Daily: WeatherDataBlock{
			Data: dailyPoints,
		},
		Alerts: alerts,
	}

	xray.AddMetadata(ctx, "OpenWeatherResult", retval)

	// Close the segment
	seg.Close(nil)

	return retval, nil
}
