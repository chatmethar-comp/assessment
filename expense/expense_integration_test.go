//go:build integration

package expense

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

func seedExpense(t *testing.T) Expense {
	var c Expense
	body := bytes.NewBufferString(`{
		"title": "New Stuff",
		"amount": 70,
		"note": "Bunny",
		"tags": [
			"rabbit",
			"young"
		]
	}`)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&c)
	if err != nil {
		t.Fatal("can't create uomer:", err)
	}
	return c
}

func TestGetAllExpense(t *testing.T) {
	seedExpense(t)
	var e []Expense
	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(e), 0)
}

func TestGetExpenseById(t *testing.T) {
	c := seedExpense(t)
	var latest Expense
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(c.ID)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, c.ID, latest.ID)
	assert.NotEmpty(t, latest.Title)
	assert.NotEmpty(t, latest.Amount)
	assert.NotEmpty(t, latest.Note)
	assert.NotEmpty(t, latest.Tags)
}

func TestCreateExpense(t *testing.T) {

	body := bytes.NewBufferString(`{
		"title": "New Stuff",
		"amount": 70,
		"note": "Bunny",
		"tags": [
			"rabbit",
			"young"
		]
	}`)
	var e Expense
	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, e.ID, 0)
	assert.Equal(t, e.Title, "New Stuff")
	assert.Equal(t, e.Amount, 70)
	assert.Equal(t, e.Note, "Bunny")
	assert.Equal(t, e.Tags, []string{"rabbit", "young"})
}

func TestUpdateExpense(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "Stuff",
		"amount": 70,
		"note": "Bunny",
		"tags": [
			"rabbit",
			"young"
		]
	}`)
	var e Expense
	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)
	body = bytes.NewBufferString(`{
		"title": "New Stuff",
		"amount": 100,
		"note": "Loony",
		"tags": [
			"rabbit",
			"young"
		]
	}`)
	res = request(http.MethodPut, uri("expenses", strconv.Itoa(e.ID)), body)
	err = res.Decode(&e)
	assert.Nil(t, err)
	assert.NotEqual(t, e.ID, 0)
	assert.Equal(t, e.Title, "New Stuff")
	assert.Equal(t, e.Amount, 100)
	assert.Equal(t, e.Note, "Loony")
	assert.Equal(t, e.Tags, []string{"rabbit", "young"})
}
