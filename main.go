package main

import (
	appConfigs "github.com/DeltaDemand/configs-server/configs"
	_ "github.com/DeltaDemand/configs-server/docs"
	"github.com/DeltaDemand/configs-server/global"
	"github.com/DeltaDemand/configs-server/internal/inputArgs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Agent配置修改系统
// @version 1.0
// @description 通过更新etcd上的Agent配置修改管理对应Agent
func main() {
	//加载configs/config.json下配置
	configs := appConfigs.LoadingConfigs()
	//检测用户输入配置
	inputArgs.Parse(&configs)
	//启动初始化本Server的一些变量
	ConfigSvc := configs.ConfigServer               //本服务
	global.Cli = configs.EtcdConfig.GetEtcdClient() //连接etcd的客户端
	router := gin.Default()

	// swagger接口文档：localhost:8887/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", ConfigSvc.EtcdGetAll)
	router.GET("/struct", global.DisplayGroups)
	router.GET("/group", ConfigSvc.EtcdGetGroup)
	router.PUT("/group", ConfigSvc.EtcdPutGroup)
	router.GET("/agent", ConfigSvc.EtcdGetAgent)
	router.PUT("/agent", ConfigSvc.EtcdPutAgent)
	router.DELETE("/agent", ConfigSvc.EtcdDel)

	router.Run(":" + ConfigSvc.ListenPort)

}
