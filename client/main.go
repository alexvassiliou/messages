package main

import (
	"context"
	"flag"
	"log"
	"time"

	api "github.com/alexvassiliou/messages/service"
	"google.golang.org/grpc"
)

func main() {
	// create a message using flags
	title := flag.String("title", "defaultTitle", "the title of your note")
	body := flag.String("body", "this is a default mesage", "the note you wish to save")
	flag.Parse()

	// Setup connection to the server, no credentials needed
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// instantiate client and pass in the server connection
	c := api.NewMessageServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// build request from the flags
	request := api.CreateRequest{
		Message: &api.Message{
			Title:   *title,
			Content: *body,
		},
	}

	// Call create on the server
	response, err := c.Create(ctx, &request)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", response)

}
