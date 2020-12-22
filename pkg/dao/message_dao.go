package dao

import (
	"context"

	"getitqec.com/server/chat/pkg/commons"
	"getitqec.com/server/chat/pkg/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageDAO struct {
	mongodb commons.MongoDB
}

func InitMessageDAO(m commons.MongoDB) IMessageDAO {
	return &MessageDAO{mongodb: m}
}

func (v *MessageDAO) InitIndex(ctx context.Context) error {
	mod := mongo.IndexModel{
		Keys: bson.M{
			"timestamp": -1, // index in ascending order
		},
		// Options: options.Index().SetExpireAfterSeconds(259200), // 72 hours
	}

	collection := v.mongodb.Client().Database(commons.ChatDatabase).Collection(commons.MessageCollection)
	_, err := collection.Indexes().CreateOne(ctx, mod)
	return err
}

func (v *MessageDAO) Create(ctx context.Context, message *dto.Message) error {
	return v.mongodb.Create(ctx, commons.ChatDatabase, commons.MessageCollection, message)
}
func (v *MessageDAO) Get(ctx context.Context, messageId string) (*dto.Message, error) {
	result := v.mongodb.Read(ctx, commons.ChatDatabase, commons.MessageCollection, bson.D{
		{Key: "_id", Value: messageId},
	})
	if result.Err() != nil {
		return nil, result.Err()
	}
	message := &dto.Message{}
	err := result.Decode(message)
	return message, err
}
func (v *MessageDAO) Update(ctx context.Context, message *dto.Message) error {
	return v.mongodb.Update(ctx, commons.ChatDatabase, commons.MessageCollection, bson.M{"_id": message.Id}, message)
}
func (v *MessageDAO) Delete(ctx context.Context, messageId string) (*dto.Message, error) {
	result := v.mongodb.Delete(ctx, commons.ChatDatabase, commons.MessageCollection, bson.D{
		{Key: "_id", Value: messageId},
	})
	if result.Err() != nil {
		return nil, result.Err()
	}
	message := &dto.Message{}
	err := result.Decode(message)
	return message, err
}

func (v *MessageDAO) Search(ctx context.Context, searchString string) ([]*dto.Message, error) {
	return []*dto.Message{}, nil
}

func (v *MessageDAO) GetMessagesByUser(ctx context.Context, receiverId string) ([]*dto.Message, error) {
	_, raws, err := v.mongodb.Query(ctx, commons.ChatDatabase, commons.MessageCollection, nil, nil, &commons.FilterData{Item: "receiverId", Value: receiverId})
	if err != nil {
		return nil, err
	}
	messages := []*dto.Message{}
	for _, raw := range raws {
		message := &dto.Message{}
		err = bson.Unmarshal(*raw, message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
