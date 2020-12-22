package dao

import (
	"context"

	"getitqec.com/server/chat/pkg/dto"
)

type IChatDAO interface {
	//InitIndex(ctx context.Context) error

	Create(ctx context.Context, chat *dto.Chat) error
	Get(ctx context.Context, chatId string) (*dto.Chat, error)
	Update(ctx context.Context, chat *dto.Chat) error
	Delete(ctx context.Context, chatId string) (*dto.Chat, error)
	GetByUser(ctx context.Context, users []string) ([]*dto.Chat, error)
	GetDirectByUser(ctx context.Context, user1, user2 string) (*dto.Chat, error)
}

type IMessageDAO interface {
	InitIndex(ctx context.Context) error

	Create(ctx context.Context, message *dto.Message) error
	Get(ctx context.Context, messageId string) (*dto.Message, error)
	Update(ctx context.Context, message *dto.Message) error
	Delete(ctx context.Context, messageId string) (*dto.Message, error)
	Search(ctx context.Context, searchString string) ([]*dto.Message, error)
	GetMessagesByUser(ctx context.Context, receiverId string) ([]*dto.Message, error)
}
