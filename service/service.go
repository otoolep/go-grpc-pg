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

	pb.RegisterDBProviderServer(s.g, (*gprcService)(&s))
	return &s
}

func (s *Service) Addr() string {
	return s.addr.String()
}

type gprcService Service

func (g *gprcService) Query(c context.Context, q *pb.QueryRequest) (*pb.QueryResponse, error) {
	_, err := g.db.Query(q.Stmt)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (g *gprcService) Exec(c context.Context, e *pb.ExecRequest) (*pb.ExecResponse, error) {
	_, err := g.db.Exec(e.Stmt)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
