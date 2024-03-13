package config

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sitetester/sequence-api/api"
	"github.com/sitetester/sequence-api/api/controller"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
)

// SetupFileLogger https://github.com/gin-gonic/gin#how-to-write-log-file
func SetupFileLogger() {
	// not needed when writing the logs to file
	gin.DisableConsoleColor()
	f, err := os.Create("logs/gin.log")
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(f)
}

func DotEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return os.Getenv(key)
}

// SetupDb TODO: Make db engine connection dynamic based on environment setting
func SetupDb(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&api.Sequence{})
	db.AutoMigrate(&api.SequenceStep{})

	return db
}

const ApiVersion = "/v1"

func SetupRouter(db *gorm.DB) *gin.Engine {
	engine := gin.Default()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	engine.Use(gin.Recovery())

	sequenceController := controller.NewSequenceController(db)
	sequenceStepsController := controller.NewSequenceStepsController(db)

	// WARNING! Currently, there is no authentication/authorization for this API
	// Some kind of token/key must be provided to avoid data loss
	v1 := engine.Group(ApiVersion)
	{
		// http://localhost:8081/api/v1/
		// Or http://127.0.0.1:8081/api/v1/
		v1.GET("/", func(ctx *gin.Context) { ctx.String(200, "It works!") })

		// Sequences
		v1.POST("/sequences", sequenceController.Create)
		v1.PUT("/sequences/:id", sequenceController.Update)
		v1.GET("/sequences/:id", sequenceController.ViewWithSteps)

		// Steps
		v1.POST("/sequence-steps", sequenceStepsController.Create)
		v1.PUT("/sequence-steps/:id", sequenceStepsController.Update)
		v1.DELETE("/sequence-steps/:id", sequenceStepsController.Delete)
		v1.GET("/sequence-steps/:id", sequenceStepsController.View)

	}

	return engine
}
