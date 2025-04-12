package api

import (
	"komari/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		session, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "Unauthorized"})
			c.Abort()
			return
		}

		// Komari is a single user system
		_, err = database.GetSession(session)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "Unauthorized."})
			c.Abort()
			return
		}

		c.Next()
	}
}
