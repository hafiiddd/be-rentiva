package middleware

import (
	"back-end/domain/model"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
)

func JWTMiddlewareConfig() echojwt.Config {
	secretFromEnv := strings.TrimSpace(os.Getenv("JWT_SECRET_KEY"))
	jwtSecretKey := []byte(secretFromEnv)
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(model.Auth)
		},
		SigningKey: jwtSecretKey,
		ErrorHandler: func(c echo.Context, err error) error {
			log.Printf("ERROR JWT: Terjadi error pada middleware: %v", err)
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Token tidak valid atau Anda tidak memiliki izin.",
			})
		},
	}

}
func SkipOptions(m echo.MiddlewareFunc) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if c.Request().Method == http.MethodOptions {
                return next(c)
            }
            return m(next)(c)
        }
    }
}