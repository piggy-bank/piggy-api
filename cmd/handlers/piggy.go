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

func GetAllPiggies(db *gorm.DB, ctx *gin.Context) {
	piggies := []entities.Piggy{}

	if err := db.Find(&piggies).Error; err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
	}

	ctx.IndentedJSON(http.StatusOK, piggies)
}

func GetPiggy(db *gorm.DB, ctx *gin.Context) {
	id := ctx.Param("piggy_id")
	piggy := getPiggy(db, id)
	if piggy == nil {
		ctx.IndentedJSON(http.StatusNotFound, "Piggy not found")
		return
	}
	ctx.IndentedJSON(http.StatusOK, piggy)
}

func CreatePiggy(db *gorm.DB, ctx *gin.Context, flowconfig *configuration.FlowConfig, profile string, log *log.Logger) {
	model := entities.Piggy{}
	piggy := entities.Piggy{}
	if err := ctx.BindJSON(&piggy); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	piggyId, err := blockchainservices.CreateBlockchainPiggy(piggy.UserAddress, piggy.Name, piggy.Description, flowconfig, profile, ctx, log)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	piggy.ID = uint(piggyId)
	if err := db.FirstOrCreate(&model, piggy).Error; err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.IndentedJSON(http.StatusCreated, model)
}

func UpdatePiggy(db *gorm.DB, ctx *gin.Context) {
	id := ctx.Param("piggy_id")
	piggy := getPiggy(db, id)
	if piggy == nil {
		ctx.IndentedJSON(http.StatusNotFound, "Piggy not found")
		return
	}

	if err := ctx.BindJSON(&piggy); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Save(&piggy).Error; err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.IndentedJSON(http.StatusOK, piggy)
}

func DeletePiggy(db *gorm.DB, ctx *gin.Context) {
	id := ctx.Param("piggy_id")
	piggy := getPiggy(db, id)
	if piggy == nil {
		ctx.IndentedJSON(http.StatusNotFound, "Piggy not found")
	}
	if err := db.Delete(&piggy).Error; err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
	}
	ctx.IndentedJSON(http.StatusOK, id)
}

func getPiggy(db *gorm.DB, id string) *entities.Piggy {
	piggy := entities.Piggy{}
	if err := db.First(&piggy, id).Error; err != nil {
		return nil
	}
	return &piggy
}
