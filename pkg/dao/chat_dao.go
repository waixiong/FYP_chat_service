package dao

import (
	"context"

	"getitqec.com/server/chat/pkg/commons"
	"getitqec.com/server/chat/pkg/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatDAO struct {
	mongodb commons.MongoDB
}

func InitChatDAO(m commons.MongoDB) IChatDAO {
	return &ChatDAO{mongodb: m}
}

// func (v *ChatDAO) InitIndex(ctx context.Context) error {
// 	mod := mongo.IndexModel{
// 		Keys: bson.M{
// 			: 1, // index in ascending order
// 		}, Options: options.Index().SetExpireAfterSeconds(1),
// 	}

// 	collection := v.client.Database(constants.Cosmos).Collection(constants.AuthTokens)
// 	_, err := collection.Indexes().CreateOne(ctx, mod)
// 	return err
// }

func (v *ChatDAO) Create(ctx context.Context, chat *dto.Chat) error {
	option := options.Find()
	option = option.SetLimit(1)
	c, err := v.mongodb.Client().Database(commons.ChatDatabase).Collection(commons.ChatCollection).CountDocuments(ctx, bson.D{{Key: "chat_id", Value: chat.ChatId}})
	if err != nil {
		return err
	}
	if c == 0 {
		return v.mongodb.Create(ctx, commons.ChatDatabase, commons.ChatCollection, chat)
	}
	// created inventory for every outlets in business
	return commons.AlreadyExist
}

func (v *ChatDAO) Get(ctx context.Context, chatId string) (*dto.Chat, error) {
	result := v.mongodb.Read(ctx, commons.ChatDatabase, commons.ChatCollection, bson.D{{Key: "chat_id", Value: chatId}})
	if result.Err() != nil {
		return nil, result.Err()
	}
	item := &dto.Chat{}
	err := result.Decode(item)
	return item, err
}

func (v *ChatDAO) Update(ctx context.Context, chat *dto.Chat) error {
	return v.mongodb.Update(ctx, commons.ChatDatabase, commons.ChatCollection, bson.D{{Key: "chat_id", Value: chat.ChatId}}, bson.D{{"$set", chat}})
}

func (v *ChatDAO) Delete(ctx context.Context, chatId string) (*dto.Chat, error) {
	result := v.mongodb.Delete(ctx, commons.ChatDatabase, commons.ChatCollection, bson.D{{Key: "chat_id", Value: chatId}})
	if result.Err() != nil {
		return nil, result.Err()
	}
	chat := &dto.Chat{}
	err := result.Decode(chat)
	// delete inventory for every outlets in business
	return chat, err
}

func (v *ChatDAO) GetByUser(ctx context.Context, users []string) ([]*dto.Chat, error) {
	//{ '$exists' : true }
	query := bson.M{}
	for _, user := range users {
		query["participants."+user] = bson.M{"$exists": true}
	}
	query["count"] = len(users)
	// fmt.Printf("\t%o", query)
	raws, err := v.mongodb.BatchRead(ctx, commons.ChatDatabase, commons.ChatCollection, query)
	if err != nil {
		return nil, err
	}
	chats := []*dto.Chat{}
	for _, raw := range raws {
		chat := &dto.Chat{}
		err = bson.Unmarshal(*raw, chat)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (v *ChatDAO) GetDirectByUser(ctx context.Context, user1 string, user2 string) (*dto.Chat, error) {
	//{ '$exists' : true }
	query := bson.M{}
	query["participants."+user1] = bson.M{"$exists": true}
	query["participants."+user2] = bson.M{"$exists": true}
	query["count"] = 2
	query["type"] = true
	// fmt.Printf("\t%o", query)
	raws, err := v.mongodb.BatchRead(ctx, commons.ChatDatabase, commons.ChatCollection, query)
	if err != nil {
		return nil, err
	}
	if len(raws) == 0 {
		return nil, nil
	}
	chat := &dto.Chat{}
	err = bson.Unmarshal(*raws[0], chat)
	if err != nil {
		return nil, err
	}
	return chat, nil
}
