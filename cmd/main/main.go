package main

import (
	"e-wallet/config"
	"e-wallet/internal/lib/env"
	"e-wallet/internal/logger"
	"e-wallet/internal/storage"
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
	db, err := storage.New(dbPath)
	if err != nil{
		log.Error("Some problem")
		os.Exit(1)
	}
	_ = db
	log.Info("Storage was init")
	//TODO: init router

	//TODO: init server
}