package core

import (
	"bytes"
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
	"request_otp":      "/v2/auth/public/generateOTP",
	"confirm_otp":      "/v2/auth/public/confirmOTP",
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

func (c *CoWinAPI) setHeaders(req *http.Request, optHeaders map[string]string) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept-Language", c.language)

	if len(optHeaders) != 0 {
		for headerName, headerValue := range optHeaders {
			req.Header.Set(headerName, headerValue)
		}
	}
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

func (c *CoWinAPI) getter(routeCode string, urlExt string, queryParams map[string]string, optHeaders map[string]string) ([]byte, error) {
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
	c.setHeaders(req, optHeaders)

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

	defer resp.Body.Close()

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

func (c *CoWinAPI) poster(routeCode string, urlExt string, postData []byte, optHeaders map[string]string) ([]byte, error) {
	routeURI, _ := Routes[routeCode]
	url := fmt.Sprintf("%s%s", DefaultServiceURL, routeURI)
	if urlExt != "" {
		url = fmt.Sprintf("%s%s", url, urlExt)
	}

	// Make POST request:
	optHeaders["Content-Type"] = "application/json"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
	if err != nil {
		return nil, err
	}

	c.setHeaders(req, optHeaders)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// read response:
	err = c.handleErrorStatusCode(resp)
	if err != nil {
		return nil, err
	}

	// read the body into a byte array
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// RequestOTP Requests an OTP to be sent to given number:
func (c *CoWinAPI) RequestOTP(mobile string) (*OTPResponse, error) {
	otpRequest := OTPRequest{Mobile: mobile}
	payload, _ := json.Marshal(&otpRequest)

	body, err := c.poster("request_otp", "", payload, map[string]string{})
	if err != nil {
		return nil, err
	}

	otpResponse := OTPResponse{}
	err = json.Unmarshal(body, &otpResponse)
	if err != nil {
		return &otpResponse, err
	}

	return &otpResponse, nil
}

// ConfirmSignedOTP Confirm an OTP sent to given number by passing SHA-256 hashed otp string:
func (c *CoWinAPI) ConfirmSignedOTP(otpSHA256 string, txnID string) (*OTPConfirmResponse, error) {
	otpConfirmRequest := OTPConfirmRequest{OTP: otpSHA256, TxnID: txnID}
	payload, _ := json.Marshal(&otpConfirmRequest)

	body, err := c.poster("confirm_otp", "", payload, map[string]string{})
	if err != nil {
		return nil, err
	}

	otpConfirmResponse := OTPConfirmResponse{}
	err = json.Unmarshal(body, &otpConfirmResponse)
	if err != nil {
		return &otpConfirmResponse, err
	}

	return &otpConfirmResponse, nil
}

// ConfirmRawOTP Confirm an OTP sent to given number by passing SHA-256 raw otp string
func (c *CoWinAPI) ConfirmRawOTP(otp string, txnID string) (*OTPConfirmResponse, error) {
	otpSHA256 := SignOTP(otp)
	otpConfirmRequest := OTPConfirmRequest{OTP: otpSHA256, TxnID: txnID}
	payload, _ := json.Marshal(&otpConfirmRequest)

	body, err := c.poster("confirm_otp", "", payload, map[string]string{})
	if err != nil {
		return nil, err
	}

	otpConfirmResponse := OTPConfirmResponse{}
	err = json.Unmarshal(body, &otpConfirmResponse)
	if err != nil {
		return &otpConfirmResponse, err
	}

	return &otpConfirmResponse, nil
}

// GetStates get all the states and their IDs listed by CoWIN
func (c *CoWinAPI) GetStates() (*StateResp, error) {
	body, err := c.getter("get_states", "", map[string]string{}, map[string]string{})
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
		map[string]string{},
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
	}, map[string]string{})

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
	body, err := c.getter("find_by_district", "", map[string]string{
		"district_id": fmt.Sprintf("%d", districtID),
		"date":        date,
	}, map[string]string{})

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

// GetCentersByLatLong Get the vaccination sessions available by lat long
func (c *CoWinAPI) GetCentersByLatLong(lat float64, long float64) (*VaccinationCentersResp, error) {
	body, err := c.getter("find_by_lat_lan", "", map[string]string{
		"lat":  fmt.Sprintf("%f", lat),
		"long": fmt.Sprintf("%f", long),
	}, map[string]string{})

	if err != nil {
		return nil, err
	}

	// de-serialize the body:
	centersResp := VaccinationCentersResp{}
	err = json.Unmarshal(body, &centersResp)
	if err != nil {
		return nil, err
	}

	return &centersResp, err
}

// GetCalendarByPIN Get the calendar by PIN
func (c *CoWinAPI) GetCalendarByPIN(pincode string, date string) (*CentersCalendarResponse, error) {
	body, err := c.getter("cal_by_pin", "", map[string]string{
		"pincode": pincode,
		"date":    date,
	}, map[string]string{})

	if err != nil {
		return nil, err
	}

	calendarsResponse := CentersCalendarResponse{}
	err = json.Unmarshal(body, &calendarsResponse)
	if err != nil {
		return nil, err
	}

	return &calendarsResponse, err
}

// GetCalendarByDistrict Get the calendar by district ID
func (c *CoWinAPI) GetCalendarByDistrict(districtID int, date string) (*CentersCalendarResponse, error) {
	body, err := c.getter("cal_by_district", "", map[string]string{
		"districtID": fmt.Sprintf("%d", districtID),
		"date":       date,
	}, map[string]string{})

	if err != nil {
		return nil, err
	}

	calendarsResponse := CentersCalendarResponse{}
	err = json.Unmarshal(body, &calendarsResponse)
	if err != nil {
		return nil, err
	}

	return &calendarsResponse, err
}

// GetCalendarByCenter Get the calendar by center ID
func (c *CoWinAPI) GetCalendarByCenter(centerID int, date string) (*CenterCalendar, error) {
	body, err := c.getter("cal_by_district", "", map[string]string{
		"districtID": fmt.Sprintf("%d", centerID),
		"date":       date,
	}, map[string]string{})

	if err != nil {
		return nil, err
	}

	centerCalendar := CenterCalendar{}
	err = json.Unmarshal(body, &centerCalendar)
	if err != nil {
		return nil, err
	}

	return &centerCalendar, nil
}

// GetCertificate Get the certificate in binary blob format
func (c *CoWinAPI) GetCertificate(beneficiaryID string, token string) ([]byte, error) {
	body, err := c.getter(
		"download_cert", "", map[string]string{
			"beneficiary_reference_id": beneficiaryID,
		},
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)},
	)

	if err != nil {
		return nil, err
	}

	return body, nil
}
