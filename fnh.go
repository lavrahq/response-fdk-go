package fnh

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/fnproject/fdk-go"
)

// Fnh is the Fn Helper utility.
type Fnh struct {
	Context fdk.Context
	Client  *http.Client
}

// QueryRequest is a request sent to
type QueryRequest struct {
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
}

// QueryErrorResponse is the response for a Query request when it fails.
type QueryErrorResponse struct {
	Path  string `json:"path"`
	Error string `json:"error"`
}

// QueryResponse is the response for a Query request.
type QueryResponse struct {
	ResultType string     `json:"result_type"`
	Result     [][]string `json:"result"`
}

// Create creates a new Fnh instance from the Fn context.
func Create(ctx context.Context) *Fnh {
	return &Fnh{
		Context: fdk.GetContext(ctx),
		Client:  http.DefaultClient,
	}
}

// Query runs a query request to the GraphQL service via the
// /query endpoint.
func (f *Fnh) Query(qreq *QueryRequest) (*QueryResponse, error) {
	data, err := json.Marshal(qreq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "http://graphql:8080/v1/query", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Hasura-Admin-Secret", f.Context.Config()["admin_secret"])

	res, err := f.Client.Do(req)
	if err != nil {
		return nil, err
	}

	qrs := new(QueryResponse)
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	json.Unmarshal(buf.Bytes(), qrs)

	return qrs, nil
}
