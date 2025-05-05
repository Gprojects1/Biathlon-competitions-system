package biathlon

import (
	"Biathlon-competitions-system/config"
	"Biathlon-competitions-system/model"
	"Biathlon-competitions-system/utils"
	"fmt"
	"sort"
	"time"
)

type BiathlonRace interface {
	ProcessGame()
}

type PursuitRace struct {
	config      config.Config
	events      []model.Event
	competitors map[int]*model.Competitor
	results     []model.Result
}

func NewPursuitRace(cfg config.Config, events []model.Event) *PursuitRace {
	return &PursuitRace{
		config:      cfg,
		events:      events,
		competitors: make(map[int]*model.Competitor),
	}
}

func (purs *PursuitRace) ProcessGame() {
	for _, event := range purs.events {
		if _, exists := purs.competitors[event.CompetitorID]; !exists {
			purs.competitors[event.CompetitorID] = &model.Competitor{
				ID:        event.CompetitorID,
				Hits:      make([]int, purs.config.FiringLines*purs.config.Laps),
				Shots:     make([]int, purs.config.FiringLines*purs.config.Laps),
				LapTimes:  make([]time.Duration, 0),
				LapSpeeds: make([]float64, 0),
			}
		}
	}

	for _, event := range purs.events {
		comp := purs.competitors[event.CompetitorID]
		comp.LastEventTime = event.Time

		switch event.EventID {
		case 1:
			comp.Registered = true
			fmt.Printf("[%s] Участник %d зарегистрирован\n", event.Time.Format("15:04:05.000"), comp.ID)

		case 2:
			startTime, err := time.Parse("15:04:05.000", event.ExtraParams)
			if err != nil {
				fmt.Printf("Ошибка парсинга времени старта: %v\n", err)
				continue
			}
			comp.ScheduledStart = startTime
			fmt.Printf("[%s] Установлено время старта для участника %d: %s\n",
				event.Time.Format("15:04:05.000"), comp.ID, event.ExtraParams)

		case 3:
			fmt.Printf("[%s] Участник %d на стартовой линии\n", event.Time.Format("15:04:05.000"), comp.ID)

		case 4:
			comp.ActualStart = event.Time
			comp.LastLapTime = event.Time
			fmt.Printf("[%s] Участник %d начал гонку\n", event.Time.Format("15:04:05.000"), comp.ID)

		case 5:
			var rangeNum int
			fmt.Sscanf(event.ExtraParams, "%d", &rangeNum)
			comp.CurrentRange = rangeNum - 1
			fmt.Printf("[%s] Участник %d на стрелковом рубеже %d\n",
				event.Time.Format("15:04:05.000"), comp.ID, rangeNum)

		case 6:
			var target int
			fmt.Sscanf(event.ExtraParams, "%d", &target)
			comp.Shots[comp.CurrentRange]++
			comp.Hits[comp.CurrentRange]++
			fmt.Printf("[%s] Участник %d поразил мишень %d на рубеже %d\n",
				event.Time.Format("15:04:05.000"), comp.ID, target, comp.CurrentRange+1)

		case 7:
			fmt.Printf("[%s] Участник %d покинул стрелковый рубеж %d\n",
				event.Time.Format("15:04:05.000"), comp.ID, comp.CurrentRange+1)

		case 8:
			comp.PenaltyStartTime = event.Time
			comp.InPenalty = true
			fmt.Printf("[%s] Участник %d начал штрафные круги\n",
				event.Time.Format("15:04:05.000"), comp.ID)

		case 9:
			if comp.InPenalty {
				penaltyTime := event.Time.Sub(comp.PenaltyStartTime)
				comp.PenaltyTime += penaltyTime
				comp.InPenalty = false
				fmt.Printf("[%s] Участник %d завершил штрафные круги (время: %s)\n",
					event.Time.Format("15:04:05.000"), comp.ID, utils.FormatDuration(penaltyTime))
			}

		case 10:
			comp.CurrentLap++
			lapTime := event.Time.Sub(comp.LastLapTime)
			comp.LapTimes = append(comp.LapTimes, lapTime)

			speed := float64(purs.config.LapLen) / lapTime.Seconds()
			comp.LapSpeeds = append(comp.LapSpeeds, speed)

			comp.LastLapTime = event.Time
			fmt.Printf("[%s] Участник %d завершил круг %d (время: %s, скорость: %.1f м/с)\n",
				event.Time.Format("15:04:05.000"), comp.ID, comp.CurrentLap,
				utils.FormatDuration(lapTime), speed)

		case 11:
			comp.NotFinished = true
			comp.Comment = event.ExtraParams
			fmt.Printf("[%s] Участник %d не может продолжить: %s\n",
				event.Time.Format("15:04:05.000"), comp.ID, event.ExtraParams)
		}
	}

	purs.prepareResults()
	purs.printResultsTable()
}

func (purs *PursuitRace) prepareResults() {
	for _, comp := range purs.competitors {
		result := model.Result{
			ID: comp.ID,
		}

		switch {
		case comp.Disqualified:
			result.Status = "Disqualified"
			result.TotalTime = "DQ"
		case comp.NotFinished:
			result.Status = "NotFinished"
			if !comp.ActualStart.IsZero() {
				totalTime := comp.LastEventTime.Sub(comp.ActualStart)
				result.TotalTime = utils.FormatDuration(totalTime)
			} else {
				result.TotalTime = "DNS"
			}
		case !comp.ActualStart.IsZero():
			result.Status = "Finished"
			totalTime := comp.LastEventTime.Sub(comp.ActualStart)
			result.TotalTime = utils.FormatDuration(totalTime)
		default:
			result.Status = "NotStarted"
			result.TotalTime = "DNS"
		}

		totalHits := 0
		totalShots := 0
		for i := range comp.Hits {
			totalHits += comp.Hits[i]
			totalShots += comp.Shots[i]
		}
		result.HitRate = fmt.Sprintf("%d/%d", totalHits, totalShots)

		for i, lapTime := range comp.LapTimes {
			result.LapTimes = append(result.LapTimes, utils.FormatDuration(lapTime))
			result.LapSpeeds = append(result.LapSpeeds, comp.LapSpeeds[i])
		}

		if comp.PenaltyTime > 0 {
			result.PenaltyTime = utils.FormatDuration(comp.PenaltyTime)
			result.PenaltySpeed = float64(purs.config.PenaltyLen) / comp.PenaltyTime.Seconds()
		}

		purs.results = append(purs.results, result)
	}

	sort.Slice(purs.results, func(i, j int) bool {
		return purs.results[i].TotalTime < purs.results[j].TotalTime
	})
}

func (purs *PursuitRace) printResultsTable() {
	fmt.Println("\nИтоговые результаты:")
	fmt.Println("----------------------------------------------------------------------------------------------------------")
	fmt.Println("| ID | Статус       | Общее время | Попадания | Круги (время/скорость)         | Штрафные круги          |")
	fmt.Println("----------------------------------------------------------------------------------------------------------")

	for _, res := range purs.results {
		lapsInfo := ""
		for i := 0; i < len(res.LapTimes); i++ {
			if i > 0 {
				lapsInfo += ", "
			}
			lapsInfo += fmt.Sprintf("%s (%.1f м/с)", res.LapTimes[i], res.LapSpeeds[i])
		}

		penaltyInfo := ""
		if res.PenaltyTime != "" {
			penaltyInfo = fmt.Sprintf("%s (%.1f м/с)", res.PenaltyTime, res.PenaltySpeed)
		} else {
			penaltyInfo = "Нет"
		}

		fmt.Printf("| %2d | %-12s | %-12s | %-9s | %-30s | %-20s |\n",
			res.ID, res.Status, res.TotalTime, res.HitRate, lapsInfo, penaltyInfo)
	}
	fmt.Println("----------------------------------------------------------------------------------------------------------")
}
