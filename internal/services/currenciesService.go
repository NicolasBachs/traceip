package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"traceip/internal/models"
	"traceip/internal/restclient"

	"github.com/go-redis/redis/v8"
)

//CurrenciesService service to obtain information about currencies
type CurrenciesService struct {
	RedisConn *redis.Client
}

//Sync allows us to keep currencies synchronized on redis database
func (cs *CurrenciesService) Sync() (interface{}, error) {
	for {
		fmt.Println("Syncronizing currencies...")
		if _, err := cs.update(); err != nil {
			time.Sleep(15 * time.Second)
			continue
		}
		time.Sleep(45 * time.Minute)
	}
}

//Get allows us to obtain last currencies that we stored on redis database
func (cs *CurrenciesService) Get() (currencies *models.Currencies, err error) {
	currencies = &models.Currencies{}
	val, err := cs.RedisConn.Get(context.Background(), "currencies").Result()
	if err == redis.Nil {
		return nil, nil
	}
	err = json.Unmarshal([]byte(val), currencies)
	return currencies, err
}

func (cs *CurrenciesService) update() (currencies *models.Currencies, err error) {
	currencies = &models.Currencies{}
	if err := getCurrenciesFromBase(currencies); err != nil {
		return nil, err
	}

	b, err := json.Marshal(currencies)
	val := fmt.Sprintf("%s", b)
	_, err = cs.RedisConn.Set(context.Background(), "currencies", val, 0).Result()

	return currencies, err
}

func getCurrenciesFromBase(currencies *models.Currencies) error {
	err := restclient.Get("http://data.fixer.io/api/latest?access_key=04d37504d1e4a10c38dc291f7485d9e8&format=1", currencies)
	if err != nil || currencies == nil {
		return errors.New("Error getting currencies")
	}
	return nil
}
