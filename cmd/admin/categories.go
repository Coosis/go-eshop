package main

import (
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/Coosis/go-eshop/internal/catalog"
)

func init() {
	addCategoryCmd.Flags().StringVar(
		&categoryName,
		"name",
		"Add Category",
		"the name of the category to create",
	)
	addCategoryCmd.MarkFlagRequired("name")

	addCategoryCmd.Flags().StringVar(
		&categorySlug,
		"slug",
		"category",
		"the slug of the category to create",
	)
	addCategoryCmd.MarkFlagRequired("slug")

	addCategoryCmd.Flags().Int32Var(
		&categoryParentID,
		"parent-id",
		-1,
		"the parent ID of the category to create",
	)
	rootCmd.AddCommand(addCategoryCmd)

	updateCategoryCmd.Flags().Int32Var(
		&categoryID,
		"id",
		-1,
		"the category ID to update",
	)
	updateCategoryCmd.MarkFlagRequired("id")

	updateCategoryCmd.Flags().StringVar(
		&categoryName,
		"name",
		"",
		"the name of the category to update",
	)
	updateCategoryCmd.MarkFlagRequired("name")

	updateCategoryCmd.Flags().StringVar(
		&categorySlug,
		"slug",
		"",
		"the slug of the category to update",
	)
	updateCategoryCmd.MarkFlagRequired("slug")

	updateCategoryCmd.Flags().Int32Var(
		&categoryParentID,
		"parent-id",
		-1,
		"the parent ID of the category to update",
	)
	rootCmd.AddCommand(updateCategoryCmd)
}

var (
	categoryID int32
	categoryName string
	categorySlug string
	categoryParentID int32
)

var addCategoryCmd = &cobra.Command{
	Use:  "add-category",
	Short: "add categories to the catalog",
	Long: `Add categories to the catalog service.`,
	RunE: addCategory,
}

func addCategory(cmd *cobra.Command, args []string) error {
	var pid *int32
	pid = nil
	if categoryParentID != -1 {
		pid = &categoryParentID
	}
	orig_req, err := json.Marshal(&catalog.CreateCategoryRequest{
		CategoryProperties: catalog.CategoryProperties{
			Name:       categoryName,
			Slug:       categorySlug,
			ParentID:   pid,
		},
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s%s", serverURL, adminPath, "/catalog/categories"),
		bytes.NewBuffer(orig_req),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return reqAndPrint(req)
}

var updateCategoryCmd = &cobra.Command{
	Use:  "update-category",
	Short: "update categories in the catalog",
	Long: `Update categories in the catalog service.`,
	RunE: updateCategory,
}

func updateCategory(cmd *cobra.Command, args []string) error {
	var pid *int32
	if categoryParentID != -1 {
		pid = &categoryParentID
	}
	orig_req, err := json.Marshal(&catalog.UpdateCategoryRequest{
		ID: categoryID,
		CategoryProperties: catalog.CategoryProperties{
			Name:     categoryName,
			Slug:     categorySlug,
			ParentID: pid,
		},
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s%s%s/%d", serverURL, adminPath, "/catalog/categories", categoryID),
		bytes.NewBuffer(orig_req),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return reqAndPrint(req)
}
