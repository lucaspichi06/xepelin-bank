package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Logger logs the deposits greater than $10k
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: verify the request
		log.Printf("Depósito mayor a $10,000 realizado")

		c.Next()
	}
}
