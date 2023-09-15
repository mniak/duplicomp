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
	flag.Parse()

	stop := lo.Must(samples.RunServer(samples.WithPort(port)))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	stop.Stop()
}
