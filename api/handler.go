package api

import (
	"cloud-run-challenge-go/internal"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

const (
	urlViaCep     = "https://viacep.com.br/"
	urlWeatherApi = "https://api.weatherapi.com/v1/"
	weatherApiKey = "465c66df5be547d790a181453242405"
)

var FetchData = func(c *fiber.Ctx, url string) (response []byte, err error) {
	client := fasthttp.Client{}

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(url)

	if err := client.DoTimeout(req, res, 30*time.Second); err != nil {
		return nil, err
	}

	if res.StatusCode() != fiber.StatusOK {
		return nil, errors.New("invalid status code")
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	return res.Body(), nil
}

func HandleRequest(c *fiber.Ctx) error {
	cep := c.Params("cep")
	cep = strings.Replace(cep, "-", "", -1)
	if len(cep) != 8 {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": "invalid zipcode"})
	}

	url := urlViaCep + "ws/" + cep + "/json"
	response, err := FetchData(c, url)
	if err != nil || response == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "can not find zipcode"})
	}
	cepResponse := ViaCepResponse{}
	err = json.Unmarshal(response, &cepResponse)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing zipcode data"})
	}

	if cepResponse.Erro {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "can not found zipcode"})
	}

	city := internal.RemoveAccents(cepResponse.Localidade)
	state := cepResponse.Uf

	url = urlWeatherApi + "current.json?key=" + weatherApiKey + "&q=" + city + " - " + state + " - Brazil&aqi=no&tides=no"
	url = strings.Replace(url, " ", "%20", -1)

	response, err = FetchData(c, url)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching weather data"})
	}

	weatherResponse := WeatherApiResponse{}
	err = json.Unmarshal(response, &weatherResponse)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing weather data"})
	}

	tempC := strconv.FormatFloat(weatherResponse.Current.TempC, 'f', -1, 64)
	tempF := strconv.FormatFloat(weatherResponse.Current.TempF, 'f', -1, 64)
	tempK := strconv.FormatFloat(weatherResponse.Current.TempC+273.15, 'f', -1, 64)

	response = []byte(`{ "temp_C": ` + tempC + `, "temp_F": ` + tempF + `, "temp_K": ` + tempK + ` }`)

	return c.Send(response)
}
