package services

import (
	"context"
	"encoding/json"
	"fmt"
	"traceip/internal/models"

	"github.com/go-redis/redis/v8"
)

//StatisticsService service to obtain statistics
type StatisticsService struct {
	RedisConn *redis.Client
	Events    chan *statisticsEvent
}

type statisticsEvent struct {
	countryCode string
	distance    float64
}

type statisticsCountryEntry struct {
	Invocations int     `json:"invocations"`
	Distance    float64 `json:"distance"`
}

//Get get precomputed statistics
func (ss *StatisticsService) Get() (*models.PrecomputedStatistics, error) {
	precomputedStatistics := models.PrecomputedStatistics{}
	val, err := ss.RedisConn.Get(context.Background(), "statistics").Result()
	if err == nil {
		err = json.Unmarshal([]byte(val), &precomputedStatistics)
		if err != nil {
			return nil, err
		}
	} else if err == redis.Nil {
		return nil, nil
	} else {
		return nil, err
	}
	return &precomputedStatistics, nil
}

//InitEventListener initialize event listener channel
func (ss *StatisticsService) InitEventListener() {
	ss.Events = make(chan *statisticsEvent)
}

//AddEvent add event to queue
func (ss *StatisticsService) AddEvent(countryCode string, distance float64) {
	if countryCode != "AR" {
		event := &statisticsEvent{
			countryCode: countryCode,
			distance:    distance,
		}
		ss.Events <- event
	}
}

//Listen listen for events
func (ss *StatisticsService) Listen() {
	for {
		e := <-ss.Events
		if e != nil {
			//Increment invocation
			invocationsMap := make(map[string]statisticsCountryEntry)
			val, err := ss.RedisConn.Get(context.Background(), "invocations").Result()
			if err == nil {
				err = json.Unmarshal([]byte(val), &invocationsMap)
				if err != nil {
					continue
				}
			}

			if val, ok := invocationsMap[e.countryCode]; ok {
				invocationsMap[e.countryCode] = statisticsCountryEntry{
					Distance:    e.distance,
					Invocations: val.Invocations + 1,
				}
			} else {
				invocationsMap[e.countryCode] = statisticsCountryEntry{
					Distance:    e.distance,
					Invocations: 1,
				}
			}

			b, err := json.Marshal(invocationsMap)
			if err != nil {
				fmt.Println(err.Error())
			}
			_, err = ss.RedisConn.Set(context.Background(), "invocations", b, 0).Result()

			//Precompute statistics
			precomputedStatistics := models.PrecomputedStatistics{}
			val, err = ss.RedisConn.Get(context.Background(), "statistics").Result()
			if err == nil {
				err = json.Unmarshal([]byte(val), &precomputedStatistics)
				if err != nil {
					continue
				}
			}

			var avg float64 = 0
			var totalInvocations int
			for countryCode, statisticsEntry := range invocationsMap {
				avg += float64(statisticsEntry.Invocations) * statisticsEntry.Distance
				totalInvocations += statisticsEntry.Invocations
				if statisticsEntry.Distance > precomputedStatistics.Max {
					precomputedStatistics.Max = statisticsEntry.Distance
					precomputedStatistics.MaxCountry = countryCode
				}

				if statisticsEntry.Distance < precomputedStatistics.Min || precomputedStatistics.Min == 0 {
					precomputedStatistics.Min = statisticsEntry.Distance
					precomputedStatistics.MinCountry = countryCode
				}
			}
			if totalInvocations > 0 {
				precomputedStatistics.Avg = avg / float64(totalInvocations)
			}

			b, err = json.Marshal(precomputedStatistics)
			if err != nil {
				fmt.Println(err.Error())
			}
			_, err = ss.RedisConn.Set(context.Background(), "statistics", b, 0).Result()

		}
	}
}
