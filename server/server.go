package server

import (
	"app/common"
	"app/config"
	"app/controller"
	"app/database"
	_ "app/docs"
	"app/middleware"
	"app/repository"
	"app/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/version"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const Language = "zh"

type Server struct {
	engine         *gin.Engine
	config         *config.Config
	userController *controller.UserController
	//redis          *redis.Client
	logger *logrus.Logger

	authContoller *controller.AuthController
	//containerController: containerController,
	authMiddleware gin.HandlerFunc

	loggerMiddleware   gin.HandlerFunc
	recoveryMiddleware gin.HandlerFunc
	repository         repository.Repository

	//casbinMiddleware gin.HandlerFunc

	controllers []controller.Controller

	transContoller *controller.TransController

	casbinController *controller.CasbinController

	uploadController *controller.UploadController
}

func New(conf *config.Config) (*Server, error) {

	// rateLimit
	rateLimitMiddleware := middleware.RateLimitMiddleware(&conf.Server.RateLimitsConfigs)

	// mysql
	db, err := database.NewMysql(&conf.DB)
	if err != nil {
		return nil, err
	}

	// redis
	rdb, err := database.NewRedis(&conf.Redis)
	if err != nil {
		return nil, err
	}

	//new initable
	repository := repository.NewRepository(db, rdb)
	if conf.DB.Migrate {
		if err = repository.Migrate(); err != nil {
			return nil, err
		}
	}

	// user
	//userRepository := repository.NewUserRepository(db,rdb)
	//if err := userRepository.Migrate(); err != nil {
	//	return nil, err
	//}

	// user
	userService := service.NewUserService(repository.User())
	userController := controller.NewUserController(userService)

	jwtService := service.NewJWTService()
	authContoller := controller.NewAuthController(userService, jwtService)

	transContoller := controller.NewTransController()
	transContoller.Trans(Language)

	// casbin
	//casbinRepository := repository.NewCasbinRepository()
	//if err := casbinRepository.Migrate(); err != nil {
	//	return nil, err
	//}

	//casbinService := service.NewCasbinService(casbinRepository)
	//casbinController := controller.NewCasbinController(casbinService)

	// upload file
	uploadService := service.NewUploadService(&conf.Storage)
	uploadController := controller.NewUploadController(uploadService)

	//post
	postService := service.NewPostService(repository.Post())
	postController := controller.NewPostController(postService)

	controllers := []controller.Controller{userController, postController}

	//logger
	logs := service.NewLoggerService(&conf.Logger)
	if err := logs.WriteLog(); err != nil {
		return nil, err
	}

	gin.SetMode(conf.Server.ENV)
	e := gin.New()
	e.Use(
		rateLimitMiddleware,
		middleware.CORSMiddleware(),
		middleware.LoggerMiddleWare(),
		middleware.Recovery(),
		//middleware.CasbinMiddleware(),
	)

	//e.LoadHTMLFiles("")

	return &Server{
		engine:           e,
		config:           conf,
		userController:   userController,
		authContoller:    authContoller,
		transContoller:   transContoller,
		uploadController: uploadController,
		authMiddleware:   middleware.AuthMiddleware(jwtService),

		//casbinMiddleware: middleware.CasbinMiddleware(),

		// logger
		loggerMiddleware:   middleware.LoggerMiddleWare(),
		recoveryMiddleware: middleware.Recovery(),
	}, nil
}

func (s *Server) Run() {
	defer s.Close()
	s.Routers()

	addr := fmt.Sprintf("%s:%d", s.config.Server.Address, s.config.Server.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("failed to start server, %v", err)
		}
	}()

	// 平滑关闭进程
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.Server.GracefulShutdownPeriod)*time.Second)
	defer cancel()

	ch := <-sig
	log.Fatalf("Receive signal: %s", ch)

	server.Shutdown(ctx)
}

func (s *Server) Close() {
	if err := s.repository.Close(); err != nil {
		s.logger.Warnf("failed to close repository, %v", err)
	}
}

func (s *Server) Routers() {
	root := s.engine
	//router.GET("/")

	// register non-resource routers
	root.GET("/", common.WrapFunc(s.getRoutes))
	root.GET("/index", controller.Index)
	root.GET("/healthz", common.WrapFunc(s.Ping))
	root.GET("/version", common.WrapFunc(version.Get))
	root.GET("/metrics", gin.WrapH(promhttp.Handler()))
	root.Any("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))

	// swagger doc
	if gin.Mode() != gin.ReleaseMode {
		root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	//active health check
	//router.GET("/health", s.Health)
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//
	//router.POST("/api/auth/token", s.authContoller.Login)
	////router.DELETE("/api/auth/token", s.authContoller.Logout)
	//router.POST("/api/auth/user", s.authContoller.Register)

	api := root.Group("/api/v1")
	controllers := make([]string, 0, len(s.controllers))
	for _, router := range s.controllers {
		router.RegisterRoute(api)
		controllers = append(controllers, router.Name())
	}
	//s.casbinMiddleware
	//s.authMiddleware, s.logMiddleware, s.recoveryMiddleware
	api.Use()
	// user api
	api.GET("/users", s.userController.List)
	api.GET("/user/:id", s.userController.Get)
	api.POST("/user", s.userController.Create)
	api.DELETE("/user/:id", s.userController.Delete)
	api.PUT("/user/:id", s.userController.Update)
	//api.GET("/user/download", s.userController.ExportUserList)

	// upload api
	api.POST("/upload", s.uploadController.Upload)

	//casbin api
	//api.POST("/casbin", s.casbinController.Create)
	//api.POST("/casbin/list", s.casbinController.List)
}

// @Summary Healthz
// @Produce json
// @Tags healthz
// @Success 200 {string}  string    "ok"
// @Router /healthz [get]
func (s *Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
