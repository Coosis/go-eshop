package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	placeOrderCmd.Flags().Int64VarP(
		&cartVersion,
		"cart-version",
		"c",
		0,
		"the version of the cart to place the order for",
	)
	placeOrderCmd.MarkFlagRequired("cart-version")
	rootCmd.AddCommand(placeOrderCmd)

	getOrdersCmd.PersistentFlags().Int32VarP(
		&getOrdersPage,
		"page",
		"p",
		1,
		"the page number to retrieve",
	)
	getOrdersCmd.PersistentFlags().Int32Var(
		&getOrdersPerPage,
		"per-page",
		250,
		"number of orders per page",
	)
	getOrdersCmd.PersistentFlags().StringVarP(
		&getOrdersBefore,
		"before",
		"b",
		"",
		`retrieve orders placed before this timestamp, in RFC3339 format.
		Example: 2020-01-01T00:00:00Z`,
	)
	getOrdersCmd.PersistentFlags().StringVarP(
		&getOrdersAfter,
		"after",
		"a",
		"",
		`retrieve orders placed after this timestamp, in RFC3339 format. 
		Example: 2020-01-01T00:00:00Z`,
	)
	getOrdersCmd.PersistentFlags().StringVar(
		&getOrdersStatus,
		"status",
		"",
		"filter orders by status (e.g., pending, shipped, delivered)",
	)
	rootCmd.AddCommand(getOrdersCmd)

	payOrderCmd.Flags().Int32VarP(
		&payOrderID,
		"id",
		"i",
		-1,
		"the ID of the order to pay for",
	)
	payOrderCmd.MarkFlagRequired("id")
	payOrderCmd.Flags().StringVar(
		&paymentIntentID,
		"payment-intent-id",
		"",
		"the payment intent ID to use for payment",
	)
	rootCmd.AddCommand(payOrderCmd)

	cancelOrderCmd.Flags().Int32VarP(
		&cancelOrderID,
		"id",
		"i",
		-1,
		"the ID of the order to cancel",
	)
	cancelOrderCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(cancelOrderCmd)

	refundOrderCmd.Flags().Int32VarP(
		&refundOrderID,
		"id",
		"i",
		-1,
		"the ID of the order to refund",
	)
	refundOrderCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(refundOrderCmd)

	rootCmd.PersistentFlags().StringVar(
		&orderPath,
		"order-path",
		"/v1/orders",
		"the base path of the order service",
	)
}

var (
	orderPath string

	cartVersion int64
)

var placeOrderCmd = &cobra.Command{
	Use:   "place-order",
	Short: "place an order",
	Long:  `Place an order for products in the cart.`,
	RunE:  placeOrder,
}

func placeOrder(cmd *cobra.Command, args []string) error {
	orig_req, err := json.Marshal(map[string]any{
		"cart_version":    cartVersion,
		"idempotency_key": generateChars(8),
	})
	if err != nil {
		return err
	}
	// needs 1. cart version 2. idempotency key
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", serverURL, orderPath),
		bytes.NewReader(orig_req),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return reqAndPrint(req)
}

var (
	getOrdersPage    int32
	getOrdersPerPage int32
	getOrdersBefore  string
	getOrdersAfter   string
	getOrdersStatus  string
)

var getOrdersCmd = &cobra.Command{
	Use:   "get-orders",
	Short: "get orders for the current user",
	Long:  `Get all orders placed by the current user.`,
	RunE:  getOrders,
}

func getOrders(cmd *cobra.Command, args []string) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s%s", serverURL, orderPath),
		nil,
	)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	if getOrdersBefore != "" {
		bef, err := time.Parse(time.RFC3339, getOrdersBefore)
		if err == nil {
			getOrdersBefore = fmt.Sprintf("%d", bef.UnixMilli())
		} else {
			return fmt.Errorf("invalid before timestamp: %v", err)
		}
		q.Add("before", getOrdersBefore)
	}
	if getOrdersAfter != "" {
		aft, err := time.Parse(time.RFC3339, getOrdersAfter)
		if err == nil {
			getOrdersAfter = fmt.Sprintf("%d", aft.UnixMilli())
		} else {
			return fmt.Errorf("invalid after timestamp: %v", err)
		}
		q.Add("after", getOrdersAfter)
	}
	q.Add("page", fmt.Sprintf("%d", getOrdersPage))
	q.Add("per_page", fmt.Sprintf("%d", getOrdersPerPage))
	if getOrdersStatus != "" {
		q.Add("status", getOrdersStatus)
	}
	req.URL.RawQuery = q.Encode()

	return reqAndPrint(req)
}

var (
	payOrderID int32
	paymentIntentID string
)

var payOrderCmd = &cobra.Command{
	Use: "pay-order --id <order_id>",
	Short: "pay for an order",
	Long: `Pay for an order with the specified order ID.`,
	RunE: payOrder,
}

func payOrder(cmd *cobra.Command, args []string) error {
	if payOrderID <= 0 {
		return fmt.Errorf("invalid order ID: %d", payOrderID)
	}
	orig_req, err := json.Marshal(map[string]any{
		"payment_intent_id": paymentIntentID,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s/%d/pay", serverURL, orderPath, payOrderID),
		bytes.NewReader(orig_req),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return reqAndPrint(req)
}

var (
	cancelOrderID int32
)

var cancelOrderCmd = &cobra.Command{
	Use: "cancel-order --id <order_id>",
	Short: "cancel an order",
	Long: `Cancel an order with the specified order ID.`,
	RunE: cancelOrder,
}

func cancelOrder(cmd *cobra.Command, args []string) error {
	if cancelOrderID <= 0 {
		return fmt.Errorf("invalid order ID: %d", cancelOrderID)
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s/%d/cancel", serverURL, orderPath, cancelOrderID),
		nil,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return reqAndPrint(req)
}

var (
	refundOrderID int32
)

var refundOrderCmd = &cobra.Command{
	Use: "refund-order --id <order_id>",
	Short: "refund an order",
	Long: `Refund an order with the specified order ID.`,
	RunE: refundOrder,
}

func refundOrder(cmd *cobra.Command, args []string) error {
	if refundOrderID <= 0 {
		return fmt.Errorf("invalid order ID: %d", refundOrderID)
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s/%d/refund", serverURL, orderPath, refundOrderID),
		nil,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return reqAndPrint(req)
}
