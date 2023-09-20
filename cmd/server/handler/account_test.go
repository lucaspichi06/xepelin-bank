package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	custom_errors "github.com/lucaspichi06/xepelin-bank/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type accountServiceMock struct {
	create func(account domain.Account) error
	read   func(id uuid.UUID) (domain.Account, error)
	update func(account domain.Account) error
}

func (a accountServiceMock) Create(account domain.Account) error {
	return a.create(account)
}
func (a accountServiceMock) Read(id uuid.UUID) (domain.Account, error) {
	return a.read(id)
}
func (a accountServiceMock) Update(account domain.Account) error {
	return a.update(account)
}

func TestAccountCreate(t *testing.T) {
	t.Run("account create success", func(t *testing.T) {
		serviceMock := accountServiceMock{
			create: func(account domain.Account) error {
				return nil
			},
		}
		a := NewAccountHandler(serviceMock)

		r := gin.Default()
		r.POST("/test", a.Create())

		body := []byte(`{"name":"test"}`)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/test", bytes.NewBuffer(body))
		if err != nil {
			t.Fail()
		}
		r.ServeHTTP(w, req)

		responseMap := make(map[string]interface{})
		responseBody, err := io.ReadAll(w.Body)
		if err != nil {
			t.Fail()
		}
		if err = json.Unmarshal(responseBody, &responseMap); err != nil {
			t.Fail()
		}

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "test", responseMap["data"].(map[string]interface{})["name"])
		assert.Equal(t, float64(0), responseMap["data"].(map[string]interface{})["balance"])
	})
	t.Run("account create invalid JSON", func(t *testing.T) {
		a := NewAccountHandler(nil)

		r := gin.Default()
		r.POST("/test", a.Create())

		body := []byte(`{"name": 10}`)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/test", bytes.NewBuffer(body))
		if err != nil {
			t.Fail()
		}
		r.ServeHTTP(w, req)

		responseMap := make(map[string]interface{})
		responseBody, err := io.ReadAll(w.Body)
		if err != nil {
			t.Fail()
		}
		if err = json.Unmarshal(responseBody, &responseMap); err != nil {
			t.Fail()
		}

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, responseMap["message"], "invalid json")
	})
	t.Run("account create internal server error", func(t *testing.T) {
		serviceErrorMock := accountServiceMock{
			create: func(account domain.Account) error {
				return errors.New("test error")
			},
		}
		aError := NewAccountHandler(serviceErrorMock)

		rError := gin.Default()
		rError.POST("/test", aError.Create())
		body := []byte(`{"name": "test"}`)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/test", bytes.NewBuffer(body))
		if err != nil {
			t.Fail()
		}
		rError.ServeHTTP(w, req)

		responseMap := make(map[string]interface{})
		responseBody, err := io.ReadAll(w.Body)
		if err != nil {
			t.Fail()
		}
		if err = json.Unmarshal(responseBody, &responseMap); err != nil {
			t.Fail()
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, responseMap["message"], "test error")
	})
}

func TestAccountGetBalance(t *testing.T) {
	t.Run("account balance success", func(t *testing.T) {
		serviceMock := accountServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{
					Balance: 1000.00,
				}, nil
			},
		}
		a := NewAccountHandler(serviceMock)

		r := gin.Default()
		r.GET("/test/:id/balance", a.GetBalance())

		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/test/7dab3e13-02c7-455e-845a-13cb8c70ae8c/balance", nil)
		if err != nil {
			t.Fail()
		}
		r.ServeHTTP(w, req)

		responseMap := make(map[string]interface{})
		responseBody, err := io.ReadAll(w.Body)
		if err != nil {
			t.Fail()
		}
		if err = json.Unmarshal(responseBody, &responseMap); err != nil {
			t.Fail()
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, 1000.00, responseMap["data"].(map[string]interface{})["balance"])
	})
	t.Run("account balance invalid ID", func(t *testing.T) {
		a := NewAccountHandler(nil)

		r := gin.Default()
		r.GET("/test/:id/balance", a.GetBalance())

		body := []byte(`{"name": 10}`)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/test/10/balance", bytes.NewBuffer(body))
		if err != nil {
			t.Fail()
		}
		r.ServeHTTP(w, req)

		responseMap := make(map[string]interface{})
		responseBody, err := io.ReadAll(w.Body)
		if err != nil {
			t.Fail()
		}
		if err = json.Unmarshal(responseBody, &responseMap); err != nil {
			t.Fail()
		}

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, responseMap["message"], "invalid param id")
	})
	t.Run("account balance not found error", func(t *testing.T) {
		serviceErrorMock := accountServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{}, custom_errors.ErrNotFound
			},
		}
		aError := NewAccountHandler(serviceErrorMock)

		rError := gin.Default()
		rError.GET("/test/:id/balance", aError.GetBalance())
		body := []byte(`{"name": "test"}`)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/test/7dab3e13-02c7-455e-845a-13cb8c70ae8c/balance", bytes.NewBuffer(body))
		if err != nil {
			t.Fail()
		}
		rError.ServeHTTP(w, req)

		responseMap := make(map[string]interface{})
		responseBody, err := io.ReadAll(w.Body)
		if err != nil {
			t.Fail()
		}
		if err = json.Unmarshal(responseBody, &responseMap); err != nil {
			t.Fail()
		}

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, responseMap["message"], custom_errors.ErrNotFound.Error())
	})
	t.Run("account balance internal server error", func(t *testing.T) {
		serviceErrorMock := accountServiceMock{
			read: func(id uuid.UUID) (domain.Account, error) {
				return domain.Account{}, errors.New("test error")
			},
		}
		aError := NewAccountHandler(serviceErrorMock)

		rError := gin.Default()
		rError.GET("/test/:id/balance", aError.GetBalance())
		body := []byte(`{"name": "test"}`)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/test/7dab3e13-02c7-455e-845a-13cb8c70ae8c/balance", bytes.NewBuffer(body))
		if err != nil {
			t.Fail()
		}
		rError.ServeHTTP(w, req)

		responseMap := make(map[string]interface{})
		responseBody, err := io.ReadAll(w.Body)
		if err != nil {
			t.Fail()
		}
		if err = json.Unmarshal(responseBody, &responseMap); err != nil {
			t.Fail()
		}

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, responseMap["message"], "test error")
	})
}
