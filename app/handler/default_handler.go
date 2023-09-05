package handler

import (
	"fmt"
	"sismo-datagroup-service/app/form"
	"sismo-datagroup-service/app/model"
	"sismo-datagroup-service/app/response"
)

var _ Handler = new(DefaultHandler)

var DefaultHandlerName = "Default"

type DefaultHandler struct{}

func (*DefaultHandler) ValidateRecord(recordForm form.RecordForm, groupMeta *model.DataGroupMate) (string, error) {
	return recordForm.Account, response.NewError(500, 50011, fmt.Sprint("ShouldNotAddRecordToThisGroup"))
}

func (*DefaultHandler) HandlerName() string {
	return DefaultHandlerName
}
