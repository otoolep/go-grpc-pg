# Generating new protobuf code

```
export PATH=$PATH:grpc/protoc/bin
export PATH=$PATH:$GOPATH/bin
cd $GOPATH/src/github.com/otoolep/go-grpc-pg/proto
protoc service.proto --go_out=plugins=grpc:.
```
