package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/emvi/hide"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

type AuthContext struct {
	UserID hide.ID
	Secure bool
}

// UserClaim structure
type UserClaim struct {
	UserID hide.ID `json:"userId"`
	jwt.StandardClaims
}

// Middleware decodes the authorization header
func Middleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userJSON := c.Request.Header.Get("user")
		if userJSON == "" {
			c.Next()
			return
		}

		userClaim := &UserClaim{
			StandardClaims: jwt.StandardClaims{},
		}
		err := json.Unmarshal([]byte(userJSON), userClaim)
		if err != nil {
			c.Next()
			return
		}

		log.Println(int64(userClaim.UserID))

		// get user from database

		// var user model.User
		// err = db.Model(&model.User{}).Where(int64(userID)).First(&user).Error
		// if err != nil || user.ID == 0 {
		// 	c.Next()
		// 	return
		// }

		ctx := context.WithValue(c.Request.Context(), userCtxKey, AuthContext{
			UserID: userClaim.UserID,
		})

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// For find the user from the context. Middleware must have run
func For(ctx context.Context) AuthContext {
	raw, _ := ctx.Value(userCtxKey).(AuthContext)
	return raw
}

func splitToken(header string) (string, error) {
	splitToken := strings.Split(header, "Bearer")

	if len(splitToken) != 2 || len(splitToken[1]) < 2 {
		return "", fmt.Errorf("bad token format")
	}

	return strings.TrimSpace(splitToken[1]), nil
}
