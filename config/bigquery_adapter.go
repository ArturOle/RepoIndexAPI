package config

import (
	"context"
    "time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func InsertRepositoryToBigQuery(ctx context.Context, repo Repository) error {
    inserter := BQClient.Dataset("repositories").Table("repos").Inserter()
    if err := inserter.Put(ctx, repo); err != nil {
        return err
    }
    return nil
}

func BatchUpsertRepositoriesToBigQuery(ctx context.Context, repos []Repository) error {
    queryStr := `
        MERGE ` + "`repos-450312.repositories.repos`" + ` T
        USING UNNEST(@rows) AS S
        ON T.ID = S.ID
        WHEN MATCHED THEN
          UPDATE SET Name = S.Name, Stars = S.Stars, URL = S.URL, Description = S.Description
        WHEN NOT MATCHED THEN
          INSERT (ID, Name, Stars, URL, Description)
          VALUES (S.ID, S.Name, S.Stars, S.URL, S.Description)
    `

    type RepoRow struct {
        ID          int64  `bigquery:"ID"`
        Name        string `bigquery:"Name"`
        Stars       int64  `bigquery:"Stars"`
        URL         string `bigquery:"URL"`
        Description string `bigquery:"Description"`
    }
    
    rows := make([]RepoRow, len(repos))
    for i, r := range repos {
        rows[i] = RepoRow{
            ID:          int64(r.ID),
            Name:        r.Name,
            Stars:       int64(r.Stars),
            URL:         r.URL,
            Description: r.Description,
        }
    }

    query := BQClient.Query(queryStr)
    query.Parameters = []bigquery.QueryParameter{
        {
            Name:  "rows",
            Value: rows,
        },
    }
    job, err := query.Run(ctx)
    if err != nil {
        return err
    }
    status, err := job.Wait(ctx)
    if err != nil {
        return err
    }
    if err := status.Err(); err != nil {
        return err
    }
    return nil
}

func UpsertRepositoryToBigQuery(ctx context.Context, repo Repository) error {
    queryStr := `
        MERGE ` + "`repos-450312.repositories.repos`" + ` T
        USING (
            SELECT @id AS ID,
                   @name AS Name,
                   @stars AS Stars,
                   @url AS URL,
                   @description AS Description
        ) S
        ON T.ID = S.ID
        WHEN MATCHED THEN
          UPDATE SET Name = S.Name, Stars = S.Stars, URL = S.URL, Description = S.Description
        WHEN NOT MATCHED THEN
          INSERT (ID, Name, Stars, URL, Description)
          VALUES (S.ID, S.Name, S.Stars, S.URL, S.Description)
    `
    query := BQClient.Query(queryStr)
    query.Parameters = []bigquery.QueryParameter{
        {Name: "id", Value: int64(repo.ID)},
        {Name: "name", Value: repo.Name},
        {Name: "stars", Value: int64(repo.Stars)},
        {Name: "url", Value: repo.URL},
        {Name: "description", Value: repo.Description},
    }
    job, err := query.Run(ctx)
    if err != nil {
        return err
    }
    status, err := job.Wait(ctx)
    if err != nil {
        return err
    }
    if status.Err() != nil {
        return status.Err()
    }
    return nil
}

func FetchRepositoriesFromBigQuery(ctx context.Context) ([]Repository, error) {
    query := BQClient.Query("SELECT * FROM `repos-450312.repositories.repos`")
    it, err := query.Read(ctx)
    if err != nil {
        return nil, err
    }

    var repos []Repository
    for {
        var row map[string]bigquery.Value
        err := it.Next(&row)
        if err == iterator.Done {
            break
        }
        if err != nil {
            return nil, err
        }
        var repo Repository
        repo.FromBigQueryRow(row)
        repos = append(repos, repo)
    }
    return repos, nil
}

func InsertRepoHistoryToBigQuery(ctx context.Context, repoHistory RepoHistory) error {
    inserter := BQClient.Dataset("repositories").Table("repo_history").Inserter()
    if err := inserter.Put(ctx, repoHistory); err != nil {
        return err
    }
    return nil
}

func BatchInsertRepoHistoryToBigQuery(ctx context.Context, repoHistory []RepoHistory) error {
    inserter := BQClient.Dataset("repositories").Table("repo_history").Inserter()
    if err := inserter.Put(ctx, repoHistory); err != nil {
        return err
    }
    return nil
}

func BatchUpsertRepoHistoriesToBigQuery(ctx context.Context, histories []RepoHistory) error {
    queryStr := `
        MERGE ` + "`repos-450312.repositories.repo_history`" + ` T
        USING UNNEST(@rows) AS S
        ON T.RepoID = S.RepoID AND T.CreatedAt = S.CreatedAt
        WHEN MATCHED THEN
          UPDATE SET Stars = S.Stars
        WHEN NOT MATCHED THEN
          INSERT (HistoryID, RepoID, Stars, CreatedAt)
          VALUES (S.HistoryID, S.RepoID, S.Stars, S.CreatedAt)
    `

    type HistoryRow struct {
        HistoryID int64  `bigquery:"HistoryID"`
        RepoID    int64  `bigquery:"RepoID"`
        Stars     int64  `bigquery:"Stars"`
        CreatedAt string `bigquery:"CreatedAt"`
    }

    rows := make([]HistoryRow, len(histories))
    for i, h := range histories {
        rows[i] = HistoryRow{
            HistoryID: time.Now().UnixNano(),
            RepoID:    int64(h.RepoID),
            Stars:     int64(h.Stars),
            CreatedAt: h.CreatedAt,
        }
    }

    query := BQClient.Query(queryStr)
    query.Parameters = []bigquery.QueryParameter{
        {
            Name:  "rows",
            Value: rows,
        },
    }

    job, err := query.Run(ctx)
    if err != nil {
        return err
    }
    status, err := job.Wait(ctx)
    if err != nil {
        return err
    }
    if status.Err() != nil {
        return status.Err()
    }
    return nil
}

func FetchRepoHistoriesFromBigQuery(ctx context.Context) ([]RepoHistory, error) {
    query := BQClient.Query("SELECT * FROM `repos-450312.repositories.repo_history`")
    it, err := query.Read(ctx)
    if err != nil {
        return nil, err
    }

    var repoHistories []RepoHistory
    for {
        var row map[string]bigquery.Value
        err := it.Next(&row)
        if err == iterator.Done {
            break
        }
        if err != nil {
            return nil, err
        }
        var repoHistory RepoHistory
        repoHistory.FromBigQueryRow(row)
        repoHistories = append(repoHistories, repoHistory)
    }
    return repoHistories, nil
}
