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
	"github.com/google/go-github/v45/github"
)

type Event struct {
	Date        *time.Time
	Description string
	Person      string
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
		{
			Date: prd.Pull.CreatedAt,
			Description: fmt.Sprintf("#%d opened by %s %q (%s)",
				*prd.Pull.Number, getName(prd.Pull.User), *prd.Pull.Title,
				*prd.Pull.HTMLURL),
			Person: getName(prd.Pull.User),
		},
	}
	if prd.Pull.ClosedAt != nil {
		daysOpen := int(prd.Pull.ClosedAt.Sub(*prd.Pull.CreatedAt).Hours() / 24)
		results = append(results, &Event{
			Date: prd.Pull.ClosedAt,
			Description: fmt.Sprintf("#%d %s after %d days %q (%s)",
				*prd.Pull.Number, prd.State, daysOpen, *prd.Pull.Title,
				*prd.Pull.HTMLURL),
		})
	} else {
		daysOpen := int(time.Since(*prd.Pull.CreatedAt).Hours() / 24)
		now := time.Now()
		results = append(results, &Event{
			Date: &now,
			Description: fmt.Sprintf("#%d %s %d days %q (%s)",
				*prd.Pull.Number, prd.State, daysOpen, *prd.Pull.Title,
				*prd.Pull.HTMLURL),
		})
	}

	for _, commit := range prd.Commits {
		results = append(results, &Event{
			Date: commit.Commit.Author.Date,
			Description: fmt.Sprintf("#%d updated by %s",
				*prd.Pull.Number, *commit.Commit.Author.Name),
			Person: *commit.Commit.Author.Name,
		})
	}

	for _, review := range prd.Reviews {
		results = append(results, &Event{
			Date: review.SubmittedAt,
			Description: fmt.Sprintf("#%d review by %s", *prd.Pull.Number,
				getName(review.User)),
			Person: getName(review.User),
		})
	}

	for _, comment := range prd.PullRequestComments {
		results = append(results, &Event{
			Date: comment.CreatedAt,
			Description: fmt.Sprintf("#%d comment by %s", *prd.Pull.Number,
				getName(comment.User)),
			Person: getName(comment.User),
		})
	}

	for _, comment := range prd.IssueComments {
		results = append(results, &Event{
			Date: comment.CreatedAt,
			Description: fmt.Sprintf("#%d comment by %s", *prd.Pull.Number,
				getName(comment.User)),
			Person: getName(comment.User),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Date.Before(*results[j].Date)
	})

	return results
}
