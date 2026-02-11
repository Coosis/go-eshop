package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/spf13/cobra"
)

func reqAndPrint(req *http.Request) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	fmt.Println("response status:", resp.Status)
	msg, _ := io.ReadAll(resp.Body)
	fmt.Println("response body:", string(msg))

	return nil
}

func generateChars(l int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, l)
	for i := range b {
		id := rand.Intn(len(charset))
		b[i] = charset[id]
	}
	return string(b)
}

func addPagingFlags(cmd *cobra.Command, page *int32, perPage *int32, defaultPerPage int32) {
	cmd.Flags().Int32VarP(
		page,
		"page",
		"p",
		1,
		"the page number to retrieve",
	)
	cmd.Flags().Int32Var(
		perPage,
		"per-page",
		defaultPerPage,
		"number of items per page",
	)
}
