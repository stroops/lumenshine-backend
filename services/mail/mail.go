package main

import (
	"crypto/tls"

	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	context "golang.org/x/net/context"

	"github.com/sirupsen/logrus"
	mail "gopkg.in/mail.v2"
)

//SendMail sends the mails
//one may either specify ToMultiple (if more than one reciver) or To. If ToMultiple specified, that one will be used
func (s *server) SendMail(c context.Context, r *pb.SendMailRequest) (*pb.Empty, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	from := cnf.MailConfig.Sender
	if r.From != "" {
		from = r.From
	}

	if r.ToMultiple == nil {
		go intSendMail(log, r.Base, from, r.To, r.Subject, r.Body, r.IsHtml)
	} else {
		for _, to := range r.ToMultiple {
			go intSendMail(log, r.Base, from, to, r.Subject, r.Body, r.IsHtml)
		}
	}

	return &pb.Empty{}, nil
}

//intSendMail is called as a goroutine
//it sends the mail and stores errors and status in the database
func intSendMail(log *logrus.Entry, b *pb.BaseRequest, from string, to string, subject string, body string, isHTML bool) {
	msg := mail.NewMessage()
	c := context.Background()

	m := &pb.SaveMailRequest{
		Base:           b,
		MailTo:         to,
		MailSubject:    subject,
		MailBody:       body,
		InternalStatus: 1,
		MailFrom:       from,
	}

	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)

	if isHTML {
		msg.SetBody("text/html", body)
	} else {
		msg.SetBody("text/plain", body)
	}

	d := mail.NewDialer(cnf.MailConfig.Host, cnf.MailConfig.Port, cnf.MailConfig.User, cnf.MailConfig.Password)
	if cnf.MailConfig.Insecure {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	d.SSL = cnf.MailConfig.UseSSL

	err := d.DialAndSend(msg)
	if err != nil {
		log.WithError(err).WithField("to", to).Error("Error sending mail")
		m.ExternalStatus = err.Error()
		m.InternalStatus = -1
	}

	_, err = dbClient.SaveMail(c, m)
	if err != nil {
		log.WithError(err).WithField("to", to).Error("Error saving mail to db")
	}
}
