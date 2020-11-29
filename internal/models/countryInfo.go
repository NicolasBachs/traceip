package models

//CountryInfo Country information
type CountryInfo struct {
	Name         string              `json:"name"`
	Alpha2Code   string              `json:"alpha2Code"`
	Alpha3Code   string              `json:"alpha3Code"`
	LatLng       []float64           `json:"latlng"`
	TimeZones    []string            `json:"timezones"`
	Currencies   []map[string]string `json:"currencies"`
	Languages    []map[string]string `json:"languages"`
	Translations map[string]string   `json:"translations"`
	Capital      string              `json:"capital"`
}
