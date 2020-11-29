package main

import (
	"log"
	"traceip/internal/infrastructure"
	"traceip/internal/router"
	"traceip/internal/services"
)

func main() {
	redis := new(infrastructure.Redis)
	redis.Init("redis:6379", "")

	statisticsService := &services.StatisticsService{
		RedisConn: redis.Conn,
	}

	r := router.InitRouter(redis.Conn, statisticsService)

	currenciesService := services.CurrenciesService{
		RedisConn: redis.Conn,
	}

	go currenciesService.Sync()

	statisticsService.InitEventListener()
	go statisticsService.Listen()

	log.Fatal(r.Start(":3000"))
}
