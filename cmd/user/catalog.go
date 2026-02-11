package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	productsPath string
)

func init() {
	getProductCmd.Flags().Int32Var(
		&getProductID,
		"id",
		-1,
		"the product id of the product to retrieve",
	)
	getProductCmd.Flags().StringVar(
		&getProductSlug,
		"slug",
		"",
		"the slug of the product to retrieve",
	)
	getProductCmd.MarkFlagsMutuallyExclusive("id", "slug")
	getProductCmd.MarkFlagsOneRequired("id", "slug")

	listProductsCmd.PersistentFlags().Int32Var(
		&listPage,
		"page",
		1,
		"the page number to retrieve",
	)
	listProductsCmd.PersistentFlags().Int32Var(
		&listPerPage,
		"per-page",
		250,
		"number of products per page",
	)
	listProductsCmd.PersistentFlags().Int32Var(
		&listMinPrice,
		"min-price",
		0,
		"minimum price (in cents) of products to list",
	)
	listProductsCmd.PersistentFlags().Int32Var(
		&listMaxPrice,
		"max-price",
		100000,
		"maximum price (in cents) of products to list",
	)
	listProductsCmd.PersistentFlags().Int32Var(
		&listCategoryID,
		"category-id",
		0,
		"category ID to filter products by",
	)

	rootCmd.AddCommand(listProductsCmd)
	rootCmd.AddCommand(getProductCmd)
	rootCmd.PersistentFlags().StringVar(
		&productsPath,
		"products-path",
		"/v1/catalog/products",
		"the base path for catalog product endpoints",
	)
}
var (
	listPage int32
	listPerPage int32
	listMinPrice int32
	listMaxPrice int32
	listCategoryID int32
)

var listProductsCmd = &cobra.Command{
	Use:   "list-products",
	Short: "list products in the catalog",
	RunE: listProducts,
}

func listProducts(cmd *cobra.Command, args []string) error {
	req, err := http.NewRequest(
		"GET",
		serverURL+productsPath,
		nil,
	)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("page", fmt.Sprintf("%d", listPage))
	q.Add("per_page", fmt.Sprintf("%d", listPerPage))
	q.Add("min_price", fmt.Sprintf("%d", listMinPrice))
	q.Add("max_price", fmt.Sprintf("%d", listMaxPrice))
	q.Add("category_id", fmt.Sprintf("%d", listCategoryID))
	req.URL.RawQuery = q.Encode()

	return reqAndPrint(req)
}

var (
	getProductID int32
	getProductSlug string
)

var getProductCmd = &cobra.Command{
	Use:   "get-product",
	Short: "get a product by ID or slug",
	RunE: getProduct,
}

func getProduct(cmd *cobra.Command, args []string) error {
	if getProductID != -1 {
		return getProductByID(getProductID)
	} else if getProductSlug != "" {
		return getProductBySlug(getProductSlug)
	} else {
		return fmt.Errorf("either id or slug must be provided!")
	}
}

func getProductByID(id int32) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s%s/%d", serverURL, productsPath, id),
		nil,
	)
	if err != nil {
		return err
	}
	return reqAndPrint(req)
}	

func getProductBySlug(slug string) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s%s/slug/%s", serverURL, productsPath, slug),
		nil,
	)
	if err != nil {
		return err
	}
	return reqAndPrint(req)
}
