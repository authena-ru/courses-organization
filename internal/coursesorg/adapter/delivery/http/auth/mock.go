package auth

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	"github.com/authena-ru/courses-organization/internal/coursesorg/adapter/delivery/http/httperr"
)

func MockAuthHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var claims jwt.MapClaims
		token, err := request.ParseFromRequest(
			r,
			request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte("mock_secret"), nil
			},
			request.WithClaims(&claims),
		)
		if err != nil {
			httperr.BadRequest("unable-to-get-jwt", err, w, r)
			return
		}

		if !token.Valid {
			httperr.Unauthorized("invalid-jwt", nil, w, r)
			return
		}

		id, idIsOk := claims["id"].(string)
		role, roleIsOk := claims["role"].(string)

		if !idIsOk || !roleIsOk {
			httperr.Unauthorized("invalid-jwt-claims", nil, w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey, User{
			ID:   id,
			Role: Role(role),
		})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
