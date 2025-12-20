package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:  "session_id",
			Value: "abc123",
		})
		return nil
	})
	g := e.Group("/secret")
	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session_id")
			if err != nil || cookie.Value != "abc123" {
				return c.String(http.StatusUnauthorized, "Unauthorized")
			}
			return next(c)
		}
	})
	g.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the secret area!")
	})
	if err := e.Start(":8144"); err != nil {
		panic(err)
	}
}
