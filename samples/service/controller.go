//Copyright (c) 2021 Hunan Antvsion Technology Co., Ltd.. All rights reserved.
//版权所有(c)2021湖南蚁景科技有限公司。保留所有权利。
//Author: lizhi
//CreateTime: 2021-11-02 1:49 下午
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lizzz49/sflow/instance"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Api interface {
	ListProcess(ctx *gin.Context)
	GetProcess(ctx *gin.Context)
	CreateProcess(ctx *gin.Context)
	StartProcess(ctx *gin.Context)
	TerminateProcess(ctx *gin.Context)
	ListActivities(ctx *gin.Context)
	GetActivity(ctx *gin.Context)
	FinishActivity(ctx *gin.Context)
	ListActions(ctx *gin.Context)
	GetAction(ctx *gin.Context)
	InvokeAction(ctx *gin.Context)
}

func NewService() Api {
	var db *gorm.DB
	dialet := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local",
		"root", "", "localhost", 3306, "sflow")
	db, err := gorm.Open(mysql.Open(dialet), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		AllowGlobalUpdate: false,
	})
	if err != nil {
		panic(err.Error())
	}
	return &service{db: db, pim: instance.NewProcessInstanceManager(db)}
}

type service struct {
	db  *gorm.DB
	pim *instance.ProcessInstanceManager
}

func (s *service) ListProcess(ctx *gin.Context) {

}

func (s *service) GetProcess(ctx *gin.Context) {
	panic("implement me")
}

func (s *service) CreateProcess(ctx *gin.Context) {
	panic("implement me")
}

func (s *service) StartProcess(ctx *gin.Context) {
	panic("implement me")
}

func (s *service) TerminateProcess(ctx *gin.Context) {
	panic("implement me")
}

func (s *service) ListActivities(ctx *gin.Context) {
	panic("implement me")
}

func (s *service) GetActivity(ctx *gin.Context) {
	panic("implement me")
}

func (s *service) FinishActivity(ctx *gin.Context) {
	panic("implement me")
}

func (s *service) ListActions(ctx *gin.Context) {
	panic("implement me")
}

func (s *service) GetAction(ctx *gin.Context) {
	panic("implement me")
}

func (s *service) InvokeAction(ctx *gin.Context) {
	panic("implement me")
}
