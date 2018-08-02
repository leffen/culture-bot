package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/plutov/culture-bot/dashboard/pkg/handlers"
	"github.com/plutov/culture-bot/dashboard/pkg/stats"
)

func main() {
	stats.Init()

	if err := handlers.ListenAndServe(); err != nil {
		log.WithError(err).Error("unable to start server")
	}
}
