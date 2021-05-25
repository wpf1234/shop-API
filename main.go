package main

import (
	"testAPI/confInit"
	"testAPI/middleware"

	"testAPI/handle"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/web"
	log "github.com/sirupsen/logrus"
)

func main() {

	service := web.NewService(
		web.Name("test.micro.api.v1.test"),
		web.Version("latest"),
		web.Address(":12345"),
	)

	err := service.Init()
	if err != nil {
		log.Error("服务初始化失败: ", err)
		return
	}
	confInit.Init()

	gin.SetMode(gin.ReleaseMode)

	g := new(handle.Gin)

	router := gin.Default()
	router.Use(middleware.Cors())

	noAuth := router.Group("/v1/test")
	noAuth.POST("/login", g.Login)

	// 需要验证 token 的路由组
	auth := router.Group("/v1/test/auth")
	auth.Use(middleware.JWTAuth())
	{
		auth.GET("/products", g.GetList)
		auth.GET("/products/one", g.GetInfoByID)
		auth.PUT("/products", g.ModifyProd)
		auth.POST("/products", g.AddNewProd)
		auth.DELETE("/products", g.DeleteProd)
	}

	service.Handle("/", router)

	err = service.Run()
	if err != nil {
		log.Error("服务启动失败：", err)
		return
	}
}
