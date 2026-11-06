/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"os"

	"github.com/fredericobenevides/pge-stress-test/internal/stresstest"
	"github.com/fredericobenevides/pge-stress-test/internal/summary"
	"github.com/spf13/cobra"
)

var (
	url         string
	requests    int
	concurrency int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pge-stress-test",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		config, err := stresstest.NewConfig(url, requests, concurrency)
		if err != nil {
			os.Exit(1)
		}

		stressTest := stresstest.NewStressTest(config)
		operationDuration, results := stressTest.Run()

		summary := summary.NewSummary(operationDuration, results)
		summary.Print()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL to stress test")
	rootCmd.Flags().IntVarP(&requests, "requests", "r", 1, "Number of requests to perform")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent requests")
}
