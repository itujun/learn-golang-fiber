package learn_golang_fiber

import (
	"bytes"
	_ "embed"
	"io"
	"mime/multipart"
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

//go:embed source/contoh.txt
var contohFile []byte

func TestMultipartForm(t *testing.T) {
	app.Post("/upload", func(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("file") // Get file from form
		if err != nil {
			// return c.Status(fiber.StatusBadRequest).SendString("File not found")
			return err
		}

		err = ctx.SaveFile(file, "../target/" + file.Filename) // Save file to disk
		if err != nil {
			return err
		}

		return ctx.SendString("Upload Success")
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, err := writer.CreateFormFile("file", "contoh.txt") // Create form file
	assert.Nil(t, err)
	file.Write(contohFile) // Write file content
	writer.Close() // Close writer to finalize the form data
	
	request := httptest.NewRequest("POST", "/upload", body)
	request.Header.Set("Content-Type", writer.FormDataContentType()) // Set content type for form data
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err :=io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Upload Success", string(bytes))
}