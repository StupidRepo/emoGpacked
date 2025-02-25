package models

import "time"

type DefaultOptions struct {
	Token    string `json:"token"`
	ServerId string `json:"server_id"`

	LastStoredToken time.Time `json:"last_stored_token"`
}
