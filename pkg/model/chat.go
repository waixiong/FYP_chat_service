package model

import (
	"context"
	"fmt"

	pb "getitqec.com/server/chat/pkg/api/v1"
	"getitqec.com/server/chat/pkg/dto"
	"getitqec.com/server/chat/pkg/logger"
)

func (m *ChatModel) RegisterStream(ctx context.Context, stream pb.ChatService_ConnectServer, userId string) {
	// fmt.Println(stream)
	// fmt.Println(userId)
	m.ClientStream[userId] = stream
	fmt.Println(m.ClientStream)

	messages, err := m.GetMessagesByUser(ctx, userId)
	if err != nil {
		logger.Log.Error("\tFail to get Messages")
	}
	fmt.Printf("\t%d new messages\n", len(messages))
	for _, message := range messages {
		pbMessage := dto.Message2PBMessage(message)
		err = stream.Send(pbMessage)
		if err != nil {
			logger.Log.Error("\tMessage Fail To Send")
		}
	}
}

func (m *ChatModel) UnregisterStream(ctx context.Context, userId string) {
	// m.ClientStream[userId] = nil
	if m.ClientStream[userId] == nil {
		return
	}
	m.ClientStream[userId].Context().Done()
	delete(m.ClientStream, userId)
}

func (m *ChatModel) UnregisterAllStream(ctx context.Context) {
	// m.ClientStream[userId] = nil
	for _, stream := range m.ClientStream {
		stream.Context().Done()
	}
	for id := range m.ClientStream {
		delete(m.ClientStream, id)
	}
}
