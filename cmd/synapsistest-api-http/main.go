package main

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/synapsis-test/cmd/synapsistest-api-http/server"
)

func main() {
	godotenv.Load()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		os.Exit(server.Run())
		defer wg.Done()
	}()
	wg.Wait()
}
