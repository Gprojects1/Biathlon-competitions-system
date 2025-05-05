package main

import (
	"Biathlon-competitions-system/biathlon"
	"Biathlon-competitions-system/config"
	"Biathlon-competitions-system/utils"
	"fmt"
	"os"
)

func main() {
	config := config.LoadConfig()
	parsedEvents, err := utils.LoadEvents()
	if err != nil {
		fmt.Printf("Ошибка парсинга событий: %v\n", err)
		os.Exit(1)
	}
	race := biathlon.NewPursuitRace(config, parsedEvents)
	race.ProcessGame()
}
