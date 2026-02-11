package main

import (
	"io"
	"fmt"
	"net/http"
	"math/rand/v2"
	"strings"

	"github.com/Coosis/go-eshop/internal/catalog"
)

func makeProduct(name string, desc string) catalog.ProductProperties {
	price := rand.Int32N(3000);
	return catalog.ProductProperties{
		Name:        name,
		Slug:        strings.ToLower(strings.ReplaceAll(name, " ", "-")),
		Description: &desc,
		PriceCents:  price,
		CategoryIDs: []int32{},
	}
}

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
