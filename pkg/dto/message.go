package dto

import "time"

type Message struct {
	Id         string    `json:"id" bson:"_id"`
	SenderId   string    `json:"senderId" bson:"senderId"`
	Text       string    `json:"text" bson:"text"`
	Timestamp  time.Time `json:"timestamp" bson:"timestamp"`
	Img        string    `json:"img" bson:"img"`
	Attachment string    `json:"attachment" bson:"attachment"`
	ReceiverId string    `json:"receiverId" bson:"receiverId"`
}

/*
{
	"id":"d32fe01b-0198-4f7f-a499-b2580f052e5a",
	"text":"aloha",
	"image":null,
	"video":null,
	"createdAt":1589773362924,
	"user":""
	"quickReplies":null,
	"customProperties":null
}
*/
