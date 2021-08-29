package data_test

import (
	"context"
	"testing"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/danesparza/weather/data"
)

func TestOpenWeather_GetWeatherReport_ReturnsValidData(t *testing.T) {
	//	Arrange
	service := data.OpenWeatherService{}
	lat := "34.016410"
	long := "-83.906870"
	ctx := context.Background()
	ctx, seg := xray.BeginSegment(ctx, "unit-test")
	defer seg.Close(nil)

	//	Act
	response, err := service.GetWeatherReport(ctx, lat, long)

	//	Assert
	if err != nil {
		t.Errorf("Error calling GetWeatherReport: %v", err)
	}

	t.Logf("Returned object: %+v", response)

}

func TestOpenWeather_GetWeatherReport_HumidityIsPresent(t *testing.T) {
	//	Arrange
	service := data.OpenWeatherService{}
	lat := "34.016410"
	long := "-83.906870"
	ctx := context.Background()
	ctx, seg := xray.BeginSegment(ctx, "unit-test")
	defer seg.Close(nil)

	//	Act
	response, err := service.GetWeatherReport(ctx, lat, long)

	//	Assert
	if err != nil {
		t.Errorf("Error calling GetWeatherReport: %v", err)
	}

	if response.Currently.Humidity < 1 {
		t.Errorf("Humidity should be greater than 1, but wasn't")
	} else {
		t.Logf("Current humidity: %v", response.Currently.Humidity)
	}
}
