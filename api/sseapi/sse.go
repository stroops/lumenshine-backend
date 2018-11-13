package main

import (
	"time"

	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
)

type SSEListener struct {
	log *logrus.Entry
}

func NewListenSSE() *SSEListener {
	return &SSEListener{
		log: helpers.GetDefaultLog(ServiceName, "SSEListener"),
	}
}

//Run runs a loop and gathers the latest data to be send to the clients
//should be called as a go routine
func (s *SSEListener) Run() {
	ctx := context.Background()

	baseRequest := &pb.BaseRequest{RequestId: "0", UpdateBy: ServiceName}

	for {
		data, err := sseClient.GetData(ctx, &pb.SSEGetDataRequest{
			Base:          baseRequest,
			SourceReciver: "sse",
			Count:         20,
		})

		if err != nil {
			s.log.WithError(err).Error("Error reading data")
			time.Sleep(5 * time.Second)
		} else {
			if data != nil && data.Data != nil {
				for _, d := range data.Data {
					hub.send <- &WsMessage{
						Account: d.StellarAccount,
						//MessageType: d.OperationType,
						//Message:     []byte(d.OperationData),
					}
				}
			} else {
				//max 10 messages per second
				time.Sleep(1 * time.Second)
			}
		}
	}
}
