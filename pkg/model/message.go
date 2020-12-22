package model

import (
	"context"
	"fmt"
	"strconv"

	mailnotificationproto "getitqec.com/server/chat/pkg/api/clients/mailnotification/v1"
	"getitqec.com/server/chat/pkg/commons"
	"getitqec.com/server/chat/pkg/dto"
	"getitqec.com/server/chat/pkg/logger"
	"getitqec.com/server/chat/pkg/protocol/grpcClient"
)

func (m *ChatModel) SendMessage(ctx context.Context, message *dto.Message) (*dto.Message, error) {
	logger.Log.Debug("\tModel send message")
	message.Id = GenerateMessageId()
	message.Timestamp = commons.MalaysiaTimeNow()
	err := m.MessageDAO.Create(ctx, message)

	// Notification
	if err != nil {
		// logger.Log.Debug(fmt.Sprintf("\tModel %o", err))
		return nil, err
	}
	logger.Log.Debug("\tAdded to database")
	client, conn, err := grpcClient.GetMailNotificationClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	data := map[string]string{}
	data["action"] = "receiveChat"
	data["params.senderId"] = message.SenderId
	data["params.content"] = message.Text
	data["params.timestamp"] = strconv.FormatInt(message.Timestamp.UnixNano()/1000000, 10)
	_, err = client.PushToUser(ctx, &mailnotificationproto.PushRequest{
		Target:       message.ReceiverId,
		Notification: true,
		Title:        "ImageChat Message",
		Body:         "A message have been receive",
		ImgUrl:       "",
		Data:         data,
	})
	if err != nil {
		logger.Log.Debug(fmt.Sprintf("\tChat send Notification %s : %o", message.ReceiverId, err))
	}
	// if m.ClientStream[message.ReceiverId] != nil {
	// 	err = m.ClientStream[message.ReceiverId].Send(dto.Message2PBMessage(message))
	// 	if err != nil {
	// 		logger.Log.Debug(fmt.Sprintf("\tChat send Message %s : %o", message.ReceiverId, err))
	// 	}
	// }
	if stream, ok := m.ClientStream[message.ReceiverId]; ok {
		pbMessage := dto.Message2PBMessage(message)
		err = stream.Send(pbMessage)
		if err != nil {
			logger.Log.Error("\tStream exist but Message Fail To Send")
		}
	}
	return message, nil
}

func (m *ChatModel) GetMessage(ctx context.Context, messageId string) (*dto.Message, error) {
	token, err := commons.VerifyGoogleAccessToken(ctx)
	if err != nil {
		return nil, commons.ErrInvalidToken
	}
	message, err := m.MessageDAO.Get(ctx, messageId)
	if message.ReceiverId != token.UserId || message.SenderId != token.UserId {
		return nil, commons.NotAuthorized
	}
	return message, nil
}

func (m *ChatModel) UpdateMessage(ctx context.Context, message *dto.Message) (*dto.Message, error) {
	token, err := commons.VerifyGoogleAccessToken(ctx)
	if err != nil {
		return nil, commons.ErrInvalidToken
	}
	messageO, err := m.MessageDAO.Get(ctx, message.Id)
	if messageO.SenderId != token.UserId {
		return nil, commons.NotAuthorized
	}
	return message, m.MessageDAO.Update(ctx, message)
	// return nil, nil
}

func (m *ChatModel) DeleteMessage(ctx context.Context, messageId string) (*dto.Message, error) {
	token, err := commons.VerifyGoogleAccessToken(ctx)
	if err != nil {
		return nil, commons.ErrInvalidToken
	}
	message, err := m.MessageDAO.Get(ctx, messageId)
	if err != nil {
		return nil, err
	}
	if message.ReceiverId != token.UserId || message.SenderId != token.UserId {
		return nil, commons.NotAuthorized
	}
	return m.MessageDAO.Delete(ctx, messageId)
}

func (m *ChatModel) GetMessagesByUser(ctx context.Context, receiverId string) ([]*dto.Message, error) {
	return m.MessageDAO.GetMessagesByUser(ctx, receiverId)
}
