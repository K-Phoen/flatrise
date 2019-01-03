package main

import (
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/K-Phoen/flatrise/flatrise/scheduler"
	"github.com/jaffee/commandeer"
	"github.com/pkg/errors"
)

type App struct {
	Cron   bool   `help:"Launch in cron-mode"`
	Engine string `help:"Trigger a search for the given engine"`
}

func (app App) Run() error {
	rabbitMqUrl, ok := os.LookupEnv("RABBITMQ_URL")
	if !ok {
		return errors.New("The DSN to use to connect to RabbitMq must be specified by the RABBITMQ_URL environment variable.")
	}

	if !app.Cron && len(app.Engine) == 0 {
		return errors.New("Choose between cron-mode or search-mode (see -h)")
	}

	if app.Cron && len(app.Engine) > 0 {
		return errors.New("Can not launch both a cron and trigger a search.")
	}

	queuesMap := map[string]string{
		"leboncoin":   "leboncoin_searchs",
		"boligportal": "boligportal_searchs",
		"blocket":     "blocket_searchs",
	}

	s := scheduler.NewGoCronScheduler(rabbitMqUrl, queuesMap)

	if app.Cron {
		s.Run()
	} else {
		return s.Search(app.Engine)
	}

	return nil
}

func main() {
	err := commandeer.Run(&App{})
	if err != nil {
		log.Fatal(err)
	}
}
