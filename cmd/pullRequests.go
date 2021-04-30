/*
Copyright © 2021 Doug Hellmann <doug@doughellmann.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/dhellmann/gh-review-stats/stats"
	"github.com/dhellmann/gh-review-stats/util"

	"github.com/google/go-github/v32/github"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const dateFmt = "2006-01-02"

// pullRequestsCmd represents the pullRequests command
var pullRequestsCmd = &cobra.Command{
	Use:   "pull-requests",
	Short: "List pull requests and some characteristics in CSV format",
	Long:  `Produce a CSV list of pull requests suitable for import into a spreadsheet.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if orgName == "" {
			cobra.CheckErr(errors.New("Missing required option --org"))
		}
		if repoName == "" {
			cobra.CheckErr(errors.New("Missing required option --repo"))
		}
		if githubToken() == "" {
			cobra.CheckErr(errors.New("Missing GitHub token"))
		}

		query := &util.PullRequestQuery{
			Org:     orgName,
			Repo:    repoName,
			DevMode: devMode,
			Client:  util.NewGithubClient(context.Background(), githubToken()),
		}

		all := stats.Bucket{
			Rule: func(prd *stats.PullRequestDetails) bool {

				// FIXME: Add an option to include all PRs
				if prd.State != "merged" {
					return false
				}

				return true
			},
			Cascade: false,
		}

		earliestDate := time.Now().AddDate(0, 0, daysBack*-1)
		theStats := &stats.Stats{
			Query:        query,
			EarliestDate: earliestDate,
			Buckets:      []*stats.Bucket{&all},
		}
		err := theStats.Populate()
		if err != nil {
			return errors.Wrap(err, "could not generate stats")
		}

		out := csv.NewWriter(os.Stdout)
		out.Write([]string{
			"ID",
			"Title",
			"State",
			"Author",
			"URL",
			"Created",
			"Closed",
			"Days to Merge",
		})

		for _, prd := range all.Requests {

			var (
				createdAt, closedAt string
				daysToMerge         int = -1
			)

			if prd.Pull.CreatedAt != nil {
				createdAt = prd.Pull.CreatedAt.Format(dateFmt)
			}
			if prd.Pull.ClosedAt != nil {
				closedAt = prd.Pull.ClosedAt.Format(dateFmt)
			}
			if prd.State == "merged" && prd.Pull.CreatedAt != nil && prd.Pull.ClosedAt != nil {
				daysToMerge = int(prd.Pull.ClosedAt.Sub(*prd.Pull.CreatedAt).Hours() / 24)
			}

			user := getName(prd.Pull.User)

			out.Write([]string{
				fmt.Sprintf("%d", *prd.Pull.Number),
				*prd.Pull.Title,
				prd.State,
				user,
				*prd.Pull.HTMLURL,
				createdAt,
				closedAt,
				fmt.Sprintf("%d", daysToMerge),
			})
		}

		out.Flush()

		return nil
	},
}

func getName(user *github.User) string {
	if user == nil {
		return "unnamed"
	}
	if user.Name != nil {
		return *user.Name
	}
	if user.Login != nil {
		return *user.Login
	}
	return "unnamed"
}

func init() {
	pullRequestsCmd.Flags().StringVarP(&orgName, "org", "o", "", "github org")
	pullRequestsCmd.Flags().StringVarP(&repoName, "repo", "r", "", "github repository")
	pullRequestsCmd.Flags().IntVar(&daysBack, "days-back", 90,
		"how many days back to query, defaults to 90")

	rootCmd.AddCommand(pullRequestsCmd)
}
