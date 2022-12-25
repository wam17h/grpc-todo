package main

import (
	"context"
	"testing"

	pb "example.com/grpc-todo/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestCreateTodo(t *testing.T) {
	testCases := []struct {
		name            string
		req             *pb.NewToDo
		headers         map[string]string
		expectedErrCode codes.Code
	}{
		{
			name: "No trace id",
			req:  &pb.NewToDo{Name: "test one"},
			headers: map[string]string{
				"authorization": "myverysecretkey",
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "no auth",
			req:  &pb.NewToDo{Name: "test one"},
			headers: map[string]string{
				"x-trace_id": "1225-two-words",
			},
			expectedErrCode: codes.Unauthenticated,
		},
		{
			name: "test ok",
			req:  &pb.NewToDo{Name: "test one"},
			headers: map[string]string{
				"x-trace_id":    "1225-two-words",
				"authorization": "myverysecretkey",
			},
			expectedErrCode: codes.OK,
		},
	}

	s := server{}
	for _, tc := range testCases {
		tCase := tc
		t.Run(tCase.name, func(t *testing.T) {
			//	t.Parallel()
			t.Logf("Testing %s ..", tCase.name)
			newHeaders := metadata.New(tCase.headers)
			ctx := metadata.NewIncomingContext(context.Background(), newHeaders)
			// call
			response, err := s.CreateTodo(ctx, tCase.req)
			if err != nil {
				t.Log("Error: ", err)

				s, ok := status.FromError(err)
				if !ok {
					// The error is not a gRPC error, handle it as a regular error
					t.Error(err)
				}
				if s.Code() != tCase.expectedErrCode {
					t.Errorf("expected %v, got %v", tCase.expectedErrCode, s.Code())
				}
			}

			if tCase.expectedErrCode > 0 && err == nil {
				t.Errorf("expected %v, got no error", tCase.expectedErrCode)
			}

			t.Log("Resp: ", response)
		})
	}

}
