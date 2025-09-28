package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/consumers"
	"github.com/tienhai2808/ecom_go/internal/container"
	"github.com/tienhai2808/ecom_go/internal/initialization"
	"github.com/tienhai2808/ecom_go/internal/router"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	Config     *config.Config
	DB         *gorm.DB
	Redis      *redis.Client
	RabbitConn *amqp091.Connection
	RabbitChan *amqp091.Channel
	Container  *container.Container
	Router     *gin.Engine
}

type Server struct {
	cfg        *config.Config
	httpServer *http.Server
	db         *initialization.DB
	rdb        *redis.Client
	rmq        *initialization.RabbitMQConn
}

func NewServer(cfg *config.Config) (*Server, error) {
	db, err := initialization.InitMySQL(cfg)
	if err != nil {
		return nil, err
	}

	rdb, err := initialization.InitRedis(cfg)
	if err != nil {
		return nil, err
	}

	rmq, err := initialization.InitRabbitMQ(cfg)
	if err != nil {
		return nil, err
	}

	ctn := container.NewContainer(db.Gorm, rdb, cfg, rmq.Chann)

	go consumers.StartSendEmailConsumer(rmq, ctn.AuthModule.SMTPService)

	r := gin.Default()

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return nil, fmt.Errorf("thiết lập Proxy thất bại: %w", err)
	}

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(corsConfig))

	api := r.Group(cfg.App.ApiPrefix)

	router.NewUserRouter(api, cfg, ctn.UserModule.UserRepository, ctn.UserModule.UserHandler)
	router.NewAuthRouter(api, cfg, ctn.UserModule.UserRepository, ctn.AuthModule.AuthHandler)
	router.NewAddressRouter(api, cfg, ctn.UserModule.UserRepository, ctn.AddressModule.AddressHandler)
	router.NewProductRouter(api, cfg, ctn.UserModule.UserRepository, ctn.ProductModule.ProductHandler)
	router.NewImageRouter(api, cfg, ctn.UserModule.UserRepository, ctn.ImageModule.ImageHandler)

	addr := fmt.Sprintf(":%d", cfg.App.Port)

	httpServer := &http.Server{
		Addr:           addr,
		Handler:        r,
		MaxHeaderBytes: 5 * 1024 * 1024,
	}

	return &Server{
		cfg,
		httpServer,
		db,
		rdb,
		rmq,
	}, nil
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	if s.db != nil {
		s.db.Close()
	}

	if s.rdb != nil {
		s.rdb.Close()
	}

	if s.rmq != nil {
		s.rmq.Close()
	}

	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Printf("Shutdown http server thất bại: %v", err)
			return
		}
	}

	log.Println("Dừng server thành công")
}

func (s *Server) GracefulShutdown(ch <- chan error) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-ch:
		log.Printf("Chạy service thất bại: %v", err)
	case <-ctx.Done():
		log.Println("Có tín hiệu dừng server")
	}

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.Shutdown(shutdownCtx)
}