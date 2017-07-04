package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/otoolep/go-grpc-pg/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Command line defaults
const (
	DefaultgRPCAddr = "localhost:11000"
)

// Command line parameters
var gRPCAddr string

func init() {
	flag.StringVar(&gRPCAddr, "grpc-addr", DefaultgRPCAddr, "gRPC connection address")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(gRPCAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %s", err.Error())
	}
	defer conn.Close()

	client := pb.NewDBProviderClient(conn)
	_ = client

	r, err := client.Query(context.Background(), &pb.QueryRequest{Stmt: "SELECT * FROM foo"})
	if err != nil {
		log.Fatalf("failed to query: %s", err.Error())
	}
	fmt.Println(">>>>", r.Columns, r.Rows)
}
