package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sitetester/sequence-api/config"
)

func main() {
	envGinMode := config.DotEnvVar("EnvGinMode")
	if envGinMode != "" {
		gin.SetMode(envGinMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	if gin.Mode() == gin.ReleaseMode {
		config.SetupFileLogger()
	}

	db := config.SetupDb("./db/sequences.db")
	engine := config.SetupRouter(db)

	err := engine.Run(":8081")
	if err != nil {
		panic(err)
	}
}
