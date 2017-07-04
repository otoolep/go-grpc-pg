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

// Service represents a gRPC service that communicates with a database backend.
type Service struct {
	grpc *grpc.Server
	db   *sql.DB

	ln   net.Listener
	addr string

	logger *log.Logger
}

// New returns an instantiated service.
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

// Addr returns the bind address of the gRPC service.
func (s *Service) Addr() string {
	return s.addr
}

// Open opens the service, starting it listening on the configured address.
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

// Close closes the service.
func (s *Service) Close() error {
	s.grpc.GracefulStop()
	s.ln = nil
	s.logger.Println("gRPC server stopped")
	return nil
}

// gprcService is an unexported type, that is the same type as Service.
//
// Having the methods that the gRPC service requires on this type means that even though
// the methods are exported, since the type is not, these methods are not visible outside
// this package.
type gprcService Service

// Query implements the Query interface of the gRPC service.
func (g *gprcService) Query(c context.Context, q *pb.QueryRequest) (*pb.QueryResponse, error) {
	rows, err := g.db.Query(q.Stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get the column names.
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	response := pb.QueryResponse{
		Columns: cols,
	}

	// Iterate through each row returned by the query.
	for rows.Next() {
		row := make([]string, len(cols))
		// Get a set of pointers to the strings allocated above.
		row_i := make([]interface{}, len(cols))
		for i, _ := range row {
			row_i[i] = &row[i]
		}

		if err := rows.Scan(row_i...); err != nil {
			return nil, err
		}

		// Add the latest rows to existing rows.
		response.Rows = append(response.Rows, &pb.Row{Values: row})
	}

	return &response, nil
}

// Exec implements the Exec interface of the gRPC service.
func (g *gprcService) Exec(c context.Context, e *pb.ExecRequest) (*pb.ExecResponse, error) {
	_, err := g.db.Exec(e.Stmt)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
