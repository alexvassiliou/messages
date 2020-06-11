package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	api "github.com/alexvassiliou/messages/service"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	stream := flag.Bool("stream", false, "set true for chat app")
	port := flag.String("port", "8080", "select the port")
	flag.Parse()

	// check if streaming or unary service is selected
	if *stream {
		var connections []*api.Connection
		s := &api.StreamServer{connections}

		// Launch ther gRPC server
		log.Fatal(runStreamingServer(ctx, *port, s))

	} else {
		var db []api.Message
		s := api.NewServer(db)

		// Launch ther gRPC server
		log.Fatal(runUnaryServer(ctx, *port, s))
	}

}

func runUnaryServer(ctx context.Context, port string, s api.MessageServiceServer) error {
	// Open a TCP listener at the given port address
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("run: %s", err)
	}

	// Open a grpc server
	server := grpc.NewServer()

	// register the relevant service to the grpc server
	api.RegisterMessageServiceServer(server, s)

	// create channel to use to shut down the server when the app is interrupted
	c := make(chan os.Signal, 1)
	// signal the channel when the specified interruption is recieved
	signal.Notify(c, os.Interrupt)
	// go routine ranges over the channel and shuts down the server when the signal is received
	go func() {
		for range c {
			log.Println("shutting down gRPC streaming server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	log.Println("starting gRPC unary server...")

	// Open the server to incoming requests
	return server.Serve(listen)
}

func runStreamingServer(ctx context.Context, port string, s api.StreamingMessageServiceServer) error {
	// Open a TCP listener at the given port address
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("runStreamingServer: %s", err)
	}

	// Open a grpc server
	server := grpc.NewServer()

	// register the relevant service to the grpc server
	api.RegisterStreamingMessageServiceServer(server, s)

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

	log.Println("starting gRPC streaming server...")

	// Open the server to incoming requests
	return server.Serve(listen)
}
