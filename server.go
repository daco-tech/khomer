package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/common-nighthawk/go-figure"
)

func main() {
	myFigure := figure.NewFigure("kHomer", "", true)
	myFigure.Print()
	
	e := echo.New()
	e.HideBanner = true
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}	
