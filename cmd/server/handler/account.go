package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	acc "github.com/lucaspichi06/xepelin-bank/internal/account"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	"github.com/lucaspichi06/xepelin-bank/internal/events"
	custom_errors "github.com/lucaspichi06/xepelin-bank/pkg/errors"
	"github.com/lucaspichi06/xepelin-bank/pkg/web"
	"net/http"
)

type Accounts interface {
	Create() gin.HandlerFunc
	GetBalance() gin.HandlerFunc
}

type account struct {
	s acc.Service
}

func NewAccountHandler(s acc.Service) Accounts {
	return &account{
		s: s,
	}
}

// Create	godoc
// @Summary	Creates a new account
// @Tags	Account
// @Description	creates a new account with the received parameters
// @Accept	json
// @Produce	json
// @Param	token	header	string	true	"token"
// @Param	transaction		body 	domain.AccountRequest	true	"Account to create"
// @Success 201	{object}	web.Response
// @Failure	400	{object}	web.ErrorResponse
// @Failure	500	{object}	web.ErrorResponse
// @Router	/accounts	[post]
func (a account) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var acc domain.AccountRequest
		err := c.ShouldBindJSON(&acc)
		if err != nil {
			web.Failure(c, http.StatusBadRequest, custom_errors.ErrInvalidJSON)
			return
		}

		event := events.NewCreateAccountEvent(acc.Name, a.s)

		newAcc, err := event.Process()
		if err != nil {
			web.Failure(c, http.StatusInternalServerError, err)
			return
		}

		web.Success(c, http.StatusCreated, newAcc)
	}
}

// Read	godoc
// @Summary	Get the balance from an account
// @Tags	Account
// @Description	get the balance from an account
// @Accept	json
// @Produce	json
// @Param	id		path	int		true	"Account ID"
// @Success	200	{object}	web.Response
// @Failure	400	{object}	web.ErrorResponse
// @Failure	404	{object}	web.ErrorResponse
// @Failure	500	{object}	web.ErrorResponse
// @Router	/accounts/{id}/balance	[get]
func (a account) GetBalance() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			web.Failure(c, http.StatusBadRequest, custom_errors.ErrInvalidID)
			return
		}

		event := events.NewBalanceEvent(id, a.s)

		res, err := event.Process()
		if err != nil {
			if errors.Is(err, custom_errors.ErrNotFound) {
				web.Failure(c, http.StatusNotFound, fmt.Errorf("account not found: %v", err))
				return
			}
			web.Failure(c, http.StatusInternalServerError, err)
			return
		}
		web.Success(c, http.StatusOK, res)
	}
}
