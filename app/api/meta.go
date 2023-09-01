package api

import (
	"github.com/gin-gonic/gin"
	"sismo-datagroup-service/app/config"
	"sismo-datagroup-service/app/form"
	"sismo-datagroup-service/app/repository"
	"sismo-datagroup-service/app/response"
	"sismo-datagroup-service/db"
)

func ApplyMetaAPI(app *gin.RouterGroup, resource *db.Resource) {
	metaEntity := repository.NewMetaEntity(resource)
	metaRoute := app.Group("/meta")

	metaRoute.GET("/:group_name", getGroupMetaByName(metaEntity))
	metaRoute.POST("", createGroupMeta(metaEntity))
}

//func getAllMeta(metaEntity repository.IMeta) func(ctx *gin.Context) {
//	return func(ctx *gin.Context) {
//		list, code, err := metaEntity.GetAll()
//		response := map[string]interface{}{
//			"metaList": list,
//			"err":      err2.GetErrorMessage(err),
//		}
//		ctx.JSON(code, response)
//	}
//}

func createGroupMeta(metaEntity repository.IMeta) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if !validateApiKey(ctx) {
			response.WithError(ctx, response.ErrorBadRequest(500000, "InvalidApiKey"))
			return
		}

		metaReq := form.MetaForm{}
		if err := ctx.Bind(&metaReq); err != nil {
			response.WithError(ctx, response.ErrorBadRequest(500000, "BindRequestErr"))
			return
		}

		meta, err := metaEntity.CreateOne(metaReq)
		if err != nil {
			response.WithError(ctx, err)
			return
		}

		response.Success(ctx, meta)
	}
}

func getGroupMetaByName(metaEntity repository.IMeta) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if !validateApiKey(ctx) {
			response.WithError(ctx, response.ErrorBadRequest(500000, "InvalidApiKey"))
			return
		}

		name := ctx.Param("group_name")
		meta, err := metaEntity.GetOneByGroupName(name)
		if err != nil {
			response.WithError(ctx, err)
			return
		}

		response.Success(ctx, meta)
	}
}

func validateApiKey(ctx *gin.Context) bool {
	apiKey := ctx.GetHeader("X-API-KEY")
	return apiKey == config.API_KEY
}
