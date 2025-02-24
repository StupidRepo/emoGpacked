package models

import (
	"encoding/json"
)

type Error struct {
	Location string `json:"location"`
	Type     string `json:"type"`
}

func (e Error) MakeErrorFromResponse(body []byte) (error, *Error) {
	err := json.Unmarshal(body, &e)
	if err != nil {
		return err, nil
	}

	return nil, &e
}
