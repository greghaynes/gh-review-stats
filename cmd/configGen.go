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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configGenCmd represents the configGen command
var configGenCmd = &cobra.Command{
	Use:   "config-gen",
	Short: "Generate the default config file",
	Long: `Generate the default config file, if it does not exist.

If it does exist, update it to include any settings that are missing.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		// Write the configuration file with any values we read in and
		// the defaults for values that are missing. If the user has
		// provided a config file name, try to use that. Otherwise use
		// the default path.
		if cfgFile != "" {
			err = viper.WriteConfigAs(cfgFile)
		} else {
			err = viper.WriteConfig()
		}
		cobra.CheckErr(errors.Wrap(err, "failed to create configuration file"))

		// Now re-read the config. This automatically uses the name
		// from the config file option, if present, because
		// initConfig() has set up viper.
		err = viper.ReadInConfig()
		cobra.CheckErr(errors.Wrap(err, "failed to re-read configuration file"))

		// Now we can figure out what name was actually used and
		// report that we wrote to it.
		filename := viper.ConfigFileUsed()
		fmt.Printf("wrote %q\n", filename)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configGenCmd)
}
