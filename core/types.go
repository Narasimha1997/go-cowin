package core

// State the state type
type State struct {
	ID    int    `json:"state_id"`
	Name  string `json:"state_name"`
	NameL string `json:"state_name_l"`
}

// District the district type
type District struct {
	StateID int    `json:"state_id"`
	ID      int    `json:"district_id"`
	Name    string `json:"district_name"`
	NameL   string `json:"district_name_l"`
}

// StateResp the API response returned by /location/states
type StateResp struct {
	States []State `json:"states"`
	TTL    int     `json:"ttl"`
}

// DistrictResp the API response returned by /location/districts
type DistrictResp struct {
	Districts []District `json:"districts"`
	TTL       int        `json:"ttl"`
}

// VaccinationSession the Session type
type VaccinationSession struct {
	CenterID           int      `json:"center_id"`
	Name               string   `json:"name"`
	NameL              string   `json:"name_l"`
	Address            string   `json:"address"`
	AddressL           string   `json:"address_l"`
	StateName          string   `json:"state_name"`
	StateNameL         string   `json:"state_name_l"`
	DistrictName       string   `json:"district_name"`
	DistrictNameL      string   `json:"district_name_l"`
	BlockName          string   `json:"block_name"`
	BlockNameL         string   `json:"block_name_l"`
	Pincode            int      `json:"pincode"`
	Lat                float32  `json:"lat"`
	Long               float32  `json:"long"`
	FromTime           string   `json:"from"`
	ToTime             string   `json:"to"`
	FeeType            string   `json:"fee_type"`
	Fee                string   `json:"fee"`
	SessionID          string   `json:"session_id"`
	Date               string   `json:"date"`
	Capacity           int      `json:"available_capacity"`
	CapacityFirstDose  int      `json:"available_capacity_dose1"`
	CapacitySecondDose int      `json:"available_capacity_dose2"`
	MinimumAge         int      `json:"min_age_limit"`
	VaccineName        string   `json:"vaccine"`
	Slots              []string `json:"slots"`
}

// VaccinationSessionResp the response returned by all find* APIs
type VaccinationSessionResp struct {
	Sessions []VaccinationSession `json:"sessions"`
}

// VaccinationCenter the type that represents a vaccination center
type VaccinationCenter struct {
	CenterID      int    `json:"center_id"`
	Name          string `json:"name"`
	NameL         string `json:"name_l"`
	Address       string `json:"address"`
	AddressL      string `json:"address_l"`
	StateName     string `json:"state_name"`
	StateNameL    string `json:"state_name_l"`
	DistrictName  string `json:"district_name"`
	DistrictNameL string `json:"district_name_l"`
	BlockName     string `json:"block_name"`
	BlockNameL    string `json:"block_name_l"`
	Pincode       string `json:"pincode"`
	Lat           string `json:"lat"`
	Long          string `json:"long"`
	FromTime      string `json:"from"`
	ToTime        string `json:"to"`
	FeeType       string `json:"fee_type"`
}

// VaccinationCentersResp the type that represents vaccination response
type VaccinationCentersResp struct {
	Centers []VaccinationCenter `json:"centers"`
	TTL     int                 `json:"ttl"`
}

// ErrorType the API response returned by /certificate/public/download upon error
type ErrorType struct {
	ErrorCode string `json:"errorCode"`
	ErrString string `json:"error"`
}
