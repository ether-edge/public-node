package jobs

import (
	"fmt"
	"public-node/internal/handlers"
	"sync"

	"gopkg.in/robfig/cron.v2"
)

var (
	cronScheduler *cron.Cron
	once          sync.Once
)

func RunCronJobsStarted() {
	once.Do(func() { 
		cronScheduler = cron.New()

		cronScheduler.AddFunc("@every 1m", func() {
			response, err := handlers.SendPostRequest()
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Response:", response)
			}
		})

		cronScheduler.Start()
		fmt.Println("Cron jobs started.")
	})
}
