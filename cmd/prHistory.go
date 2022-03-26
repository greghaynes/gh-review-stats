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
	"sort"
	"strconv"
	"strings"

	"github.com/dhellmann/gh-review-stats/events"
	"github.com/dhellmann/gh-review-stats/stats"
	"github.com/dhellmann/gh-review-stats/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type keyCount struct {
	Key   string
	Count int
}

// prHistoryCmd represents the prHistory command
var prHistoryCmd = &cobra.Command{
	Use:       "pr-history pull-request-id...",
	Short:     "Summarize the history of one pull request",
	Long:      `Produce stats and a history log of one pull request`,
	ValidArgs: []string{"pull-request"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("expecting at least 1 pull-request-id argument, got %d",
				len(args))
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
				{
					Rule: func(*stats.PullRequestDetails) bool {
						return true
					},
				},
			},
		}

		toIgnore := reviewersToIgnore()

		// fetch all of the event data for all pull requests
		for _, arg := range args {
			prID, err := strconv.Atoi(arg)
			if err != nil {
				return errors.Wrap(err, "pull-request-id must be a number")
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
		}

		// merge the events into a single stream
		allEvents := []*events.Event{}
		for _, prd := range prStats.Buckets[0].Requests {
			events := events.GetOrderedEvents(prd)
			allEvents = append(allEvents, events...)
		}
		sort.Slice(allEvents, func(i, j int) bool {
			return allEvents[i].Date.Before(*allEvents[j].Date)
		})

		// prepare to summarize activity of participants
		// (maps user names to unique dates)
		personActivityDates := map[string]map[string]bool{}
		// (maps dates to activity count)
		dateActivity := map[string]int{}

		// show the event log
		var previous *events.Event
		for _, e := range allEvents {
			if _, ok := toIgnore[e.Person]; ok {
				continue
			}

			if previous != nil {
				delay := int(math.Floor(e.Date.Sub(*previous.Date).Hours() / 24))
				if delay > 1 {
					fmt.Printf("%d days\n", delay)
				}
			}

			fmt.Printf("%s: %s\n", e.Date.Format("Mon Jan _2"), e.Description)

			if _, ok := personActivityDates[e.Person]; !ok {
				personActivityDates[e.Person] = map[string]bool{}
			}
			dateKey := e.Date.Format(dateFmt)
			personActivityDates[e.Person][dateKey] = true

			if _, ok := dateActivity[dateKey]; !ok {
				dateActivity[dateKey] = 0
			}
			dateActivity[dateKey]++

			previous = e
		}

		// show the number of dates each reviewer was active
		pairs := []keyCount{}
		for person, dates := range personActivityDates {
			if person == "" {
				continue
			}
			if _, ok := toIgnore[person]; ok {
				continue
			}
			pairs = append(pairs, keyCount{
				Key:   person,
				Count: len(dates),
			})
		}
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Count > pairs[j].Count
		})

		fmt.Printf("\nNumber of Engaged Days\n")
		for _, p := range pairs {
			fmt.Printf("%s: %d\n", p.Key, p.Count)
		}

		// show the amount of activity on each day
		pairs = []keyCount{}
		maxDailyActivity := 0
		for date, count := range dateActivity {
			pairs = append(pairs, keyCount{
				Key:   date,
				Count: count,
			})
			if count > maxDailyActivity {
				maxDailyActivity = count
			}
		}
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Key > pairs[j].Key
		})

		fmt.Printf("\nEngagement by Day\n")
		for _, p := range pairs {
			//barLength := int(math.Floor(float64(p.Count) / 100 * 25))
			barLength := int(math.Floor((float64(p.Count) / float64(maxDailyActivity)) * 60))
			bar := strings.Repeat("*", barLength)
			fmt.Printf("%s: %3d %s\n", p.Key, p.Count, bar)
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
