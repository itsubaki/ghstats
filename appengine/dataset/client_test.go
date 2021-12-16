package dataset_test

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghstats/appengine/dataset"
)

func TestCreateIfNotExists(t *testing.T) {
	if _, err := os.Stat("../../credentials.json"); os.IsNotExist(err) {
		return
	}

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../../credentials.json")

	cases := []struct {
		name string
		meta bigquery.TableMetadata
	}{
		{"test", dataset.CommitsMeta},
		{"test", dataset.PullReqsMeta},
		{"test", dataset.PullReqCommitsMeta},
		{"test", dataset.WorkflowRunsMeta},
		{"test", dataset.WorkflowJobsMeta},
	}

	for _, c := range cases {
		ctx := context.Background()
		client, err := dataset.New(ctx)
		if err != nil {
			t.Fatalf("new bigquery client: %v", err)
		}

		if err := client.CreateIfNotExists(ctx, c.name, []bigquery.TableMetadata{c.meta}); err != nil {
			t.Errorf("create if not exists: %v", err)
		}
	}
}
