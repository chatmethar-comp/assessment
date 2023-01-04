package expense

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func PutExpenseHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	e := Expense{}
	err = c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	stmt, err := db.Prepare(`UPDATE expenses SET title=$2 , amount=$3 , note=$4 , tags=$5 WHERE id=$1;`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can not Prepare Statement" + err.Error()})
	}
	if _, err := stmt.Exec(2, e.Title, e.Amount, e.Note, pq.Array(e.Tags)); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "Can not execute" + err.Error()})
	}
	e.ID = id
	return c.JSON(http.StatusCreated, e)
}
