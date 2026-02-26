package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	serverURL string
	adminPath string
)

var rootCmd = &cobra.Command{
	Use:   "eshop-admin [command]",
	Short: "admin cli tool for interacting with the e-shop server",
	Long: `
A command line tool used by admin to interact with the e-shop
server's various services, like restocking, schedule seckill sales.
	`,
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&adminPath,
		"admin-path",
		"/v1/admin",
		"the base path for admin endpoints",
	)
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
