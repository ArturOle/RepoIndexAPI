package config

import (
    "context"
    "fmt"
    "log"

    "cloud.google.com/go/bigquery"
)

var BQClient *bigquery.Client
const repoTableID = "repos"
const repoHistoryTableID = "repo_history"
const projectID = "repos-450312"
const datasetID = "repositories"


func InitBigQuery() error {
    ctx := context.Background()
    var err error
    BQClient, err = bigquery.NewClient(ctx, projectID)
    if err != nil {
        return err
    }
    if err := createDataset(ctx, datasetID); err != nil {
        log.Printf("Dataset creation: %v", err)
    }
    if err := createTable(ctx, datasetID, repoTableID); err != nil {
        log.Printf("Table creation for %q: %v", repoTableID, err)
    }
    if err := createTable(ctx, datasetID, repoHistoryTableID); err != nil {
        log.Printf("Table creation for %q: %v", repoHistoryTableID, err)
    }
    return nil
}

func ReinitBigQuery() error {
    ctx := context.Background()
    var err error

    BQClient, err = bigquery.NewClient(ctx, projectID)
    if err != nil {
        log.Fatal("Failed to connect to BigQuery:", err)
    }

    if err := BQClient.Dataset(datasetID).DeleteWithContents(ctx); err != nil {
        log.Println("Failed to delete dataset:", err)
    }

    if err := createDataset(ctx, datasetID); err != nil {
        log.Println("Dataset creation:", err)
    }

    if err := createTable(ctx, datasetID, repoTableID); err != nil {
        log.Println("Table creation for repos:", err)
    }

    if err := createTable(ctx, datasetID, repoHistoryTableID); err != nil {
        log.Println("Table creation for repo_history:", err)
    }

    return nil
}

func createDataset(ctx context.Context, datasetID string) error {
    ds := BQClient.Dataset(datasetID)

    if err := ds.Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
        log.Printf("Dataset %q creation error (possibly already exists): %v", datasetID, err)
    }

    return nil
}

func createTable(ctx context.Context, datasetID, tableID string) error {
    table := BQClient.Dataset(datasetID).Table(tableID)
    var meta *bigquery.TableMetadata

    switch tableID { 
    case repoTableID:
        meta = &bigquery.TableMetadata{
            Schema: RepoSchema,
            Clustering: &bigquery.Clustering{Fields: []string{"ID"}},
        }
    case repoHistoryTableID:
        meta = &bigquery.TableMetadata{
            Schema: RepoHistorySchema,
        }
    default:
        return fmt.Errorf("unknown tableID %s", tableID)
    }

    if err := table.Create(ctx, meta); err != nil {
        return err
    }
    return nil
}