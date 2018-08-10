package main

import (
	"github.com/Soneso/lumenshine-backend/pb"

	"github.com/Soneso/lumenshine-backend/services/db/models"

	_ "github.com/lib/pq"
	context "golang.org/x/net/context"
)

//Save the mail data to db
func (s *server) SaveMail(ctx context.Context, r *pb.SaveMailRequest) (*pb.Empty, error) {
	var m models.Mail
	m.MailFrom = r.MailFrom
	m.MailTo = r.MailTo
	m.MailSubject = r.MailSubject
	m.MailBody = r.MailBody

	m.ExternalStatus = r.ExternalStatus
	m.ExternalStatusID = r.ExternalStatusId
	m.InternalStatus = r.InternalStatus
	m.UpdatedBy = r.Base.UpdateBy

	err := m.Insert(db)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
