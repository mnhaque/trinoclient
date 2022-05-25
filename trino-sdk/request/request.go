package request

import (
	"database/sql"

	"github.com/mnhaque/trinoclient/trino-sdk/config"
	"github.com/mnhaque/trinoclient/trino-sdk/trinodriver"
)

type Response struct {
	QueryId string
	Data    *sql.Rows
	Error   error
}

func New(c *config.Configuration, query string) *Response {
	data, queryId, err := getDataByQuery(query, c)
	r := &Response{
		Data:    data,
		QueryId: queryId,
		Error:   err,
	}

	return r
}

func getDataByQuery(query string, c *config.Configuration) (*sql.Rows, string, error) {
	dsn, _ := c.GetDSN()
	db, err := sql.Open("trino", dsn)
	if err != nil {
		return nil, trinodriver.QueryId, err
	}

	rows, queryException := db.Query(query)
	if queryException != nil {
		return nil, trinodriver.QueryId, err
	}

	return rows, trinodriver.QueryId, nil
}
