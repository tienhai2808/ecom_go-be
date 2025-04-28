package server

import (
	"backend/internal/admin"
	"backend/internal/auth"
	"backend/internal/common"
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/cache"
	"backend/internal/product"
	"backend/internal/user"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type Application struct {
	AppCtx *common.AppContext
	Router *gin.Engine
}

func NewApplication() *Application {
	appConfig, err := config.LoadAppConfig()
	if err != nil {
		log.Fatalf("❤️ Lỗi khi load cấu hình app: %v", err)
	}

	db, err := database.ConnectToDatabase()
	if err != nil {
		log.Fatalf("❤️ Lỗi kết nối tới database: %v", err)
	}

	redisClient, err := cache.ConnectToRedis()
	if err != nil {
		log.Fatalf("❤️ Lỗi kết nối tới redis: %v", err)
	}

	appCtx := &common.AppContext{
		DB:     db,
		Redis:  redisClient,
		Config: appConfig,
	}

	r := gin.Default()
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("❤️ Could not set trusted proxies: %v", err)
	}

	application := &Application{
		AppCtx: appCtx,
		Router: r,
	}

	application.initRoutes()

	return application
}

func (app *Application) initRoutes() {
	api := app.Router.Group(app.AppCtx.Config.App.ApiPrefix)
	{
		auth.AuthRouter(api, app.AppCtx)
		admin.AdminRouter(api, app.AppCtx)
		user.UserRouter(api, app.AppCtx)
		product.ProductRouter(api, app.AppCtx)
	}
}

func (app *Application) Run() {
	fmt.Println("💚 Kết nối MySQL thành công")
	fmt.Println("💚 Kết nối Redis thành công")
	addr := app.AppCtx.Config.App.Host + ":" + app.AppCtx.Config.App.Port
	if err := app.Router.Run(addr); err != nil {
		log.Fatalf("❤️ Không thể khởi động server: %v", err)
	}
}

func (app *Application) Close() {
	if sqlDB, err := app.AppCtx.DB.DB(); err == nil {
		sqlDB.Close()
	}
	if app.AppCtx.Redis != nil {
		app.AppCtx.Redis.Close()
	}
}
