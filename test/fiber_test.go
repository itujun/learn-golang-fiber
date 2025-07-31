package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var	app = fiber.New()

func TestRoutingHelloWorld(t *testing.T) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	request := httptest.NewRequest("GET", "/", nil)
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err :=io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, World!", string(bytes))
}

func TestCtx(t *testing.T) {
	app.Get("/hello", func(c *fiber.Ctx) error {
		name := c.Query("name", "guest") // Default to "World" if no name is provided
		return c.SendString("Hello, " + name + "!")
	})

	request := httptest.NewRequest("GET", "/hello?name=Lev", nil)
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err :=io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, Lev!", string(bytes))

	request = httptest.NewRequest("GET", "/hello", nil)
	response, err = app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err =io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, guest!", string(bytes))
}

func TestHttpRequest(t *testing.T) {
	app.Get("/request", func(c *fiber.Ctx) error {
		first := c.Get("firstname") 	// headers
		last := c.Cookies("lastname")   // cookies
		return c.SendString("Hello, " + first + " " + last + "!")
	})

	request := httptest.NewRequest("GET", "/request", nil)
	request.Header.Set("firstname", "Lev") 								// Set header
	request.AddCookie(&http.Cookie{Name: "lastname", Value: "Tempest"}) // Set cookie
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err :=io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, Lev Tempest!", string(bytes))
}