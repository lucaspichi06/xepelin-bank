package handler

import (
	"errors"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	tran "github.com/lucaspichi06/xepelin-bank/internal/transaction"
	custom_errors "github.com/lucaspichi06/xepelin-bank/pkg/errors"
	"github.com/lucaspichi06/xepelin-bank/pkg/web"

	"github.com/gin-gonic/gin"
	"net/http"
)

type Transactions interface {
	Process() gin.HandlerFunc
}

type transaction struct {
	s tran.Service
}

func NewTransactionsHandler(s tran.Service) Transactions {
	return &transaction{
		s: s,
	}
}

// Process	godoc
// @Summary	Process a received transaction
// @Tags	Transaction
// @Description	process a received transaction
// @Accept	json
// @Produce	json
// @Param	token	header	string	true	"token"
// @Param	transaction		body 	domain.Transaction	true	"Transaction to process"
// @Success 201	{object}	web.Response
// @Failure	400	{object}	web.ErrorResponse
// @Failure	500	{object}	web.ErrorResponse
// @Router	/transactions	[post]
func (t transaction) Process() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tr domain.Transaction
		err := c.ShouldBindJSON(&tr)
		if err != nil {
			web.Failure(c, http.StatusBadRequest, custom_errors.ErrInvalidJSON)
			return
		}

		c.Set("transaction_type", tr.Type)
		c.Set("transaction_amount", tr.Amount)

		if err = t.s.Create(&tr); err != nil {
			if errors.Is(err, custom_errors.ErrInvalidTransactionType) {
				web.Failure(c, http.StatusBadRequest, custom_errors.ErrInvalidTransactionType)
				return
			}
			web.Failure(c, http.StatusInternalServerError, err)
			return
		}

		web.Success(c, http.StatusCreated, tr)
	}
}
