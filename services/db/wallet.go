package main

import (
	"errors"
	"github.com/Soneso/lumenshine-backend/pb"

	"github.com/Soneso/lumenshine-backend/services/db/models"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/queries/qm"
	context "golang.org/x/net/context"
)

//CheckWallet checks whether the given name and fedname are ok
func (s *server) CheckWalletData(ctx context.Context, r *pb.CheckWalletRequest) (*pb.CheckWalletResponse, error) {
	resp := new(pb.CheckWalletResponse)

	if r.PublicKey_0 != "" {
		exists, err := models.UserWallets(db, qm.Where("user_id=? and public_key_0 ilike ?", r.UserId, r.PublicKey_0)).Exists()
		if err != nil {
			return nil, err
		}
		resp.PublicKey_0Ok = !exists
	}

	if r.WalletName != "" {
		exists, err := models.UserWallets(db, qm.Where("user_id=? and wallet_name ilike ?", r.UserId, r.WalletName)).Exists()
		if err != nil {
			return nil, err
		}
		resp.NameOk = !exists
	}

	if r.FederationAddress != "" {
		exists, err := models.UserWallets(db, qm.Where("federation_address ilike ?", r.FederationAddress)).Exists()
		if err != nil {
			return nil, err
		}
		resp.FederationAddressOk = !exists
	}

	return resp, nil
}

func (s *server) AddWallet(ctx context.Context, r *pb.AddWalletRequest) (*pb.IDResponse, error) {
	uExists, err := models.UserProfiles(db, qm.Where("id=?", r.UserId)).Exists()
	if err != nil {
		return nil, err
	}
	if !uExists {
		return nil, errors.New("User does not exist")
	}

	//check wallet does not exists
	exists, err := models.UserWallets(db, qm.Where("user_id=? and wallet_name ilike ?", r.UserId, r.WalletName)).Exists()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("Wallet for user already exists")
	}

	//if fedname specified, check that fedname does not exist for other user
	if r.FederationAddress != "" {
		exists, err := models.UserWallets(db, qm.Where("user_id<>? and federation_address ilike ?", r.UserId, r.FederationAddress)).Exists()
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("FederationName already exists for other user")
		}
	}

	//add the wallet for the user
	w := &models.UserWallet{
		UserID:            int(r.UserId),
		PublicKey0:        r.PublicKey_0,
		WalletName:        r.WalletName,
		FederationAddress: r.FederationAddress,
		ShowOnHomescreen:  r.ShowOnHomescreen,
		UpdatedBy:         r.Base.UpdateBy,
	}

	err = w.Insert(db)
	if err != nil {
		return nil, err
	}

	return &pb.IDResponse{Id: int64(w.ID)}, nil
}

func (s *server) RemoveWallet(ctx context.Context, r *pb.RemoveWalletRequest) (*pb.Empty, error) {
	w, err := models.UserWallets(db, qm.Where("id=? and user_id=?", r.Id, r.UserId)).One()
	if err != nil {
		return nil, err
	}

	err = w.Delete(db)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) WalletChangeName(ctx context.Context, r *pb.WalletChangeNameRequest) (*pb.Empty, error) {
	w, err := models.UserWallets(db, qm.Where("id=? and user_id=?", r.Id, r.UserId)).One()
	if err != nil {
		return nil, err
	}

	w.WalletName = r.Name
	w.UpdatedBy = r.Base.UpdateBy

	err = w.Update(db, "wallet_name", "updated_at", "updated_by")
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) WalletChangeFederationAddress(ctx context.Context, r *pb.WalletChangeFederationAddressRequest) (*pb.Empty, error) {
	w, err := models.UserWallets(db, qm.Where("id=? and user_id=?", r.Id, r.UserId)).One()
	if err != nil {
		return nil, err
	}

	w.FederationAddress = r.FederationAddress
	w.UpdatedBy = r.Base.UpdateBy

	err = w.Update(db, "federation_address", "updated_at", "updated_by")
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) GetUserWallets(ctx context.Context, r *pb.GetWalletsRequest) (*pb.GetWalletsResponse, error) {
	wallets, err := models.UserWallets(db, qm.Where("user_id=?", r.UserId)).All()
	if err != nil {
		return nil, err
	}

	ret := new(pb.GetWalletsResponse)
	for _, w := range wallets {
		ret.Wallets = append(ret.Wallets, &pb.Wallet{
			Id:                int64(w.ID),
			UserId:            int64(w.UserID),
			PublicKey_0:       w.PublicKey0,
			WalletName:        w.WalletName,
			ShowOnHomescreen:  w.ShowOnHomescreen,
			FederationAddress: w.FederationAddress,
		})
	}

	return ret, nil
}

func (s *server) WalletIsLast(ctx context.Context, r *pb.WalletIsLastRequest) (*pb.BoolResponse, error) {
	wallets, err := models.UserWallets(db, qm.Where("user_id=?", r.UserId)).All()
	if err != nil {
		return nil, err
	}

	ret := new(pb.BoolResponse)

	if len(wallets) == 1 && wallets[0].ID == int(r.Id) {
		ret.Value = true
	}

	return ret, nil
}

func (s *server) WalletSetHomescreen(ctx context.Context, r *pb.WalletSetHomescreenRequest) (*pb.Empty, error) {
	wallet, err := models.UserWallets(db, qm.Where("id=? and user_id=?", r.Id, r.UserId)).One()
	if err != nil {
		return nil, err
	}
	wallet.ShowOnHomescreen = r.Visible
	wallet.UpdatedBy = r.Base.UpdateBy
	err = wallet.Update(db,
		models.UserWalletColumns.ShowOnHomescreen,
		models.UserWalletColumns.UpdatedAt,
		models.UserWalletColumns.UpdatedBy,
	)

	return &pb.Empty{}, err
}
