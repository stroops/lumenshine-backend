package main

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Soneso/lumenshine-backend/pb"

	"github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/volatiletech/sqlboiler/queries/qm"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/boil"

	context "golang.org/x/net/context"
)

func (s *server) AddUserContact(ctx context.Context, r *pb.AddUserContactRequest) (*pb.IDResponse, error) {
	contact := &models.UserContact{}
	contact.UserID = int(r.UserId)
	contact.ContactName = r.ContactName
	contact.StellarAddress = r.StellarAddress
	contact.PublicKey = r.PublicKey
	contact.UpdatedBy = r.Base.UpdateBy
	err := contact.Insert(db, boil.Infer())

	if err != nil {
		return nil, err
	}

	return &pb.IDResponse{Id: int64(contact.ID)}, nil
}

func (s *server) UpdateUserContact(ctx context.Context, r *pb.UpdateUserContactRequest) (*pb.Empty, error) {
	contact, err := models.UserContacts(qm.Where(models.UserContactColumns.ID+"=?", r.Id)).One(db)
	if err != nil {
		return nil, err
	}
	if contact == nil {
		return nil, errors.New("user contact not found")
	}

	contact.ContactName = r.ContactName
	contact.StellarAddress = r.StellarAddress
	contact.PublicKey = r.PublicKey
	contact.UpdatedAt = time.Now()
	contact.UpdatedBy = r.Base.UpdateBy
	_, err = contact.Update(db, boil.Whitelist(
		models.UserContactColumns.ContactName,
		models.UserContactColumns.StellarAddress,
		models.UserContactColumns.PublicKey,
		models.UserContactColumns.UpdatedAt,
		models.UserPushtokenColumns.UpdatedBy))

	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) DeleteUserContact(ctx context.Context, r *pb.IDRequest) (*pb.Empty, error) {
	contact, err := models.UserContacts(qm.Where(models.UserContactColumns.ID+"=?", r.Id)).One(db)
	if err != nil {
		return nil, err
	}
	if contact == nil {
		return nil, errors.New("user contact not found")
	}

	_, err = contact.Delete(db)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) GetUserContacts(ctx context.Context, r *pb.IDRequest) (*pb.GetUserContactsResponse, error) {
	dbContacts, err := models.UserContacts(qm.Where(models.UserContactColumns.UserID+"=?", r.Id),
		qm.OrderBy(models.UserContactColumns.ContactName)).All(db)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	var contacts []*pb.UserContact
	for _, dbContact := range dbContacts {
		contact := &pb.UserContact{
			Id:             int64(dbContact.ID),
			ContactName:    dbContact.ContactName,
			StellarAddress: dbContact.StellarAddress,
			PublicKey:      dbContact.PublicKey,
		}
		contacts = append(contacts, contact)
	}

	return &pb.GetUserContactsResponse{UserContacts: contacts}, nil
}
