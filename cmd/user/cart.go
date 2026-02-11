package main

import (
	"bytes"
	// "context"
	"encoding/json"
	"fmt"

	// "math/rand/v2"
	"net/http"
	// "os"
	// "strings"

	// "github.com/Coosis/go-eshop/internal/catalog"
	"github.com/Coosis/go-eshop/internal/cart"
	"github.com/spf13/cobra"
)

var (
	cartPath string
)

func init() {
	GetCartCmd.PersistentFlags().StringVar(
		&cartPath,
		"cart-path",
		"/v1/cart",
		"the base path of the cart service",
	)
	GetCartCmd.Flags().Int32VarP(
		&getCartPage,
		"page",
		"p",
		1,
		"the page number to retrieve",
	)
	GetCartCmd.Flags().Int32Var(
		&getCartPerPage,
		"per-page",
		250,
		"number of items per page",
	)

	AddItemCmd.Flags().Int32Var(
		&addItemProductID,
		"product-id",
		-1,
		"the product id of the item to add to the cart",
	)
	AddItemCmd.MarkFlagRequired("product-id")
	AddItemCmd.Flags().Int32VarP(
		&addItemQuantity,
		"quantity",
		"q",
		1,
		"the quantity of the item to add to the cart",
	)
	AddItemCmd.Flags().Int32VarP(
		&addItemPage,
		"page",
		"p",
		1,
		"the page number to retrieve",
	)
	AddItemCmd.Flags().Int32Var(
		&addItemPerPage,
		"per-page",
		250,
		"number of items per page",
	)

	updateCartItemCmd.Flags().Int32Var(
		&updateItemProductID,
		"product-id",
		-1,
		"the product id of the item to update in the cart",
	)
	updateCartItemCmd.MarkFlagRequired("product-id")
	updateCartItemCmd.Flags().Int32VarP(
		&updateItemQuantity,
		"quantity",
		"q",
		1,
		"the new quantity of the item in the cart",
	)
	updateCartItemCmd.Flags().Int32VarP(
		&updateItemPage,
		"page",
		"p",
		1,
		"the page number to retrieve",
	)
	updateCartItemCmd.Flags().Int32Var(
		&updateItemPerPage,
		"per-page",
		250,
		"number of items per page",
	)

	changeCartItemQtyCmd.Flags().Int32Var(
		&changeItemProductID,
		"product-id",
		-1,
		"the product id of the item to update in the cart",
	)
	changeCartItemQtyCmd.MarkFlagRequired("product-id")
	changeCartItemQtyCmd.Flags().Int32Var(
		&changeItemDelta,
		"delta",
		0,
		"the delta to apply to item quantity",
	)
	changeCartItemQtyCmd.MarkFlagRequired("delta")
	changeCartItemQtyCmd.Flags().Int32VarP(
		&changeItemPage,
		"page",
		"p",
		1,
		"the page number to retrieve",
	)
	changeCartItemQtyCmd.Flags().Int32Var(
		&changeItemPerPage,
		"per-page",
		250,
		"number of items per page",
	)

	removeCartItemCmd.Flags().Int32Var(
		&removeItemProductID,
		"product-id",
		-1,
		"the product id of the item to remove from the cart",
	)
	removeCartItemCmd.MarkFlagRequired("product-id")
	removeCartItemCmd.Flags().Int32VarP(
		&removeItemPage,
		"page",
		"p",
		1,
		"the page number to retrieve",
	)
	removeCartItemCmd.Flags().Int32Var(
		&removeItemPerPage,
		"per-page",
		250,
		"number of items per page",
	)

	rootCmd.AddCommand(GetCartCmd)
	rootCmd.AddCommand(AddItemCmd)
	rootCmd.AddCommand(clearCartCmd)
	rootCmd.AddCommand(updateCartItemCmd)
	rootCmd.AddCommand(changeCartItemQtyCmd)
	rootCmd.AddCommand(removeCartItemCmd)
}

var (
	getCartPage int32
	getCartPerPage int32
)

var GetCartCmd = &cobra.Command{
	Use:   "get-cart",
	Short: "retrieve the current user's cart",
	RunE:  GetCart,
}

func GetCart(cmd *cobra.Command, args []string) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s%s", serverURL, cartPath),
		nil,
	)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("page", fmt.Sprintf("%d", getCartPage))
	q.Add("per_page", fmt.Sprintf("%d", getCartPerPage))
	req.URL.RawQuery = q.Encode()

	return reqAndPrint(req)
}

var (
	addItemProductID int32
	addItemQuantity  int32
	addItemPage      int32
	addItemPerPage   int32
)

var AddItemCmd = &cobra.Command{
	Use:   "add-cart-item",
	Short: "add an item to the current user's cart",
	RunE:  AddCartItem,
}

func AddCartItem(cmd *cobra.Command, args []string) error {
	b, err := json.Marshal(cart.AddCartItemRequest{
		ProductID: addItemProductID,
		Quantity:  addItemQuantity,
		CartPaging: cart.CartPaging{
			Page:    addItemPage,
			PerPage: addItemPerPage,
		},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s/items", serverURL, cartPath),
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return reqAndPrint(req)
}

var clearCartCmd = &cobra.Command{
	Use:   "clear-cart",
	Short: "clear all items from the current user's cart",
	RunE: clearCart,
}

func clearCart(cmd *cobra.Command, args []string) error {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s%s", serverURL, cartPath),
		nil,
	)
	if err != nil {
		return err
	}
	return reqAndPrint(req)
}

var (
	updateItemProductID int32
	updateItemQuantity  int32
	updateItemPage      int32
	updateItemPerPage   int32
)

var updateCartItemCmd = &cobra.Command{
	Use:   "update-cart-item",
	Short: "update an item in the current user's cart",
	RunE:  UpdateCartItem,
}

func UpdateCartItem(cmd *cobra.Command, args []string) error {
	orig_req, err := json.Marshal(cart.UpdateCartItemRequest{
		ProductID: updateItemProductID,
		Quantity:  updateItemQuantity,
		CartPaging: cart.CartPaging{
			Page:    updateItemPage,
			PerPage: updateItemPerPage,
		},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s%s/items/%d", serverURL, cartPath, updateItemProductID),
		bytes.NewReader(orig_req),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return reqAndPrint(req)
}

var (
	changeItemProductID int32
	changeItemDelta int32
	changeItemPage int32
	changeItemPerPage int32
)

var changeCartItemQtyCmd = &cobra.Command{
	Use:   "change-cart-item-qty",
	Short: "change quantity of an item in the cart by delta",
	RunE:  changeCartItemQty,
}

func changeCartItemQty(cmd *cobra.Command, args []string) error {
	req, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s%s/items/%d", serverURL, cartPath, changeItemProductID),
		nil,
	)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("delta", fmt.Sprintf("%d", changeItemDelta))
	q.Add("page", fmt.Sprintf("%d", changeItemPage))
	q.Add("per_page", fmt.Sprintf("%d", changeItemPerPage))
	req.URL.RawQuery = q.Encode()
	return reqAndPrint(req)
}

var (
	removeItemProductID int32
	removeItemPage int32
	removeItemPerPage int32
)

var removeCartItemCmd = &cobra.Command{
	Use:   "remove-cart-item",
	Short: "remove an item from the cart",
	RunE:  removeCartItem,
}

func removeCartItem(cmd *cobra.Command, args []string) error {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s%s/items/%d", serverURL, cartPath, removeItemProductID),
		nil,
	)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("page", fmt.Sprintf("%d", removeItemPage))
	q.Add("per_page", fmt.Sprintf("%d", removeItemPerPage))
	req.URL.RawQuery = q.Encode()
	return reqAndPrint(req)
}
