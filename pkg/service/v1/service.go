package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	pb "getitqec.com/server/chat/pkg/api/v1"
	"getitqec.com/server/chat/pkg/commons"
	"getitqec.com/server/chat/pkg/dto"
	"getitqec.com/server/chat/pkg/handlers"
	"getitqec.com/server/chat/pkg/logger"
	"getitqec.com/server/chat/pkg/model"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/metadata"
	//pb "./proto"
)

var httpClient = &http.Client{}

// logger is to mock a sophisticated logging system. To simplify the example, we just print out the content.
// func logger(format string, a ...interface{}) {
// 	fmt.Printf("LOG:\t"+format+"\n", a...)
// }

// var (
// 	//port = flag.Int("port", 50051, "the port to serve on")

// 	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
// 	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
// )

// Server class
type Server struct {
	model model.ChatModelI

	pb.UnimplementedChatServiceServer
}

func (s *Server) SendMessage(ctx context.Context, req *pb.Message) (*pb.Message, error) {
	// return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Println(md)
	handler := &handlers.SendMessageHandler{Model: s.model}
	return handler.SendMessage(ctx, req)
}
func (s *Server) UpdateMessage(ctx context.Context, m *pb.Message) (*pb.Message, error) {
	token, err := commons.VerifyGoogleAccessToken(ctx)
	if err != nil {
		return nil, commons.ErrInvalidToken
	}
	message, err := s.model.UpdateMessage(ctx, dto.PBMessage2Message(m, token.UserId))
	if err != nil {
		return nil, err
	}
	return dto.Message2PBMessage(message), err
}
func (s *Server) DeleteMessage(ctx context.Context, m *pb.Message) (*empty.Empty, error) {
	// return nil, status.Errorf(codes.Unimplemented, "method DeleteMessage not implemented")
	_, err := s.model.DeleteMessage(ctx, m.Id)
	return &empty.Empty{}, err
}

func (s *Server) Connect(stream pb.ChatService_ConnectServer) error {
	// md, _ := metadata.FromIncomingContext(stream.Context())
	// fmt.Println(md)
	token, err := commons.VerifyGoogleAccessToken(stream.Context())
	if err != nil {
		return commons.ErrInvalidToken
	}
	// token := &oauth2.Tokeninfo{}
	// token.UserId, _ = commons.GetUserID(stream.Context())
	s.model.RegisterStream(stream.Context(), stream, token.UserId)
	// messages, err := s.model.GetMessagesByUser(stream.Context(), token.UserId)
	// if err == nil {
	// 	for _, message := range messages {
	// 		e := stream.Send(dto.Message2PBMessage(message))
	// 		if e == nil {
	// 			fmt.Println("error getting messages on connect [for loop]")
	// 		}
	// 	}
	// } else {
	// 	fmt.Println("error getting messages on connect")
	// }
	for {
		readMessage, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Close connection")
			s.model.UnregisterStream(stream.Context(), token.UserId)
			return nil
		}
		if err != nil {
			fmt.Println("Err receive")
			s.model.UnregisterStream(stream.Context(), token.UserId)
			return err
		}
		if readMessage.Id == "" {
			return nil
		}
		fmt.Println("Receive " + readMessage.Id)
		s.model.DeleteMessage(stream.Context(), readMessage.Id)
		// for _, note := range s.routeNotes[key] {
		// 	if err := stream.Send(note); err != nil {
		// 		return err
		// 	}
		// }
	}
}

// NewServer return new auth server service
func NewServer(model model.ChatModelI, ctx context.Context) *Server {
	server := &Server{}
	server.model = model
	// close stream when server shutting down
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Log.Warn("closing down gRPC service...")

			server.CloseAllConnectStream(ctx)

			<-ctx.Done()
		}
	}()
	return server
}

func (s *Server) CloseAllConnectStream(ctx context.Context) {
	s.model.UnregisterAllStream(ctx)
}
