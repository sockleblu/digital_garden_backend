package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/sockleblu/digital_garden_backend/graph/helpers"
	"github.com/sockleblu/digital_garden_backend/graph/model"
	"gorm.io/gorm"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// Middleware decodes the share session cookie and packs the session into context
func Middleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if header == "" || header == "null" {
				next.ServeHTTP(w, r)
				return
			}

			// Validate jwt token
			tokenStr := strings.Split(header, "Bearer ")[1]
			if tokenStr == "" || tokenStr == "null" {
				log.Printf("tokenStr equals %s", tokenStr)
				// next.ServeHTTP(w,r)
				// next.ServeHTTP()
				return
			}
			log.Printf("tokenStr: %s\n", tokenStr)
			username, err := helpers.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// get the user from the database
			id, err := helpers.GetUserIdByUsername(db, username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user := model.User{
				ID:       id,
				Username: username,
			}

			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *model.User {
	raw, _ := ctx.Value(userCtxKey).(*model.User)
	return raw
}
