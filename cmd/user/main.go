package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	serverURL string
)

var rootCmd = &cobra.Command{
	Use:   "eshop [command]",
	Short: "user cli tool for interacting with the e-shop server",
	Long:  "A command line tool used by users to interact with the e-shop server.",
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&serverURL,
		"server-url",
		"http://localhost:8144",
		"the base URL of the e-shop server",
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
