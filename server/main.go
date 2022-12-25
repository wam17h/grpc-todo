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
	log.Printf("Received: %v", in.GetName())
	//create a random ID for this todo

	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "Cannot read headers")
	}
	// fmt.Printf("%+v\n", headers)
	// Get the value for the "authorization" key from the headers
	authorization := headers.Get("authorization")
	fmt.Println("Authorization header:", authorization)

	// ctype := headers.Get("content-type")
	// fmt.Println("content-type:", ctype)

	// todo test how an error is handled
	if in.GetName() == "error" {
		// return one of these codes https://pkg.go.dev/google.golang.org/grpc/codes
		return nil, status.Error(codes.Unauthenticated, "Not authenticated")
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
