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
	adjustStockCmd.Flags().Int32Var(
		&stockProductID,
		"product-id",
		-1,
		"the product ID to adjust stock for",
	)
	adjustStockCmd.MarkFlagRequired("product-id")
	adjustStockCmd.Flags().Int32Var(
		&stockDelta,
		"delta",
		0,
		"the stock delta to apply",
	)
	adjustStockCmd.MarkFlagRequired("delta")
	adjustStockCmd.Flags().StringVar(
		&stockReason,
		"reason",
		"",
		"the reason for adjustment",
	)
	adjustStockCmd.Flags().StringVar(
		&stockCreatedBy,
		"created-by",
		"admin-cli",
		"the user creating the adjustment",
	)
	rootCmd.AddCommand(adjustStockCmd)

	listStockAdjustmentsCmd.Flags().Int32Var(
		&stockProductID,
		"product-id",
		-1,
		"the product ID to filter adjustments",
	)
	listStockAdjustmentsCmd.MarkFlagRequired("product-id")
	listStockAdjustmentsCmd.Flags().StringVar(
		&stockCreatedAfter,
		"created-after",
		"",
		"RFC3339 timestamp (inclusive)",
	)
	listStockAdjustmentsCmd.Flags().StringVar(
		&stockCreatedBefore,
		"created-before",
		"",
		"RFC3339 timestamp (inclusive)",
	)
	listStockAdjustmentsCmd.Flags().StringVar(
		&stockCreatedBy,
		"created-by",
		"",
		"filter by creator",
	)
	listStockAdjustmentsCmd.Flags().Int32Var(
		&stockDeltaMin,
		"delta-min",
		0,
		"minimum delta",
	)
	listStockAdjustmentsCmd.Flags().Int32Var(
		&stockDeltaMax,
		"delta-max",
		0,
		"maximum delta",
	)
	listStockAdjustmentsCmd.Flags().Int32VarP(
		&stockPage,
		"page",
		"p",
		1,
		"the page number to retrieve",
	)
	listStockAdjustmentsCmd.Flags().Int32Var(
		&stockPerPage,
		"per-page",
		250,
		"number of items per page",
	)
	rootCmd.AddCommand(listStockAdjustmentsCmd)

	getStockAdjustmentCmd.Flags().Int64Var(
		&stockAdjustmentID,
		"id",
		-1,
		"the stock adjustment ID to retrieve",
	)
	getStockAdjustmentCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(getStockAdjustmentCmd)
}

var (
	stockProductID int32
	stockDelta int32
	stockReason string
	stockCreatedBy string

	stockCreatedAfter string
	stockCreatedBefore string
	stockDeltaMin int32
	stockDeltaMax int32
	stockPage int32
	stockPerPage int32

	stockAdjustmentID int64
)

var adjustStockCmd = &cobra.Command{
	Use:   "adjust-stock",
	Short: "adjust stock for a product",
	RunE:  adjustStock,
}

func adjustStock(cmd *cobra.Command, args []string) error {
	body := map[string]any{
		"product_id": stockProductID,
		"delta": stockDelta,
		"reason": stockReason,
		"created_by": stockCreatedBy,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s%s", serverURL, adminPath, "/stock/adjustments"),
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return reqAndPrint(req)
}

var listStockAdjustmentsCmd = &cobra.Command{
	Use:   "list-stock-adjustments",
	Short: "list stock adjustments for a product",
	RunE:  listStockAdjustments,
}

func listStockAdjustments(cmd *cobra.Command, args []string) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s%s%s", serverURL, adminPath, "/stock/adjustments"),
		nil,
	)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("product_id", fmt.Sprintf("%d", stockProductID))
	if stockCreatedAfter != "" {
		t, err := time.Parse(time.RFC3339, stockCreatedAfter)
		if err != nil {
			return fmt.Errorf("invalid created-after: %v", err)
		}
		q.Add("created_after", fmt.Sprintf("%d", t.UnixMilli()))
	}
	if stockCreatedBefore != "" {
		t, err := time.Parse(time.RFC3339, stockCreatedBefore)
		if err != nil {
			return fmt.Errorf("invalid created-before: %v", err)
		}
		q.Add("created_before", fmt.Sprintf("%d", t.UnixMilli()))
	}
	if stockCreatedBy != "" {
		q.Add("created_by", stockCreatedBy)
	}
	if stockDeltaMin != 0 {
		q.Add("delta_min", fmt.Sprintf("%d", stockDeltaMin))
	}
	if stockDeltaMax != 0 {
		q.Add("delta_max", fmt.Sprintf("%d", stockDeltaMax))
	}
	q.Add("page", fmt.Sprintf("%d", stockPage))
	q.Add("per_page", fmt.Sprintf("%d", stockPerPage))
	req.URL.RawQuery = q.Encode()
	return reqAndPrint(req)
}

var getStockAdjustmentCmd = &cobra.Command{
	Use:   "get-stock-adjustment",
	Short: "get a stock adjustment by ID",
	RunE:  getStockAdjustment,
}

func getStockAdjustment(cmd *cobra.Command, args []string) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s%s%s/%d", serverURL, adminPath, "/stock/adjustments", stockAdjustmentID),
		nil,
	)
	if err != nil {
		return err
	}
	return reqAndPrint(req)
}
