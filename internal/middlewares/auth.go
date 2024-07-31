package middlewares

import (
	"attendance/pkg"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		cookie, err := c.Request.Cookie("token")

		tokenString := ""
		if err != nil {

			if authHeader != "" {
				tokenString = authHeader[len("Bearer "):]
			} else {
				response := pkg.BuildResponse(http.StatusUnauthorized, pkg.Unauthorized, pkg.Null())
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		} else {
			tokenString = cookie.Value
		}

		claims, err := pkg.VerifyToken(tokenString, os.Getenv("JWT_SECRET_KEY"))
		// fmt.Println(tokenString, "token")
		if err != nil {
			// fmt.Println(err.Error(), "claims err")
			response := pkg.BuildResponse(http.StatusUnauthorized, pkg.Unauthorized, pkg.Null())
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		user := &pkg.AuthUser{
			ID:       claims.ID,
			Username: claims.Username,
		}

		c.Set("user", user)
		c.Next()
	}
}
