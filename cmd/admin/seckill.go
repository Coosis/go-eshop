package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Coosis/go-eshop/internal/seckill"
	"github.com/spf13/cobra"
)

func init() {
	addSeckillCmd.Flags().Int32Var(
		&addSeckillProductID,
		"product-id",
		-1,
		"the product ID for the seckill event",
	)
	addSeckillCmd.MarkFlagRequired("product-id")

	addSeckillCmd.Flags().StringVarP(
		&addSeckillStartTime,
		"start-time",
		"s",
		"",
		"the start time for the seckill event (Unix timestamp in milliseconds)",
	)
	addSeckillCmd.MarkFlagRequired("start-time")

	addSeckillCmd.Flags().StringVarP(
		&addSeckillEndTime,
		"end-time",
		"e",
		"",
		"the end time for the seckill event (Unix timestamp in milliseconds)",
	)
	addSeckillCmd.MarkFlagRequired("end-time")

	addSeckillCmd.Flags().Int32Var(
		&addSeckillPriceCents,
		"seckill-price-cents",
		-1,
		"the seckill price in cents for the seckill event",
	)
	addSeckillCmd.MarkFlagRequired("seckill-price-cents")

	addSeckillCmd.Flags().Int32Var(
		&addSeckillStock,
		"seckill-stock",
		-1,
		"the seckill stock for the seckill event",
	)
	addSeckillCmd.MarkFlagRequired("seckill-stock")
	rootCmd.AddCommand(addSeckillCmd)

	updateSeckillCmd.Flags().Int32Var(
		&updateSeckillID,
		"id",
		-1,
		"the seckill event ID to update",
	)
	updateSeckillCmd.MarkFlagRequired("id")
	updateSeckillCmd.Flags().Int32Var(
		&updateSeckillProductID,
		"product-id",
		-1,
		"the product ID for the seckill event",
	)
	updateSeckillCmd.MarkFlagRequired("product-id")
	updateSeckillCmd.Flags().StringVar(
		&updateSeckillStartTime,
		"start-time",
		"",
		"the start time for the seckill event (RFC3339)",
	)
	updateSeckillCmd.MarkFlagRequired("start-time")
	updateSeckillCmd.Flags().StringVar(
		&updateSeckillEndTime,
		"end-time",
		"",
		"the end time for the seckill event (RFC3339)",
	)
	updateSeckillCmd.MarkFlagRequired("end-time")
	updateSeckillCmd.Flags().Int32Var(
		&updateSeckillPriceCents,
		"seckill-price-cents",
		-1,
		"the seckill price in cents for the seckill event",
	)
	updateSeckillCmd.MarkFlagRequired("seckill-price-cents")
	updateSeckillCmd.Flags().Int32Var(
		&updateSeckillStock,
		"seckill-stock",
		-1,
		"the seckill stock for the seckill event",
	)
	updateSeckillCmd.MarkFlagRequired("seckill-stock")
	rootCmd.AddCommand(updateSeckillCmd)
}

var (
	addSeckillProductID int32
	addSeckillStartTime string
	addSeckillEndTime string
	addSeckillPriceCents int32
	addSeckillStock int32
	updateSeckillID int32
	updateSeckillProductID int32
	updateSeckillStartTime string
	updateSeckillEndTime string
	updateSeckillPriceCents int32
	updateSeckillStock int32
)

var addSeckillCmd = &cobra.Command{
	Use:   "add-seckill",
	Short: "Add a new seckill event",
	Long:  `Add a new seckill event to the system.`,
	RunE: addSeckill,
}

func addSeckill(cmd *cobra.Command, args []string) error {
	sttime, err := time.Parse(time.RFC3339, addSeckillStartTime)
	if err != nil {
		return fmt.Errorf("invalid start time format: %v", err)
	}
	tstart := sttime.UnixMilli()

	edtime, err := time.Parse(time.RFC3339, addSeckillEndTime)
	if err != nil {
		return fmt.Errorf("invalid end time format: %v", err)
	}
	tend := edtime.UnixMilli()
	orig_req, err := json.Marshal(&seckill.SeckillEventInfo{
		ProductID: addSeckillProductID,
		StartTime: tstart,
		EndTime: tend,
		SeckillPriceCents: addSeckillPriceCents,
		SeckillStock: addSeckillStock,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s%s", serverURL, adminPath, "/seckill/events"),
		bytes.NewReader(orig_req),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return reqAndPrint(req)
}

var updateSeckillCmd = &cobra.Command{
	Use:   "update-seckill",
	Short: "Update a seckill event",
	Long:  `Update a seckill event in the system.`,
	RunE:  updateSeckill,
}

func updateSeckill(cmd *cobra.Command, args []string) error {
	sttime, err := time.Parse(time.RFC3339, updateSeckillStartTime)
	if err != nil {
		return fmt.Errorf("invalid start time format: %v", err)
	}
	tstart := sttime.UnixMilli()

	edtime, err := time.Parse(time.RFC3339, updateSeckillEndTime)
	if err != nil {
		return fmt.Errorf("invalid end time format: %v", err)
	}
	tend := edtime.UnixMilli()
	orig_req, err := json.Marshal(&seckill.SeckillEventInfo{
		ProductID: updateSeckillProductID,
		StartTime: tstart,
		EndTime: tend,
		SeckillPriceCents: updateSeckillPriceCents,
		SeckillStock: updateSeckillStock,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s%s%s/%d", serverURL, adminPath, "/seckill/events", updateSeckillID),
		bytes.NewReader(orig_req),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return reqAndPrint(req)
}
