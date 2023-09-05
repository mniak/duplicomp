package main

import (
	"flag"

	"github.com/mniak/duplicomp/internal/samples"
	"github.com/samber/lo"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 9000, "TCP port to connect")
	flag.Parse()

	lo.Must0(samples.RunSendPing(samples.WithPort(port)))
}
