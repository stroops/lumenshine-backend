package main

import (
	"time"

	"github.com/Soneso/lumenshine-backend/pb"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Soneso/lumenshine-backend/services/db/models"

	_ "github.com/lib/pq"
	context "golang.org/x/net/context"
)

func (s *server) GetAllJwtKeys(ctx context.Context, r *pb.Empty) (*pb.KeyListResponse, error) {
	jwts, err := models.JWTKeys().All(db)
	if err != nil {
		return nil, err
	}

	var ret pb.KeyListResponse
	for _, j := range jwts {
		ret.KeyValues = append(ret.KeyValues, &pb.KeyListResponse_KeyValue{
			KeyName:   j.KeyName,
			KeyValue1: j.KeyValue1,
			KeyValue2: j.KeyValue2,
			Valid1To:  j.Valid1To.Unix(),
			Valid2To:  j.Valid2To.Unix(),
		})
	}

	return &ret, nil
}

func (s *server) SetJwtKey(ctx context.Context, r *pb.JwtSetKeyRequest) (*pb.Empty, error) {
	jwt, err := models.JWTKeys(qm.Where(
		models.JWTKeyColumns.KeyName+"=?", r.Key,
	)).One(db)
	if err != nil {
		return nil, err
	}

	jwt.KeyValue1 = r.Value1
	jwt.KeyValue2 = r.Value2
	jwt.Valid1To = time.Unix(r.Expiry1, 0)
	jwt.Valid2To = time.Unix(r.Expiry2, 0)
	_, err = jwt.Update(db, boil.Whitelist(
		models.JWTKeyColumns.KeyValue1, models.JWTKeyColumns.KeyValue2,
		models.JWTKeyColumns.Valid1To, models.JWTKeyColumns.Valid2To,
		models.JWTKeyColumns.UpdatedAt),
	)
	return &pb.Empty{}, err
}
