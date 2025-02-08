package main

import (
    "log"
    "time"

    "repo_api/config"
    "repo_api/api"
)

func main() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered from panic: %v", r)
            time.Sleep(5 * time.Second)
            main()
        }
    }()

    if err := config.InitBigQuery(); err != nil {
        log.Fatalf("Failed to initialize BigQuery: %v", err)
    }

    api.StartBackgroundTask()

    router := api.SetupRouter()
    log.Println("Server is running on :8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}