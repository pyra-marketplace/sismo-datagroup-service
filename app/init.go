package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"sismo-datagroup-service/app/api"
	"sismo-datagroup-service/app/cache"
	"sismo-datagroup-service/app/config"
	"sismo-datagroup-service/app/handler"
	"sismo-datagroup-service/db"
	"sismo-datagroup-service/middlewares"
)

type Routes struct {
}

func (app Routes) StartGin() {
	loadEnv()
	cache.InitCache()

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(middlewares.NewRecovery())
	r.Use(middlewares.NewCors([]string{"*"}))
	r.GET("swagger/*any", middlewares.NewSwagger())

	handler.InitHandler()
	publicRoute := r.Group("/api/v1")
	resource, err := db.InitResource()
	if err != nil {
		logrus.Error(err)
	}
	defer resource.Close()

	r.Static("/template/css", "./template/css")
	r.Static("/template/images", "./template/images")
	//r.Static("/template", "./template")

	r.NoRoute(func(context *gin.Context) {
		//context.File("./template/route_not_found.html")
		context.File("./template/index.html")
	})

	api.ApplyRecordAPI(publicRoute, resource)
	api.ApplyMetaAPI(publicRoute, resource)
	//api.ApplyUserAPI(publicRoute, resource)
	r.Run(":" + os.Getenv("PORT"))
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print(err)
	}

	apiKey := os.Getenv("API_KEY")
	fmt.Println("apiKey: ", apiKey)
	config.InitApiKey(apiKey)
}
