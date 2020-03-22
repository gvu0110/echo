package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from the Echo Web Server!")
}

func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catAge := c.QueryParam("age")

	dataType := c.Param("data")
	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("Your cat's name is %s and cat's age is %s", catName, catAge))
	}

	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"age":  catAge,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "Data type must be 'string' or 'json'",
	})

}

func main() {
	fmt.Println("Welcome to the Echo Web Server!")

	e := echo.New()

	e.GET("/", hello)
	e.GET("/cats/:data", getCats)

	e.Start(":8000")
}
