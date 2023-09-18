package main

import (
	"flag"
	"log"
	"syscall"

	"github.com/mniak/duplicomp"
)

func main() {
	var listenAddress string
	var primaryTarget string
	var shadowTarget string
	flag.StringVar(&listenAddress, "listen-address", ":9091", "TCP address to listen")
	flag.StringVar(&primaryTarget, "target", ":9001", "Connection target")
	flag.StringVar(&shadowTarget, "shadow-target", "", "Shadow connection target")
	flag.Parse()

	cmp := LogComparator{}

	stopGw, err := duplicomp.StartNewGateway(listenAddress, primaryTarget, shadowTarget, cmp)
	if err != nil {
		log.Fatalln(err)
	}

	wait(syscall.SIGTERM, syscall.SIGINT)
	stopGw.GracefulStop()
}
