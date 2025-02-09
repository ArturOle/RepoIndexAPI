package api

import (
    "context"
    "encoding/csv"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

    "repo_api/config"
)

func StoreRepositories() {
    resp, err := http.Get(config.GITHUB_API)
    if err != nil {
        log.Println("Failed to fetch GitHub repositories:", err)
        return
    }
    defer resp.Body.Close()

    repos, err := decodeResponse(resp)
    if err != nil {
        log.Println("Failed to decode response:", err)
        return
    }

    if err := saveToCSV(repos); err != nil {
        log.Println("Failed to save to CSV:", err)
    }

    if err := updateBigQuery(repos); err != nil {
        log.Println("Failed to update BigQuery:", err)
    }
}

func decodeResponse(resp *http.Response) ([]config.Repository, error) {
    var data struct {
        Items []config.Repository `json:"items"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, err
    }
    return data.Items, nil
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

        history := config.RepoHistory{
            RepoID:    repo.ID,
            Stars:     repo.Stars,
            CreatedAt: time.Now().Format(time.RFC3339),
        }

        if err := config.InsertRepoHistoryToBigQuery(ctx, history); err != nil {
            return err
        }
    }
    return nil
}