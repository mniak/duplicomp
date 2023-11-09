package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/mniak/duplicomp/internal/samples"
	"github.com/samber/lo"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 9000, "TCP port to listen")
	flag.StringVar(&samples.ApplicationVersion, "version", "v1.0", "Set application version")
	flag.Parse()

	p := samples.WithPort(port)
	stop := lo.Must(samples.RunServer(p))

	wait(syscall.SIGINT, syscall.SIGTERM)
	stop.Stop()
}

func wait(signals ...os.Signal) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, signals...)
	<-sigs
}
