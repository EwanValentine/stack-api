// Authware
// 
// A generic authenticate middleware for validating request tokens
// for stack-api.
//
// Pass in any generic authentication service, which has a `VaidateToken` method
// which matches the interface `Authenticator`, to wrap around authenticated API endpoints.
// 
// @copyright Stack API - 2017
// @author    Ewan Valentine <ewan.valentine89@gmail.com>

package authware

import (
	"os"
	"strings"

	api "github.com/ewanvalentine/stack-api"
)

// Authenticator - 
type Authenticator interface {
	ValidateToken(string) (bool, error)
}

// Authware - 
type Authware struct {
	Service Authenticator
}

// RequireAuth - 
func (auth *Authware) RequireAuth(next api.Handler) api.Handler {
	return func(c *api.Context) {

		if os.Getenv("AUTH_ENABLED") == "false" {
			next(c)
			return
		}

		authHeader := c.Header("Authorization")

		// If blank, check lower case value
		if authHeader == "" {
			authHeader = c.Header("authorization")
		}

		// If still blank, return 401
		if authHeader == "" {
			c.JSON(api.D{"_message": "No token detected."}, 401)
			return
		}

		// Split header value
		token := strings.Split(authHeader, " ")

		// Validate token against service
		isValid, err := auth.Service.ValidateToken(token[1])

		if err != nil || isValid == false {
			c.JSON(
				api.D{"_message": "Unauthorized."},
				401,
			)
			return
		}

		next(c)
	}
}
