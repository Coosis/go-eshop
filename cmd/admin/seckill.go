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
}

var (
	addSeckillProductID int32
	addSeckillStartTime string
	addSeckillEndTime string
	addSeckillPriceCents int32
	addSeckillStock int32
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
