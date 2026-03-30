package middleware

import (
	"context"
	"net/http"
	"strconv"

	"awesomeProject/internal/models"
	"awesomeProject/internal/repository"

	"github.com/sirupsen/logrus"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

func AuthMiddleware(userRepo *repository.UserRepository, logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDStr := r.Header.Get("X-User-Id")
			if userIDStr == "" {
				if r.URL.Path == "/login" || r.URL.Path == "/health" {
					next.ServeHTTP(w, r)
					return
				}
				http.Error(w, "X-User-Id header is required", http.StatusUnauthorized)
				return
			}

			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}

			user, err := userRepo.GetByID(userID)
			if err != nil {
				logger.WithError(err).Error("Failed to get user")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if user == nil {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserContextKey).(*models.User)
		if !ok || user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if user.Role != "admin" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func GetUserFromContext(r *http.Request) *models.User {
	user, ok := r.Context().Value(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
