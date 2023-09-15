package main

import (
	"flag"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/duplicomp/internal/samples"
	"github.com/samber/lo"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 9000, "TCP port to connect")
	flag.Parse()

	msg := gofakeit.Sentence(8)
	lo.Must0(samples.RunSendPing(msg, samples.WithPort(port)))
}
