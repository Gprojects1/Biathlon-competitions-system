package utils

import (
	"Biathlon-competitions-system/model"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func LoadEvents() ([]model.Event, error) {
	eventsFile, err := ioutil.ReadFile("sunny_5_skiers/events")
	if err != nil {
		log.Fatal(err)
	}

	events := strings.Split(string(eventsFile), "\n")
	var parsedEvents []model.Event

	for _, line := range events {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		timeEnd := strings.Index(line, "]")
		if timeEnd == -1 {
			return nil, fmt.Errorf("неверный формат строки: отсутствует ']'")
		}

		timeStr := strings.TrimSpace(line[1:timeEnd])
		rest := strings.TrimSpace(line[timeEnd+1:])

		eventTime, err := time.Parse("15:04:05.000", timeStr)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга времени '%s': %v", timeStr, err)
		}

		var eventID, competitorID int
		var extraParams string

		n, err := fmt.Sscanf(rest, "%d %d %s", &eventID, &competitorID, &extraParams)
		if err != nil && n < 2 {
			if _, err := fmt.Sscanf(rest, "%d %d", &eventID, &competitorID); err != nil {
				return nil, fmt.Errorf("ошибка парсинга параметров '%s': %v", rest, err)
			}
		}

		parsedEvents = append(parsedEvents, model.Event{
			Time:         eventTime,
			EventID:      eventID,
			CompetitorID: competitorID,
			ExtraParams:  extraParams,
		})
	}

	return parsedEvents, nil
}
