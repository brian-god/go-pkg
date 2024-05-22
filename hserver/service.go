package hserver

import (
	"context"
	"fmt"
	"github.com/brian-god/go-pkg/configs"
	"github.com/brian-god/go-pkg/hserver/i18n"
	"github.com/brian-god/go-pkg/hserver/log"
	"github.com/brian-god/go-pkg/hserver/middleware/cors"
	"github.com/brian-god/go-pkg/hserver/middleware/ratelimit"
	"github.com/brian-god/go-pkg/token"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/gzip"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	service *Service
	once    sync.Once
)

type Service struct {
	option    *Option
	routers   []Router
	handlers  []app.HandlerFunc
	Tokenizer *token.Token
	config    *configs.Bootstrap
	hertz     *server.Hertz
}

func NewService(config *configs.Bootstrap) *Service {
	once.Do(func() {
		opt := Option{
			RateQPS:         config.Server.RateQPS,
			CryptoKey:       config.Crypto.AppKey,
			TokenIssuer:     config.JWT.Issuer,
			TokenSigningKey: config.JWT.SigningKey,
			ReleaseMode:     configs.Mode != configs.Development,
		}
		port := config.Server.Port
		if port <= 0 {
			port = 8888
		}
		tracePort := config.Server.TracerPort
		if port <= 0 {
			tracePort = 8881
		}
		//Configuration Log
		log.Config(config.Log)
		addr := fmt.Sprintf(":%d", port)
		h := server.Default(server.WithHostPorts(addr), server.WithTracer(prometheus.NewServerTracer(fmt.Sprintf(":%d", tracePort), "/hertz")))
		service = &Service{
			option:    &opt,
			hertz:     h,
			Tokenizer: token.New(opt.TokenIssuer, opt.TokenSigningKey),
			config:    config,
		}
	})
	return service
}

// RegisterRouters 注册路由
func (s *Service) RegisterRouters(routers ...Router) {
	s.routers = append(s.routers, routers...)
}

// Use 使用中间件
func (s *Service) Use(handlers ...app.HandlerFunc) {
	s.handlers = append(s.handlers, handlers...)
}
func (s *Service) GetHertz() *server.Hertz {
	return s.hertz
}

// Run 运行服务
func (s *Service) Run() {
	// Set up cross domain and flow limiting middleware
	s.hertz.Use(cors.Handler())
	s.hertz.Use(ratelimit.WithTimeoutHandler(s.option.RateQPS))
	//Use compression
	s.hertz.Use(gzip.Gzip(gzip.DefaultCompression))
	//internationalization
	s.hertz.Use(i18n.Handler())
	//Register custom processors
	if len(s.handlers) > 0 {
		s.hertz.Use(s.handlers...)
	}
	// Set the routing of each module in sequence
	for _, r := range s.routers {
		r.ConfigRoutes(s.hertz, s.Tokenizer)
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		s.hertz.Spin()
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	hlog.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.hertz.Shutdown(ctx); err != nil {
		hlog.Fatal("Server forced to shutdown:", err)
	}
	hlog.Fatal("Server exiting")
}
