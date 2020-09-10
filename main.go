package main

import (
	"flag"
	"fmt"
	"net/http"
	"log"
	"strconv"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func main() {

	bind := flag.String("bind", "127.0.0.1:8090", "bind address [host]:[port]")
	selfTest := flag.Bool("self-test", false, "test gpio pins")
	flag.Parse()

	pins := []int{9, 10, 22, 27, 17, 4, 3, 2}
	err := rpio.Open()
	if err != nil {
		log.Printf("ERROR: Unable to open gpio: %s", err)
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

	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "pong\n")
	})


	http.HandleFunc("/api/relay", func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		action := query.Get("action")
		relay, err := strconv.Atoi(query.Get("relay"))
		if err != nil || relay >= len(pins) {
			fmt.Fprintf(w, "invalid relay\n")
			return
		}

		pin := rpio.Pin(pins[relay])
		pin.Output()

		if action == "on" {
			pin.Low()
			fmt.Fprintf(w, "on\n")
			return
		} else if action == "off" {
			pin.High()
			fmt.Fprintf(w, "off\n")
			return
		}
		fmt.Fprintf(w, "invalid action\n")
	})

	log.Printf("INFO: Listening on: %s", *bind)
	err = http.ListenAndServe(*bind, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
