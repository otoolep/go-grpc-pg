# Generating new protobuf code

```
export PATH=$PATH:$GOPATH/bin
export PATH=$PATH:<path to decompressed protobuf toolchain>/protoc/bin
cd $GOPATH/src/github.com/otoolep/go-grpc-pg/proto
protoc service.proto --go_out=plugins=grpc:.
```

Once a new `.go` file has been generated, this file should be committed.
