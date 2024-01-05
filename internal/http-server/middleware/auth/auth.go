package middlewareAuth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	restJwt "url-shortener/internal/lib/jwt"
	"url-shortener/internal/lib/logger/sl"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrFailedIsAdminCheck = errors.New("failed to check if user is admin")
)

type PermissionProvider interface {
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"uid"`
}

// New creates new auth middleware.
func New(
	log *slog.Logger,
	appSecret string,
	permProvider PermissionProvider,
) func(next http.Handler) http.Handler {
	const op = "middleware.auth.New"

	log = log.With(slog.String("op", op))

	// Возвращаем функцию-обработчик
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем JWT-токен из запроса
			tokenStr := restJwt.ExtractBearerToken(r)
			if tokenStr == "" {
				ctx := context.WithValue(r.Context(), "status", http.StatusUnauthorized)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			// Парсим и валидируем токен, использая appSecret
			parsedToken, err := Parse(tokenStr, appSecret)
			if err != nil {
				log.Warn("failed to parse token", sl.Err(err))
				// But if token is invalid, we shouldn't handle request
				ctx := context.WithValue(r.Context(), "status", http.StatusUnauthorized)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}
			claims, ok := parsedToken.Claims.(*tokenClaims)
			if !ok {
				log.Warn("failed to parse token claims")
				ctx := context.WithValue(r.Context(), "status", http.StatusUnauthorized)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			log.Info("user authorized", slog.Any("claims", claims))

			// // Отправляем запрос для проверки, является ли пользователь админов
			isAdmin, err := permProvider.IsAdmin(r.Context(), int64(claims.UserId))
			if err != nil {
				log.Error("failed to check if user is admin", sl.Err(err))

				ctx := context.WithValue(r.Context(), "isAdmin", "error")
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			// Полученны данные сохраняем в контекст,
			// откуда его смогут получить следующие хэндлеры.
			ctx := context.WithValue(r.Context(), "uid", claims.UserId)
			ctx = context.WithValue(r.Context(), "isAdmin", isAdmin)
			ctx = context.WithValue(r.Context(), "status", http.StatusOK)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Parse(tokenStr string, appSecret string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(appSecret), nil
	})
}
