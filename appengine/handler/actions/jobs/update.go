package jobs

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/pkg/actions/jobs"
)

func Update(c *gin.Context) {
	owner := c.Param("owner")
	repository := c.Param("repository")
	traceID := c.GetString("trace_id")

	ctx := context.Background()
	dsn := dataset.Name(owner, repository)
	log := logf.New(traceID, c.Request)

	list, err := ListJobs(ctx, projectID, dsn)
	if err != nil {
		log.ErrorReport("list jobs: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Debug("len(jobs)=%v", len(list))

	for _, j := range list {
		job, err := jobs.Get(ctx, &jobs.GetInput{
			Owner:      owner,
			Repository: repository,
			PAT:        os.Getenv("PAT"),
			JobID:      j.JobID,
		})
		if err != nil {
			log.ErrorReport("get jobID=%v: %v", j.JobID, err)
			c.AbortWithStatus(http.StatusInternalServerError)
			continue
		}

		if err := UpdateJob(ctx, projectID, dsn, job); err != nil {
			log.Info("update job(%v): %v", j.JobID, err)
			continue
		}
		log.Debug("updated. jobID=%v", j.JobID)
	}

	c.JSON(http.StatusOK, gin.H{
		"path": c.Request.URL.Path,
	})
}

func ListJobs(ctx context.Context, projectID, dsn string) ([]dataset.WorkflowJob, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.WorkflowJobsMeta.Name)
	query := fmt.Sprintf("select job_id from `%v` where status != \"completed\"", table)

	out := make([]dataset.WorkflowJob, 0)
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 1 {
			return
		}

		if values[0] == nil {
			return
		}

		out = append(out, dataset.WorkflowJob{
			JobID: values[0].(int64),
		})
	}); err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	return out, nil
}

func UpdateJob(ctx context.Context, projectID, dsn string, j *github.WorkflowJob) error {
	if j.GetStatus() != "completed" {
		return nil
	}

	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.WorkflowJobsMeta.Name)
	query := fmt.Sprintf("update %v set status = \"%v\", conclusion = \"%v\", completed_at = \"%v\" where job_id = %v",
		table,
		j.GetStatus(),
		j.GetConclusion(),
		j.GetCompletedAt().Format("2006-01-02 15:04:05 UTC"),
		j.GetID(),
	)

	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		return
	}); err != nil {
		return fmt.Errorf("query(%v): %v", query, err)
	}

	return nil
}
