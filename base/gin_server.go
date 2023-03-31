package base

import (
	"context"
	"github.com/xiazhe-x/basis"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var GinApplication *gin.Engine

func Gin() *gin.Engine {
	Check(GinApplication)
	return GinApplication
}

type GinServerStarter struct {
	basis.BaseStarter
}

func (i *GinServerStarter) Init(ctx basis.StarterContext) {
	//配置运行环境和测试环境
	if SystemConf.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
		GinApplication = gin.New()
	} else {
		gin.SetMode(gin.DebugMode)
		GinApplication = gin.Default()
	}
}

func (i *GinServerStarter) Start(c basis.StarterContext) {
	port := c.Props().GetDefault("app.port", "8080")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: GinApplication,
	}
	defer func() {
		if err := recover(); err != nil {
			logrus.Error(err)
		}
	}()
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Debug("listen: %s\n", err)
		}
	}()
	logrus.Println("listen:", port)
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("服务器关闭 :", err)
	}
	logrus.Println("服务器正在退出 ")
}
