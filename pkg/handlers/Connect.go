package handlers

import (
	pb "getitqec.com/server/chat/pkg/api/v1"
	"getitqec.com/server/chat/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConnectHandler struct {
	Model model.ChatModelI
}

func (s *ConnectHandler) Connect(srv pb.ChatService_ConnectServer) error {
	// userID, err := commons.GetUserID(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	return status.Errorf(codes.Unimplemented, "method Connect not implemented")
	// return s.Model.CreateItem(ctx, PBItem2Item(item))
}
