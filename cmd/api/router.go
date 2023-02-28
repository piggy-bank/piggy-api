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
	a.Router.POST("/users/on-boarding", a.UserSignup)
	a.Router.DELETE("/users/:user_id", a.DeleteUser)
	a.Router.PUT("/users/:user_id/disable", a.DisableUser)
	a.Router.PUT("/users/:user_id/enable", a.EnableUser)
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
