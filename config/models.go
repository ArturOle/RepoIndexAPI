package config

import (
    "cloud.google.com/go/bigquery"
    "gorm.io/gorm"
)

type Repository struct {
    gorm.Model
    Name        string `json:"name"`
    Stars       int    `json:"stargazers_count"`
    URL         string `json:"url"`
    Description string `json:"description"`
}

func (repo Repository) ToBigQueryRow() map[string]bigquery.Value {
    return map[string]bigquery.Value{
        "ID":          repo.ID,
        "Name":        repo.Name,
        "Stars":       repo.Stars,
        "URL":         repo.URL,
        "Description": repo.Description,
    }
}

func (repo *Repository) FromBigQueryRow(row map[string]bigquery.Value) {
    repo.ID = uint(row["ID"].(int64))
    repo.Name = row["Name"].(string)
    repo.Stars = int(row["Stars"].(int64))
    repo.URL = row["URL"].(string)
    repo.Description = row["Description"].(string)
}

func (repo Repository) Save() (map[string]bigquery.Value, string, error) {
    return repo.ToBigQueryRow(), "", nil
}

type RepoHistory struct {
    gorm.Model
    RepoID uint `json:"repo_id"`
    Stars  int  `json:"stargazers_count"`
    CreatedAt string `json:"created_at"`
}

func (repoHistory RepoHistory) ToBigQueryRow() map[string]bigquery.Value {
    return map[string]bigquery.Value{
        "RepoID":    repoHistory.RepoID,
        "Stars":     repoHistory.Stars,
        "CreatedAt": repoHistory.CreatedAt,
    }
}

func (repoHistory *RepoHistory) FromBigQueryRow(row map[string]bigquery.Value) {
    repoHistory.RepoID = uint(row["RepoID"].(int64))
    repoHistory.Stars = int(row["Stars"].(int64))
    repoHistory.CreatedAt = row["CreatedAt"].(string)
}

func (repoHistory RepoHistory) Save() (map[string]bigquery.Value, string, error) {
    return repoHistory.ToBigQueryRow(), "", nil
}
