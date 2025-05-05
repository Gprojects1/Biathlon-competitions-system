package model

import "time"

type Event struct {
	Time         time.Time
	EventID      int
	CompetitorID int
	ExtraParams  string
}
