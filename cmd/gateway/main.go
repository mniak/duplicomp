package main

import (
	"flag"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/mniak/duplicomp"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).
		Level(zerolog.TraceLevel).With().
		Timestamp().
		Caller().
		Logger()

	var listenAddress string
	var primaryTarget string
	var shadowTarget string
	flag.StringVar(&listenAddress, "listen-address", ":9091", "TCP address to listen")
	flag.StringVar(&primaryTarget, "target", ":9001", "Connection target")
	flag.StringVar(&shadowTarget, "shadow-target", "", "Shadow connection target")
	flag.Parse()

	cmp := LogComparator{
		logger: logger,
	}

	stopGw, err := duplicomp.StartNewGateway(listenAddress, primaryTarget, shadowTarget, cmp)
	if err != nil {
		log.Fatalln(err)
	}

	wait(syscall.SIGTERM, syscall.SIGINT)
	stopGw.GracefulStop()
}
