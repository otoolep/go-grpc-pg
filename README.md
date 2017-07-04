go-grpc-pg [![Circle CI](https://circleci.com/gh/otoolep/go-grpc-pg/tree/master.svg?style=svg)](https://circleci.com/gh/otoolep/go-grpc-pg/tree/master) [![GoDoc](https://godoc.org/github.com/otoolep/go-grpc-pg?status.svg)](https://godoc.org/github.com/otoolep/go-grpc-pg) [![Go Report Card](https://goreportcard.com/badge/github.com/otoolep/go-grpc-pg)](https://goreportcard.com/report/github.com/otoolep/go-grpc-pg)
======

A simple service demonstrating Go, gRPC, and PostgreSQL. Integration with [CircleCI](http://www.circleci.com) included.

## Building go-grpc-pg
*Building go-httpd requires Go 1.8 or later. [gvm](https://github.com/moovweb/gvm) is a great tool for installing and managing your versions of Go.*

Download and build it like so:
```
mkdir go-grpc-pg # Or any directory of your choice
cd go-grpc-pg/
export GOPATH=$PWD
go get -t github.com/otoolep/go-grpc-pg
cd src/github.com/otoolep/go-grpc-pg
go install ./...
```
Some people consider using a distinct `GOPATH` environment variable for each project _doing it wrong_. In practise I, and many other Go programmers, find this actually most convenient.

### Optional step to speed up testing
Unit testing actually uses SQLite, which is built as part of the test suite -- there is no need to install SQLite separately. However the compilation of SQLite is slow, and quickly becomes tedious if continually repeated. To avoid continual compilation every test run, execute the following command:
```
cd $GOPATH
go install github.com/mattn/go-sqlite3
```

## Running go-grpc-pg
Once built as per above, launch the server as follows:
```
$GOPATH/bin/server
```
This assumes PostgreSQL is listening on `localhost`, port 5432. Run `$GOPATH/bin/server -h` to learn the full configuration the server expects of PostgreSQL.

### Generating queries
Assuming the server is up and running, execute the client as follows.
```
$GOPATH/bin/client
```