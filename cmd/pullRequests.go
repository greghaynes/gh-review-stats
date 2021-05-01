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
	"strings"
	"time"

	"github.com/dhellmann/gh-review-stats/stats"
	"github.com/dhellmann/gh-review-stats/util"

	"github.com/google/go-github/v32/github"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const dateFmt = "2006-01-02"

// newPullRequestsCmd creates a pullRequests command
func newPullRequestsCommand() *cobra.Command {
	var outputFileName string

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

			var earliestDate time.Time
			if daysBack > 0 {
				earliestDate = time.Now().AddDate(0, 0, daysBack*-1)
				fmt.Fprintf(os.Stderr, "including data since %s\n",
					earliestDate.Format("2006-01-02"))
			}

			theStats := &stats.Stats{
				Query:        query,
				EarliestDate: earliestDate,
				Buckets:      []*stats.Bucket{&all},
			}
			err := theStats.Populate()
			if err != nil {
				return errors.Wrap(err, "could not generate stats")
			}

			var out *csv.Writer
			if outputFileName == "" {
				out = csv.NewWriter(os.Stdout)
			} else {
				outFile, err := os.Create(outputFileName)
				cobra.CheckErr(errors.Wrap(err, "could not create output file"))
				defer outFile.Close()
				fmt.Fprintf(os.Stderr, "writing to %s\n", outputFileName)
				out = csv.NewWriter(outFile)
			}

			out.Write([]string{
				"ID",
				"Title",
				"State",
				"Author",
				"URL",
				"Created",
				"Closed",
				"Days to Merge",
				"Review Activity",
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
					strings.TrimSpace(*prd.Pull.Title),
					prd.State,
					user,
					*prd.Pull.HTMLURL,
					createdAt,
					closedAt,
					fmt.Sprintf("%d", daysToMerge),
					fmt.Sprintf("%d", prd.AllActivityCount),
				})
			}

			out.Flush()

			return nil
		},
	}

	pullRequestsCmd.Flags().StringVarP(&outputFileName, "output", "O", "",
		"output file to create (defaults to stdout)")

	return pullRequestsCmd
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
	rootCmd.AddCommand(newPullRequestsCommand())
}
