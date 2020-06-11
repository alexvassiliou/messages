package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	api "github.com/alexvassiliou/messages/service"
	"google.golang.org/grpc"
)

func main() {
	// Create the context and specify a port for the listener
	ctx := context.Background()
	port := "8080"

	// Create the db store and pass to the new server
	var db []api.Message
	s := api.NewServer(db)

	// Launch ther gRPC server
	log.Fatal(runServer(ctx, port, s))

}

func runServer(ctx context.Context, port string, s api.MessageServiceServer) error {
	// Open a TCP listener at the given port address
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("runServer: %s", err)
	}

	// Open a grpc server
	server := grpc.NewServer()
	// Register the message server passed through to the grpc server
	api.RegisterMessageServiceServer(server, s)

	// create channel to use to shut down the server when the app is interrupted
	c := make(chan os.Signal, 1)
	// signal the channel when the specified interruption is recieved
	signal.Notify(c, os.Interrupt)
	// go routine ranges over the channel and shuts down the server when the signal is received
	go func() {
		for range c {
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	log.Println("starting gRPC server...")

	// Open the server to incoming requests
	return server.Serve(listen)

}
