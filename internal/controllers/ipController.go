package controllers

import (
	"fmt"
	"net/http"
	"traceip/internal/helpers"
	"traceip/internal/models"
	"traceip/internal/services"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
)

//IPController IP endpoints handlers
type IPController struct {
	RedisConn         *redis.Client
	StatisticsService *services.StatisticsService
}

//Trace handler
func (ipc *IPController) Trace(c echo.Context) error {
	params := make(map[string]interface{})
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ErrorResponse(err.Error()))
	}

	ipAddr := fmt.Sprintf("%v", params["ip"])

	ipService := services.IPService{
		RedisConn: ipc.RedisConn,
	}
	countryInfo, err := ipService.Trace(ipAddr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ErrorResponse(err.Error()))
	}

	currenciesService := services.CurrenciesService{
		RedisConn: ipc.RedisConn,
	}
	currencies, err := currenciesService.Get()
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.ErrorResponse(err.Error()))
	}

	distance := helpers.DistanceToBuenosAires(countryInfo.LatLng[0], countryInfo.LatLng[1]) / 1000.0
	response := traceResponse(countryInfo, currencies, distance, "USD")

	go ipc.StatisticsService.AddEvent(countryInfo.Alpha2Code, distance)

	return c.JSON(http.StatusOK, response)
}

func traceResponse(countryInfo *models.CountryInfo, currencies *models.Currencies, distance float64, currencyBase string) (response map[string]interface{}) {
	response = make(map[string]interface{})

	//Monedas
	currenciesRate := make(map[string]string)
	var ref float64 = 1.0
	if currencies.Base != currencyBase {
		ref = currencies.Rates[currencyBase]
	}
	for _, currency := range countryInfo.Currencies {
		if currencyCode, ok := currency["code"]; ok {
			currency := ref / currencies.Rates[currencyCode]
			conversionEntry := fmt.Sprintf("%f %s", currency, currencyBase)
			currenciesRate[currencyCode] = conversionEntry
		}
	}

	//Horarios
	timezoneHour := make(map[string]string)
	for _, timezone := range countryInfo.TimeZones {
		timezoneHour[timezone] = helpers.GetTime(timezone)
	}

	//Idiomas
	idiomas := []map[string]string{}
	for _, language := range countryInfo.Languages {
		idiomas = append(idiomas, map[string]string{
			"nombre":       language["name"],
			"nombreNativo": language["nativeName"],
			"isoCode":      language["iso639_1"],
		})
	}

	response["monedas"] = currenciesRate
	response["idiomas"] = idiomas
	response["distanciaEstimada"] = distance
	response["pais"] = countryInfo.Name
	response["isoCode"] = countryInfo.Alpha2Code
	response["horarios"] = timezoneHour
	return response
}
