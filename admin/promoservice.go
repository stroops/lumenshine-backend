package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	context "golang.org/x/net/context"

	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"

	"github.com/Soneso/lumenshine-backend/admin/config"
	"github.com/Soneso/lumenshine-backend/admin/db"
)

//GetPromos returns all active promos
func (s *server) GetPromos(c context.Context, r *pb.Empty) (*pb.GetPromosResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	promos, err := db.GetActivePromos()
	if err != nil {
		//c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading existing promos", cerr.GeneralError))
		log.WithError(err).Error("Error reading existing promos")
		return nil, err
	}

	var promoItems []*pb.Promo
	for _, promo := range promos {
		var buttons []*pb.PromoButton
		err = json.Unmarshal([]byte(promo.Buttons), &buttons)
		if err != nil {
			log.WithError(err).WithField("ID", promo.ID).Error("Error deserializing buttons")
			return nil, err
		}

		item := &pb.Promo{
			Id:      int64(promo.ID),
			Name:    promo.Name,
			Title:   promo.Title,
			Text:    promo.PromoText,
			Type:    promo.PromoType,
			Buttons: buttons,
		}
		image, err := getImageResponse(promo.ID, promo.ImageType)
		if err != nil {
			log.WithError(err).WithField("ID", promo.ID).Error("Could not read image")
			return nil, err
		}
		item.Image = image

		promoItems = append(promoItems, item)
	}

	return &pb.GetPromosResponse{Promos: promoItems}, nil
}

func getImageResponse(id int, imageType string) (*pb.PromoImage, error) {
	fileName := fmt.Sprintf("%d.%v", id, imageType)
	filePath := filepath.Join(config.Cnf.Promo.ImagesPath, fileName)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	str := base64.StdEncoding.EncodeToString(content)
	mimeType := db.MimeTypes[imageType]

	return &pb.PromoImage{
		Content:  str,
		MimeType: mimeType}, nil
}
