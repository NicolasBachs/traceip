package controllers

import (
	"net/http"
	"traceip/internal/helpers"
	"traceip/internal/services"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
)

//StatisticsController Statistics endpoints handlers
type StatisticsController struct {
	RedisConn *redis.Client
}

//Get Returns statistics
func (sc *StatisticsController) Get(c echo.Context) error {
	statisticsService := services.StatisticsService{
		RedisConn: sc.RedisConn,
	}
	precomputedStatistics, err := statisticsService.Get()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err.Error()))
	}
	if precomputedStatistics == nil {
		return c.JSON(http.StatusNoContent, nil)
	}
	return c.JSON(http.StatusOK, precomputedStatistics)
}
