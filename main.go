package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	username = "adam"
	password = "12345"

	cookieName  = "sessionID"
	cookieValue = "8500RfpFDt&S"
)

type Cat struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

type Dog struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

type Hamster struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

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

func addCat(c echo.Context) error {
	cat := Cat{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Read the request body FAIL: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &cat)
	if err != nil {
		log.Printf("Unmarshal addCat FAIL: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("There is a new cat %#v", cat)
	return c.String(http.StatusOK, fmt.Sprintf("We got your cat %s", cat.Name))
}

func addDog(c echo.Context) error {
	dog := Dog{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("Process addDog request FAIL: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("There is a new dog %#v", dog)
	return c.String(http.StatusOK, fmt.Sprintf("We got your dog %s", dog.Name))
}

func addHamster(c echo.Context) error {
	hamster := Hamster{}

	err := c.Bind(&hamster)
	if err != nil {
		log.Printf("Process addHamster request FAIL: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("There is a new hamster %#v", hamster)
	return c.String(http.StatusOK, fmt.Sprintf("We got your hamster %s", hamster.Name))
}

func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "You are in the main Admin page!")
}

func mainCookie(c echo.Context) error {
	return c.String(http.StatusOK, "You are in the main Cookie page!")
}

func authValidator(user, pass string, c echo.Context) (bool, error) {
	if user == username && pass == password {
		return true, nil
	}
	return false, nil
}

func logIn(c echo.Context) error {
	user := c.QueryParam("username")
	pass := c.QueryParam("password")

	if user == username && pass == password {
		cookie := &http.Cookie{}
		// cookie := new(http.Cookie)

		cookie.Name = cookieName
		cookie.Value = cookieValue
		cookie.Expires = time.Now().Add(8 * time.Hour)

		c.SetCookie(cookie)

		return c.String(http.StatusOK, "You are logged in!")
	}
	return c.String(http.StatusUnauthorized, "Your username or password is invalid!")
}

////////// Middlewares section //////////
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "AdamVu/1.0")
		c.Response().Header().Set("Custom-Header", "No-Value")
		return next(c)
	}
}

func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(cookieName)

		if err != nil {
			log.Println(err)
			return err
		}

		if cookie.Value == cookieValue {
			return next(c)
		}

		return c.String(http.StatusUnauthorized, "You don't have the right cookie")
	}
}

func main() {
	fmt.Println("Welcome to the Echo Web Server!")

	e := echo.New()

	e.Use(ServerHeader)

	adminGroup := e.Group("/admin")
	cookieGroup := e.Group("/cookie")

	cookieGroup.Use(checkCookie)

	// Use middleware to log server interaction
	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `[${time_custom}]  ${status}  ${method}  ${host}${path}  ${latency_human}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	// Use middleware for basic authentication
	adminGroup.Use(middleware.BasicAuth(authValidator))

	adminGroup.GET("/main", mainAdmin)

	cookieGroup.GET("/main", mainCookie)

	e.GET("/login", logIn)
	e.GET("/", hello)
	e.GET("/cats/:data", getCats)
	e.POST("/addcat", addCat)
	e.POST("/adddog", addDog)
	e.POST("/addhamster", addHamster)
	e.Start(":8000")
}
