package api

import (
    "os"
    "strconv"

    "encoding/csv"
    "encoding/json"
    "log"
    "net/http"
    "repo_api/config"
    "context"
)

func StoreRepositories() {
    resp, err := http.Get(config.GITHUB_API)
    if err != nil {
        log.Println("Failed to fetch GitHub repositories:", err)
        return
    }
    defer resp.Body.Close()

    data, err := decodeResponse(resp)
    if err != nil {
        log.Println("Failed to decode response:", err)
        return
    }

    if err := saveToCSV(data.Items); err != nil {
        log.Println("Failed to save to CSV:", err)
    }

    if err := updateBigQuery(data.Items); err != nil {
        log.Println("Failed to update BigQuery:", err)
    }
}

func decodeResponse(resp *http.Response) (*struct {
    Items []config.Repository `json:"items"`
}, error) {
    var data struct {
        Items []config.Repository `json:"items"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, err
    }
    return &data, nil
}

func saveToCSV(repos []config.Repository) error {
    csvFile, err := os.Create("repositories.csv")
    if err != nil {
        return err
    }
    defer csvFile.Close()

    writer := csv.NewWriter(csvFile)
    defer writer.Flush()

    writer.Write([]string{"id", "Name", "Stars", "URL"})

    for _, repo := range repos {
        writer.Write([]string{
            strconv.Itoa(int(repo.ID)),
            repo.Name,
            strconv.Itoa(int(repo.Stars)),
            repo.URL,
        })
    }

    return nil
}

func updateBigQuery(repos []config.Repository) error {
    ctx := context.Background()
    for _, repo := range repos {
        if err := config.InsertRepositoryToBigQuery(ctx, repo); err != nil {
            return err
        }
    }
    return nil
}