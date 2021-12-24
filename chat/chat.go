package chat

import (
	"log"
	"golang.org/x/net/context"
)

type Server struct{
	UnimplementedChatServiceServer
}

func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error){
	log.Printf("Received message body from client: %v", message)
	return &Message{Body: "Hello From the Server!", 
		MessageNumber: message.MessageNumber,
		Timestamps: message.Timestamps,
		}, nil
}