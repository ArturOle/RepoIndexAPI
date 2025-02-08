package api

import (
    "context"
    "repo_api/config"
    "cloud.google.com/go/bigquery"
    "google.golang.org/api/iterator"
)

func GetRepositories(limit int) ([]config.Repository, error) {
    ctx := context.Background()

    query := config.BQClient.Query(
        `SELECT ID, Name, Stars, URL, Description
        FROM ` + "`repos-450312.repositories.repos`" + `
        LIMIT @limit`,
    )
    query.Parameters = []bigquery.QueryParameter{
        {Name: "limit", Value: limit},
    }

    it, err := query.Read(ctx)
    if err != nil {
        return nil, err
    }

    var repos []config.Repository
    for {
        var row map[string]bigquery.Value
        err := it.Next(&row)
        if err == iterator.Done {
            break
        }
        if err != nil {
            return nil, err
        }
        var repo config.Repository
        repo.FromBigQueryRow(row)
        repos = append(repos, repo)
    }

    return repos, nil
}