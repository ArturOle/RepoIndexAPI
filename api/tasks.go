package api

import (
    "log"
    "github.com/robfig/cron/v3"
)

func StartBackgroundTask() {
    c := cron.New()

    _, err := c.AddFunc("@hourly", func() {
        log.Println("Running cron job to fetch and store repositories")
        StoreRepositories()
    })
    if err != nil {
        log.Fatal("Failed to schedule cron job:", err)
    }

    c.Start()
}