package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/churilovmn1/workout-tracker/internal/models"
	"github.com/churilovmn1/workout-tracker/internal/service"
)

type contextKey string

const (
	ctxUserID contextKey = "user_id"
	ctxRole   contextKey = "role"
)

// AuthMiddleware проверяет JWT-токен из заголовка Authorization.
//
// Схема работы:
//  1. Извлекаем токен из "Bearer <token>"
//  2. Парсим и валидируем подпись через AuthService (HMAC-SHA256)
//  3. Кладём user_id и role в context — handler достаёт их через getUserID()
//
// Если токен невалиден — возвращаем 401 и прерываем цепочку middleware.
func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				writeError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok {
				writeError(w, http.StatusUnauthorized, "invalid authorization format")
				return
			}

			claims, err := authService.ParseToken(token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), ctxUserID, claims.UserID)
			ctx = context.WithValue(ctx, ctxRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AdminOnly проверяет роль из context — пропускает только admin.
// Всегда идёт после AuthMiddleware: к этому моменту role уже в context.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value(ctxRole).(models.Role)
		if role != models.RoleAdmin {
			writeError(w, http.StatusForbidden, "admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// getUserID извлекает ID аутентифицированного пользователя из context.
func getUserID(r *http.Request) int {
	id, _ := r.Context().Value(ctxUserID).(int)
	return id
}
