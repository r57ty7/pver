/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"net/http"
	"os"

	"github.com/r57ty7/pver/infra"
	"github.com/r57ty7/pver/service"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile        string
	conf           service.Config
	gitRepository  GitRepository
	pomFvm         FileVersionManager
	npmFvm         FileVersionManager
	jiraService    JiraService
	jiraRepository service.JiraRepository
)

func NewCmdRoot(version string) *cobra.Command {

	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:     "pver",
		Version: version,
		Short:   "project versioning",
		// 		Long: `A longer description that spans multiple lines and likely contains
		// examples and usage of using your application. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("%+v", conf)
		},
	}

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pver.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(newPomCmd())
	rootCmd.AddCommand(newNpmCmd())
	rootCmd.AddCommand(newJiraCmd())

	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".pver" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".pver")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Printf("unmarshal error: %v\n", err)
		return
	}

	pomFvm = service.NewMavenProject()
	npmFvm = service.NewNpmProject()
	gitRepository = service.NewRepository("./")

	jiraRepository, err = infra.NewJiraRepository(
		http.DefaultClient,
		conf.Jira.BaseURL,
		conf.Jira.Username,
		conf.Jira.Password,
	)
	if err != nil {
		fmt.Printf("initialize jiraRepository error: %v\n", err)
		return
	}
	jiraService = service.NewJiraService(jiraRepository)

}
