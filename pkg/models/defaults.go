package models

import "time"

type DefaultOptions struct {
	Token string `json:"token"`

	LastStoredToken time.Time `json:"last_stored_token"`
}
