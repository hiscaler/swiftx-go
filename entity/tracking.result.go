package entity

// TrackingResult 物流跟踪结果
type TrackingResult struct {
	Result struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	} `json:"result"`
	TrackingNo        string `json:"trackingNo"`
	TrackingEventList []struct {
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
	} `json:"trackingEventList"`
}
