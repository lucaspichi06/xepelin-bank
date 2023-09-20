package main

import (
	"database/sql"
	"fmt"
	"github.com/lucaspichi06/xepelin-bank/cmd/server/handler"
	"github.com/lucaspichi06/xepelin-bank/docs"
	"github.com/lucaspichi06/xepelin-bank/internal/account"
	"github.com/lucaspichi06/xepelin-bank/internal/transaction"
	"github.com/lucaspichi06/xepelin-bank/pkg/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

// @title Xepelin Bank
// @version 1.0
// @description This API Handles Accounts and Transactions.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
func main() {
	// opening the DB
	//db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/my_db", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT")))
	os.Setenv("TOKEN", "my-secret-token")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/my_db", "root", "rootpass", "localhost", "3306"))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	//storage := store.NewSqlStore(db)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	// account section
	accountRepository := account.NewRepository(db)
	accountService := account.NewService(accountRepository)
	accountHandler := handler.NewAccountHandler(accountService)

	acc := r.Group("/accounts")
	{
		acc.GET(":id/balance", accountHandler.GetBalance())
		acc.POST("", middleware.Authentication(), accountHandler.Create())
	}

	// transaction section
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository, accountService)
	transactionHandler := handler.NewTransactionsHandler(transactionService)

	tran := r.Group("/transactions")
	{
		tran.POST("", middleware.Authentication(), middleware.Logger(), transactionHandler.Process())
	}

	// documentation section
	docs.SwaggerInfo.Host = os.Getenv("HOST")
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
