package models

import "time"

type Url struct {
	Url        string
	Clicks     int
	Expiration time.Time
	Hash       string
}
