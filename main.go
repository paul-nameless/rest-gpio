package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/stianeikeland/go-rpio"
)


func main() {
	pins := [8]int{9, 10, 22, 27, 17, 4, 3, 2}
	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	defer rpio.Close()

	e := echo.New()

	e.GET("/api/relay", func(c echo.Context) error {
		action := c.QueryParam("action")
		relay, err := strconv.Atoi(c.QueryParam("relay"))
		if err != nil || relay >= len(pins) {
			return c.String(http.StatusBadRequest, "invalid relay")
		}

		pin := rpio.Pin(pins[relay])
		pin.Output()

		if action == "on" {
			pin.Low()
			return c.String(http.StatusOK, "on")
		} else if action == "off" {
			pin.High()
			return c.String(http.StatusOK, "off")
		}
		return c.String(http.StatusBadRequest, "invalid action")
	})

	e.Logger.Fatal(e.Start(":1323"))

}
