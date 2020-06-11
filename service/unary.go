package api

import (
	context "context"
	"fmt"

	"github.com/google/uuid"
)

//Server is the message server it holds a collection of messages
type server struct {
	db []Message
}

// NewServer instantiates a message server
func NewServer(db []Message) MessageServiceServer {
	return &server{
		db: db,
	}
}

func (s *server) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	// Create the id for the new message
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("Create: %s", err)
	}

	// Create the message from the request
	m := Message{
		Id:      id.String(),
		Content: req.Message.Content,
		Time:    req.Message.Time,
	}

	// Add the message to the db store
	s.db = append(s.db, m)

	// Success message that the message has been added to the database
	success := fmt.Sprintf("New entry created: %v", m.Id)
	fmt.Println(success)

	// Return the response with the message id
	return &CreateResponse{
		Id: m.Id,
	}, nil
}

func (s *server) Read(ctx context.Context, req *ReadRequest) (*ReadResponse, error) {
	var result Message

	for _, message := range s.db {
		if message.Id == req.Id {
			result = message
		}
	}

	return &ReadResponse{
		Message: &result,
	}, nil
}
