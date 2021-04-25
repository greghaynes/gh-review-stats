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
	"fmt"

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
		fmt.Println("reviewers called", reviewersToIgnore())
		return nil
	},
}

func reviewersToIgnore() map[string]interface{} {
	result := map[string]interface{}{}
	for _, i := range ignoredReviewers {
		result[i] = nil
	}
	for _, i := range viper.GetStringSlice(ignoreConfigOptionName) {
		result[i] = nil
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
}
