package model

type Result struct {
	ID           int
	Status       string
	TotalTime    string
	HitRate      string
	LapTimes     []string
	LapSpeeds    []float64
	PenaltyTime  string
	PenaltySpeed float64
}
