package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Coosis/go-eshop/internal/seckill"
	"github.com/spf13/cobra"
)

func init() {
	purchaseSeckillCmd.Flags().Int32Var(
		&seckillEventID,
		"event-id",
		-1,
		"the seckill event ID to purchase from",
	)
	purchaseSeckillCmd.MarkFlagRequired("event-id")

	purchaseSeckillCmd.Flags().Int64Var(
		&seckillQuantity,
		"quantity",
		-1,
		"the quantity to purchase",
	)
	purchaseSeckillCmd.MarkFlagRequired("quantity")

	purchaseSeckillCmd.Flags().StringVar(
		&seckillKey,
		"seckill-key",
		"",
		"the idempotency key for the purchase attempt (optional, will be randomly generated if not provided)",
	)

	listSeckillsCmd.Flags().Int32VarP(
		&listSeckillsPage,
		"page",
		"p",
		1,
		"the page number to retrieve",
	)
	listSeckillsCmd.Flags().Int32VarP(
		&listSeckillsPerPage,
		"per-page",
		"n",
		250,
		"number of seckill events per page",
	)

	getSeckillStatusCmd.Flags().StringVar(
		&seckillStatusKey,
		"idempotency-key",
		"",
		"the idempotency key for the purchase attempt",
	)

	rootCmd.Flags().StringVar(
		&seckillPath,
		"seckill-path",
		"/v1/seckill",
		"the base path for seckill API endpoints",
	)

	rootCmd.AddCommand(purchaseSeckillCmd)
	rootCmd.AddCommand(listSeckillsCmd)
	rootCmd.AddCommand(getSeckillCmd)
	rootCmd.AddCommand(getSeckillStatusCmd)
}

var (
	seckillPath string
	seckillKey  string
	seckillStatusKey string

	seckillEventID int32
	seckillQuantity int64
)

var purchaseSeckillCmd = &cobra.Command{
	Use:   "purchase-seckill",
	Short: "Purchase a seckill product",
	Long:  `This command allows a user to purchase a seckill product.`,
	RunE: runPurchaseSeckill,
}

func runPurchaseSeckill(cmd *cobra.Command, args []string) error {
	if seckillKey == "" {
		seckillKey = generateChars(16)
	}
	orig_req, err := json.Marshal(&seckill.SeckillAttempt{
		EventID: seckillEventID,
		Quantity: seckillQuantity,
		IdempotencyKey: seckillKey,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Purchasing seckill product with EventID=%d, Quantity=%d, IdempotencyKey=%s\n",
		seckillEventID,
		seckillQuantity,
		seckillKey,
	)
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s%s%s/attempt", serverURL, seckillPath, "/events", fmt.Sprintf("/%d", seckillEventID)),
		bytes.NewReader(orig_req),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return reqAndPrint(req)
}

var (
	listSeckillsPage int32
	listSeckillsPerPage int32
)

var listSeckillsCmd = &cobra.Command{
	Use:   "list-seckill",
	Short: "List seckill events",
	Long:  `This command retrieves the list of seckill events.`,
	RunE: listAllSeckills,
}

func listAllSeckills(cmd *cobra.Command, args []string) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s%s%s", serverURL, seckillPath, "/events"),
		nil,
	)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("page", fmt.Sprintf("%d", listSeckillsPage))
	q.Add("per_page", fmt.Sprintf("%d", listSeckillsPerPage))
	req.URL.RawQuery = q.Encode()
	return reqAndPrint(req)
}

var getSeckillCmd = &cobra.Command{
	Use:   "get-seckill [event_id]",
	Short: "Get details of a seckill event",
	Long:  `This command retrieves the details of a specific seckill event by its ID.`,
	RunE: getSeckillByID,
}

func getSeckillByID(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("seckill event ID is required")
	}
	eventID := args[0]
	if _, err := strconv.ParseInt(eventID, 10, 32); err != nil {
		return fmt.Errorf("invalid seckill event ID: %s", eventID)
	}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s%s%s%s", serverURL, seckillPath, "/events/", eventID),
		nil,
	)
	if err != nil {
		return err
	}
	return reqAndPrint(req)
}

var getSeckillStatusCmd = &cobra.Command{
	Use:   "get-seckill-status",
	Short: "Get the status of a seckill event",
	Long:  `This command retrieves the current status of a specific seckill purchase by idempotency key.`,
	RunE: getSeckillStatus,
}

func getSeckillStatus(cmd *cobra.Command, args []string) error {
	if seckillStatusKey == "" {
		if len(args) < 1 {
			return fmt.Errorf("idempotency key is required")
		}
		seckillStatusKey = args[0]
	}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s%s%s%s/status", serverURL, seckillPath, "/attempts/", seckillStatusKey),
		nil,
	)
	if err != nil {
		return err
	}
	return reqAndPrint(req)
}
