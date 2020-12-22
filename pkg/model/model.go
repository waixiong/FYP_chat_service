package model

import (

	//pb "./proto"

	// "github.com/syndtr/goleveldb/leveldb"

	"context"

	pb "getitqec.com/server/chat/pkg/api/v1"
	"getitqec.com/server/chat/pkg/commons"
	"getitqec.com/server/chat/pkg/dao"
)

type ChatModel struct {
	// signdb        *leveldb.DB
	// authdb        *leveldb.DB
	// sessiondb     *leveldb.DB
	// usersessiondb *leveldb.DB
	// db      *dynamodb.DynamoDB
	ChatDAO      dao.IChatDAO
	MessageDAO   dao.IMessageDAO
	ClientStream map[string]pb.ChatService_ConnectServer
}

// InitModel ...
func InitModel(m commons.MongoDB) ChatModelI {
	// dao := &dao.UserDAO{}
	_messageDao := dao.InitMessageDAO(m)
	ctx := context.TODO()
	_messageDao.InitIndex(ctx)
	ctx.Done()
	_chatDao := dao.InitChatDAO(m)
	return &ChatModel{MessageDAO: _messageDao, ChatDAO: _chatDao, ClientStream: map[string]pb.ChatService_ConnectServer{}}
}

// need interface
