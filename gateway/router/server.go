/**
* @file: server.go ==> gateway/router
* @package: router
* @author: jingxiu
* @since: 2022/11/8
* @desc: 网络服务入口
 */

package router

import (
	_ "gateway/docs"

	"context"
	"fmt"
	"gateway/config"
	"gateway/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GinApplication() *gin.Engine {
	//设置gin模式
	//gin.SetMode(global.VP.GetString("RunMode"))
	engine := gin.New()
	// 使用日志打印
	engine.Use(gin.Logger())
	//定义404中间件
	engine.NoRoute(middleware.NoRoute())
	//	 路由挂载
	loadRouter(engine)
	// 开启swagger
	Swagger(engine)
	return engine
}

func Server(engine *gin.Engine) {
	fmt.Println("welcome to JINGXIU-CLI!!!")
	fmt.Printf("%s is starting...\n", config.C.Gateway.ServerName)
	fmt.Printf("version: %s\n", config.C.Gateway.Version)
	server := http.Server{
		Addr:    config.C.Gateway.Listen,
		Handler: engine,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server listen err:%s", err)
		}
		fmt.Printf("listening in 127.0.0.1: %s\n", config.C.Gateway.Listen)
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 在此阻塞
	<-quit
	ctx, channel := context.WithTimeout(context.Background(), 5*time.Second)
	defer channel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown error")
	}
	log.Printf("%s exiting...", config.C.Gateway.ServerName)
}

func Swagger(e *gin.Engine) {
	//e.Static("/html","./public")
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
