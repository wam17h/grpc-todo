package main

import (
	"context"
	"log"
	"time"

	pb "example.com/grpc-todo/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	ADDRESS = "localhost:50051"
)

type TodoTask struct {
	Name        string
	Description string
	Done        bool
}

func main() {
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewTodoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	todos := []TodoTask{
		{Name: "code review", Description: "Review grpc code", Done: false},
		{Name: "add auth to ctx", Description: "Add auth headers in the context", Done: false},
		{Name: "write medium tutorial", Description: "Write a better grpc tutorial", Done: false},
		{Name: "Go To Gym", Description: "Review grpc code", Done: false},
		{Name: "", Description: "This will produce an error", Done: false},
	}

	//ctx = context.WithValue(ctx, "authorization", "Bearer <your-jwt-token>")

	newHeaders := metadata.New(map[string]string{
		"x-trace_id":    "1225-two-words",
		"authorization": "myverysecretkey",
	})
	ctx = metadata.NewOutgoingContext(ctx, newHeaders)

	for _, todo := range todos {
		res, err := c.CreateTodo(ctx, &pb.NewToDo{Name: todo.Name, Description: todo.Description, Done: todo.Done})
		if err != nil {
			// log.Printf("%+v\n", err)
			log.Fatalf("could not create todo: %v", err)
		}

		log.Printf(`ID: %s, Name: %s`, res.GetId(), res.GetName())
	}

}
