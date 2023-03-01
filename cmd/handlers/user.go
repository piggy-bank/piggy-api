package handlers

import (
	"fmt"
	"math/rand"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/manubidegain/piggy-api/cmd/api/configuration"
	"github.com/manubidegain/piggy-api/cmd/entities"

	"github.com/jinzhu/gorm"
)

func GetAllUsers(db *gorm.DB, ctx *gin.Context) {
	value, find := ctx.GetQuery("email")
	if find {
		user, err := findUserByMail(db, value)
		if err != nil {
			ctx.IndentedJSON(http.StatusNotFound, "Mail not found")
			return
		}
		ctx.IndentedJSON(http.StatusOK, user)
	} else {
		users := []entities.User{}
		db.Find(&users)
		ctx.IndentedJSON(http.StatusOK, users)
	}

}

func UserSignup(db *gorm.DB, ctx *gin.Context) {
	model := entities.User{}
	user := entities.User{}
	if err := ctx.BindJSON(&user); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
	}
	id := ctx.GetString("UUID")
	user.ID = id

	if err := db.FirstOrCreate(&model, user).Error; err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
	}
	if model.Status {
		ctx.IndentedJSON(http.StatusUnprocessableEntity, "user already exist")
		return
	}
	ctx.IndentedJSON(http.StatusCreated, model)
}

func GetUser(db *gorm.DB, ctx *gin.Context) {
	id := ctx.Param("user_id")
	user := getUser(db, id)
	if user == nil {
		ctx.IndentedJSON(http.StatusNotFound, "User not found")
		return
	}
	ctx.IndentedJSON(http.StatusOK, user)
}

// TODO
func UpdateUser(db *gorm.DB, ctx *gin.Context) {
	id := ctx.Param("user_id")
	user := getUser(db, id)
	if user == nil {
		ctx.IndentedJSON(http.StatusNotFound, "User not found")
		return
	}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Save(&user).Error; err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.IndentedJSON(http.StatusOK, user)
}

func DeleteUser(db *gorm.DB, ctx *gin.Context) {
	id := ctx.Param("user_id")
	user := getUser(db, id)
	if user == nil {
		ctx.IndentedJSON(http.StatusNotFound, "User not found")
	}
	if err := db.Delete(&user).Error; err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
	}
	ctx.IndentedJSON(http.StatusOK, id)
}

func DisableUser(db *gorm.DB, ctx *gin.Context) {
	id := ctx.Param("user_id")
	user := getUser(db, id)
	if user == nil {
		ctx.IndentedJSON(http.StatusNotFound, "User not found")
	}
	user.Disable()
	if err := db.Save(&user).Error; err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
	}
	ctx.IndentedJSON(http.StatusOK, id)
}

func EnableUser(db *gorm.DB, ctx *gin.Context) {
	id := ctx.Param("user_id")
	user := getUser(db, id)
	if user == nil {
		ctx.IndentedJSON(http.StatusNotFound, "User not found")
	}
	user.Enable()
	if err := db.Save(&user).Error; err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
	}
	ctx.IndentedJSON(http.StatusOK, id)
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func ForgotPassword(db *gorm.DB, client *auth.Client, ctx *gin.Context, config *configuration.Config) {
	forgotRequest := ForgotPasswordRequest{}
	if err := ctx.BindJSON(&forgotRequest); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err := PasswordResetEmail(client, ctx, config, forgotRequest.Email); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.IndentedJSON(http.StatusOK, nil)

}

func getUser(db *gorm.DB, id string) *entities.User {
	user := entities.User{}
	if err := db.First(&user, id).Error; err != nil {
		return nil
	}
	return &user
}

func getUserByToken(db *gorm.DB, id string) *entities.User {
	user := entities.User{ID: id}
	if err := db.First(&user, user).Error; err != nil {
		return nil
	}
	return &user
}

func findUserByModel(db *gorm.DB, user entities.User) (*entities.User, error) {
	model := entities.User{}
	if err := db.Raw("SELECT * FROM `users` WHERE `email` = (?)", user.Email).Scan(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func findUserByMail(db *gorm.DB, email string) (*entities.User, error) {
	model := entities.User{}
	user := entities.User{Email: email}
	if err := db.First(&model, user).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func randomNumberToString() string {
	random := fmt.Sprint(rand.Intn(999999))
	random = addLeftZeros(random)
	return random
}

func addLeftZeros(random string) string {
	if len(random) < 6 {
		random = "0" + random
		return addLeftZeros(random)
	}
	return random
}
