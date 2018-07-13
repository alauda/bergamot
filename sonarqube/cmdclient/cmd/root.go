package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "sonarclient",
	Short: "Sonarqube client",
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var (
	workDir    string
	sonarHost  string
	sonarToken string
)

func init() {
	RootCmd.AddCommand(taskMonitorCmd)
	taskMonitorCmd.Flags().StringVar(&workDir, "w", "./", "sonar scanner work directory")
	RootCmd.PersistentFlags().StringVar(&sonarHost, "host", "localhost:9000", "sonarqube server url")
	RootCmd.PersistentFlags().StringVar(&sonarToken, "token", "", "sonarqube api token")
}
