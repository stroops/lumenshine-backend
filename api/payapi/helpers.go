package main

import (
	mw "github.com/Soneso/lumenshine-backend/api/middleware"
	"github.com/Soneso/lumenshine-backend/pb"
)

//NewBaseRequest returns a new base for all requests
func NewBaseRequest(uc *mw.IcopContext) *pb.BaseRequest {
	return &pb.BaseRequest{RequestId: uc.RequestID, UpdateBy: ServiceName}
}
