package router

import (
	"encoding/json"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	ginMiddleware "github.com/oapi-codegen/gin-middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	"gorm.io/gorm"

	"go-api-arch-clean-template/adapter/controller/gin/handler"
	"go-api-arch-clean-template/adapter/controller/gin/middleware"
	"go-api-arch-clean-template/adapter/controller/gin/presenter"
	"go-api-arch-clean-template/adapter/gateway"
	"go-api-arch-clean-template/pkg"
	"go-api-arch-clean-template/pkg/logger"
	"go-api-arch-clean-template/usecase"
)

// Swaggerの設定をする
func setupSwagger(router *gin.Engine) (*openapi3.T, error) {
	swagger, err := presenter.GetSwagger()
	if err != nil {
		return nil, err
	}

	env := pkg.GetEnvDefault("APP_ENV", "development")
	if env == "development" {
		swaggerJson, _ := json.Marshal(swagger)
		var SwaggerInfo = &swag.Spec{
			InfoInstanceName: "swagger",
			SwaggerTemplate:  string(swaggerJson),
		}
		swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	return swagger, nil
}

func NewGinRouter(db *gorm.DB, corsAllowOrigins []string) (*gin.Engine, error) {
	router := gin.Default()

	router.Use(middleware.CorsMiddleware(corsAllowOrigins))
	swagger, err := setupSwagger(router)
	if err != nil {
		logger.Warn(err.Error())
		return nil, err
	}

	router.Use(middleware.GinZap())
	router.Use(middleware.RecoveryWithZap())

	// ViewのHTMLの設定です。
	router.LoadHTMLGlob("./adapter/presenter/html/*")
	router.GET("/", handler.Index)

	// Healthチェック用のAPIです。
	router.GET("/health", handler.Health)

	apiGroup := router.Group("/api")
	{
		apiGroup.Use(middleware.TimeoutMiddleware(2 * time.Second))
		v1 := apiGroup.Group("/v1")
		{
			v1.Use(ginMiddleware.OapiRequestValidator(swagger))
			// Album APIを追加します。
			albumRepository := gateway.NewAlbumRepository(db)
			albumUseCase := usecase.NewAlbumUseCase(albumRepository)
			albumHandler := handler.NewAlbumHandler(albumUseCase)
			presenter.RegisterHandlers(v1, albumHandler)
		}
	}
	return router, err
}
