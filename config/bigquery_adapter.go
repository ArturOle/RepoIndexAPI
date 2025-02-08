package config

import (
    "context"
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
