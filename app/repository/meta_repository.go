package repository

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"sismo-datagroup-service/app/cache"
	"sismo-datagroup-service/app/form"
	"sismo-datagroup-service/app/handler"
	"sismo-datagroup-service/app/model"
	"sismo-datagroup-service/app/response"
	"sismo-datagroup-service/db"
	"time"
)

var MetaEntity IMeta

type metaEntity struct {
	resource *db.Resource
	repo     *mongo.Collection
}

type IMeta interface {
	GetAll() ([]model.DataGroupMate, response.Error)
	GetOneByGroupName(groupName string) (*model.DataGroupMate, response.Error)
	CreateOne(metaForm form.MetaForm) (*model.DataGroupMate, response.Error)
}

func NewMetaEntity(resource *db.Resource) IMeta {
	userRepo := resource.DB.Collection("meta")
	MetaEntity = &metaEntity{resource: resource, repo: userRepo}
	return MetaEntity
}

func (entity *metaEntity) GetAll() ([]model.DataGroupMate, response.Error) {
	metaList := []model.DataGroupMate{}
	ctx, cancel := initContext()
	defer cancel()
	cursor, err := entity.repo.Find(ctx, bson.M{})

	if err != nil {
		logrus.Print(err)
		return []model.DataGroupMate{}, response.NewError(http.StatusBadRequest, 500000, "QueryMongodbErr")
	}

	for cursor.Next(ctx) {
		var meta model.DataGroupMate
		err = cursor.Decode(&meta)
		if err != nil {
			logrus.Print(err)
		}
		metaList = append(metaList, meta)
	}
	return metaList, nil
}

func (entity *metaEntity) GetOneByGroupName(groupName string) (*model.DataGroupMate, response.Error) {
	cacheMeta, ok := cache.MemCache.GetCachedMeta(groupName)
	var meta *model.DataGroupMate
	if !ok {
		ctx, cancel := initContext()
		defer cancel()
		err :=
			entity.repo.FindOne(ctx, bson.M{"group_name": groupName}).Decode(&meta)
		if err != nil {
			logrus.Print(err)
			return nil, response.NewError(http.StatusBadRequest, 50000, "DataGroupNotExist")
		}
		return meta, nil
	}

	meta, ok = cacheMeta.(*model.DataGroupMate)
	if ok {
		fmt.Println("cached result....")
		return meta, nil
	}
	return nil, response.NewError(http.StatusBadRequest, 50000, fmt.Sprint("DataGroupCacheErr"))
}

func (entity *metaEntity) CreateOne(metaForm form.MetaForm) (*model.DataGroupMate, response.Error) {
	ctx, cancel := initContext()
	defer cancel()

	metaForm.Handler = handler.GetHandlerName(metaForm.Handler)

	if metaForm.Handler == handler.TwitterFollowerHandlerName && metaForm.TwitterConfig.Followers == 0 {
		return nil, response.NewError(http.StatusBadRequest, 50000, "FollowersNotSet")
	}

	groupMeta := model.DataGroupMate{
		Id:                primitive.NewObjectID(),
		GroupName:         metaForm.Name,
		Description:       metaForm.Description,
		Spec:              metaForm.Spec,
		StartAt:           time.Now(),
		UpdatedAt:         time.Now(),
		GenerateFrequency: metaForm.GenerateFrequency,
		Handler:           metaForm.Handler,
		TwitterConfig:     metaForm.TwitterConfig,
	}

	found, _ := entity.GetOneByGroupName(groupMeta.GroupName)
	if found != nil {
		return nil, response.NewError(http.StatusBadRequest, 50000, "GroupNameIsTaken")
	}
	_, err := entity.repo.InsertOne(ctx, groupMeta)
	if err != nil {
		logrus.Print(err)
		return nil, response.NewError(http.StatusBadRequest, 50000, "InsertRecordErr")
	}

	cache.MemCache.Cached(&groupMeta)
	return &groupMeta, nil
}
