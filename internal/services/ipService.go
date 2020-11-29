package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"traceip/internal/models"
	"traceip/internal/restclient"

	"github.com/go-redis/redis/v8"
)

//IPService service to obtain information about an IP address
type IPService struct {
	RedisConn *redis.Client
}

//Trace Obtain information about the country where an IP comes from
func (ips *IPService) Trace(ipAddr string) (countryInfo *models.CountryInfo, err error) {
	countryInfo = &models.CountryInfo{}
	val, err := ips.RedisConn.Get(context.Background(), "ip"+ipAddr).Result()
	if err == redis.Nil {
		//Non cached
		countrySummary := &models.CountrySummary{}
		if err := getCountrySummaryFromIP(ipAddr, countrySummary); err != nil {
			return nil, err
		}
		if err := getCountryInfoFromCode(countrySummary.CountryCode, countryInfo); err != nil {
			return nil, err
		}

		b, err := json.Marshal(countryInfo)
		if err != nil {
			return nil, err
		}
		val := fmt.Sprintf("%s", b)
		_, err = ips.RedisConn.Set(context.Background(), "ip"+ipAddr, val, 0).Result()
	} else if err != nil {
		return nil, err
	} else {
		//Cached
		err = json.Unmarshal([]byte(val), countryInfo)
		if err != nil {
			return nil, err
		}
	}

	return countryInfo, nil
}

func getCountrySummaryFromIP(ipAddr string, countrySummary *models.CountrySummary) error {
	countrySummaryURL := fmt.Sprintf("https://api.ip2country.info/ip?%s", ipAddr)
	err := restclient.Get(countrySummaryURL, countrySummary)
	if err != nil || countrySummary == nil {
		return errors.New("Error getting country summary from IP - Error: " + err.Error())
	}
	return nil
}

func getCountryInfoFromCode(countryCode string, countryInfo *models.CountryInfo) error {
	countryInfoURL := fmt.Sprintf("https://restcountries.eu/rest/v2/alpha/%s?fields=name;alpha2Code;alpha3Code;timezones;languages;capital;currencies;translations;latlng", countryCode)
	err := restclient.Get(countryInfoURL, countryInfo)
	if err != nil || countryInfo == nil {
		return errors.New("Error getting country info from IP - Error: " + err.Error())
	}
	return nil
}
