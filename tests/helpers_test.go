package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-secret-key")
}

func setupRouter() *gin.Engine {
	return gin.New()
}

func makeRequest(method, url string, body interface{}) *http.Request {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(
		context.Background(),
		method,
		url,
		bytes.NewBuffer(b),
	)
	req.Header.Set("Content-Type", "application/json")
	return req
}
