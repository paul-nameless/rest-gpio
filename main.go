package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stianeikeland/go-rpio"
)

func selfTest() {

}

func main() {

	bind := flag.String("bind", "127.0.0.1:8090", "bind address [host]:[port]")
	selfTest := flag.Bool("self-test", false, "test gpio pins")
	flag.Parse()

	pins := [...]int{9, 10, 22, 27, 17, 4, 3, 2}
	err := rpio.Open()
	if err != nil {
		fmt.Println("ERROR: unable to open gpio:", err.Error())
	}

	defer rpio.Close()

	if *selfTest {
		for _, pin := range pins {
			pin := rpio.Pin(pin)
			pin.Output()
			pin.Low()
			time.Sleep(500 * time.Millisecond)
		}
		for _, pin := range pins {
			pin := rpio.Pin(pin)
			pin.Output()
			pin.High()
			time.Sleep(500 * time.Millisecond)
		}
		return
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${remote_ip} | ${method} ${uri} ${status} - ${latency_human}\n",
	}))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

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

	e.Logger.Fatal(e.Start(*bind))

}
