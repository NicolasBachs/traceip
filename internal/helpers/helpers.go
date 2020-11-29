package helpers

import (
	"math"
	"strconv"
	"time"
)

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

//DistanceToBuenosAires calculate the distance of LatLng to Buenos Aires, return unit is meters
func DistanceToBuenosAires(lat, lon float64) float64 {

	var la1, lo1, la2, lo2, r float64
	la1 = -34.0 * math.Pi / 180
	lo1 = -64.0 * math.Pi / 180
	la2 = lat * math.Pi / 180
	lo2 = lon * math.Pi / 180

	r = 6378100 // Earth radius in meters

	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

//GetTime from timezone in UTC-XX:00 format
func GetTime(tz string) string {
	currentTime := time.Now().UTC()
	if tz[0:3] == "UTC" {
		if len(tz) > 3 {
			factor := 1
			if tz[3:4] == "-" {
				factor = -1
			}
			h, err := strconv.Atoi(tz[4:6])
			if err == nil {
				currentTime = currentTime.Add(time.Duration(h*factor) * time.Hour)
			}
		}
	}
	return currentTime.Format("2006-01-02 15:04:05")
}

//ErrorResponse error response format
func ErrorResponse(err string) (response map[string]interface{}) {
	return map[string]interface{}{
		"Error": err,
	}
}
