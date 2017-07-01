package service

import (
	"database/sql"
	"log"
	"net"
	"os"

	pb "github.com/otoolep/go-grpc-pg/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Service struct {
	g  *grpc.Server
	db *sql.DB

	ln   net.Listener
	addr net.Addr

	logger *log.Logger
}

func NewService(ln net.Listener, db *sql.DB) *Service {
	s := Service{
		g:      grpc.NewServer(),
		db:     db,
		ln:     ln,
		addr:   ln.Addr(),
		logger: log.New(os.Stderr, "[service] ", log.LstdFlags),
	}

	pb.RegisterDBProviderServer(s.g, &s)
	return &s
}

func (s *Service) Addr() string {
	return s.addr.String()
}

func (s *Service) Query(context.Context, *pb.QueryRequest) (*pb.QueryResponse, error) {
	return nil, nil
}

func (s *Service) Exec(context.Context, *pb.ExecRequest) (*pb.ExecResponse, error) {
	return nil, nil
}
