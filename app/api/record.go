package api

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"sismo-datagroup-service/app/csv"
	"sismo-datagroup-service/app/form"
	"sismo-datagroup-service/app/handler"
	"sismo-datagroup-service/app/repository"
	"sismo-datagroup-service/app/response"
	"sismo-datagroup-service/db"
)

func ApplyRecordAPI(app *gin.RouterGroup, resource *db.Resource) {
	recordEntity := repository.NewRecordEntity(resource)
	recordRoute := app.Group("/record")

	recordRoute.GET("/:group_name", getRecordsByGroupName(recordEntity))
	recordRoute.POST("", addRecord(recordEntity))
	recordRoute.POST("/submitList", uploadCsvFile(recordEntity))
}

//func getAllRecord(recordEntity repository.IRecord) func(ctx *gin.Context) {
//	return func(ctx *gin.Context) {
//		list, code, err := recordEntity.GetAll()
//		response := map[string]interface{}{
//			"metaList": list,
//			"err":      err2.GetErrorMessage(err),
//		}
//		ctx.JSON(code, response)
//	}
//}

func addRecord(recordEntity repository.IRecord) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		recordReq := form.RecordForm{}
		if err := ctx.Bind(&recordReq); err != nil {
			response.WithError(ctx, response.ErrorBadRequest(500000, err.Error()))
			return
		}
		groupMeta, ok := recordEntity.IsDataGroupExist(recordReq.GroupName)
		if !ok {
			response.WithError(ctx, response.ErrorBadRequest(500000, "DataGroupNotExist"))
			return
		}

		hd := handler.HandlerMap[groupMeta.Handler]
		record, err := recordEntity.CreateOne(recordReq, hd.ValidateRecord)
		if err != nil {
			response.WithError(ctx, err)
			return
		}
		response.Success(ctx, record)
	}
}

func getRecordsByGroupName(recordEntity repository.IRecord) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		groupName := ctx.Param("group_name")
		dataGroupList, err := recordEntity.GetDataGroupList(groupName)
		if err != nil {
			response.WithError(ctx, err)
			return
		}
		response.Success(ctx, dataGroupList)
	}
}

func uploadCsvFile(recordEntity repository.IRecord) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		groupName := ctx.Request.FormValue("group_name")
		file, header, err := ctx.Request.FormFile("file")
		if err != nil {
			response.WithError(ctx, response.ErrorBadRequest(500000, "NotFileProvided"))
			return
		}
		defer file.Close()

		// Create a new file in the uploads directory
		path := "./uploads/" + header.Filename
		out, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		// Copy the file content to the new file
		_, err = io.Copy(out, file)
		if err != nil {
			response.WithError(ctx, response.ErrorBadRequest(500000, err.Error()))
			return
		}

		list := csv.ParseCSV(path)
		insertErr := recordEntity.InsertRecordsOneByOne(list, groupName)
		if err != nil {
			response.WithError(ctx, insertErr)
			return
		}

		response.Success(ctx, "File uploaded successfully")
	}
}
