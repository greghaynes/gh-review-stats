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

	"github.com/dhellmann/gh-review-stats/util"
	"github.com/spf13/cobra"

	"github.com/google/go-github/v45/github"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	githubTokenConfigOptionName               = "github.token"
	githubEnterpriseBaseURLConfigOptionName   = "github.enterprise.base-url"
	githubEnterpriseUploadURLConfigOptionName = "github.enterprise.upload-url"
)

var cfgFile string

// devMode is a flag telling us whether we are in developer mode
var devMode bool

// orgName and repoName are the GitHub organization and repository to query
var orgName, repoName string

// daysBack is the number of days of history to examine (older items are ignored)
var daysBack int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-review-stats",
	Short: "GitHub Review Statistics",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func githubToken() string {
	return viper.GetString(githubTokenConfigOptionName)
}

func newGithubClient(ctx context.Context) (*github.Client, error) {
	return util.NewGithubClient(
		ctx,
		githubToken(),
		viper.GetString(githubEnterpriseBaseURLConfigOptionName),
		viper.GetString(githubEnterpriseUploadURLConfigOptionName),
	)
}

func addHistoryArgs(theCommand *cobra.Command) {
	theCommand.PersistentFlags().StringVarP(&orgName, "org", "o", "",
		"github org")
	theCommand.PersistentFlags().StringVarP(&repoName, "repo", "r", "",
		"github repository")
	theCommand.PersistentFlags().IntVar(&daysBack, "days-back", 90,
		"how many days back to query")
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	viper.SetDefault(githubTokenConfigOptionName, "")
	viper.SetDefault(githubEnterpriseBaseURLConfigOptionName, "")
	viper.SetDefault(githubEnterpriseUploadURLConfigOptionName, "")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.gh-review-stats.yml)")
	rootCmd.PersistentFlags().BoolVar(&devMode, "dev", false,
		"enable developer mode, shortcutting some queries")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gh-review-stats".
		viper.AddConfigPath(home)
		viper.SetConfigName(".gh-review-stats")
		viper.SetConfigType("yml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
