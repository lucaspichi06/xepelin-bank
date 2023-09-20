package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lucaspichi06/xepelin-bank/internal/domain"
	custom_errors "github.com/lucaspichi06/xepelin-bank/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type transactionServiceMock struct {
	create func(tr *domain.Transaction) error
}

func (t transactionServiceMock) Create(tr *domain.Transaction) error {
	return t.create(tr)
}

func TestTransactionCreate(t *testing.T) {
	t.Run("transaction create success", func(t *testing.T) {
		serviceMock := transactionServiceMock{
			create: func(tr *domain.Transaction) error {
				return nil
			},
		}
		tr := NewTransactionsHandler(serviceMock)

		r := gin.Default()
		r.POST("/test", tr.Process())

		body := []byte(`{"account_id":"d70d0a95-af7f-4098-8d81-caca1934e94d","type":"withdraw","amount":10000.33}`)

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
		assert.Equal(t, "d70d0a95-af7f-4098-8d81-caca1934e94d", responseMap["data"].(map[string]interface{})["account_id"])
		assert.Equal(t, "withdraw", responseMap["data"].(map[string]interface{})["type"])
		assert.Equal(t, 10000.33, responseMap["data"].(map[string]interface{})["amount"])
		assert.NotNil(t, responseMap["data"].(map[string]interface{})["transaction_id"])
	})
	t.Run("transaction create invalid JSON", func(t *testing.T) {
		tr := NewTransactionsHandler(nil)

		r := gin.Default()
		r.POST("/test", tr.Process())

		body := []byte(`{"account_id":"d70d0a91231235-af7f-4098-8d81-caca1934e94d","type":"withdraw","amount":10000.33}`)

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
	t.Run("transaction create invalid transaction type", func(t *testing.T) {
		serviceMock := transactionServiceMock{
			create: func(tr *domain.Transaction) error {
				return custom_errors.ErrInvalidTransactionType
			},
		}
		tr := NewTransactionsHandler(serviceMock)

		r := gin.Default()
		r.POST("/test", tr.Process())

		body := []byte(`{"account_id":"d70d0a95-af7f-4098-8d81-caca1934e94d","type":"invalid","amount":10000.33}`)

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
		assert.Contains(t, responseMap["message"], custom_errors.ErrInvalidTransactionType.Error())
	})
	t.Run("transaction create internal server error", func(t *testing.T) {
		serviceMock := transactionServiceMock{
			create: func(tr *domain.Transaction) error {
				return errors.New("test error")
			},
		}
		tr := NewTransactionsHandler(serviceMock)

		r := gin.Default()
		r.POST("/test", tr.Process())

		body := []byte(`{"account_id":"d70d0a95-af7f-4098-8d81-caca1934e94d","type":"withdraw","amount":10000.33}`)

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

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, responseMap["message"], "test error")
	})
}
