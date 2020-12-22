package model

import (
	"context"

	pb "getitqec.com/server/chat/pkg/api/v1"
	"getitqec.com/server/chat/pkg/dto"
)

type ChatModelI interface {
	// // Chat

	// GetChat(ctx context.Context, chatId string) (*dto.Chat, error)
	// GetChatByUser(ctx context.Context, users []string) ([]*dto.Chat, error)
	// GetDirectChatByUser(ctx context.Context, user1, user2 string) (*dto.Chat, error)
	// UpdateChat(ctx context.Context, chat *dto.Chat) error
	// DeleteChat(ctx context.Context, chatId string) (*dto.Chat, error)
	RegisterStream(ctx context.Context, stream pb.ChatService_ConnectServer, userId string)
	UnregisterStream(ctx context.Context, userId string)
	UnregisterAllStream(ctx context.Context)
	// Message

	SendMessage(ctx context.Context, message *dto.Message) (*dto.Message, error)
	GetMessage(ctx context.Context, messageId string) (*dto.Message, error)
	UpdateMessage(ctx context.Context, message *dto.Message) (*dto.Message, error)
	DeleteMessage(ctx context.Context, messageId string) (*dto.Message, error)
	GetMessagesByUser(ctx context.Context, receiverId string) ([]*dto.Message, error)
}
