package main

import (
	"github.com/Soneso/lumenshine-backend/pb"
	context "golang.org/x/net/context"
)

func (s *server) ListenFor(ctx context.Context, r *pb.ListenForRequest) (*pb.Empty, error) {
	//log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	return nil, nil
}

//RemoveListening
//GetData
