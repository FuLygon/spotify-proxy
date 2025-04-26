package models

type PlayerQueue struct {
	Queue []map[string]interface{} `json:"queue"`
}

type PlayerRecentTrack struct {
	Items []struct {
		Track map[string]interface{} `json:"track"`
	} `json:"items"`
	Limit int `json:"limit"`
}
