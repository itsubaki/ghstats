package prlist

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/prstats/pkg/prstats"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	in := prstats.GetStatsInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		State:   c.String("state"),
		PerPage: c.Int("perpage"),
	}

	list, err := prstats.GetList(context.Background(), &in, time.Now(), time.Unix(0, 0))
	if err != nil {
		return fmt.Errorf("list PR: %v", err)
	}

	format := strings.ToLower(c.String("format"))
	if err := print(format, list); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func print(format string, list []*github.PullRequest) error {
	if format == "json" {
		for _, r := range list {
			fmt.Println(r)
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("id, title, created_at, merged_at, lead_time(hours), ")

		for _, r := range list {
			fmt.Printf("%v, %v, %v, %v, ", *r.ID, strings.ReplaceAll(*r.Title, ",", ""), r.CreatedAt, r.MergedAt)
			if r.MergedAt != nil {
				fmt.Printf("%.4f, ", r.MergedAt.Sub(*r.CreatedAt).Hours())
			}

			fmt.Println()
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}
