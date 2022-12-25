package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	pb "example.com/grpc-todo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	port = 50051
)

type server struct {
	pb.UnimplementedTodoServiceServer
}

func (s *server) CreateTodo(ctx context.Context, in *pb.NewToDo) (*pb.Todo, error) {
	//create a random ID for this todo

	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "Cannot read headers")
	}
	if len(headers.Get("x-trace_id")) == 0 {
		log.Println("ERROR: missing trace_id")
		return nil, status.Error(codes.InvalidArgument, "no trace_id")
	}
	trace_id := headers.Get("x-trace_id")[0]
	log.Printf("%s: New Req: %v", trace_id, in.GetName())

	// fmt.Printf("%+v\n", headers)
	// Get the value for the "authorization" key from the headers
	authorization := headers.Get("authorization")
	log.Printf("Authorization header: %v\n", authorization)
	if len(authorization) == 0 || authorization[0] != "myverysecretkey" {
		log.Println("ERROR: bad auth")
		return nil, status.Error(codes.Unauthenticated, "Not authenticated")
	}

	if in.GetName() == "" {
		// return one of these codes https://pkg.go.dev/google.golang.org/grpc/codes
		log.Println("ERROR: invalid arg: name")
		return nil, status.Error(codes.InvalidArgument, "Invalid argument: Name")
	}

	//todo: add the other bits of the todo - desc and done
	return &pb.Todo{Id: strconv.Itoa(rand.Intn(100-1) + 1), Name: in.GetName()}, nil
}

func main() {
	//used for creating random IDs for the todos
	rand.Seed(time.Now().UnixNano())

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, &server{})
	reflection.Register(s)

	log.Printf("server listening on%v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
