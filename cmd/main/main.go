package main

import (
	"context"
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
	dotenv := env.Env_reader()
	//TODO: init logger
	log := logger.Logger_init(cfg.Env)
	log.Info("Start new logger!")
	
	//TODO: init storage
	s, err := storage.New(dotenv.DbURL)
	if err != nil{
		log.Error("Some problem", ErrorWrapper(err))
		os.Exit(1)
	}
	if err != nil{
		log.Error("Db error", ErrorWrapper(err))
	}
	userid, err := s.GetUserIDByEmail(context.Background(), "test@mail.ru")
	if err != nil{
		log.Error("Cannot get userid", ErrorWrapper(err))
		os.Exit(1)
	}
	err = s.ConvertCurrency(
		context.Background(),
		userid,
		"RUB", "USD",
		100,
		0.012734,
	)
	if err != nil{
		log.Error("Cannot convert", ErrorWrapper(err))
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

