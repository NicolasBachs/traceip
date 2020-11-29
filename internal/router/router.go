package router

import (
	"traceip/internal/controllers"
	"traceip/internal/services"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//InitRouter initialize routes to access resources
func InitRouter(redisConn *redis.Client, statisticsService *services.StatisticsService) *echo.Echo {
	r := echo.New()

	initRoutes(r, redisConn, statisticsService)
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	r.Use(middleware.Logger())

	r.Static("/docs", "static/docs")
	r.Static("/rapidoc", "static/rapidoc")

	return r
}

func initRoutes(r *echo.Echo, redisConn *redis.Client, statisticsService *services.StatisticsService) {
	ipController := controllers.IPController{
		RedisConn:         redisConn,
		StatisticsService: statisticsService,
	}
	r.POST("/traceip", ipController.Trace)

	statisticsController := controllers.StatisticsController{
		RedisConn: redisConn,
	}
	r.GET("/statistics", statisticsController.Get)
}
