package main

import (
	"os"
	"os/signal"
	"sync"
	"context"
	"log"
	"itout/go-ethereum-lightclient/client"
	"itout/go-ethereum-lightclient/configs"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := client.Client{
		BeaconBaseURL: os.Getenv("BEACON_BASE_URL"),
		TrustedRoot: os.Getenv("TRUSTED_ROOT"),
		Spec: *configs.Mainnet,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	
	go func() {
		defer wg.Done()
		if err := client.StartClient(ctx); err != nil {
			log.Printf("Client stopped with error: %v", err)
		}
	}()
	
	wg.Wait()
}