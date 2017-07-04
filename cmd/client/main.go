package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	pb "github.com/otoolep/go-grpc-pg/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Command line defaults
const (
	DefaultgRPCAddr = "localhost:11000"
)

var Prompt = []byte(`>> `)

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

	s := bufio.NewScanner(os.Stdin)
	f := bufio.NewWriter(os.Stdout)

	for {
		f.Write(Prompt)
		f.Flush()

		for s.Scan() {
			line := strings.Trim(s.Text(), " ")
			if line == "\n" {
				continue
			}

			_, err = client.Query(context.Background(), &pb.QueryRequest{s.Text()})
			if err != nil {
				fmt.Printf("failed to query: %s\n", err.Error())
			}

			f.Write(Prompt)
			f.Flush()
		}
	}
}
