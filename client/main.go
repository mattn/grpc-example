package main

import (
	"io"
	"strconv"

	pb "github.com/mattn/grpc-example/proto"
	"github.com/mattn/sc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const address = "127.0.0.1:11111"

func add(client pb.CustomerServiceClient, name string, age int) error {
	person := &pb.Person{
		Name: name,
		Age:  int32(age),
	}
	_, err := client.AddPerson(context.Background(), person)
	return err
}

func list(client pb.CustomerServiceClient) error {
	stream, err := client.ListPerson(context.Background(), new(pb.RequestType))
	if err != nil {
		return err
	}
	for {
		person, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		grpclog.Println(person)
	}
	return nil
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewCustomerServiceClient(conn)

	(&sc.Cmds{
		{
			Name: "list",
			Desc: "list: listing person",
			Run: func(c *sc.C, args []string) error {
				return list(client)
			},
		},
		{
			Name: "add",
			Desc: "add [name] [age]: add person",
			Run: func(c *sc.C, args []string) error {
				if len(args) != 2 {
					return sc.UsageError
				}
				name := args[0]
				age, err := strconv.Atoi(args[1])
				if err != nil {
					return err
				}
				return add(client, name, age)
			},
		},
	}).Run(&sc.C{
		Desc: "grpc example",
	})
}
