package main

import (
	"e-wallet/config"
	"e-wallet/internal/lib/env"
	"e-wallet/internal/logger"
)

func main(){
	//TODO: init config
	cfg := config.Config_init()
	//TODO: init .env
	env := env.Env_reader()
	
	//TODO: init logger
	log := logger.Logger_init(cfg.Env)
	log.Info("Start new logger!")
	
	//TODO: init storage
	
	//TODO: init router

	//TODO: init server
}