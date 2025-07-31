package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/USA-RedDragon/dinner-elo/internal/config"
	"github.com/USA-RedDragon/dinner-elo/internal/store"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	server        *http.Server
	metricsServer *http.Server
	pprofServer   *http.Server
	stopped       bool
	config        *config.Config
}

const defTimeout = 5 * time.Second

func NewServer(config *config.Config, store store.Store) *Server {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	writeTimeout := defTimeout
	if config.PProf.Enabled {
		writeTimeout = 60 * time.Second
	}

	applyMiddleware(r, config)
	r.Use(func(ctx *gin.Context) {
		ctx.Set("store", store)
		ctx.Set("config", config)
	})

	applyRoutes(r, config)

	var metricsServer *http.Server
	var pprofServer *http.Server

	if config.Metrics.Enabled {
		metricsRouter := gin.New()
		applyMiddleware(metricsRouter, config)

		metricsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))
		metricsServer = &http.Server{
			Addr:              fmt.Sprintf("%s:%d", config.Metrics.Address, config.Metrics.Port),
			ReadHeaderTimeout: defTimeout,
			WriteTimeout:      writeTimeout,
			Handler:           metricsRouter,
		}
	}

	if config.PProf.Enabled {
		pprofRouter := gin.New()
		applyMiddleware(pprofRouter, config)
		pprof.Register(pprofRouter)
		pprofServer = &http.Server{
			Addr:              fmt.Sprintf("%s:%d", config.PProf.Address, config.PProf.Port),
			ReadHeaderTimeout: defTimeout,
			WriteTimeout:      writeTimeout,
			Handler:           pprofRouter,
		}
	}

	return &Server{
		server: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", config.HTTP.Address, config.HTTP.Port),
			ReadHeaderTimeout: defTimeout,
			WriteTimeout:      writeTimeout,
			Handler:           r,
		},
		metricsServer: metricsServer,
		pprofServer:   pprofServer,
		config:        config,
	}
}

func (s *Server) Start() error {
	waitGrp := sync.WaitGroup{}
	if s.server != nil {
		listener, err := net.Listen("tcp", s.server.Addr)
		if err != nil {
			return err
		}
		waitGrp.Add(1)
		go func() {
			defer waitGrp.Done()
			if err := s.server.Serve(listener); err != nil && !s.stopped {
				slog.Error("HTTP server error", "error", err.Error())
			}
		}()
	}
	slog.Info("HTTP server started", "address", s.config.HTTP.Address, "port", s.config.HTTP.Port)

	if s.config.Metrics.Enabled {
		if s.metricsServer != nil {
			metricsListener, err := net.Listen("tcp", s.metricsServer.Addr)
			if err != nil {
				return err
			}
			waitGrp.Add(1)
			go func() {
				defer waitGrp.Done()
				if err := s.metricsServer.Serve(metricsListener); err != nil && !s.stopped {
					slog.Error("Metrics server error", "error", err.Error())
				}
			}()
		}

		slog.Info("Metrics server started", "address", s.config.Metrics.Address, "port", s.config.Metrics.Port)
	}

	if s.config.PProf.Enabled {
		if s.pprofServer != nil {
			pprofListener, err := net.Listen("tcp", s.pprofServer.Addr)
			if err != nil {
				return err
			}
			waitGrp.Add(1)
			go func() {
				defer waitGrp.Done()
				if err := s.pprofServer.Serve(pprofListener); err != nil && !s.stopped {
					slog.Error("PProf server error", "error", err.Error())
				}
			}()
		}

		slog.Info("PProf server started", "address", s.config.PProf.Address, "port", s.config.PProf.Port)
	}

	go func() {
		waitGrp.Wait()
	}()
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.stopped = true

	errGrp := errgroup.Group{}
	if s.server != nil {
		errGrp.Go(func() error {
			return s.server.Shutdown(ctx)
		})
	}
	if s.metricsServer != nil {
		errGrp.Go(func() error {
			return s.metricsServer.Shutdown(ctx)
		})
	}
	if s.pprofServer != nil {
		errGrp.Go(func() error {
			return s.pprofServer.Shutdown(ctx)
		})
	}

	return errGrp.Wait()
}
