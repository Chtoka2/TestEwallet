package main

import (
	"e-wallet/config"
	"e-wallet/internal/http_router"
	"e-wallet/internal/lib/env"
	"e-wallet/internal/logger"
	"e-wallet/internal/storage"
	"log/slog"
	"os"
	"net/http"
)

func main(){
	//TODO: init config
	cfg := config.Config_init()
	//TODO: init .env
	dotenv := env.Env_reader()
	//TODO: init logger
	log := logger.Logger_init(cfg.Env)
	log.Info("Start new logger!")
	
	//TODO: init storage
	s, err := storage.New(dotenv.DbURL)
	if err != nil{
		log.Error("Database error", ErrorWrapper(err))
		os.Exit(1)
	}
	log.Info("Storage was init")
	//TODO: init router
	r := http_router.Init_router(log, s)
	//TODO: init server
	serv := &http.Server{
		Addr: cfg.HttpServer.Address,
		Handler: r,
		ReadTimeout: cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout: cfg.HttpServer.IdleTimeout,
	}
	log.Info("Address of server", slog.String("Address", serv.Addr))
	if err := serv.ListenAndServe(); err != nil{
		log.Info("Server cannot run", ErrorWrapper(err))
	}
	log.Info("Server stopped")
}

func ErrorWrapper(err error) slog.Attr{
	return slog.Attr{
		Key: "Error",
		Value: slog.StringValue(err.Error()),
	}
}