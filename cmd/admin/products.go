package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Coosis/go-eshop/internal/catalog"
	"github.com/spf13/cobra"
)

func init() {
	addProductCmd.Flags().StringVar(
		&productName,
		"name",
		"",
		"the name of the product",
	)
	addProductCmd.MarkFlagRequired("name")

	addProductCmd.Flags().StringVar(
		&productSlug,
		"slug",
		"",
		"the slug of the product",
	)
	addProductCmd.MarkFlagRequired("slug")

	addProductCmd.Flags().StringVar(
		&productDescription,
		"description",
		"",
		"the description of the product",
	)

	addProductCmd.Flags().Int32Var(
		&productPriceCents,
		"price-cents",
		-1,
		"the price in cents",
	)
	addProductCmd.MarkFlagRequired("price-cents")

	addProductCmd.Flags().StringVar(
		&productCategoryIDs,
		"category-ids",
		"",
		"comma-separated list of category IDs",
	)

	updateProductCmd.Flags().Int32Var(
		&productID,
		"id",
		-1,
		"the ID of the product to update",
	)
	updateProductCmd.MarkFlagRequired("id")

	updateProductCmd.Flags().StringVar(
		&productName,
		"name",
		"",
		"the name of the product",
	)
	updateProductCmd.MarkFlagRequired("name")

	updateProductCmd.Flags().StringVar(
		&productSlug,
		"slug",
		"",
		"the slug of the product",
	)
	updateProductCmd.MarkFlagRequired("slug")

	updateProductCmd.Flags().StringVar(
		&productDescription,
		"description",
		"",
		"the description of the product",
	)

	updateProductCmd.Flags().Int32Var(
		&productPriceCents,
		"price-cents",
		-1,
		"the price in cents",
	)
	updateProductCmd.MarkFlagRequired("price-cents")

	updateProductCmd.Flags().StringVar(
		&productCategoryIDs,
		"category-ids",
		"",
		"comma-separated list of category IDs",
	)

	rootCmd.AddCommand(addProductCmd)
	rootCmd.AddCommand(updateProductCmd)
}

var (
	productID int32
	productName string
	productSlug string
	productDescription string
	productPriceCents int32
	productCategoryIDs string
)

var addProductCmd = &cobra.Command{
	Use:   "add-product",
	Short: "add a product to the catalog",
	RunE:  addProductCmdRun,
}

func addProductCmdRun(cmd *cobra.Command, args []string) error {
	cats, err := parseInt32List(productCategoryIDs)
	if err != nil {
		return err
	}
	var desc *string
	if productDescription != "" {
		desc = &productDescription
	}
	body, err := json.Marshal(&catalog.CreateProductRequest{
		ProductProperties: catalog.ProductProperties{
			Name:        productName,
			Slug:        productSlug,
			Description: desc,
			PriceCents:  productPriceCents,
			CategoryIDs: cats,
		},
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s%s", serverURL, adminPath, "/catalog/products"),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return reqAndPrint(req)
}

var updateProductCmd = &cobra.Command{
	Use:   "update-product",
	Short: "update a product in the catalog",
	RunE:  updateProductCmdRun,
}

func updateProductCmdRun(cmd *cobra.Command, args []string) error {
	cats, err := parseInt32List(productCategoryIDs)
	if err != nil {
		return err
	}
	var desc *string
	if productDescription != "" {
		desc = &productDescription
	}
	body, err := json.Marshal(&catalog.UpdateProductRequest{
		ID: productID,
		ProductProperties: catalog.ProductProperties{
			Name:        productName,
			Slug:        productSlug,
			Description: desc,
			PriceCents:  productPriceCents,
			CategoryIDs: cats,
		},
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s%s%s/%d", serverURL, adminPath, "/catalog/products", productID),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return reqAndPrint(req)
}
