package handlers

import (
	"context"

	pb "getitqec.com/server/chat/pkg/api/v1"
	"getitqec.com/server/chat/pkg/commons"
	"getitqec.com/server/chat/pkg/dto"
	"getitqec.com/server/chat/pkg/logger"
	"getitqec.com/server/chat/pkg/model"
)

type SendMessageHandler struct {
	Model model.ChatModelI
}

func (s *SendMessageHandler) SendMessage(ctx context.Context, req *pb.Message) (*pb.Message, error) {
	// userID, err := commons.GetUserID(ctx)
	token, err := commons.VerifyGoogleAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	logger.Log.Debug("\tHandler send message")
	m, err := s.Model.SendMessage(ctx, dto.PBMessage2Message(req, token.UserId))
	if err != nil {
		return nil, err
	}
	return dto.Message2PBMessage(m), err
}
