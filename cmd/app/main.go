package main

import (
    "log"
    "time"

    "saucer_api/config"
    "saucer_api/api"
)

func main() {
    // Catch any unexpected panic.
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered from panic: %v", r)
            time.Sleep(5 * time.Second)
            main() // Restart the main function.
        }
    }()

    // Expected errors are handled here.
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