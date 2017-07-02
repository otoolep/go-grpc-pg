package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/otoolep/go-grpc-pg/service"
)

// Command line defaults
const (
	DefaultgRPCAddr       = "localhost:11000"
	DefaultPostgreSQLAddr = "localhost:5432"
)

// Command line parameters
var gRPCAddr string
var pgAddr string

func init() {
	flag.StringVar(&gRPCAddr, "grpc-addr", DefaultgRPCAddr, "Set the gRPC bind address")
	flag.StringVar(&pgAddr, "pg-addr", DefaultPostgreSQLAddr, "Set PostgreSQL address")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	// Create the service.
	srv := service.New(gRPCAddr, nil)

	// Start the service.
	if err := srv.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start service: %s", err.Error())
		os.Exit(1)
	}
	log.Println("service started successfully")

	// Block until a signal is received.
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	log.Println("service exiting")
}
