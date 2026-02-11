package auth

import "github.com/labstack/echo/v4"

const (
	UserIDKey = "user_id"
)

func RequireUserID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uid, ok := c.Get(UserIDKey).(int32)
		if !ok || uid <= 0 {
			return echo.NewHTTPError(401, "unauthorized: missing or invalid user ID")
		}
		return next(c)
	}
}

func DevWithUserID(userID int32) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(UserIDKey, userID)
			return next(c)
		}
	}
}
