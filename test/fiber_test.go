package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestRouteParameter(t *testing.T) {
	app.Get("/users/:userId/orders/:orderId", func(c *fiber.Ctx) error {
		userId := c.Params("userId")	
		orderId := c.Params("orderId")
		return c.SendString("Get user " + userId + " orders " + orderId)
	})

	request := httptest.NewRequest("GET", "/users/lev/orders/10", nil)
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err :=io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Get user lev orders 10", string(bytes))
}

func TestFormRequest(t *testing.T) {
	app.Post("/hello", func(c *fiber.Ctx) error {
		name := c.FormValue("name") // Get form value
		return c.SendString("Hello, " + name + "!")
	})

	body := strings.NewReader("name=Lev")
	request := httptest.NewRequest("POST", "/hello", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Set content type for form data
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err :=io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, Lev!", string(bytes))
}