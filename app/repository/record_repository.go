package repository

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sismo-datagroup-service/app/cache"
	"sismo-datagroup-service/app/form"
	"sismo-datagroup-service/app/handler"
	"sismo-datagroup-service/app/model"
	"sismo-datagroup-service/app/response"
	"sismo-datagroup-service/db"
	"strings"
	"time"
)

var RecordEntity IRecord

type recordEntity struct {
	resource *db.Resource
	repo     *mongo.Collection
}

type (
	IRecord interface {
		GetAll(groupName string) ([]model.DataGroupRecord, response.Error)
		GetOneByAccount(recordForm form.RecordForm) (*model.DataGroupRecord, response.Error)
		CreateOne(recordForm form.RecordForm, groupMeta *model.DataGroupMate, recordHandler handler.HandlerFunc) (*model.DataGroupRecord, response.Error)
		GetDataGroupList(groupName string) (map[string]string, response.Error)
		IsDataGroupExist(groupName string) (*model.DataGroupMate, bool)
		InsertRecords(records []model.DataGroupRecord, groupName string) response.Error
		InsertRecordsOneByOne(records []model.DataGroupRecord, groupName string) response.Error
	}
)

func NewRecordEntity(resource *db.Resource) IRecord {
	userRepo := resource.DB.Collection("record")
	RecordEntity = &recordEntity{resource: resource, repo: userRepo}
	return RecordEntity
}

func (entity *recordEntity) GetAll(groupName string) ([]model.DataGroupRecord, response.Error) {
	accountList := []model.DataGroupRecord{}

	_, ok := entity.IsDataGroupExist(groupName)
	if !ok {
		return nil, response.ErrorBadRequest(500000, "DataGroupNotExist")
	}

	entity.repo = entity.resource.DB.Collection(groupName)

	ctx, cancel := initContext()
	defer cancel()
	cursor, err := entity.repo.Find(ctx, bson.M{})

	if err != nil {
		logrus.Print(err)
		return []model.DataGroupRecord{}, response.ErrorBadRequest(500000, err.Error())
	}

	for cursor.Next(ctx) {
		var meta model.DataGroupRecord
		err = cursor.Decode(&meta)
		if err != nil {
			logrus.Print(err)
		}
		accountList = append(accountList, meta)
	}
	return accountList, nil
}

func (entity *recordEntity) GetDataGroupList(groupName string) (map[string]string, response.Error) {
	groupMembers := make(map[string]string)
	groupMeta, ok := entity.IsDataGroupExist(groupName)
	if !ok {
		return nil, response.ErrorBadRequest(500000, "DataGroupNotExist")
	}

	cacheGroupMembers, ok := cache.MemCache.GetGroupMembers(groupName)
	if !ok {
		entity.repo = entity.resource.DB.Collection(groupName)

		ctx, cancel := initContext()
		defer cancel()
		cursor, err := entity.repo.Find(ctx, bson.M{})

		if err != nil {
			logrus.Print(err)
			return groupMembers, response.ErrorBadRequest(500000, err.Error())
		}

		for cursor.Next(ctx) {
			var meta model.DataGroupRecord
			err = cursor.Decode(&meta)
			if err != nil {
				logrus.Print(err)
			}
			groupMembers[meta.Account] = meta.Value
		}

		expiredAt := calExpireTime(groupMeta)
		cache.MemCache.CacheGroupMembers(groupName, groupMembers, expiredAt)
		return groupMembers, nil
	}

	groupMembers, ok = cacheGroupMembers.(map[string]string)
	if !ok {
		return groupMembers, nil
	}
	return groupMembers, nil

}

func (entity *recordEntity) GetOneByAccount(recordForm form.RecordForm) (*model.DataGroupRecord, response.Error) {
	_, ok := entity.IsDataGroupExist(recordForm.GroupName)
	if !ok {
		return nil, response.ErrorBadRequest(500000, "DataGroupNotExist")
	}

	var record *model.DataGroupRecord
	cacheRecord, ok := cache.MemCache.GetCachedRecordByGroupName(recordForm.GroupName, recordForm.Account)
	if !ok {
		ctx, cancel := initContext()
		defer cancel()
		entity.repo = entity.resource.DB.Collection(recordForm.GroupName)

		err :=
			entity.repo.FindOne(ctx, bson.M{"account": recordForm.Account}).Decode(&record)

		if err != nil {
			logrus.Print(err)
			return nil, response.ErrorBadRequest(500000, err.Error())
		}
		cache.MemCache.CacheRecord(recordForm.GroupName, recordForm.Account, record)
		return record, nil
	}

	record, ok = cacheRecord.(*model.DataGroupRecord)
	if !ok {
		return nil, response.ErrorBadRequest(500000, "ParseCacheErr")
	}

	return record, nil
}

func (entity *recordEntity) CreateOne(recordForm form.RecordForm, groupMeta *model.DataGroupMate, recordHandler handler.HandlerFunc) (*model.DataGroupRecord, response.Error) {
	// validate and return account
	account, err := recordHandler(recordForm, groupMeta)
	if err != nil {
		logrus.Print(err.Error())
		return nil, response.ErrorBadRequest(500000, err.Error())
	}
	if account == "" {
		logrus.Print("EmptyAccount")
		return nil, response.ErrorBadRequest(500000, "EmptyAccount")
	}

	recordForm.Account = account

	groupMeta, ok := entity.IsDataGroupExist(recordForm.GroupName)
	if !ok {
		return nil, response.ErrorBadRequest(500000, "DataGroup not exist")
	}

	ctx, cancel := initContext()
	defer cancel()

	logrus.Print("GroupName:", recordForm.GroupName)
	entity.repo = entity.resource.DB.Collection(recordForm.GroupName)

	expiredAt := calExpireTime(groupMeta)

	groupRecord := model.DataGroupRecord{
		Id:        primitive.NewObjectID(),
		Account:   recordForm.Account,
		Value:     "1",
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	found, _ := entity.GetOneByAccount(recordForm)
	if found != nil {
		return nil, response.ErrorBadRequest(500000, "UsernameIsTaken")
	}
	_, insertErr := entity.repo.InsertOne(ctx, groupRecord)

	if insertErr != nil {
		logrus.Print(err)
		return nil, response.ErrorBadRequest(500000, "UsernameIsTaken")
	}

	return &groupRecord, nil
}

func (entity *recordEntity) InsertRecords(records []model.DataGroupRecord, groupName string) response.Error {
	ctx, cancel := initContext()
	defer cancel()

	_, ok := entity.IsDataGroupExist(groupName)
	if !ok {
		return response.ErrorBadRequest(500000, "DataGroupNotExist")
	}

	// validate and return account
	entity.repo = entity.resource.DB.Collection(groupName)

	interfaces := make([]interface{}, len(records))
	for i, v := range records {
		interfaces[i] = v
	}

	_, err := entity.repo.InsertMany(ctx, interfaces)
	if err != nil {
		return response.ErrorBadRequest(500000, err.Error())

	}
	return nil
}

func (entity *recordEntity) InsertRecordsOneByOne(records []model.DataGroupRecord, groupName string) response.Error {
	ctx, cancel := initContext()
	defer cancel()

	groupMeta, ok := entity.IsDataGroupExist(groupName)
	if !ok {
		return response.ErrorBadRequest(500000, "DataGroupNotExist")
	}

	expiredAt := calExpireTime(groupMeta)

	entity.repo = entity.resource.DB.Collection(groupName)
	for _, record := range records {
		found, _ := entity.GetOneByAccount(form.RecordForm{Account: record.Account, GroupName: groupName})
		if found != nil {
			continue
		}

		record.ExpiredAt = expiredAt

		_, err := entity.repo.InsertOne(ctx, record)
		if err != nil {
			logrus.Print(err)
			continue
		}
		cache.MemCache.CacheRecord(groupName, record.Account, &record)
	}

	return nil
}

func (entity *recordEntity) IsDataGroupExist(groupName string) (*model.DataGroupMate, bool) {
	meta, err := MetaEntity.GetOneByGroupName(groupName)
	if err != nil {
		return nil, false
	}
	//cache.MemCache.Cached()
	return meta, true
}

func calExpireTime(groupMeta *model.DataGroupMate) time.Time {
	expiredAt := groupMeta.StartAt
	frequent := strings.ToLower(groupMeta.GenerateFrequency)
	switch {
	case frequent == "once":
		expiredAt = expiredAt.Add(24 * 365 * 10 * time.Hour)
	case frequent == "daily":
		for {
			expiredAt = expiredAt.Add(24 * time.Hour)
			if expiredAt.After(time.Now()) {
				break
			}
		}
	case frequent == "weekly":
		for {
			expiredAt = expiredAt.Add(7 * 24 * time.Hour)
			if expiredAt.After(time.Now()) {
				break
			}
		}
	default:
		expiredAt = expiredAt.Add(24 * 365 * 10 * time.Hour)
	}
	return expiredAt
}
