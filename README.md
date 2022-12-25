# GRPC TODO

A proof of concept GRPC server and client written in Golang.  
Based on the code from: https://www.golinuxcloud.com/golang-grpc/

## Installing protobuffer

* for macs: `brew install protobuff` - this will install the protoc command line tool  
* for linux: `sudo apt install -y protobuf-compiler`

### golang libs for protobuffers

```
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

NOTE: You should add the `protoc-gen-go-grpc` to your PATH

```
PATH="${PATH}:${HOME}/go/bin"

```


## Instructions:  
* check out the code
* run `make server`, this will start a server listening on port 50051
* open a new terminal window and run `make client`, this will start a client that will send todo records to the server
* if you update the `proto/todo.proto` definition file you will have to regenerate the protopuf code - this can be done by running `make protoc`
* You can also send requests to the server by using the [grpc-client-cli](https://github.com/vadimi/grpc-client-cli) tool

## Updates:  
* Added headers to the outgoing request via metadata.NewOutgoingContext(ctx, newHeaders)  
* Can read headers on the server side via metadata.FromIncomingContext(ctx)

* Added error handling eg status.Error(codes.Unauthenticated, "Not authenticated")  
To return a standard error code from: https://pkg.go.dev/google.golang.org/grpc/codes

