package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/manubidegain/piggy-api/cmd/entities"
)

func GetAllDonations(db *gorm.DB, ctx *gin.Context) {
	donations := []entities.Donation{}

	if err := db.Preload("Piggy").Find(&donations).Error; err != nil {
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

func CreateDonation(db *gorm.DB, ctx *gin.Context) {
	model := entities.Donation{}
	donation := entities.Donation{}
	if err := ctx.BindJSON(&donation); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
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
