/*
Copyright © 2022 Doug Hellmann <doug@doughellmann.com>

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
	"fmt"
	"math"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/dhellmann/gh-review-stats/events"
	"github.com/dhellmann/gh-review-stats/stats"
	"github.com/dhellmann/gh-review-stats/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// prHistoryCmd represents the prHistory command
var prHistoryCmd = &cobra.Command{
	Use:       "pr-history pull-request-id",
	Short:     "Summarize the history of one pull request",
	Long:      `Produce stats and a history log of one pull request`,
	ValidArgs: []string{"pull-request"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expecting 1 pull-request-id argument, got %d",
				len(args))
		}
		prID, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.Wrap(err, "pull-request-id must be a number")
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		query := &util.PullRequestQuery{
			Org:     orgName,
			Repo:    repoName,
			DevMode: devMode,
			Client:  util.NewGithubClient(ctx, githubToken()),
		}

		prStats := &stats.Stats{
			Query: query,
			Buckets: []*stats.Bucket{
				&stats.Bucket{
					Rule: func(*stats.PullRequestDetails) bool {
						return true
					},
				},
			},
		}

		pr, _, err := query.Client.PullRequests.Get(ctx, orgName, repoName, prID)
		if err != nil {
			return errors.Wrap(err, "failed to fetch pull request")
		}
		prStats.ProcessOne(ctx, pr)

		select {
		case <-ctx.Done():
			return nil
		default:
		}

		prd := prStats.Buckets[0].Requests[0]
		fmt.Printf("Pull request: %s\n", *prd.Pull.HTMLURL)

		var (
			createdAt, closedAt string
			daysOpen            int = -1
		)

		if prd.Pull.CreatedAt != nil {
			createdAt = prd.Pull.CreatedAt.Format(dateFmt)
		}
		if prd.Pull.ClosedAt != nil {
			closedAt = prd.Pull.ClosedAt.Format(dateFmt)
		}
		if prd.Pull.CreatedAt != nil {
			if prd.State == "merged" && prd.Pull.ClosedAt != nil {
				daysOpen = int(prd.Pull.ClosedAt.Sub(*prd.Pull.CreatedAt).Hours() / 24)
			} else {
				daysOpen = int(time.Since(*prd.Pull.CreatedAt).Hours() / 24)
			}
		}
		fmt.Printf("Created     : %s\n", createdAt)
		fmt.Printf("Closed      : %s\n", closedAt)
		fmt.Printf("Age         : %d days\n", daysOpen)
		fmt.Printf("\n")

		var previous *events.Event
		events := events.GetOrderedEvents(prd)
		for _, e := range events {
			if previous != nil {
				delay := int(math.Floor(e.Date.Sub(*previous.Date).Hours() / 24))
				if delay > 1 {
					fmt.Printf("%d days\n", delay)
				}
			}
			fmt.Printf("%s: %s\n", e.Date.Format(dateFmt), e.Description)
			previous = e
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(prHistoryCmd)
	addHistoryArgs(prHistoryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//prHistoryCmd.PersistentFlags().String("pr", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prHistoryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
