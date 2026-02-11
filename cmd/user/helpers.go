package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
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
