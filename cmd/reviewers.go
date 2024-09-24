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
	"fmt"
	"os"
	"os/signal"
	"sort"
	"time"

	"github.com/dhellmann/gh-review-stats/reviewers"
	"github.com/dhellmann/gh-review-stats/util"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const ignoreConfigOptionName = "reviewers.ignore"

// ignoredReviewers is the list of github ids to leave out of the
// stats
var ignoredReviewers = []string{}

// reviewersCmd represents the reviewers command
var reviewersCmd = &cobra.Command{
	Use:   "reviewers",
	Short: "List reviewers of PRs in a repo",
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

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		ghClient, err := newGithubClient(ctx)
		if err != nil {
			return err
		}

		query := &util.PullRequestQuery{
			Org:     orgName,
			Repo:    repoName,
			DevMode: devMode,
			Client:  ghClient,
		}

		var earliestDate time.Time
		if daysBack > 0 {
			earliestDate = time.Now().AddDate(0, 0, daysBack*-1)
		}

		reviewerStats := &reviewers.Stats{
			Query:        query,
			EarliestDate: earliestDate,
		}

		err = query.IteratePullRequests(ctx, reviewerStats.ProcessOne)
		if err != nil {
			return errors.Wrap(err, "failed to retrieve pull request details")
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}

		toIgnore := reviewersToIgnore()
		orderedReviewers := reviewerStats.ReviewersInOrder()

		for _, reviewer := range orderedReviewers {

			if _, ok := toIgnore[reviewer]; ok {
				continue
			}

			count := reviewerStats.ReviewCounts[reviewer]
			prs := reviewerStats.PRsForReviewer(reviewer)

			fmt.Printf("%d/%d: %s\n", count, len(prs), reviewer)

			sort.Slice(prs, func(i, j int) bool {
				return prs[i].ReviewCount > prs[j].ReviewCount
			})
			for _, prWithCount := range prs {
				pr := prWithCount.PR
				fmt.Printf("\t%3d: %s [%s] %q\n", prWithCount.ReviewCount,
					*pr.HTMLURL, *pr.User.Login, *pr.Title)
			}
		}

		return nil
	},
}

func reviewersToIgnore() map[string]bool {
	result := map[string]bool{}
	for _, i := range ignoredReviewers {
		result[i] = true
	}
	for _, i := range viper.GetStringSlice(ignoreConfigOptionName) {
		result[i] = true
	}
	return result
}

func init() {
	rootCmd.AddCommand(reviewersCmd)

	// Here you will define your flags and configuration settings.

	viper.SetDefault(ignoreConfigOptionName, []string{})

	reviewersCmd.Flags().StringSliceVarP(&ignoredReviewers,
		"ignore", "i", []string{},
		"ignore a reviewer (useful for bots), can be repeated")
	addHistoryArgs(reviewersCmd)
}
