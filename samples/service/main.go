//Copyright (c) 2021 Hunan Antvsion Technology Co., Ltd.. All rights reserved.
//版权所有(c)2021湖南蚁景科技有限公司。保留所有权利。
//Author: lizhi
//CreateTime: 2021-11-02 1:44 下午
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	s := NewService()
	v1 := engine.Group("/sflow/v1")
	v1.GET("/processes", s.ListProcess)
	v1.GET("/process/:pid", s.GetProcess)
	v1.POST("/process", s.CreateProcess)
	v1.POST("/process/start/:pid", s.StartProcess)
	v1.DELETE("/process/:pid", s.TerminateProcess)
	v1.GET("/activities/:pid", s.ListActivities)
	v1.GET("/activity/:pid/:aid", s.GetActivity)
	v1.POST("/activity/:pid/:aid", s.FinishActivity)
	v1.GET("/actions/:pid/:aid", s.ListActions)
	v1.GET("/action/:pid/:aid/:xid", s.GetAction)
	v1.POST("/action/:pid/:aid/:xid", s.InvokeAction)

}
