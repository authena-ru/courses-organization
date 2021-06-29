package auth

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/httperr"
	"github.com/authena-ru/courses-organization/internal/domain/course"
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

		jwtID, idIsOk := claims["id"].(string)
		jwtAcademicType, roleIsOk := claims["type"].(string)

		if !idIsOk || !roleIsOk {
			httperr.Unauthorized("invalid-jwt-claims", nil, w, r)

			return
		}

		academicType := course.NewAcademicTypeFromString(jwtAcademicType)
		academic, err := course.NewAcademic(jwtID, academicType)
		if err != nil {
			httperr.Unauthorized("invalid-academic-type", nil, w, r)
		}

		ctx := context.WithValue(r.Context(), academicCtxKey, academic)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
