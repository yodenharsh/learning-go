package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"restapi/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Bearer")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token := cookie.Value

		jwtSecret := os.Getenv("JWT_SECRET")

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})

		if errors.Is(err, jwt.ErrTokenExpired) {
			http.Error(w, "Token Expired", http.StatusUnauthorized)
			return
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			http.Error(w, "Malformed Token", http.StatusUnauthorized)
			return
		} else if err != nil {
			utils.ErrorHandler(err, "")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else if !parsedToken.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKey("role"), claims["role"])
		ctx = context.WithValue(ctx, ContextKey("expiresAt"), claims["exp"])
		ctx = context.WithValue(ctx, ContextKey("id"), claims["uid"])
		ctx = context.WithValue(ctx, ContextKey("username"), claims["username"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
