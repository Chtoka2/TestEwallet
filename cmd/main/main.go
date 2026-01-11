package main

import (
	"e-wallet/config"
	"e-wallet/internal/lib/env"
	"e-wallet/internal/logger"
	"e-wallet/internal/storage"
	"log/slog"
	"os"
)

func main(){
	//TODO: init config
	cfg := config.Config_init()
	//TODO: init .env
	dbPath := env.Env_reader()
	//TODO: init logger
	log := logger.Logger_init(cfg.Env)
	log.Info("Start new logger!")
	
	//TODO: init storage
	s, err := storage.New(dbPath)
	if err != nil{
		log.Error("Some problem", ErrorWrapper(err))
		os.Exit(1)
	}
	if err != nil{
		log.Error("Db error", ErrorWrapper(err))
	}
	userid, err := s.EnterAuth("test@mail.ru", "12345678")
	if err != nil{
		log.Error("Error", ErrorWrapper(err))
		os.Exit(1)
	}
	err = s.CreateEWallet(userid, "CNY")
	if err != nil{
		log.Error("Error", ErrorWrapper(err))
		os.Exit(1)
	}
	log.Info("Storage was init")
	//TODO: init router
	
	//TODO: init server
}

func ErrorWrapper(err error) slog.Attr{
	return slog.Attr{
		Key: "Error",
		Value: slog.StringValue(err.Error()),
	}
}

