# go-cowin
Unofficial GO SDK for Indian Government's Co-WIN API Platform. This SDK helps developers to easily integrate Co-WIN APIs with their existing eco-system.

### Key Features:
1. Pre-defined functions for all the publically available CoWIN APIs.
2. Typed definitions of all request/response types.
3. Functions for authenticating with CoWIN platform.
4. Functions for Requesting, Signing (using SHA-256) and obtaining Bearer Tokens.
5. Query vaccination centers, calendars and sessions by location, district or PIN code.
6. API to download the certificate in PDF format.

### Functions available:

```go

// RequestOTP Requests an OTP to be sent to given number:
func (c *CoWinAPI) RequestOTP(mobile string) (*OTPResponse, error

// ConfirmSignedOTP Confirm an OTP sent to given number by passing SHA-256 hashed otp string:
func (c *CoWinAPI) ConfirmSignedOTP(otpSHA256 string, txnID string) (*OTPConfirmResponse, error)

// ConfirmRawOTP Confirm an OTP sent to given number by passing SHA-256 raw otp string
func (c *CoWinAPI) ConfirmRawOTP(otp string, txnID string) (*OTPConfirmResponse, error)

// GetStates get all the states and their IDs listed by CoWIN
func (c *CoWinAPI) GetStates() (*StateResp, error)

// GetDistricts Get all the districts under a state, pass stateID as the parameter
func (c *CoWinAPI) GetDistricts(stateID int) (*DistrictResp, error)

// GetSessionsByPIN Get the vaccination sessions available by PIN
func (c *CoWinAPI) GetSessionsByPIN(pincode string, date string) (*VaccinationSessionResp, error)

// GetSessionsByDistrict Get the vaccination sessions available by district ID
func (c *CoWinAPI) GetSessionsByDistrict(districtID int, date string) (*VaccinationSessionResp, error)

// GetCentersByLatLong Get the vaccination sessions available by lat long
func (c *CoWinAPI) GetCentersByLatLong(lat float64, long float64) (*VaccinationCentersResp, error)

// GetCalendarByPIN Get the calendar by PIN
func (c *CoWinAPI) GetCalendarByPIN(pincode string, date string) (*CentersCalendarResponse, error)

// GetCalendarByDistrict Get the calendar by district ID
func (c *CoWinAPI) GetCalendarByDistrict(districtID int, date string) (*CentersCalendarResponse, error) 

// GetCalendarByCenter Get the calendar by center ID
func (c *CoWinAPI) GetCalendarByCenter(centerID int, date string) (*CenterCalendar, error)

// GetCertificate Get the certificate in binary (byte-array) format
func (c *CoWinAPI) GetCertificate(beneficiaryID string, token string) ([]byte, error)

```

### How to use?
Requires Golang to be installed on the dev system.

1. Get the package:
```
go get github.com/Narasimha1997/go-cowin
```

2. Use it in your codebase:
```go

import (
	"fmt"
     // import the module
	"github.com/Narasimha1997/go-cowin/core"
)

func main() {
    /* create a new API instance
     parameters:
        1. language: Default:  en_US, see web compatible language strings.
        2. userAgent: Default: "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.3 Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/43.4"
                    You can pass your own.
        Pass empty strings to use defaults.
    */
    api := core.NewCoWinAPI("", "")

    // call one of the APIs:
    resp, err := api.GetCentersByLatLong(12.22, 77.12)
    // err will contain the necessary information if there is an error in the API call and resp will be nil
    // if no error, resp will be a struct of the specified type.
}

```
Sometimes, CoWIN API will return empty fields in some responses in cases where the given field is not applicable. Make sure you handle empty fields properly.

All credits goes to Government of India. Check more APIs [here](https://apisetu.gov.in/public/marketplace/api/cowin).

### TODO:
1. Write tests for all APIs.

### Contributing:
All contributions are welcome. Feel free to raise issues, make PRs or suggest new features.