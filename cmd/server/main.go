package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
	"github.com/otoolep/go-grpc-pg/service"
)

// Command line defaults
const (
	DefaultgRPCAddr           = "localhost:11000"
	DefaultPostgreSQLAddr     = "localhost:5432"
	DefaultPostgreSQLDB       = "pggo"
	DefaultPostgreSQLUser     = "postgres"
	DefaultPostgreSQLPassword = "postgres"
)

// Command line parameters
var gRPCAddr string
var pgAddr string
var pgDB string
var pgUser string
var pgPassword string

func init() {
	flag.StringVar(&gRPCAddr, "grpc-addr", DefaultgRPCAddr, "Set the gRPC bind address")
	flag.StringVar(&pgAddr, "pg-addr", DefaultPostgreSQLAddr, "Set PostgreSQL address")
	flag.StringVar(&pgDB, "pg-db", DefaultPostgreSQLDB, "Set PostgreSQL database")
	flag.StringVar(&pgUser, "pg-user", DefaultPostgreSQLUser, "Set PostgreSQL user")
	flag.StringVar(&pgPassword, "pg-password", DefaultPostgreSQLPassword, "Set PostgreSQL password")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags)
	log.SetPrefix("[main] ")

	// Create database connection.
	conn, err := pgConnection()
	if err != nil {
		log.Fatal("invalid PostgreSQL connection parameters:", err.Error())
	}
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal("failed to open PostgreSQL:", err.Error())
	}

	// Verify connection.
	if err := db.Ping(); err != nil {
		log.Fatal("failed to verify connection to PostgreSQL:", err.Error())
	}
	log.Println("database connection verified")

	// Create the service.
	srv := service.New(gRPCAddr, db)

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

func pgConnection() (string, error) {
	host, port, err := net.SplitHostPort(pgAddr)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`user=%s password=%s dbname=%s host=%s port=%s`,
		pgUser, pgPassword, pgDB, host, port), nil
}
