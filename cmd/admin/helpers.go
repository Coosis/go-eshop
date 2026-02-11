package main

import (
	"io"
	"fmt"
	"net/http"
	"math/rand/v2"
	"strings"
	"strconv"

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

func parseInt32List(csv string) ([]int32, error) {
	if strings.TrimSpace(csv) == "" {
		return []int32{}, nil
	}
	parts := strings.Split(csv, ",")
	out := make([]int32, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.ParseInt(p, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid int value: %q", p)
		}
		out = append(out, int32(n))
	}
	return out, nil
}
