package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/GDH-Project/api/cmd/config"
	"github.com/GDH-Project/api/internal/grpc"
	"github.com/GDH-Project/api/internal/handler"
	m "github.com/GDH-Project/api/internal/middleware"
	"github.com/GDH-Project/api/internal/repository"
	"github.com/GDH-Project/api/internal/resource"
	"github.com/GDH-Project/api/internal/service"
	usecase "github.com/GDH-Project/api/internal/use_case"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	Version = "dev"
)

type Options struct {
	Debug bool `doc:"Enable debug logging" short:"d"`
}

func main() {

	cli := humacli.New(func(hooks humacli.Hooks, opts *Options) {
		log := config.InitLogger(opts.Debug)
		log.Info("GDH Project API SERVER", zap.String("version", Version))
		cfg := config.GetConfig(log)

		if opts.Debug {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

		var r *gin.Engine

		if opts.Debug {
			r = gin.Default()
		} else {
			r = gin.New()
			r.TrustedPlatform = gin.PlatformCloudflare
			r.Use(ginzap.Ginzap(log, time.RFC3339, true))
			r.Use(ginzap.RecoveryWithZap(log, true))
		}

		// cors 설정
		var corsHosts []string
		corsConfig := cors.DefaultConfig()
		if cfg.CorsHostList != "" {
			corsHosts = strings.Split(cfg.CorsHostList, ",")
			log.Info("CORS host list", zap.Any("hostlist", corsHosts))
		}
		corsConfig.AllowOrigins = corsHosts
		corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
		r.Use(cors.New(corsConfig))

		// huma config
		humaConfig := huma.DefaultConfig("GDH-API 서버 입니다.", Version)
		humaConfig.CreateHooks = nil
		humaConfig.SchemasPath = ""
		humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
			"bearer": {
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
			},
		}

		api := humagin.New(r, humaConfig)

		// Dependency
		db := resource.InitDB(cfg.DbUrl, log)
		// gRPC Client 생성
		grpcClientConn := grpc.NewBaseClient(log, cfg)

		userGrpcClient := grpc.NewUserClient(log, grpcClientConn)
		userService := service.NewUserService(log, userGrpcClient)
		userUseCase := usecase.NewUserUseCase(log, userService)

		authGrpcClient := grpc.NewAuthClient(log, grpcClientConn)
		authService := service.NewAuthService(log, authGrpcClient)
		authUseCase := usecase.NewAuthService(log, authService)

		metaRepository := repository.MetaRepository(log, db)
		metaService := service.NewMetaService(log, metaRepository)
		metaUseCase := usecase.NewMetaUseCase(log, metaService)

		_ = repository.MetaRepository(log, db)

		middleware := m.NewMiddleware(api, log, authUseCase)

		// gRPC 미들웨어 적용
		r.Use(middleware.WithGrpcMeta())

		// Register Handler
		handler.RegisterAuthHandler(api, log, authUseCase, userUseCase, middleware)
		handler.RegisterMetaHandler(api, log, metaUseCase)

		server := http.Server{
			Addr:    ":8080",
			Handler: r,
		}
		// 서버 시작시
		hooks.OnStart(func() {
			log.Info("서버를 시작합니다 :8080")
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal("서버를 초기화 하지 못했습니다.", zap.Error(err))
			}
		})
		// 서버 종료시
		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			stopped := make(chan struct{})

			go func() {
				_ = server.Shutdown(ctx)
				close(stopped)
			}()

			select {
			case <-stopped:
				log.Info("api서버가 정상적으로 종료되었습니다.")
			case <-ctx.Done():
				log.Warn("서버가 종료 제한시간에 도달하여 강제 종료 되었습니다.")
			}
		})
	})
	cli.Run()
}
