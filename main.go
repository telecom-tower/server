package main

//go:generate protoc -I $GOPATH/src/github.com/telecom-tower/towerapi/v1 telecomtower.proto --go_out=plugins=grpc:$GOPATH/src/github.com/telecom-tower/towerapi/v1

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime/trace"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	ws2811 "github.com/supcik/rpi_ws281x_go"
	"github.com/telecom-tower/grpc-renderer"
)

var version = "master"

func main() { // nolint: gocyclo
	debug := flag.Bool("debug", false, "Debug mode")
	verbose := flag.Bool("verbose", false, "Verbose mode")
	showVer := flag.Bool("version", false, "Show version")
	traceFile := flag.String("trace", "", "Generate tracing file")
	brightness := flag.Int("brightness", 128, "Maximum LED brightness (between 0 and 255)")
	grpcPort := flag.Int("grpc-port", 10000, "listening gRPC port")

	flag.Parse()

	if *showVer {
		fmt.Printf("Telecom Tower Server // version : %v\n", version)
		os.Exit(0)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else if *verbose {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	if *traceFile != "" {
		f, err := os.Create(*traceFile)
		if err != nil {
			log.Panic(errors.WithMessage(err, "Unable to create trace file"))
		}
		err = trace.Start(f)
		defer trace.Stop()
		if err != nil {
			log.Panic(errors.WithMessage(err, "Unable to trace"))
		}
	}

	// Create and run hub
	wsopt := ws2811.DefaultOptions
	wsopt.DmaNum = 10
	wsopt.Channels[0].Brightness = *brightness
	wsopt.Channels[0].LedCount = 1024
	ws, err := ws2811.MakeWS2811(&wsopt)
	if err != nil {
		log.Fatal(err)
	}
	if err = ws.Init(); err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		log.Fatal(renderer.Serve(grpcLis, ws))
		wg.Done()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			trace.Stop()
			log.Info("Finished")
			os.Exit(0)
		}
	}()

	wg.Wait()
}
