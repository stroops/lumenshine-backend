package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Soneso/lumenshine-backend/pb"

	"github.com/Soneso/lumenshine-backend/services/db/models"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"golang.org/x/crypto/bcrypt"
	context "golang.org/x/net/context"
)

func getUserByIDOrEmail(r *pb.GetUserByIDOrEmailRequest) (*models.UserProfile, error) {
	if r.Email == "" && r.Id == 0 {
		return nil, errors.New("need email or id")
	}

	q := []qm.QueryMod{}

	if r.Id != 0 {
		q = append(q, qm.Where("id=?", r.Id))
	} else {
		q = append(q, qm.Where(models.UserProfileColumns.Email+"=?", r.Email))
	}

	u, err := models.UserProfiles(q...).One(db)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *server) GetUserByMailtoken(ctx context.Context, r *pb.UserMailTokenRequest) (*pb.UserMailTokenResponse, error) {
	u, err := models.UserProfiles(qm.Where(models.UserProfileColumns.MailConfirmationKey+"=?", r.Token)).One(db)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.UserMailTokenResponse{
				TokenNotFound: true,
			}, nil
		}
		return nil, err
	}

	return &pb.UserMailTokenResponse{
		UserId:                 int64(u.ID),
		MailConfirmationExpiry: u.MailConfirmationExpiryDate.Unix(),
		MailConfirmed:          u.MailConfirmed,
	}, nil
}

func (s *server) GetUserDetails(ctx context.Context, r *pb.GetUserByIDOrEmailRequest) (*pb.UserDetailsResponse, error) {

	u, err := getUserByIDOrEmail(r)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.UserDetailsResponse{UserNotFound: true}, nil
		}
		return nil, err
	}

	ret := &pb.UserDetailsResponse{
		Id:                     int64(u.ID),
		MailConfirmed:          u.MailConfirmed,
		MailConfirmationKey:    u.MailConfirmationKey,
		MailConfirmationExpiry: int64(u.MailConfirmationExpiryDate.Unix()),
		TfaSecret:              u.TfaSecret,
		TfaTempSecret:          u.TfaTempSecret,
		TfaUrl:                 u.TfaURL,
		TfaConfirmed:           u.TfaConfirmed,
		MnemonicConfirmed:      u.MnemonicConfirmed,
		Password:               u.Password,
		Email:                  u.Email,
		MessageCount:           int64(u.MessageCount),
	}

	if r.WithQrData {
		ret.TfaQrcode = u.TfaQrcode
	}

	return ret, nil
}

func (s *server) GetUserProfile(ctx context.Context, r *pb.IDRequest) (*pb.UserProfileResponse, error) {
	u, err := models.UserProfiles(qm.Where(
		models.UserProfileColumns.ID+"=?", r.Id,
	)).One(db)

	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.UserProfileResponse{UserNotFound: true}, nil
		}
		return nil, err
	}

	return &pb.UserProfileResponse{
		Id:            int64(u.ID),
		Email:         u.Email,
		Salutation:    u.Salutation,
		Title:         u.Title,
		Forename:      u.Forename,
		Lastname:      u.Company,
		Company:       u.Company,
		StreetAddress: u.StreetAddress,
		StreetNumber:  u.StreetNumber,
		ZipCode:       u.ZipCode,
		City:          u.City,
		State:         u.State,
		CountryCode:   u.CountryCode,
		Nationality:   u.Nationality,
		MobileNr:      u.MobileNR,
		BirthDay:      int64(u.BirthDay.Unix()),
		BirthPlace:    u.BirthPlace,
		Password:      u.Password,
	}, nil
}

func (s *server) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.IDResponse, error) {
	if r.Email == "" {
		return nil, errors.New("need email")
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	var u models.UserProfile
	u.Email = strings.ToLower(r.Email)
	u.MailConfirmed = false
	u.MailConfirmationKey = r.MailConfirmationKey
	u.MailConfirmationExpiryDate = time.Unix(r.MailConfirmationExpiry, 0)
	u.TfaTempSecret = r.TfaTempSecret
	u.TfaConfirmed = false

	u.Salutation = r.Salutation
	u.Title = r.Title
	u.Forename = r.Forename
	u.Lastname = r.Lastname
	u.Company = r.Company
	u.StreetAddress = r.StreetAddress
	u.StreetNumber = r.StreetNumber
	u.ZipCode = r.ZipCode
	u.City = r.City
	u.State = r.State
	u.CountryCode = r.CountryCode
	u.Nationality = r.Nationality
	u.MobileNR = r.MobileNr
	u.BirthDay = time.Unix(r.BirthDay, 0)
	u.BirthPlace = r.BirthPlace
	u.Password = r.Password
	u.UpdatedBy = r.Base.UpdateBy

	err = u.Insert(tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//also insert the security data
	var ud models.UserSecurity
	ud.UserID = u.ID
	ud.KDFSalt = r.KdfSalt

	ud.MnemonicMasterKey = r.MnemonicMasterKey
	ud.MnemonicMasterIv = r.MnemonicMasterIv
	ud.WordlistMasterKey = r.WordlistMasterKey
	ud.WordlistMasterIv = r.WordlistMasterIv
	ud.Mnemonic = r.Mnemonic
	ud.MnemonicIv = r.MnemonicIv
	ud.Wordlist = r.Wordlist
	ud.WordlistIv = r.WordlistIv

	ud.PublicKey0 = r.PublicKey_0
	ud.PublicKey188 = r.PublicKey_188
	ud.UpdatedBy = r.Base.UpdateBy

	err = ud.Insert(tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//also create a initial wallet for the user
	var w models.UserWallet
	w.UserID = u.ID
	w.WalletName = "Primary Wallet"
	w.PublicKey0 = r.PublicKey_0
	w.ShowOnHomescreen = true
	err = w.Insert(tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &pb.IDResponse{Id: int64(u.ID)}, nil
}

func (s *server) ExistsEmail(ctx context.Context, r *pb.ExistsEmailRequest) (*pb.ExistsEmailResponse, error) {
	if r.Email == "" {
		return nil, errors.New("need email")
	}

	exists, err := models.UserProfiles(qm.Where(models.UserProfileColumns.Email+"=?", strings.ToLower(r.Email))).Exists(db)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.ExistsEmailResponse{Exists: false}, nil
		}

		return nil, err
	}
	return &pb.ExistsEmailResponse{Exists: exists}, nil
}

func (s *server) GetCountryList(ctx context.Context, r *pb.LanguageCodeRequest) (*pb.CountryListResponse, error) {
	if r.LanguageCode == "" {
		return nil, errors.New("need language code")
	}

	countries, err := models.Countries(qm.Where(models.CountryColumns.LangCode+"=?", r.LanguageCode)).All(db)
	if err != nil {
		return nil, err
	}

	ret := new(pb.CountryListResponse)
	for _, c := range countries {
		ret.Countries = append(ret.Countries, &pb.Country{Code: c.LangCode, Name: c.CountryName})
	}

	return ret, nil
}

func (s *server) GetSalutationList(ctx context.Context, r *pb.LanguageCodeRequest) (*pb.SalutationListResponse, error) {
	if r.LanguageCode == "" {
		return nil, errors.New("need language code")
	}

	salutations, err := models.Salutations(qm.Where(models.SalutationColumns.LangCode+"=?", r.LanguageCode)).All(db)
	if err != nil {
		return nil, err
	}

	ret := new(pb.SalutationListResponse)
	for _, s := range salutations {
		ret.Salutation = append(ret.Salutation, s.Salutation)
	}

	return ret, nil
}

func (s *server) SetUserTFAConfirmed(ctx context.Context, r *pb.SetUserTfaConfirmedRequest) (*pb.Empty, error) {
	u, err := models.UserProfiles(qm.Where("id=?", r.UserId)).One(db)
	if err != nil {
		return nil, err
	}

	u.TfaQrcode = r.TfaQrcode
	u.TfaURL = r.TfaUrl
	u.TfaSecret = u.TfaTempSecret
	u.TfaTempSecret = ""
	u.TfaConfirmed = true
	u.UpdatedAt = time.Now()
	u.UpdatedBy = r.Base.UpdateBy

	_, err = u.Update(db, boil.Whitelist(
		models.UserProfileColumns.TfaQrcode,
		models.UserProfileColumns.TfaURL,
		models.UserProfileColumns.TfaSecret,
		models.UserProfileColumns.TfaTempSecret,
		models.UserProfileColumns.TfaConfirmed,
		models.UserProfileColumns.UpdatedAt,
		models.UserProfileColumns.UpdatedBy,
	))
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) SetUserMailConfirmed(ctx context.Context, r *pb.IDRequest) (*pb.Empty, error) {
	u, err := models.UserProfiles(qm.Where("id=?", r.Id)).One(db)
	if err != nil {
		return nil, err
	}

	u.MailConfirmed = true
	u.UpdatedBy = r.Base.UpdateBy
	_, err = u.Update(db, boil.Whitelist(
		models.UserProfileColumns.MailConfirmed,
		models.UserProfileColumns.UpdatedAt,
		models.UserProfileColumns.UpdatedBy,
	))
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) SetUserMnemonicConfirmed(ctx context.Context, r *pb.IDRequest) (*pb.Empty, error) {
	u, err := models.UserProfiles(qm.Where("id=?", r.Id)).One(db)
	if err != nil {
		return nil, err
	}

	u.MnemonicConfirmed = true
	u.UpdatedBy = r.Base.UpdateBy
	_, err = u.Update(db, boil.Whitelist(
		models.UserProfileColumns.MnemonicConfirmed,
		models.UserProfileColumns.UpdatedAt,
		models.UserProfileColumns.UpdatedBy,
	))
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) GetUserSecurities(ctx context.Context, r *pb.IDRequest) (*pb.UserSecurityResponse, error) {
	ss, err := models.UserSecurities(qm.Where(
		models.UserSecurityColumns.UserID+"=?", r.Id,
	)).One(db)

	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.UserSecurityResponse{UserNotFound: true}, nil
		}
		return nil, err
	}

	return &pb.UserSecurityResponse{
		Id:      int64(ss.ID),
		UserId:  int64(ss.UserID),
		KdfSalt: ss.KDFSalt,

		MnemonicMasterKey: ss.MnemonicMasterKey,
		MnemonicMasterIv:  ss.MnemonicMasterIv,
		WordlistMasterKey: ss.WordlistMasterKey,
		WordlistMasterIv:  ss.WordlistMasterIv,
		Mnemonic:          ss.Mnemonic,
		MnemonicIv:        ss.MnemonicIv,
		Wordlist:          ss.Wordlist,
		WordlistIv:        ss.WordlistIv,
		PublicKey_0:       ss.PublicKey0,
		PublicKey_188:     ss.PublicKey188,
	}, nil
}

func (s *server) SetUserMailToken(ctx context.Context, r *pb.SetMailTokenRequest) (*pb.Empty, error) {
	u, err := models.UserProfiles(qm.Where(
		models.UserProfileColumns.ID+"=?", r.UserId,
	)).One(db)

	if err != nil {
		return nil, err
	}

	u.MailConfirmationKey = r.MailConfirmationKey
	u.MailConfirmationExpiryDate = time.Unix(r.MailConfirmationExpiry, 0)
	u.UpdatedBy = r.Base.UpdateBy
	_, err = u.Update(db, boil.Whitelist(
		models.UserProfileColumns.MailConfirmationExpiryDate,
		models.UserProfileColumns.MailConfirmationKey,
		models.UserProfileColumns.UpdatedAt,
		models.UserProfileColumns.UpdatedBy,
	))

	return &pb.Empty{}, err
}

func (s *server) SetUserSecurities(ctx context.Context, r *pb.UserSecurityRequest) (*pb.Empty, error) {

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	//need to update the user password, because pub188 could have changed
	pwd, err := bcrypt.GenerateFromPassword([]byte(r.PublicKey_188), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, err := models.UserProfiles(qm.Where("id=?", r.UserId)).One(tx)
	if err != nil {
		return nil, err
	}
	user.Password = string(pwd)
	user.UpdatedBy = r.Base.UpdateBy
	_, err = user.Update(tx, boil.Whitelist(
		models.UserProfileColumns.Password,
		models.UserProfileColumns.UpdatedAt,
		models.UserProfileColumns.UpdatedBy,
	))
	if err != nil {
		return nil, err
	}

	//update the security data
	u, err := models.UserSecurities(qm.Where(
		models.UserSecurityColumns.UserID+"=?", r.UserId,
	)).One(tx)
	if err != nil {
		return nil, err
	}

	u.KDFSalt = r.KdfSalt

	u.MnemonicMasterKey = r.MnemonicMasterKey
	u.MnemonicMasterIv = r.MnemonicMasterIv
	u.WordlistMasterKey = r.WordlistMasterKey
	u.WordlistMasterIv = r.WordlistMasterIv
	u.Mnemonic = r.Mnemonic
	u.MnemonicIv = r.MnemonicIv
	u.Wordlist = r.Wordlist
	u.WordlistIv = r.WordlistIv
	u.PublicKey0 = r.PublicKey_0
	u.PublicKey188 = r.PublicKey_188
	u.UpdatedBy = r.Base.UpdateBy

	_, err = u.Update(tx, boil.Whitelist(
		models.UserSecurityColumns.KDFSalt,
		models.UserSecurityColumns.MnemonicMasterKey,
		models.UserSecurityColumns.MnemonicMasterIv,
		models.UserSecurityColumns.WordlistMasterKey,
		models.UserSecurityColumns.WordlistMasterIv,
		models.UserSecurityColumns.Mnemonic,
		models.UserSecurityColumns.MnemonicIv,
		models.UserSecurityColumns.Wordlist,
		models.UserSecurityColumns.WordlistIv,
		models.UserSecurityColumns.PublicKey0,
		models.UserSecurityColumns.PublicKey188,

		models.UserSecurityColumns.UpdatedAt,
		models.UserSecurityColumns.UpdatedBy,
	))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &pb.Empty{}, nil
}

func (s *server) UpdateUserSecurityPassword(ctx context.Context, r *pb.UserSecurityRequest) (*pb.Empty, error) {
	u, err := models.UserSecurities(qm.Where(
		models.UserSecurityColumns.UserID+"=?", r.UserId,
	)).One(db)
	if err != nil {
		return nil, err
	}

	u.KDFSalt = r.KdfSalt
	u.MnemonicMasterKey = r.MnemonicMasterKey
	u.MnemonicMasterIv = r.MnemonicMasterIv
	u.WordlistMasterKey = r.WordlistMasterKey
	u.WordlistMasterIv = r.WordlistMasterIv
	u.UpdatedBy = r.Base.UpdateBy

	_, err = u.Update(db, boil.Whitelist(
		models.UserSecurityColumns.KDFSalt,
		models.UserSecurityColumns.MnemonicMasterKey,
		models.UserSecurityColumns.MnemonicMasterIv,
		models.UserSecurityColumns.WordlistMasterKey,
		models.UserSecurityColumns.WordlistMasterIv,
		models.UserSecurityColumns.UpdatedAt,
		models.UserSecurityColumns.UpdatedBy,
	))
	return &pb.Empty{}, err
}

func (s *server) SetTempTfaSecret(ctx context.Context, r *pb.SetTempTfaSecretRequest) (*pb.Empty, error) {
	u, err := models.UserProfiles(qm.Where("id=?", r.UserId)).One(db)
	if err != nil {
		return nil, err
	}

	u.TfaTempSecret = r.TfaSecret
	u.UpdatedAt = time.Now()
	u.UpdatedBy = r.Base.UpdateBy

	_, err = u.Update(db, boil.Whitelist(
		models.UserProfileColumns.TfaTempSecret,
		models.UserProfileColumns.UpdatedAt,
		models.UserProfileColumns.UpdatedBy,
	))
	return &pb.Empty{}, err
}

func (s *server) GetUserMessages(ctx context.Context, r *pb.UserMessageListRequest) (*pb.UserMessageListResponse, error) {
	ret := new(pb.UserMessageListResponse)

	var currentCount, archiveCount int64
	currentCount = 0
	archiveCount = 0

	if !r.Archive {
		messages, err := models.UserMessages(qm.Where("user_id=?", r.UserId)).All(db)
		if err != nil {
			return nil, err
		}

		for _, m := range messages {
			ret.MessageListItems = append(ret.MessageListItems, &pb.UserMessageItem{
				Id:          int64(m.ID),
				Title:       m.Title,
				DateCreated: m.CreatedAt.Unix(),
			})
		}

		currentCount = int64(len(messages))
		archiveCount, err = models.UserMessageArchives(qm.Where("user_id=?", r.UserId)).Count(db)
		if err != nil {
			return nil, err
		}
	} else {
		//get archive messages
		messages, err := models.UserMessageArchives(qm.Where("id=?", r.UserId)).All(db)
		if err != nil {
			return nil, err
		}

		for _, m := range messages {
			ret.MessageListItems = append(ret.MessageListItems, &pb.UserMessageItem{
				Id:          int64(m.ID),
				Title:       m.Title,
				DateCreated: m.CreatedAt.Unix(),
			})
		}
		archiveCount = int64(len(messages))
		currentCount, err = models.UserMessages(qm.Where("user_id=?", r.UserId)).Count(db)
		if err != nil {
			return nil, err
		}
	}
	ret.ArchiveCount = archiveCount
	ret.CurrentCount = currentCount
	return ret, nil
}

func (s *server) GetUserMessage(ctx context.Context, r *pb.UserMessageRequest) (*pb.UserMessageItem, error) {
	ret := &pb.UserMessageItem{}

	if !r.Archive {
		m, err := models.UserMessages(qm.Where("id=?", r.MessageId)).One(db)
		if err != nil {
			if err == sql.ErrNoRows {
				return &pb.UserMessageItem{}, nil
			}
			return nil, err
		}
		ret.Id = int64(m.ID)
		ret.UserId = int64(m.UserID)
		ret.Title = m.Title
		ret.Message = m.Message
		ret.DateCreated = m.CreatedAt.Unix()
	} else {
		m, err := models.UserMessageArchives(qm.Where("id=?", r.MessageId)).One(db)
		if err != nil {
			if err == sql.ErrNoRows {
				return &pb.UserMessageItem{}, nil
			}
			return nil, err
		}
		ret.Id = int64(m.ID)
		ret.UserId = int64(m.UserID)
		ret.Title = m.Title
		ret.Message = m.Message
		ret.DateCreated = m.CreatedAt.Unix()
	}

	return ret, nil
}

func (s *server) MoveMessageToArchive(ctx context.Context, r *pb.IDRequest) (*pb.Empty, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	m, err := models.UserMessages(qm.Where("id=?", r.Id)).One(tx)
	if err != nil {
		return nil, err
	}

	ma := new(models.UserMessageArchive)
	ma.UserID = m.UserID
	ma.Title = m.Title
	ma.Message = m.Message
	ma.CreatedAt = m.CreatedAt
	ma.UpdatedAt = time.Now()
	err = ma.Insert(tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = m.Delete(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &pb.Empty{}, nil
}

func (s *server) AddPushToken(ctx context.Context, r *pb.AddPushTokenRequest) (*pb.Empty, error) {
	pushToken, err := models.UserPushtokens(qm.Where(models.UserPushtokenColumns.PushToken+"=?", r.PushToken)).One(db)
	if err != nil {
		return nil, err
	}

	if pushToken != nil {
		if pushToken.UserID != int(r.UserId) {
			pushToken.UserID = int(r.UserId)
			pushToken.UpdatedAt = time.Now()
			pushToken.UpdatedBy = r.Base.UpdateBy
			_, err = pushToken.Update(db, boil.Whitelist(
				models.UserPushtokenColumns.UserID,
				models.UserPushtokenColumns.UpdatedAt,
				models.UserPushtokenColumns.UpdatedBy))

			if err != nil {
				return nil, err
			}
		}
	} else {
		pushToken = &models.UserPushtoken{}
		pushToken.UserID = int(r.UserId)
		pushToken.PushToken = r.PushToken
		pushToken.DeviceType = r.DeviceType.String()
		pushToken.UpdatedBy = r.Base.UpdateBy
		err = pushToken.Insert(db, boil.Infer())

		if err != nil {
			return nil, err
		}
	}

	return &pb.Empty{}, nil
}

func (s *server) UpdatePushToken(ctx context.Context, r *pb.UpdatePushTokenRequest) (*pb.Empty, error) {
	u, err := models.UserProfiles(qm.Where(models.UserProfileColumns.ID+"=?", r.UserId)).One(db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User id:%d does not exist", r.UserId)
		}
		return nil, err
	}
	pushToken, err := models.UserPushtokens(
		qm.Where(models.UserPushtokenColumns.UserID+"=?", u.ID),
		qm.Where(models.UserPushtokenColumns.PushToken+"=?", r.OldPushToken)).
		One(db)
	if err != nil {
		return nil, err
	}
	if pushToken != nil {
		_, err = pushToken.Delete(db)
		if err != nil {
			return nil, err
		}
	}

	pushToken, err = models.UserPushtokens(qm.Where(models.UserPushtokenColumns.PushToken+"=?", r.NewPushToken)).One(db)
	if err != nil {
		return nil, err
	}

	if pushToken != nil {
		if pushToken.UserID != int(r.UserId) {
			pushToken.UserID = int(r.UserId)
			pushToken.UpdatedAt = time.Now()
			pushToken.UpdatedBy = r.Base.UpdateBy
			_, err = pushToken.Update(db, boil.Whitelist(
				models.UserPushtokenColumns.UserID,
				models.UserPushtokenColumns.UpdatedAt,
				models.UserPushtokenColumns.UpdatedBy))

			if err != nil {
				return nil, err
			}
		}
	} else {
		pushToken = &models.UserPushtoken{}
		pushToken.UserID = int(r.UserId)
		pushToken.PushToken = r.NewPushToken
		pushToken.DeviceType = r.DeviceType.String()
		pushToken.UpdatedBy = r.Base.UpdateBy
		err = pushToken.Insert(db, boil.Infer())

		if err != nil {
			return nil, err
		}
	}

	return &pb.Empty{}, nil
}

func (s *server) DeletePushToken(ctx context.Context, r *pb.DeletePushTokenRequest) (*pb.Empty, error) {
	u, err := models.UserProfiles(qm.Where(models.UserProfileColumns.ID+"=?", r.UserId)).One(db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User id:%d does not exist", r.UserId)
		}
		return nil, err
	}
	pushToken, err := models.UserPushtokens(
		qm.Where(models.UserPushtokenColumns.UserID+"=?", u.ID),
		qm.Where(models.UserPushtokenColumns.PushToken+"=?", r.PushToken)).
		One(db)
	if err != nil {
		return nil, err
	}
	if pushToken != nil {
		_, err = pushToken.Delete(db)
		if err != nil {
			return nil, err
		}
	}

	return &pb.Empty{}, nil
}

func (s *server) GetPushTokens(ctx context.Context, r *pb.IDRequest) (*pb.GetPushTokensResponse, error) {
	dbPushTokens, err := models.UserPushtokens(qm.Where(models.UserPushtokenColumns.UserID+"=?", r.Id)).All(db)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	var tokens []*pb.PushToken
	for _, dbToken := range dbPushTokens {
		token := &pb.PushToken{
			PushToken:  dbToken.PushToken,
			DeviceType: pb.DeviceType(pb.DeviceType_value[dbToken.DeviceType])}
		tokens = append(tokens, token)
	}

	return &pb.GetPushTokensResponse{PushTokens: tokens}, nil
}

func (s *server) AddKycDocument(ctx context.Context, r *pb.AddKycDocumentRequest) (*pb.AddKycDocumentResponse, error) {

	document := &models.UserKycDocument{}
	document.UserID = int(r.UserId)
	document.Type = r.DocumentType.String()
	document.Format = r.DocumentFormat.String()
	document.Side = r.DocumentSide.String()
	document.IDCountryCode = r.IdCountryCode
	document.IDIssueDate = time.Unix(r.IdIssueDate, 0)
	document.IDExpirationDate = time.Unix(r.IdExpirationDate, 0)
	document.IDNumber = r.IdNumber
	document.UpdatedBy = r.Base.UpdateBy
	err := document.Insert(db, boil.Infer())

	if err != nil {
		return nil, err
	}

	return &pb.AddKycDocumentResponse{DocumentId: int64(document.ID)}, nil
}
