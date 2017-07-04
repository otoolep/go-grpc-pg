package service

import (
	"database/sql"
	"testing"

	pb "github.com/otoolep/go-grpc-pg/proto"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Test_NewService tests creation of a service in the simplest manner.
func Test_NewService(t *testing.T) {
	db := mustCreateSQLiteDatabase()
	defer db.Close()

	s := New(":0", db)
	if s == nil {
		t.Fatalf("failed to create service")
	}
}

// Test_NewServiceQueries performs a full test of the service.
func Test_NewServiceOpen(t *testing.T) {
	db := mustCreateSQLiteDatabase()
	defer db.Close()

	s := New(":0", db)
	if s == nil {
		t.Fatalf("failed to create service")
	}

	// Test opening of the service.
	if err := s.Open(); err != nil {
		t.Fatalf("failed to open service: %s", err.Error())
	}

	// Connect to the gRPC service.
	addr := s.Addr()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to connect to service: %s", err.Error())
	}
	defer conn.Close()

	// Create a client.
	client := pb.NewDBProviderClient(conn)

	// Test creation of a table.
	_, err = client.Exec(context.Background(), &pb.ExecRequest{Stmt: `CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, name TEXT)`})
	if err != nil {
		t.Fatalf("failed to create table: %s", err.Error())
	}

	// Run some queries.
	r, err := client.Query(context.Background(), &pb.QueryRequest{Stmt: "SELECT * FROM foo"})
	if err != nil {
		t.Fatalf("failed to query empty table: %s", err.Error())
	}
	if exp, got := 2, len(r.Columns); exp != got {
		t.Fatalf("wrong number of columns returned, exp: %d, got: %d ", exp, got)
	}

	// Insert some data.
	_, err = client.Exec(context.Background(), &pb.ExecRequest{Stmt: `INSERT INTO foo(name) VALUES("fiona")`})
	if err != nil {
		t.Fatalf("failed to insert a row: %s", err.Error())
	}

	// Run more queries.
	r, err = client.Query(context.Background(), &pb.QueryRequest{Stmt: "SELECT * FROM foo"})
	if err != nil {
		t.Fatalf("failed to query non-empty table: %s", err.Error())
	}
	if exp, got := 2, len(r.Columns); exp != got {
		t.Fatalf("wrong number of columns returned, exp: %d, got: %d ", exp, got)
	}
}

// mustCreateSQLiteDatabase creates an in-memory SQLite database, or panics.
func mustCreateSQLiteDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic("failed to create in-memory SQLite database")
	}
	return db
}
