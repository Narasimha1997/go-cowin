package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	client    *http.Client
}

// NewCoWinAPI Creates and returns a new CoWinAPI given language and userAgent string
func NewCoWinAPI(language string, userAgent string) *CoWinAPI {
	coWinAPI := CoWinAPI{}
	if language != "" {
		coWinAPI.language = language
	} else {
		coWinAPI.language = "en_US"
	}

	if userAgent != "" {
		coWinAPI.userAgent = userAgent
	} else {
		coWinAPI.userAgent = DefaultUserAgent
	}

	coWinAPI.client = &http.Client{}

	return &coWinAPI
}

func (c *CoWinAPI) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept-Language", c.language)
}

func (c *CoWinAPI) handleErrorStatusCode(resp *http.Response) error {
	status := resp.StatusCode
	switch status {
	case 400:
		return errors.New("Bad Request")
	case 500:
		return errors.New("Internal Server Error")
	case 401:
		return errors.New("Unauthenticated Access")
	default:
		return nil
	}
}

func (c *CoWinAPI) getter(routeCode string, urlExt string, queryParams map[string]string) ([]byte, error) {
	routeURI, _ := Routes[routeCode]
	url := fmt.Sprintf("%s%s", DefaultServiceURL, routeURI)
	if urlExt != "" {
		url = fmt.Sprintf("%s%s", url, urlExt)
	}

	// make post:
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// set headers
	c.setHeaders(req)

	// set Query parameters:
	if len(queryParams) > 0 {
		q := req.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}

		req.URL.RawQuery = q.Encode()
	}

	// make the GET request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// read response:
	err = c.handleErrorStatusCode(resp)
	if err != nil {
		return nil, err
	}

	// now read the body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetStates get all the states and their IDs listed by CoWIN
func (c *CoWinAPI) GetStates() (*StateResp, error) {
	body, err := c.getter("get_states", "", map[string]string{})
	if err != nil {
		return nil, err
	}

	// serialize the body:
	states := StateResp{}
	err = json.Unmarshal(body, &states)
	if err != nil {
		return nil, err
	}

	return &states, nil
}

// GetDistricts Get all the districts under a state, pass stateID as the parameter
func (c *CoWinAPI) GetDistricts(stateID int) (*DistrictResp, error) {
	body, err := c.getter(
		"get_districts", fmt.Sprintf("/%d", stateID), map[string]string{},
	)
	if err != nil {
		return nil, err
	}

	// de-serialize the body:
	districts := DistrictResp{}
	err = json.Unmarshal(body, &districts)
	if err != nil {
		return nil, err
	}

	return &districts, nil
}

// GetSessionsByPIN Get the vaccination sessions available by PIN
func (c *CoWinAPI) GetSessionsByPIN(pincode string, date string) (*VaccinationSessionResp, error) {
	body, err := c.getter("find_by_pin", "", map[string]string{
		"pincode": pincode,
		"date":    date,
	})

	if err != nil {
		return nil, err
	}

	// de-serialize the body:
	vaccinationResp := VaccinationSessionResp{}
	err = json.Unmarshal(body, &vaccinationResp)
	if err != nil {
		return nil, err
	}

	return &vaccinationResp, err
}

// GetSessionsByDistrict Get the vaccination sessions available by district ID
func (c *CoWinAPI) GetSessionsByDistrict(districtID int, date string) (*VaccinationSessionResp, error) {
	body, err := c.getter("find_by_pin", "", map[string]string{
		"district_id": fmt.Sprintf("%d", districtID),
		"date":        date,
	})

	if err != nil {
		return nil, err
	}

	// de-serialize the body:
	vaccinationResp := VaccinationSessionResp{}
	err = json.Unmarshal(body, &vaccinationResp)
	if err != nil {
		return nil, err
	}

	return &vaccinationResp, err
}

// GetCertificate Get the certificate in binary blob format
func (c *CoWinAPI) GetCertificate(beneficiaryID string) ([]byte, error) {
	body, err := c.getter(
		"download_cert", "", map[string]string{
			"beneficiary_reference_id": beneficiaryID,
		},
	)

	if err != nil {
		return nil, err
	}

	return body, nil
}
