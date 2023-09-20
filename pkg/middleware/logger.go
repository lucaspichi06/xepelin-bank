package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	"net/http"
)

// Logger logs the deposits greater than $10k
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		trType := c.Value("transaction_type")
		trAmount := c.Value("transaction_amount").(float64)
		if trType == string(domain.Deposit) && trAmount > 10000 {
			if c.Writer.Status() != http.StatusCreated {
				fmt.Println("Deposit greater than $10,000 failed to process")
				return
			}

			fmt.Println("Deposit greater than $10,000 processed successfully")
		}
	}
}
