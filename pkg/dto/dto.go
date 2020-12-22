package dto

import (
	pb "getitqec.com/server/chat/pkg/api/v1"
	"getitqec.com/server/chat/pkg/commons"
)

func PBMessage2Message(req *pb.Message, senderId string) *Message {
	return &Message{
		SenderId:   senderId,
		Id:         req.Id,
		Text:       req.Text,
		ReceiverId: req.ReceiverId,
		Timestamp:  commons.MilliToTime(req.Timestamp),
		Img:        req.Img,
		Attachment: req.Attachment,
	}
}

func Message2PBMessage(m *Message) *pb.Message {
	return &pb.Message{
		SenderId:   m.SenderId,
		Id:         m.Id,
		Timestamp:  commons.TimeToMilli(m.Timestamp), //milliseconds
		Text:       m.Text,
		ReceiverId: m.ReceiverId,
		Img:        m.Img,
		Attachment: m.Attachment,
	}
}
