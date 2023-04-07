package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	streamcatcher "microsomes.com/stgo/streamCatcher"
)

func main() {

	server := streamcatcher.NewStreamCatcherServer()

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go server.StartAndServe(wg, ctx)

	sigCn := make(chan os.Signal)
	signal.Notify(sigCn, syscall.SIGINT, syscall.SIGTERM)

	<-sigCn
	fmt.Println("Shutting down")
	cancel()
	wg.Wait()

}
