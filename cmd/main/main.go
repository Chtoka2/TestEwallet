package main

import (
	"e-wallet/config"
	"fmt"
)

func main(){
	//TODO: init config
	cfg := config.Config_init()
	fmt.Println(cfg)
	//TODO: init logger

	//TODO: init storage

	//TODO: init router

	//TODO: init server
}