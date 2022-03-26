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
package events

import (
	"fmt"
	"sort"
	"time"

	"github.com/dhellmann/gh-review-stats/stats"
	"github.com/google/go-github/v32/github"
)

type Event struct {
	Date        *time.Time
	Description string
}

func getName(user *github.User) string {
	if user.Name != nil {
		return *user.Name
	}
	if user.Login != nil {
		return *user.Login
	}
	return "unnamed"
}

func GetOrderedEvents(prd *stats.PullRequestDetails) []*Event {
	results := []*Event{
		&Event{
			Date:        prd.Pull.CreatedAt,
			Description: "pull request opened",
		},
	}
	if prd.Pull.ClosedAt != nil {
		desc := "pull request closed"
		if prd.State == "merged" {
			desc = "pull request merged"
		}
		results = append(results, &Event{
			Date:        prd.Pull.ClosedAt,
			Description: desc,
		})
	}

	for _, review := range prd.Reviews {
		results = append(results, &Event{
			Date:        review.SubmittedAt,
			Description: fmt.Sprintf("review by %s", getName(review.User)),
		})
	}

	for _, comment := range prd.PullRequestComments {
		results = append(results, &Event{
			Date:        comment.CreatedAt,
			Description: fmt.Sprintf("comment by %s", getName(comment.User)),
		})
	}

	for _, comment := range prd.IssueComments {
		results = append(results, &Event{
			Date:        comment.CreatedAt,
			Description: fmt.Sprintf("comment by %s", getName(comment.User)),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Date.Before(*results[j].Date)
	})

	return results
}
