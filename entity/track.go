package entity

// Track 物流轨迹
type Track struct {
	Event          string `json:"event"`
	Description    string `json:"description"`
	LocalTime      string `json:"localTime"`
	LocalGmtOffset string `json:"localGmtOffset"`
	Location       string `json:"location"`
	Iso3166Cc      string `json:"iso3166Cc"`
	Iso3166Sc      string `json:"iso3166Sc"`
	CityUppercase  string `json:"cityUppercase"`
	PostalCode     string `json:"postalCode"`
	PodImageCount  int    `json:"podImageCount"`
}
