go-grpc-pg [![Circle CI](https://circleci.com/gh/otoolep/go-grpc-pg/tree/master.svg?style=svg)](https://circleci.com/gh/otoolep/go-grpc-pg/tree/master) [![GoDoc](https://godoc.org/github.com/otoolep/go-grpc-pg?status.svg)](https://godoc.org/github.com/otoolep/go-grpc-pg) [![Go Report Card](https://goreportcard.com/badge/github.com/otoolep/go-grpc-pg)](https://goreportcard.com/report/github.com/otoolep/go-grpc-pg)
======

A simple service demonstrating Go, gRPC, and PostgreSQL. Integration with [CircleCI](http://www.circleci.com) included.

## Building and running go-grpc-pg
*Building go-httpd requires Go 1.8 or later. [gvm](https://github.com/moovweb/gvm) is a great tool for installing and managing your versions of Go.*

Starting and running go-grpc-pg is easy. Download and build it like so:
```
mkdir go-grpc-pg # Or any directory of your choice
cd go-grpc-pg/
export GOPATH=$PWD
go get github.com/otoolep/go-grpc-pg
cd src/github.com/otoolep/go-grpc-pg
go install ./...
```
Some people consider using a distinct `GOPATH` environment variable for each project _doing it wrong_. In practise I, and many other Go programmers, find this actually most convenient.

Run it like so:
```
$GOPATH/bin/server
```
