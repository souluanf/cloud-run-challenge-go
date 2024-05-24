package test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"cloud-run-challenge-go/api"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func startServer() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Get("/:cep", api.HandleRequest)
	go func() {
		err := app.Listen(":8080")
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(2 * time.Second)
}

func TestIntegration(t *testing.T) {
	startServer()

	tests := []struct {
		url          string
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			url:          "http://localhost:8080/08210010",
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"temp_C": 19.0,
				"temp_F": 66.2,
				"temp_K": 292.15,
			},
		},
		{
			url:          "http://localhost:8080/00000000",
			expectedCode: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "can not found zipcode",
			},
		},
		{
			url:          "http://localhost:8080/123",
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: map[string]interface{}{
				"error": "invalid zipcode",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			resp, err := http.Get(tt.url)
			assert.Nil(t, err)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					panic(err)
				}
			}(resp.Body)

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			assert.Nil(t, err)

			assert.Equal(t, tt.expectedBody, result)
		})
	}
}
