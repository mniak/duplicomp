package main

import (
	"flag"

	"github.com/mniak/duplicomp/internal/samples"
	"github.com/samber/lo"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 9000, "TCP port to listen")
	flag.Parse()

	lo.Must0(samples.RunServer(samples.WithPort(port)))
}
