package main

// Copyright 2016 Jacques Supcik / Bluemasters
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Towerserver implements a web server, receiving frames on a websocket, and
sending them on the LED display of the tower.
*/

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gitlab.com/geomyidae/ws2811"
	"gitlab.com/geomyidae/ws2811gw"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultBrightness = 64
)

func main() {
	// these defaults are OK for the official telecom tower.
	var debug = flag.Bool("debug", false, "set debug mode")
	var rows = flag.Int("rows", 8, "LED matrix rows")
	var columns = flag.Int("columns", 128, "LED matrix columns")
	var dmaNum = flag.Int("dma-num", 5, "DMA Number")
	var gpioPin = flag.Int("gpio-pin", 18, "GPIO Pin")
	var port = flag.Int("port", 8484, "HTTP daemon port")
	flag.Parse()

	ledsCount := *rows * *columns

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.Debugln("Starting...")
	r := mux.NewRouter()
	ws2811gw.Init(r, func(r *http.Request) ws2811gw.Ws2811Engine {
		opt := ws2811.DefaultOptions
		opt.LedCount = ledsCount
		opt.DmaNum = *dmaNum
		opt.GpioPin = *gpioPin
		opt.Brightness = defaultBrightness
		brightness := r.FormValue("brightness")
		if brightness != "" {
			val, err := strconv.Atoi(brightness)
			if err == nil {
				opt.Brightness = val
			}
		}
		res, _ := ws2811.MakeWS2811(&opt)
		return res
	})

	// Start HTTP server
	log.Infof("Starting server on port %v", *port)
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%v", *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
