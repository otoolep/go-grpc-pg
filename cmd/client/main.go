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
		prompt(f)

		for s.Scan() {
			line := strings.TrimSpace(s.Text())
			if line == "" {
				continue
			}

			if isQuery(line) {
				resp, err := client.Query(context.Background(), &pb.QueryRequest{s.Text()})
				if err != nil {
					fmt.Printf("query error: %s\n", err.Error())
					prompt(f)
					continue
				}
				for _, r := range resp.Rows {
					for _, v := range r.Values {
						fmt.Printf("%s\t", v)
					}
					fmt.Println()
				}
			} else {
				resp, err := client.Exec(context.Background(), &pb.ExecRequest{s.Text()})
				if err != nil {
					fmt.Printf("exec error: %s\n", err.Error())
					prompt(f)
					continue
				}
				fmt.Printf("Last Insert ID: %d, rows affected: %d\n", resp.LastInsertId, resp.RowsAffected)
			}

			prompt(f)
		}
	}
}

func isQuery(line string) bool {
	index := strings.Index(line, " ")
	if index >= 0 {
		return strings.ToUpper(line[:index]) == "SELECT"
	}
	return false
}

func prompt(w *bufio.Writer) {
	w.Write(Prompt)
	w.Flush()
}
