package middleware

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/Coosis/go-eshop/internal/auth"
)

func WithDebugUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set(auth.UserIDKey, int64(1))
		return next(c)
	}
}

func WithUserID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		raw := c.Request().Header.Get("X-User-ID")
		if raw == "" {
			return echo.NewHTTPError(401, "missing X-User-ID header")
		}
		uid, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || uid <= 0 {
			return echo.NewHTTPError(401, "invalid X-User-ID header")
		}
		c.Set(auth.UserIDKey, uid)
		return next(c)
	}
}
