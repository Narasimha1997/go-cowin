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

// ErrorType the API response returned by /certificate/public/download upon error
type ErrorType struct {
	ErrorCode string `json:"errorCode"`
	ErrString string `json:"error"`
}
