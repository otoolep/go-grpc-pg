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
	grpc *grpc.Server
	db   *sql.DB

	ln   net.Listener
	addr string

	logger *log.Logger
}

func New(addr string, db *sql.DB) *Service {
	s := Service{
		grpc:   grpc.NewServer(),
		db:     db,
		addr:   addr,
		logger: log.New(os.Stderr, "[service] ", log.LstdFlags),
	}

	pb.RegisterDBProviderServer(s.grpc, (*gprcService)(&s))
	return &s
}

func (s *Service) Addr() string {
	return s.addr
}

func (s *Service) Open() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.ln = ln
	s.logger.Println("listening on", s.addr)

	go func() {
		err := s.grpc.Serve(s.ln)
		if err != nil {
			s.logger.Println("gRPC Serve() returned:", err.Error())
		}
	}()

	return nil
}

func (s *Service) Close() error {
	s.grpc.GracefulStop()
	s.ln = nil
	s.logger.Println("gRPC server stopped")
	return nil
}

type gprcService Service

func (g *gprcService) Query(c context.Context, q *pb.QueryRequest) (*pb.QueryResponse, error) {
	rows, err := g.db.Query(q.Stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	response := pb.QueryResponse{
		Columns: cols,
	}

	for rows.Next() {
		row_i := make([]interface{}, len(cols))
		row := make([]string, len(cols))
		for i, _ := range row {
			row_i[i] = &row[i]
		}

		if err := rows.Scan(row_i...); err != nil {
			return nil, err
		}
		response.Rows = append(response.Rows, &pb.Row{Values: row})
	}
	return &response, nil
}

func (g *gprcService) Exec(c context.Context, e *pb.ExecRequest) (*pb.ExecResponse, error) {
	_, err := g.db.Exec(e.Stmt)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
