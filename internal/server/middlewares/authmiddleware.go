package middlewares

import (
	"api-go/internal/auth"
	"api-go/internal/utils"
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Chave de contexto para evitar colisões de nomes
type contextKey string

const userContextKey = contextKey("user")

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// AuthMiddleware verifica o token JWT e, se válido, adiciona os claims ao contexto da requisição.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Pega o cabeçalho Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header not found")
			return
		}

		// 2. Valida o formato "Bearer <token>"
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenString := headerParts[1]
		claims := &auth.Claims{}

		// 3. Parse e valida o token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Verifica se o método de assinatura é HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token: "+err.Error())
			return
		}

		if !token.Valid {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// 4. Adiciona os claims ao contexto da requisição para uso posterior
		ctx := context.WithValue(r.Context(), userContextKey, claims)

		// 5. Chama o próximo handler com o novo contexto
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext extrai os claims do usuário do contexto da requisição
func GetUserFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(userContextKey).(*auth.Claims)
	return claims, ok
}
