package dataset

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"cloud.google.com/go/bigquery"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
)

type Client struct {
	client    *bigquery.Client
	ProjectID string
}

func New(ctx context.Context) *Client {
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		panic(fmt.Sprintf("find default credentials: %v", err))
	}

	client, err := bigquery.NewClient(ctx, creds.ProjectID)
	if err != nil {
		panic(fmt.Sprintf("new bigquery client: %v", err))
	}

	return &Client{
		client:    client,
		ProjectID: creds.ProjectID,
	}
}

func (c *Client) CreateIfNotExists(ctx context.Context, datasetName string, meta []bigquery.TableMetadata) error {
	location := "US"
	if len(os.Getenv("DATASET_LOCATION")) > 0 {
		location = os.Getenv("DATASET_LOCATION")
	}

	if _, err := c.client.Dataset(datasetName).Metadata(ctx); err != nil {
		// not found then create dataset
		if err := c.client.Dataset(datasetName).Create(ctx, &bigquery.DatasetMetadata{
			Location: location,
		}); err != nil {
			return fmt.Errorf("create %v: %v", datasetName, err)
		}
	}

	for _, m := range meta {
		ref := c.client.Dataset(datasetName).Table(m.Name)
		if _, err := ref.Metadata(ctx); err == nil {
			// already exists
			continue
		}

		if err := ref.Create(ctx, &m); err != nil {
			return fmt.Errorf("create %v/%v: %v", datasetName, m.Name, err)
		}
	}

	return nil
}

func (c *Client) Insert(ctx context.Context, datasetName, tableName string, items []interface{}) error {
	if err := c.client.Dataset(datasetName).Table(tableName).Inserter().Put(ctx, items); err != nil {
		return fmt.Errorf("insert %v/%v: %v", datasetName, tableName, err)
	}

	return nil
}

func (c *Client) Query(ctx context.Context, query string, fn func(values []bigquery.Value)) error {
	it, err := c.client.Query(query).Read(ctx)
	if err != nil {
		return fmt.Errorf("query(%v): %v", query, err)
	}

	var values []bigquery.Value
	for {
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}

		if err != nil {
			return fmt.Errorf("iterator: %v", err)
		}

		fn(values)
	}

	return nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Raw() *bigquery.Client {
	return c.client
}

func ProjectID() string {
	creds, err := google.FindDefaultCredentials(context.Background())
	if err != nil {
		panic(fmt.Sprintf("find default credentials: %v", err))
	}

	return creds.ProjectID
}

func CreateIfNotExists(ctx context.Context, datasetName string, meta []bigquery.TableMetadata) error {
	client := New(ctx)
	defer client.Close()

	return client.CreateIfNotExists(ctx, datasetName, meta)
}

func Insert(ctx context.Context, datasetName, tableName string, items []interface{}) error {
	client := New(ctx)
	defer client.Close()

	return client.Insert(ctx, datasetName, tableName, items)
}

func Query(ctx context.Context, query string, fn func(values []bigquery.Value)) error {
	client := New(ctx)
	defer client.Close()

	return client.Query(ctx, query, fn)
}

var invalid = regexp.MustCompile(`[!?"'#$%&@\+\-\*/=~^;:,.|()\[\]{}<>]`)

func Name(owner, repository string) string {
	own := invalid.ReplaceAllString(owner, "_")
	rep := invalid.ReplaceAllString(repository, "_")
	return fmt.Sprintf("%v_%v", own, rep)
}
