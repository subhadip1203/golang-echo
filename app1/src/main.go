package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// -------- route functions ---------/
func homeRoute(c echo.Context) error {
	return c.String(http.StatusOK, "Home route")
}

func testRoute(c echo.Context) error {
	user := c.QueryParam("user")
	return c.String(http.StatusOK, fmt.Sprintf("The user is %s", user))
}

func testRoute2(c echo.Context) error {
	user := c.Param("data")
	return c.String(http.StatusOK, fmt.Sprintf("The user is %s", user))
}

func testRoute3(c echo.Context) error {
	user := User{}
	defer c.Request().Body.Close()
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Invalid JSON")
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Invalid Data")
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Error")
	}
	return c.String(http.StatusOK, string(userJson))
}

func homeAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "Admin home page")
}

func testAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "Admin test page")
}

func main() {
	// ------ ENV file ----------//
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	PORT := os.Getenv("PORT")

	// ----------- create routes --------/
	e := echo.New()
	e.GET("/", homeRoute)
	e.GET("/test", testRoute)
	e.GET("/test2/:data", testRoute2)
	e.POST("/test3", testRoute3)

	// ------ routing group ------//
	groupAdmin := e.Group("/admin", middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} latency=${latency}\n",
	}))
	groupAdmin.GET("", homeAdmin)
	groupAdmin.GET("/test", testAdmin)

	// ----------- running server ---------/
	e.Logger.Fatal(e.Start(":" + PORT))
}
