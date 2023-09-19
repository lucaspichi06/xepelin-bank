package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	acc "github.com/lucaspichi06/xepelin-bank/internal/account"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	custom_errors "github.com/lucaspichi06/xepelin-bank/pkg/errors"
	"github.com/lucaspichi06/xepelin-bank/pkg/web"
	"net/http"
	"strconv"
)

type Accounts interface {
	Create() gin.HandlerFunc
	Read() gin.HandlerFunc
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
// @Param	transaction		body 	domain.Account	true	"Account to create"
// @Success 201	{object}	web.Response
// @Failure	400	{object}	web.ErrorResponse
// @Failure	500	{object}	web.ErrorResponse
// @Router	/accounts	[post]
func (a account) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var account domain.Account
		err := c.ShouldBindJSON(&account)
		if err != nil {
			web.Failure(c, http.StatusBadRequest, custom_errors.ErrInvalidJSON)
			return
		}

		account.ID, err = a.s.Create(account)
		if err != nil {
			web.Failure(c, http.StatusInternalServerError, err)
			return
		}

		web.Success(c, http.StatusCreated, account)
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
func (a account) Read() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, http.StatusBadRequest, custom_errors.ErrInvalidID)
			return
		}
		res, err := a.s.Read(id)
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
