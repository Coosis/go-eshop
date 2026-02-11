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
}

var (
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
