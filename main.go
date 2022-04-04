package main

import (
	"log"
	"net/http"
)

var (
	config = &Config{}
)

func main() {
	InitConfig()
	if err := InitDB(config.dbURL, config.dbName); err != nil {
		panic("Failed to init mongo db")
	}
	if err := InitRedis(config.redisAddr); err != nil {
		panic("Failed to init mongo db")
	}
	http.HandleFunc("/", HelloWorld) // 设置访问的路由
	http.HandleFunc("/url", ScanURL)
	http.HandleFunc("/file", ScanFile)
	http.HandleFunc("/history", GetHistory)
	err := http.ListenAndServe(":9091", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
