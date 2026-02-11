package comm

import "fmt"

type Page[T any] struct {
	Items   []T   `json:"items"`
	Page    int32 `json:"page"`
	PerPage int32 `json:"per_page"`
	Total   int64 `json:"total"`
}

var (
	InternalError = fmt.Errorf("internal server error")
)
