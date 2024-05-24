package test

import (
	"cloud-run-challenge-go/api"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"net/http/httptest"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	app := fiber.New()
	app.Get("/:cep", api.HandleRequest)

	tests := []struct {
		cep          string
		expectedCode int
	}{
		{"01001000", fiber.StatusOK},
		{"invalid", fiber.StatusUnprocessableEntity},
		{"00000000", fiber.StatusNotFound},
	}
	for _, tt := range tests {
		req := httptest.NewRequest("GET", "/"+tt.cep, nil)
		resp, _ := app.Test(req)
		assert.Equal(t, tt.expectedCode, resp.StatusCode)
	}
}

func TestFetchData(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	tests := []struct {
		url         string
		expectError bool
	}{
		{"https://viacep.com.br/ws/01001000/json", false}, // valid URL
		{"https://invalidurl.com", true},                  // invalid URL
	}
	for _, tt := range tests {
		_, err := api.FetchData(ctx, tt.url)
		if tt.expectError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestHandleRequest_InvalidZipCode(t *testing.T) {
	app := fiber.New()
	app.Get("/:cep", api.HandleRequest)

	req := httptest.NewRequest("GET", "/123", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

func TestHandleRequest_ValidZipCode(t *testing.T) {
	app := fiber.New()
	app.Get("/:cep", api.HandleRequest)

	req := httptest.NewRequest("GET", "/01001000", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestFetchData_InvalidUrl(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	_, err := api.FetchData(ctx, "https://invalidurl.com")
	assert.NotNil(t, err)
}

func TestFetchData_ValidUrl(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	_, err := api.FetchData(ctx, "https://viacep.com.br/ws/01001000/json")
	assert.Nil(t, err)
}
