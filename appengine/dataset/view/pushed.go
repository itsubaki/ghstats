package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghz/appengine/dataset"
)

func PushedMeta(dsn string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pushed",
		ViewQuery: fmt.Sprintf(
			`
			SELECT
				A.owner,
				A.repository,
				A.id,
				A.login,
				B.message,
				A.head_sha,
				B.sha,
				B.date as committed_at,
				A.created_at as pushed_at,
				TIMESTAMP_DIFF(A.created_at, B.date, MINUTE) as duration
			FROM %v as A
			INNER JOIN %v as B
			ON A.sha = B.sha
			`,
			fmt.Sprintf("`%v.%v.%v`", dataset.ProjectID, dsn, dataset.EventsPushMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", dataset.ProjectID, dsn, dataset.CommitsMeta.Name),
		),
	}
}
