package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Coosis/go-eshop/internal/catalog"
	"github.com/spf13/cobra"
)

func init() {
	mockProductsCmd.PersistentFlags().StringVar(
		&mock_products_csv_path,
		"file",
		"product_sheet.csv",
		"path to the file containing mock product data",
	)
	rootCmd.AddCommand(mockProductsCmd)
}

var (
	mock_products_csv_path string
)

var mockProductsCmd = &cobra.Command{
	Use:   "mock-products",
	Short: "add mock products to the catalog",
	Long: `
Add mock products to the catalog service from a CSV file.
The CSV file should have the following columns:
name,description.
For price, a random price between $0.00 and $30.00 will be assigned.
For category IDs, no categories will be assigned.
	`,
	RunE: addMockProducts,
}

func addMockProducts(cmd *cobra.Command, args []string) error {
	_, err := os.Stat(mock_products_csv_path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("file does not exist: ", mock_products_csv_path)
			return err
		}
		return err
	}

	contents, err := os.ReadFile(mock_products_csv_path)
	if err != nil {
		return err
	}
	contents_str := strings.Split(string(contents), "\n")[1:]
	for _, line := range contents_str {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) != 2 {
			fmt.Println("skipping malformed line: ", line)
			continue
		}
		name := fields[0]
		desc := fields[1]

		prop := makeProduct(name, desc)
		pid, err := addProduct(context.Background(), prop)
		if err != nil {
			fmt.Printf("error adding product: %v, skipping...", err)
			continue
		}
		err = adjustStockLevel(context.Background(), pid, 100, "initial stock for mock product")
		if err != nil {
			fmt.Printf("error adjusting stock level for product ID %d: %v, skipping...\n", pid, err)
		}
	}
	return nil
}

func addProduct(
	ctx context.Context,
	prop catalog.ProductProperties,
) (int32, error) {
	b, err := json.Marshal(prop)
	if err != nil {
		return -1, err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s%s%s", serverURL, adminPath, "/catalog/products"),
		bytes.NewReader(b),
	)
	if err != nil {
		return -1, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	fmt.Println("response status:", resp.Status)
	msg, _ := io.ReadAll(resp.Body)
	var mp_res struct {
		ID int32 `json:"id"`
	}
	err = json.Unmarshal(msg, &mp_res)
	if err != nil {
		return -1, err
	}
	fmt.Printf("added product ID: %d\n", mp_res.ID)
	return mp_res.ID, nil
}

func adjustStockLevel(
	ctx context.Context,
	productID int32,
	delta int32,
	reason string,
) error {
	reqBody := map[string]any{
		"product_id": productID,
		"delta":  delta,
		"reason":     reason,
		"created_by": "mock-products-cmd",
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(
		ctx,
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
