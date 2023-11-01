package main

import (
	"flag"
	"log"
	"os"
	"syscall"

	"github.com/mniak/duplicomp"
	"github.com/mniak/duplicomp/log2"
)

func main() {
	var listenAddress string
	var primaryTarget string
	var shadowTarget string
	flag.StringVar(&listenAddress, "listen-address", ":9091", "TCP address to listen")
	flag.StringVar(&primaryTarget, "target", ":9001", "Connection target")
	flag.StringVar(&shadowTarget, "shadow-target", "", "Shadow connection target")
	flag.Parse()

	cmp := LogComparator{
		logger: log2.Sub(log2.FromWriter(os.Stdout), "[Comparator] "),
	}

	stopGw, err := duplicomp.StartNewGateway(listenAddress, primaryTarget, shadowTarget, cmp)
	if err != nil {
		log.Fatalln(err)
	}

	wait(syscall.SIGTERM, syscall.SIGINT)
	stopGw.GracefulStop()
}
