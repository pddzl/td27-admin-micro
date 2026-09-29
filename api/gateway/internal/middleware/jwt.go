package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"

	"td27/api/gateway/internal/svc"
	"td27/pkg/api"
	"td27/rpc/basis/types/sysManagement/permission_pb"
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
		authHeader := r.Header.Get("x-token")
		if authHeader == "" {
			api.FailWithRequest(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := m.parseToken(tokenStr)
		if err != nil {
			api.FailWithRequest(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		tokenHash := HashToken(tokenStr)
		if GetBlocklist().IsBlocklisted(tokenHash) {
			api.FailWithRequest(w, http.StatusUnauthorized, "token has been invalidated")
			return
		}

		// Optional per-request RBAC enforcement against the rpc Casbin enforcer
		if m.svcCtx.Config.Auth.EnforceRbac {
			if !m.authorize(r, claims) {
				api.FailWithRequest(w, http.StatusForbidden, "no permission to perform this operation")
				return
			}
		}

		ctx := context.WithValue(r.Context(), UserIdKey, claims["userId"])
		ctx = context.WithValue(ctx, UsernameKey, claims["username"])
		ctx = context.WithValue(ctx, RoleIdsKey, claims["roleIds"])
		next(w, r.WithContext(ctx))
	}
}

// selfServicePrefixes are routes any authenticated user may call regardless of
// RBAC policies: strictly self-scoped operations where the target identity is
// overridden server-side from the JWT (never the request body).
var selfServicePrefixes = []string{
	"/user/modify-password",
	"/role/elevation/request",
	"/role/elevation/my",
}

func isSelfService(path string) bool {
	for _, prefix := range selfServicePrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// authorize checks the request against the rpc Casbin enforcer. The user id is
// sent along so the rpc side resolves roles live from the database instead of
// trusting the (possibly stale) role ids embedded in the JWT. Fails closed: an
// rpc error denies the request rather than silently allowing it.
func (m *JwtMiddleware) authorize(r *http.Request, claims jwt.MapClaims) bool {
	// Self-service routes skip the policy check: without this, a plain user
	// could not even change their own password or file an elevation request.
	if isSelfService(r.URL.Path) {
		return true
	}

	var userId int64
	if f, ok := claims["userId"].(float64); ok {
		userId = int64(f)
	}

	req := &permission_pb.CheckPermissionReq{
		RoleIds:  roleIdsFromClaims(claims["roleIds"]),
		UserId:   userId,
		Resource: r.URL.Path,
		Action:   httpMethodToAction(r.Method, r.URL.Path),
	}

	resp, err := m.svcCtx.PermissionClient.CheckPermission(r.Context(), req)
	if err != nil {
		logx.Errorf("rbac check failed for %s %s, denying: %v", r.Method, r.URL.Path, err)
		return false
	}
	return resp.GetAllowed()
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
