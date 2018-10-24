package main

import (
	"database/sql"
	"errors"

	"github.com/Soneso/lumenshine-backend/pb"

	"github.com/Soneso/lumenshine-backend/services/db/models"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	context "golang.org/x/net/context"
)

//CheckWallet checks whether the given name and fedname are ok
func (s *server) CheckWalletData(ctx context.Context, r *pb.CheckWalletRequest) (*pb.CheckWalletResponse, error) {
	resp := new(pb.CheckWalletResponse)

	if r.PublicKey_0 != "" {
		exists, err := models.UserWallets(qm.Where("user_id=? and public_key_0 ilike ?", r.UserId, r.PublicKey_0)).Exists(db)
		if err != nil {
			return nil, err
		}
		resp.PublicKey_0Ok = !exists
	}

	if r.WalletName != "" {
		exists, err := models.UserWallets(qm.Where("user_id=? and wallet_name ilike ?", r.UserId, r.WalletName)).Exists(db)
		if err != nil {
			return nil, err
		}
		resp.NameOk = !exists
	}

	if r.FriendlyId != "" && r.Domain != "" {
		exists, err := models.UserWallets(qm.Where("friendly_id ilike ? and domain ilike ?", r.FriendlyId, r.Domain)).Exists(db)
		if err != nil {
			return nil, err
		}
		resp.FederationAddressOk = !exists
	}

	return resp, nil
}

func (s *server) AddWallet(ctx context.Context, r *pb.AddWalletRequest) (*pb.IDResponse, error) {
	uExists, err := models.UserProfiles(qm.Where("id=?", r.UserId)).Exists(db)
	if err != nil {
		return nil, err
	}
	if !uExists {
		return nil, errors.New("User does not exist")
	}

	//check wallet does not exists
	exists, err := models.UserWallets(qm.Where("user_id=? and wallet_name ilike ?", r.UserId, r.WalletName)).Exists(db)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("Wallet for user already exists")
	}

	//if fedname specified, check that fedname does not exist for other user
	if r.FriendlyId != "" && r.Domain != "" {
		exists, err := models.UserWallets(qm.Where("user_id<>? and friendly_id ilike ? and domain ilike ?", r.UserId, r.FriendlyId, r.Domain)).Exists(db)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("FederationName already exists for other user")
		}
	}
	maxOrder, err := models.UserWallets(qm.Select(models.UserWalletColumns.OrderNR),
		qm.Where(models.UserWalletColumns.UserID+"=?", r.UserId),
		qm.OrderBy(models.UserWalletColumns.OrderNR+" DESC")).One(db)
	if err != nil {
		return nil, err
	}

	//add the wallet for the user
	w := &models.UserWallet{
		UserID:           int(r.UserId),
		PublicKey0:       r.PublicKey_0,
		WalletName:       r.WalletName,
		FriendlyID:       r.FriendlyId,
		Domain:           r.Domain,
		ShowOnHomescreen: r.ShowOnHomescreen,
		UpdatedBy:        r.Base.UpdateBy,
	}
	if maxOrder != nil {
		w.OrderNR = maxOrder.OrderNR + 1
	}

	err = w.Insert(db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return &pb.IDResponse{Id: int64(w.ID)}, nil
}

func (s *server) RemoveWallet(ctx context.Context, r *pb.RemoveWalletRequest) (*pb.Empty, error) {
	wallet, err := models.UserWallets(qm.Where("id=? and user_id=?", r.Id, r.UserId),
		qm.Load(models.UserWalletRels.WalletPaymentTemplates)).One(db)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	if wallet.R != nil && wallet.R.WalletPaymentTemplates != nil {
		_, err = wallet.R.WalletPaymentTemplates.DeleteAll(tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	_, err = wallet.Delete(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	updatesql := "UPDATE user_wallet SET order_nr = order_nr - 1 WHERE user_id=$1 AND order_nr>$2"
	_, err = queries.Raw(updatesql, r.UserId, wallet.OrderNR).Exec(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &pb.Empty{}, nil
}

func (s *server) WalletChangeOrder(ctx context.Context, r *pb.WalletChangeOrderRequest) (*pb.Empty, error) {
	orderTo := int(r.OrderNr)
	if r.OrderNr < 0 {
		orderTo = 0
	}
	wallet, err := models.UserWallets(qm.Where(models.UserWalletColumns.PublicKey0+"=? and "+models.UserContactColumns.UserID+"=?", r.PublicKey_0, r.UserId)).One(db)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("Wallet not found")
	}
	if wallet.OrderNR == orderTo {
		return &pb.Empty{}, nil
	}
	max, err := models.UserWallets(qm.Select(models.UserWalletColumns.ID, models.UserWalletColumns.OrderNR),
		qm.Where(models.UserWalletColumns.UserID+"=?", r.UserId),
		qm.OrderBy(models.UserWalletColumns.OrderNR+" DESC")).One(db)
	if err != nil {
		return nil, err
	}
	if max.OrderNR == 0 {
		return &pb.Empty{}, nil
	}
	if orderTo > max.OrderNR {
		orderTo = max.OrderNR
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	orderFrom := wallet.OrderNR

	updatesql := "UPDATE user_wallet SET order_nr = order_nr + 1 WHERE user_id=$1 AND order_nr>=$2 AND order_nr<$3"
	if orderTo > orderFrom {
		updatesql = "UPDATE user_wallet SET order_nr = order_nr - 1 WHERE user_id=$1 AND order_nr<=$2 AND order_nr>$3"
	}
	_, err = queries.Raw(updatesql, r.UserId, orderTo, orderFrom).Exec(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	wallet.OrderNR = orderTo
	_, err = wallet.Update(tx, boil.Whitelist(models.UserWalletColumns.OrderNR,
		models.UserWalletColumns.UpdatedBy,
		models.UserWalletColumns.UpdatedAt))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &pb.Empty{}, nil
}

func (s *server) WalletChangeName(ctx context.Context, r *pb.WalletChangeNameRequest) (*pb.Empty, error) {
	w, err := models.UserWallets(qm.Where("id=? and user_id=?", r.Id, r.UserId)).One(db)
	if err != nil {
		return nil, err
	}

	w.WalletName = r.Name
	w.UpdatedBy = r.Base.UpdateBy

	_, err = w.Update(db, boil.Whitelist(
		models.UserWalletColumns.WalletName,
		models.UserWalletColumns.UpdatedAt,
		models.UserWalletColumns.UpdatedBy,
	))
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) WalletChangeFederationAddress(ctx context.Context, r *pb.WalletChangeFederationAddressRequest) (*pb.Empty, error) {
	w, err := models.UserWallets(qm.Where("id=? and user_id=?", r.Id, r.UserId)).One(db)
	if err != nil {
		return nil, err
	}

	w.FriendlyID = r.FriendlyId
	w.Domain = r.Domain
	w.UpdatedBy = r.Base.UpdateBy

	_, err = w.Update(db, boil.Whitelist(
		models.UserWalletColumns.FriendlyID,
		models.UserWalletColumns.Domain,
		models.UserWalletColumns.UpdatedAt,
		models.UserWalletColumns.UpdatedBy,
	))
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) GetUserWallets(ctx context.Context, r *pb.GetWalletsRequest) (*pb.GetWalletsResponse, error) {
	wallets, err := models.UserWallets(qm.Where("user_id=?", r.UserId), qm.OrderBy(models.UserWalletColumns.OrderNR)).All(db)
	if err != nil {
		return nil, err
	}

	ret := new(pb.GetWalletsResponse)
	for _, w := range wallets {
		ret.Wallets = append(ret.Wallets, &pb.Wallet{
			Id:               int64(w.ID),
			UserId:           int64(w.UserID),
			PublicKey_0:      w.PublicKey0,
			WalletName:       w.WalletName,
			ShowOnHomescreen: w.ShowOnHomescreen,
			FriendlyId:       w.FriendlyID,
			Domain:           w.Domain,
		})
	}

	return ret, nil
}

func (s *server) WalletIsLast(ctx context.Context, r *pb.WalletIsLastRequest) (*pb.BoolResponse, error) {
	wallets, err := models.UserWallets(qm.Where("user_id=?", r.UserId)).All(db)
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
	wallet, err := models.UserWallets(qm.Where("id=? and user_id=?", r.Id, r.UserId)).One(db)
	if err != nil {
		return nil, err
	}
	wallet.ShowOnHomescreen = r.Visible
	wallet.UpdatedBy = r.Base.UpdateBy
	_, err = wallet.Update(db, boil.Whitelist(
		models.UserWalletColumns.ShowOnHomescreen,
		models.UserWalletColumns.UpdatedAt,
		models.UserWalletColumns.UpdatedBy,
	))

	return &pb.Empty{}, err
}

func (s *server) AddPaymentTemplate(ctx context.Context, r *pb.AddPaymentTemplateRequest) (*pb.IDResponse, error) {
	w, err := models.UserWallets(qm.Select(models.UserWalletColumns.UserID), qm.Where("id=?", r.WalletId)).One(db)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, errors.New("Wallet does not exist")
	}
	if w.UserID != int(r.UserId) {
		return nil, errors.New("Wallet does not bellong to the user")
	}

	//add the template for the wallet
	t := &models.PaymentTemplate{
		WalletID:                int(r.WalletId),
		RecepientStellarAddress: r.RecipientStellarAddress,
		RecepientPK:             r.RecipientPublickey,
		AssetCode:               r.AssetCode,
		IssuerPK:                r.IssuerPublickey,
		Amount:                  r.Amount,
		MemoType:                r.MemoType.String(),
		Memo:                    r.Memo,
		UpdatedBy:               r.Base.UpdateBy,
	}

	err = t.Insert(db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return &pb.IDResponse{Id: int64(t.ID)}, nil
}

func (s *server) GetPaymentTemplates(ctx context.Context, r *pb.GetTemplatesRequest) (*pb.GetTemplatesResponse, error) {
	w, err := models.UserWallets(qm.Select(models.UserWalletColumns.UserID), qm.Where("id=?", r.WalletId)).One(db)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, errors.New("Wallet does not exist")
	}
	if w.UserID != int(r.UserId) {
		return nil, errors.New("Wallet does not bellong to the user")
	}

	templates, err := models.PaymentTemplates(qm.Where("wallet_id=?", r.WalletId)).All(db)
	if err != nil {
		return nil, err
	}

	ret := new(pb.GetTemplatesResponse)
	for _, t := range templates {
		template := pb.PaymentTemplate{
			Id:                      int64(t.ID),
			WalletId:                int64(t.WalletID),
			RecipientStellarAddress: t.RecepientStellarAddress,
			RecipientPublickey:      t.RecepientPK,
			AssetCode:               t.AssetCode,
			Amount:                  t.Amount,
			Memo:                    t.Memo,
		}
		if t.MemoType != "" {
			template.MemoType = pb.MemoType(pb.MemoType_value[t.MemoType])
		}
		ret.Templates = append(ret.Templates, &template)
	}

	return ret, nil
}

func (s *server) RemovePaymentTemplate(ctx context.Context, r *pb.RemovePaymentTemplateRequest) (*pb.Empty, error) {
	t, err := models.PaymentTemplates(qm.Where("id=?", r.Id), qm.Load(models.PaymentTemplateRels.Wallet)).One(db)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.New("Template does not exist")
	}
	if t.R.Wallet.UserID != int(r.UserId) {
		return nil, errors.New("Wallet does not bellong to the user")
	}
	_, err = t.Delete(db)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
