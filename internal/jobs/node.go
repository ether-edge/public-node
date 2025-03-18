package jobs

import (
	"fmt"
	"public-node/internal/handlers"

	"gopkg.in/robfig/cron.v2"
)


func RunCronJobsStarted() *cron.Cron {
	c := cron.New()

	// Schedule the task to run every 10 seconds
	_, err := c.AddFunc("*/10 * * * * *", handlers.SendPostRequest)
	if err != nil {
		fmt.Println("Error scheduling cron job:", err)
		return nil
	}

	// Start the cron scheduler
	c.Start()
	fmt.Println("Cron started...")

	return c
}