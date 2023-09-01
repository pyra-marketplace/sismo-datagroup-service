package handler

import (
	"fmt"
	"sismo-datagroup-service/app/form"
)

type Handler interface {
	ValidateRecord(record form.RecordForm) (string, error)
	HandlerName() string
}

var enabledHandlers = []Handler{
	&DefaultHandler{},
	&TwitterFollowerHandler{},
}

var HandlerMap = make(map[string]Handler)

func InitHandler() {
	for _, value := range enabledHandlers {
		fmt.Println("add handler ", value.HandlerName())
		HandlerMap[value.HandlerName()] = value
	}
}

func GetHandlerName(handlerName string) string {
	if handlerName != TwitterFollowerHandlerName {
		return DefaultHandlerName
	}
	return TwitterFollowerHandlerName
}
