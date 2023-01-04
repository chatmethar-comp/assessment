package main

import (
	"log"

	// "net/http"
	"os"

	// "strconv"

	"github.com/chatmethar-comp/assessment/expense"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// "github.com/lib/pq"
)

func main() {
	expense.InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/expenses", expense.CreateExpenseHandler)
	e.GET("/expenses", expense.GetExpenseHandler)
	e.GET("/expenses/:id", expense.GetExpenseIdHandler)
	e.PUT("/expenses/:id", expense.PutExpenseHandler)

	log.Fatal(e.Start(os.Getenv("PORT")))
}
