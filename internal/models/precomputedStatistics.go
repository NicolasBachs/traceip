package models

//PrecomputedStatistics Result of precomputed statistics
type PrecomputedStatistics struct {
	Min        float64 `json:"min"`
	Max        float64 `json:"max"`
	Avg        float64 `json:"avg"`
	MinCountry string  `json:"minCountry"`
	MaxCountry string  `json:"maxCountry"`
}
