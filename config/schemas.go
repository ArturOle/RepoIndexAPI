package config

import "cloud.google.com/go/bigquery"

// Define the schema for the "repos" table
var RepoSchema = bigquery.Schema{
    &bigquery.FieldSchema{Name: "ID", Type: bigquery.IntegerFieldType, Required: true},
    &bigquery.FieldSchema{Name: "Name", Type: bigquery.StringFieldType, Required: true},
    &bigquery.FieldSchema{Name: "Stars", Type: bigquery.IntegerFieldType, Required: true},
    &bigquery.FieldSchema{Name: "URL", Type: bigquery.StringFieldType, Required: true},
    &bigquery.FieldSchema{Name: "Description", Type: bigquery.StringFieldType, Required: false},
}
