package config

import (
    "context"
    "log"

    "cloud.google.com/go/bigquery"
)

var BQClient *bigquery.Client

func InitBigQuery() error {
    ctx := context.Background()
    var err error
    BQClient, err = bigquery.NewClient(ctx, "repos-450312")
    if err != nil {
        return err
    }
    if err := createDataset(ctx, "repositories"); err != nil {
        log.Printf("Dataset creation: %v", err)
    }
    if err := createTable(ctx, "repositories", "repos"); err != nil {
        log.Printf("Table creation: %v", err)
    }
    return nil
}


func ReinitBigQuery() {
    ctx := context.Background()
    var err error

    BQClient, err = bigquery.NewClient(ctx, "repos-450312")
    if err != nil {
        log.Fatal("Failed to connect to BigQuery:", err)
    }

    if err := BQClient.Dataset("repositories").DeleteWithContents(ctx); err != nil {
        log.Println("Failed to delete dataset:", err)
    }

    if err := createDataset(ctx, "repositories"); err != nil {
        log.Println("Dataset already exists:", err)
    }

    if err := createTable(ctx, "repositories", "repos"); err != nil {
        log.Println("Table already exists:", err)
    }
}

func createDataset(ctx context.Context, datasetID string) error {
    dataset := BQClient.Dataset(datasetID)
    if err := dataset.Create(ctx, nil); err != nil {
        return err
    }
    return nil
}

func createTable(ctx context.Context, datasetID, tableID string) error {
    table := BQClient.Dataset(datasetID).Table(tableID)
    if err := table.Create(ctx, &bigquery.TableMetadata{Schema: RepoSchema}); err != nil {
        return err
    }
    return nil
}
