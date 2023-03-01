package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"firebase.google.com/go/auth"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/manubidegain/piggy-api/cmd/api/configuration"
	"github.com/manubidegain/piggy-api/cmd/entities"
	"github.com/manubidegain/piggy-api/cmd/repository"
	"github.com/manubidegain/piggy-api/firebase"
	"github.com/manubidegain/piggy-api/utils"

	"github.com/jinzhu/gorm"
)

// App has router and db instances
type App struct {
	Router     *gin.Engine
	DB         *gorm.DB
	Uploader   repository.ClientUploader
	AuthClient *auth.Client
	Config     *configuration.Config
}

// App initialize with predefined configuration
func (a *App) Initialize() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	//Calculate and build profile from environment.
	profile := utils.CalculateProfile()
	log.Printf("Instance running in %s scope", profile)
	a.Config = utils.BuildConfig(profile)

	dbURI := getDataBaseURI(a.Config, profile)

	//Open and migrate database connection with GORM
	db, err := gorm.Open(a.Config.DB.Dialect, dbURI)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Could not connect database")
	}
	a.DB = DBMigrate(db)

	// initialize new gin engine (for server)
	a.Router = gin.Default()
	public := a.Router.Group("/public/users")
	useCorsMiddleware(public)
	public.GET("/piggy", a.GetAllPiggies)
	public.GET("/piggy/:piggy_id", a.GetPiggy)
	// configure firebase
	firebaseAuth := firebase.SetupFirebase()
	a.AuthClient = firebaseAuth

	//TO-DO Handle config for test and dev
	a.Uploader.SetupBucket(profile)

	a.Router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	// set db & firebase auth to gin context with a middleware to all incoming request
	if profile != "dev" {
		a.Router.Use(func(c *gin.Context) {
			c.Set("firebaseAuth", firebaseAuth)
		})
		// using the auth middleware to validate api requests
		a.Router.Use(firebase.AuthMiddleware)

	}
	// setting routers

	a.setRouters()

	// start server
	a.Run(fmt.Sprintf(":%s", port))

}

func getDataBaseURI(config *configuration.Config, profile string) string {
	if profile == "dev" {
		return fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8&parseTime=True&loc=UTC",
			config.DB.Username,
			config.DB.Password,
			config.DB.DatabaseName,
		)
	}
	username := os.Getenv(config.DB.Username)
	password := os.Getenv(config.DB.Password)

	return fmt.Sprintf("%s:%s@cloudsql(%s:%s:%s)/%s?charset=utf8&parseTime=True&loc=UTC",
		username,
		password,
		config.DB.Project,
		config.DB.Zone,
		config.DB.InstanceName,
		config.DB.DatabaseName,
	)
}

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&entities.User{}, &entities.Piggy{}, &entities.Donation{})
	db.LogMode(true)
	return db
}

// Set all required routers
func (a *App) setRouters() {
	a.setUserRouters()
	a.setDonationRouters()
	a.setPiggyRouters()
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
