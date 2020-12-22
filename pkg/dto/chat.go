package dto

type Chat struct {
	ChatId        string               `json:"chat_id" bson:"chat_id"`
	Type          bool                 `json:"type" bson:"type"` // true == direct, false == group
	Participants  map[string]string    `json:"-" bson:"participants"`
	MParticipants map[string]*ChatUser `json:"participants" bson:"-"`
	Count         int64                `json:"count" bson:"count"`
	Name          string               `json:"name" bson:"name"`
	// Participant   *string              `json:"-" bson:"participants"`
}

type ChatUser struct {
	Uid  string `json:"uid" bson:"uid"`
	Name string `json:"name" bson:"name"`
}
