package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	dbURL     string `json:"dbURL"`
	dbName    string `json:"dbName"`
	redisAddr string `json:"redisAddr"`
}

func InitConfig() {
	f, err := os.Open("config.json")
	if err != nil {
		panic("Failed to open config file")
	}
	defer func() {
		f.Close()
	}()
	if err = json.NewDecoder(f).Decode(config); err != nil {
		panic("Failed to decode config file")
	}
	return
}
