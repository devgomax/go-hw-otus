package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("Usage: go-telnet [--timeout] {host} {port}")
	}

	host, port := args[0], args[1]
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("Connection error: %v", err)
	}
	defer client.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		client.Close()
	}()

	go func() {
		defer wg.Done()
		defer client.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := client.Send(); err != nil {
					return
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		defer client.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := client.Receive(); err != nil {
					return
				}
			}
		}
	}()

	wg.Wait()
}
