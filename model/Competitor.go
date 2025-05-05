package model

import "time"

type Competitor struct {
	ID               int
	Registered       bool
	ScheduledStart   time.Time
	ActualStart      time.Time
	Finished         bool
	Disqualified     bool
	NotFinished      bool
	CurrentLap       int
	CurrentRange     int
	InPenalty        bool
	Hits             []int
	Shots            []int
	LapTimes         []time.Duration
	LapSpeeds        []float64
	PenaltyTime      time.Duration
	LastEventTime    time.Time
	LastLapTime      time.Time
	PenaltyStartTime time.Time
	Comment          string
}
