package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type Release struct {
	Owner           string    `bigquery:"owner"`
	Repository      string    `bigquery:"repository"`
	ID              int64     `bigquery:"id"`
	TagName         string    `bigquery:"tag_name"`
	TagSHA          string    `bigquery:"tag_sha"`
	Login           string    `bigquery:"login"`
	TargetCommitish string    `bigquery:"target_commitish"`
	Name            string    `bigquery:"name"`
	CreatedAt       time.Time `bigquery:"created_at"`
	PublishedAt     time.Time `bigquery:"published_at"`
}

var ReleasesMeta = bigquery.TableMetadata{
	Name: "releases",
	TimePartitioning: &bigquery.TimePartitioning{
		Type:  bigquery.MonthPartitioningType,
		Field: "created_at",
	},
	Schema: bigquery.Schema{
		{Name: "owner", Type: bigquery.StringFieldType, Required: true},
		{Name: "repository", Type: bigquery.StringFieldType, Required: true},
		{Name: "id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "tag_name", Type: bigquery.StringFieldType, Required: true},
		{Name: "tag_sha", Type: bigquery.StringFieldType, Required: true},
		{Name: "login", Type: bigquery.StringFieldType, Required: true},
		{Name: "target_commitish", Type: bigquery.StringFieldType, Required: true},
		{Name: "name", Type: bigquery.StringFieldType, Required: true},
		{Name: "created_at", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "published_at", Type: bigquery.TimestampFieldType, Required: true},
	},
}
