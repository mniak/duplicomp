package main

import (
	"os"
	"os/signal"
)

func wait(sigs ...os.Signal) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, sigs...)
	<-sigchan
}
