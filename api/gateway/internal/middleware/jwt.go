package middleware

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"td27/api/gateway/internal/svc"

	"td27/pkg"
)

var (
	UserIdKey   = "userId"
	UsernameKey = "username"
	RoleIdsKey  = "roleIds"
)

type JwtMiddleware struct {
	svcCtx *svc.ServiceContext
}

func NewJwtMiddleware(svcCtx *svc.ServiceContext) *JwtMiddleware {
	return &JwtMiddleware{svcCtx: svcCtx}
}

func (m *JwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJson(w, http.StatusUnauthorized, pkg.Error(401, "missing authorization header"))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			writeJson(w, http.StatusUnauthorized, pkg.Error(401, "invalid authorization format"))
			return
		}

		claims, err := m.parseToken(tokenStr)
		if err != nil {
			writeJson(w, http.StatusUnauthorized, pkg.Error(401, "invalid or expired token"))
			return
		}

		tokenHash := HashToken(tokenStr)
		if GetBlocklist().IsBlocklisted(tokenHash) {
			writeJson(w, http.StatusUnauthorized, pkg.Error(401, "token has been invalidated"))
			return
		}

		ctx := context.WithValue(r.Context(), UserIdKey, claims["userId"])
		ctx = context.WithValue(ctx, UsernameKey, claims["username"])
		ctx = context.WithValue(ctx, RoleIdsKey, claims["roleIds"])
		next(w, r.WithContext(ctx))
	}
}

func (m *JwtMiddleware) parseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(m.svcCtx.Config.Auth.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

func writeJson(w http.ResponseWriter, status int, resp *pkg.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
