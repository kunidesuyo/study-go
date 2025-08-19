package main

import (
	"awesomeProject/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/timeout"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"

	"awesomeProject/api"
	"awesomeProject/app/controllers"
	"awesomeProject/app/models"
	"awesomeProject/configs"
)

func corsMiddleware(allowOrigins []string) gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = allowOrigins
	return cors.New(config)
}

func timeoutMiddleware(duration time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(duration),
		timeout.WithResponse(func(c *gin.Context) {
			c.JSON(
				http.StatusRequestTimeout,
				api.ErrorResponse{Message: "timeout"},
			)
			c.Abort()
		}),
	)
}

func main() {
	if err := models.SetDatabase(models.InstanceMySQL); err != nil {
		logger.Fatal(err.Error())
	}

	router := gin.Default()

	router.Use(ginzap.Ginzap(logger.ZapLogger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger.ZapLogger, true))
	router.Use(corsMiddleware(configs.Config.APICorsAllowOrigins))

	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}

	if configs.Config.IsDevelopment() {
		swaggerJson, _ := json.Marshal(swagger)
		var SwaggerInfo = &swag.Spec{
			InfoInstanceName: "swagger",
			SwaggerTemplate:  string(swaggerJson),
		}
		swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	router.GET("/health", controllers.Health)

	apiGroup := router.Group("/api")
	{
		apiGroup.Use(timeoutMiddleware(2 * time.Second))
		v1 := apiGroup.Group("/v1")
		{
			v1.Use(middleware.OapiRequestValidator(swagger))
			albumHandler := &controllers.AlbumHandler{}
			api.RegisterHandlers(v1, albumHandler)
		}
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Server Shutdown: %s", err.Error()))
	}
	<-ctx.Done()
	logger.Info("Shutdown")
}
