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