package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	productsPath string
	categoriesPath string
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

	addPagingFlags(listProductsCmd, &listPage, &listPerPage, 250)
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

	addPagingFlags(listCategoriesCmd, &listCategoriesPage, &listCategoriesPerPage, 250)

	getCategoryCmd.Flags().Int32Var(
		&getCategoryID,
		"id",
		-1,
		"the category id to retrieve",
	)
	getCategoryCmd.Flags().StringVar(
		&getCategorySlug,
		"slug",
		"",
		"the slug of the category to retrieve",
	)
	getCategoryCmd.MarkFlagsMutuallyExclusive("id", "slug")
	getCategoryCmd.MarkFlagsOneRequired("id", "slug")

	listCategoryProductsCmd.Flags().Int32Var(
		&listCategoryProductsID,
		"id",
		-1,
		"the category id to list products for",
	)
	listCategoryProductsCmd.MarkFlagRequired("id")
	addPagingFlags(listCategoryProductsCmd, &listCategoryProductsPage, &listCategoryProductsPerPage, 250)

	rootCmd.AddCommand(listProductsCmd)
	rootCmd.AddCommand(getProductCmd)
	rootCmd.AddCommand(listCategoriesCmd)
	rootCmd.AddCommand(getCategoryCmd)
	rootCmd.AddCommand(listCategoryProductsCmd)
	rootCmd.PersistentFlags().StringVar(
		&productsPath,
		"products-path",
		"/v1/catalog/products",
		"the base path for catalog product endpoints",
	)
	rootCmd.PersistentFlags().StringVar(
		&categoriesPath,
		"categories-path",
		"/v1/catalog/categories",
		"the base path for catalog category endpoints",
	)
}
var (
	listPage int32
	listPerPage int32
	listMinPrice int32
	listMaxPrice int32
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

var (
	listCategoriesPage int32
	listCategoriesPerPage int32
)

var listCategoriesCmd = &cobra.Command{
	Use:   "list-categories",
	Short: "list categories in the catalog",
	RunE:  listCategories,
}

func listCategories(cmd *cobra.Command, args []string) error {
	req, err := http.NewRequest(
		"GET",
		serverURL+categoriesPath,
		nil,
	)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("page", fmt.Sprintf("%d", listCategoriesPage))
	q.Add("per_page", fmt.Sprintf("%d", listCategoriesPerPage))
	req.URL.RawQuery = q.Encode()
	return reqAndPrint(req)
}

var (
	getCategoryID int32
	getCategorySlug string
)

var getCategoryCmd = &cobra.Command{
	Use:   "get-category",
	Short: "get a category by ID or slug",
	RunE:  getCategory,
}

func getCategory(cmd *cobra.Command, args []string) error {
	if getCategoryID != -1 {
		return getCategoryByID(getCategoryID)
	} else if getCategorySlug != "" {
		return getCategoryBySlug(getCategorySlug)
	}
	return fmt.Errorf("either id or slug must be provided")
}

func getCategoryByID(id int32) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s%s/%d", serverURL, categoriesPath, id),
		nil,
	)
	if err != nil {
		return err
	}
	return reqAndPrint(req)
}

func getCategoryBySlug(slug string) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s%s/slug/%s", serverURL, categoriesPath, slug),
		nil,
	)
	if err != nil {
		return err
	}
	return reqAndPrint(req)
}

var (
	listCategoryProductsID int32
	listCategoryProductsPage int32
	listCategoryProductsPerPage int32
)

var listCategoryProductsCmd = &cobra.Command{
	Use:   "list-category-products",
	Short: "list products under a category",
	RunE:  listCategoryProducts,
}

func listCategoryProducts(cmd *cobra.Command, args []string) error {
	if listCategoryProductsID <= 0 {
		return fmt.Errorf("category id must be provided")
	}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s%s/%d/products", serverURL, categoriesPath, listCategoryProductsID),
		nil,
	)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("page", fmt.Sprintf("%d", listCategoryProductsPage))
	q.Add("per_page", fmt.Sprintf("%d", listCategoryProductsPerPage))
	req.URL.RawQuery = q.Encode()
	return reqAndPrint(req)
}
