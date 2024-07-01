package golemrouter

import (
	"context"

	clerk "github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.infratographer.com/x/echox/echozap"
	"go.uber.org/zap"
)

func New(log *zap.Logger) (*echo.Echo, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(echozap.Middleware(log))

	return e, nil
}

func ClerkAuth(secret, token string) echo.MiddlewareFunc {
	clerk.SetKey(secret)
	f := clerkhttp.RequireHeaderAuthorization(clerkhttp.AuthorizedParty(func(s string) bool {
		claims, err := jwt.Verify(context.Background(), &jwt.VerifyParams{
			Token: token,
		})
		if err != nil {
			return false
		}
		if claims.SessionID == "" {
			return false
		}
		return true
	}))
	return echo.WrapMiddleware(f)
}
