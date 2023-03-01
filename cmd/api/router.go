package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	handler "github.com/manubidegain/piggy-api/cmd/handlers"
)

func (a *App) setUserRouters() {
	a.Router.GET("/users", a.GetAllUsers)
	a.Router.GET("/users/:user_id", a.GetUser)
	a.Router.PUT("/users/:user_id", a.UpdateUser)
	a.Router.POST("/users", a.UserSignup)
	a.Router.DELETE("/users/:user_id", a.DeleteUser)
	a.Router.PUT("/users/:user_id/disable", a.DisableUser)
	a.Router.PUT("/users/:user_id/enable", a.EnableUser)
}

func (a *App) setPiggyRouters() {
	a.Router.GET("/piggy", a.GetAllPiggies)
	a.Router.GET("/piggy/:piggy_id", a.GetPiggy)
	a.Router.PUT("/piggy/:piggy_id", a.UpdatePiggy)
	a.Router.POST("/piggy", a.CreatePiggy)
	a.Router.DELETE("/piggy/:piggy_id", a.DeletePiggy)
}

func (a *App) setDonationRouters() {
	a.Router.GET("/donation", a.GetAllUserDonations)
	a.Router.GET("/donation/:donation_id", a.GetDonation)
	a.Router.PUT("/donation/:donation_id", a.UpdateDonation)
	a.Router.POST("/donation", a.CreateDonation)
	a.Router.DELETE("/donation/:donation_id", a.DeleteDonation)
}

// User Handlers.
func (a *App) GetAllUsers(ctx *gin.Context) {
	handler.GetAllUsers(a.DB, ctx)
}

func (a *App) UserSignup(ctx *gin.Context) {
	handler.UserSignup(a.DB, ctx)
}

func (a *App) GetUser(ctx *gin.Context) {
	handler.GetUser(a.DB, ctx)
}

func (a *App) UpdateUser(ctx *gin.Context) {
	handler.UpdateUser(a.DB, ctx)
}

func (a *App) DeleteUser(ctx *gin.Context) {
	handler.DeleteUser(a.DB, ctx)
}

func (a *App) DisableUser(ctx *gin.Context) {
	handler.DisableUser(a.DB, ctx)
}

func (a *App) EnableUser(ctx *gin.Context) {
	handler.EnableUser(a.DB, ctx)
}

func (a *App) ForgotPassword(ctx *gin.Context) {
	handler.ForgotPassword(a.DB, a.AuthClient, ctx, a.Config)
}

// Piggy Handlers.
func (a *App) GetAllPiggies(ctx *gin.Context) {
	handler.GetAllPiggies(a.DB, ctx)
}

func (a *App) GetPiggy(ctx *gin.Context) {
	handler.GetPiggy(a.DB, ctx)
}

func (a *App) CreatePiggy(ctx *gin.Context) {
	handler.CreatePiggy(a.DB, ctx)
}

func (a *App) UpdatePiggy(ctx *gin.Context) {
	handler.UpdatePiggy(a.DB, ctx)
}

func (a *App) DeletePiggy(ctx *gin.Context) {
	handler.DeletePiggy(a.DB, ctx)
}

// Donation Handlers.
func (a *App) GetAllUserDonations(ctx *gin.Context) {
	handler.GetAllUserDonations(a.DB, ctx)
}

func (a *App) GetDonation(ctx *gin.Context) {
	handler.GetDonation(a.DB, ctx)
}

func (a *App) CreateDonation(ctx *gin.Context) {
	handler.CreateDonation(a.DB, ctx)
}

func (a *App) UpdateDonation(ctx *gin.Context) {
	handler.UpdateDonation(a.DB, ctx)
}

func (a *App) DeleteDonation(ctx *gin.Context) {
	handler.DeleteDonation(a.DB, ctx)
}

func useCorsMiddleware(public *gin.RouterGroup) {
	public.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))
}

func Ping(c *gin.Context) {

	c.String(http.StatusOK, "pong")
}
