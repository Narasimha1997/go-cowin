package core

import (
	"errors"
	"net/http"
)

// DefaultServiceURL the default URL of the CoWIN Production API server
const DefaultServiceURL = "https://cdn-api.co-vin.in/api"

// DefaultUserAgent the default user-agent
const DefaultUserAgent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.3 Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/43.4"

// Routes all the API routes that are available in this SDK
var Routes map[string]string = map[string]string{
	"get_states":       "/v2/admin/location/states",
	"get_districts":    "/v2/admin/location/districts",
	"find_by_pin":      "/v2/appointment/sessions/public/findByPin",
	"find_by_district": "/v2/appointment/sessions/public/findByDistrict",
	"find_by_lat_lan":  "/v2/appointment/centers/public/findByLatLong",
	"cal_by_pin":       "/v2/appointment/sessions/public/calendarByPin",
	"cal_by_district":  "/v2/appointment/sessions/public/calendarByDistrict",
	"cal_by_center":    "/v2/appointment/sessions/public/calendarByCenter",
	"download_cert":    "/v2/registration/certificate/public/download",
}

func apiResponsePainc(err string) error {
	return errors.New(err)
}

// CoWinAPI The main type that implements all the CoWIN APIs
type CoWinAPI struct {
	language  string
	userAgent string
}

func NewCoWinAPI(lanugae string, userAgent string) *CoWinAPI {
	cowinApi := CoWinAPI{}
	if lanugae != "" {
		cowinApi.language = lanugae
	} else {
		cowinApi.language = "en_US"
	}

	if userAgent != "" {
		cowinApi.userAgent = userAgent
	} else {
		cowinApi.userAgent = DefaultUserAgent
	}

	return &cowinApi
}

func (c *CoWinAPI) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept-Language", c.language)
}

func (c *CoWinAPI) GetStates() {

}
