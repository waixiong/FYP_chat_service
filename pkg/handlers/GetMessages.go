package handlers

import (
	"context"

	pb "getitqec.com/server/chat/pkg/api/v1"
	"getitqec.com/server/chat/pkg/dto"
	"getitqec.com/server/chat/pkg/model"
)

type GetMessagesHandler struct {
	Model model.ChatModelI
}

func (s *GetMessagesHandler) GetMessages(ctx context.Context, receiverId string) ([]*pb.Message, error) {
	// userID, err := commons.GetUserID(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	messages, err := s.Model.GetMessagesByUser(ctx, receiverId)
	if err != nil {
		return nil, err
	}
	pbMessages := []*pb.Message{}
	for _, message := range messages {
		pbMessages = append(pbMessages, dto.Message2PBMessage(message))
	}
	return pbMessages, nil
}
