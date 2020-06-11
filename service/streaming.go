package api

import (
	"context"
	"fmt"
	"sync"
)

// Connection is the object representing the strem
type Connection struct {
	stream StreamingMessageService_CreateStreamServer
	id     string
	active bool
	error  chan error
}

// StreamServer is used to register the service in gRPC
type StreamServer struct {
	Connection []*Connection
}

// CreateStream starts the stream of messages from the client
func (s *StreamServer) CreateStream(pconn *Connect, stream StreamingMessageService_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		error:  make(chan error),
	}

	s.Connection = append(s.Connection, conn)
	return <-conn.error
}

// BroadcastMessage defines when a message is returned to the client
func (s *StreamServer) BroadcastMessage(ctx context.Context, m *Message) (*Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	// range over the connection
	for _, conn := range s.Connection {
		wait.Add(1)
		// go routine to send each message and close the connection once sent
		go func(m *Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(m)
				fmt.Println("Sending message to: ", conn.stream)

				if err != nil {
					fmt.Printf("Error with Stream: %v - Error: %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}

		}(m, conn)
	}

	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
	return &Close{}, nil
}
