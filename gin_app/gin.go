package gin_app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/artisanhe/tools/conf"
	"github.com/artisanhe/tools/env"
)

type GinApp struct {
	Name        string
	IP          string
	Port        int
	SwaggerPath string
	WithCORS    bool
	app         *gin.Engine
}

func (a GinApp) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"Port":     80,
		"WithCORS": false,
	}
}

func (a GinApp) MarshalDefaults(v interface{}) {
	if g, ok := v.(*GinApp); ok {
		if g.Name == "" {
			g.Name = os.Getenv("PROJECT_NAME")
		}

		if g.SwaggerPath == "" {
			g.SwaggerPath = "./swagger.json"
		}

		if g.Port == 0 {
			g.Port = 80
		}
	}
}

func (a *GinApp) Init() {
	if env.IsOnline() {
		gin.SetMode(gin.ReleaseMode)
	}

	a.app = gin.New()

	if a.WithCORS {
		a.app.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "HEAD", "DELETE", "PATCH"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "AppToken", "AccessKey", "Cookie"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	a.app.Use(gin.Recovery(), WithServiceName(a.Name), Logger())
}

type GinEngineFunc func(router *gin.Engine)

func (a *GinApp) Register(ginEngineFunc GinEngineFunc) {
	ginEngineFunc(a.app)
}

func (a *GinApp) Start() {
	a.MarshalDefaults(a)
	a.app.GET("/healthz", func(c *gin.Context) {})

	srv := &http.Server{
		Addr:    a.getAddr(),
		Handler: a.app,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server run failed[%s]\n", err.Error())
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown on failed[%+v]\n", err)
	}

	time.Sleep(2 * time.Second)
	fmt.Printf("Server exiting\n")
}

func (a GinApp) getAddr() string {
	return fmt.Sprintf("%s:%d", a.IP, a.Port)
}
