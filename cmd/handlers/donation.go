package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/manubidegain/piggy-api/cmd/api/configuration"
	blockchainservices "github.com/manubidegain/piggy-api/cmd/blockchain-services"
	"github.com/manubidegain/piggy-api/cmd/entities"
)

func GetAllUserDonations(db *gorm.DB, ctx *gin.Context) {
	userId := ctx.GetString("userID")
	donations := []entities.Donation{}

	if err := db.Preload("Piggy").Find(&donations, "sender_id = ?", userId).Error; err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
	}

	ctx.IndentedJSON(http.StatusOK, donations)
}

func GetDonation(db *gorm.DB, ctx *gin.Context) {
	id := ctx.Param("donation_id")
	donation := getDonation(db, id)
	if donation == nil {
		ctx.IndentedJSON(http.StatusNotFound, "Donation not found")
		return
	}
	ctx.IndentedJSON(http.StatusOK, donation)
}

func CreateDonation(db *gorm.DB, ctx *gin.Context, flowconfig *configuration.FlowConfig, profile string, log *log.Logger, projectConfig *configuration.ProjectConfig) {
	model := entities.Donation{}
	donation := entities.Donation{}
	if err := ctx.BindJSON(&donation); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	donationId, err := blockchainservices.MintDonation(donation.SenderID, donation.Comment, donation.PiggyID, flowconfig, profile, ctx, log, projectConfig)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	donation.ID = uint(donationId)
	if err := db.FirstOrCreate(&model, donation).Error; err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.IndentedJSON(http.StatusCreated, model)
}

func UpdateDonation(db *gorm.DB, ctx *gin.Context) {
	id := ctx.Param("donation_id")
	donation := getDonation(db, id)
	if donation == nil {
		ctx.IndentedJSON(http.StatusNotFound, "Donation not found")
		return
	}

	if err := ctx.BindJSON(&donation); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Save(&donation).Error; err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.IndentedJSON(http.StatusOK, donation)
}

func DeleteDonation(db *gorm.DB, ctx *gin.Context) {
	id := ctx.Param("donation_id")
	donation := getDonation(db, id)
	if donation == nil {
		ctx.IndentedJSON(http.StatusNotFound, "Donation not found")
	}
	if err := db.Delete(&donation).Error; err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
	}
	ctx.IndentedJSON(http.StatusOK, id)
}

func getDonation(db *gorm.DB, id string) *entities.Donation {
	donation := entities.Donation{}
	if err := db.First(&donation, id).Error; err != nil {
		return nil
	}
	return &donation
}
