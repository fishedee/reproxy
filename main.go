package main

import (
	"fmt"
	"runtime"
	. "github.com/fishedee/reverse-proxy/module"
)

func main(){
	runtime.GOMAXPROCS(2)

	config,err := GetConfigFromFile("config")
	if err != nil{
		fmt.Println("配置文件失败 "+err.Error())
		return
	}

	err = InitLogger(config.Log)
	if err != nil{
		fmt.Println("启动日志失败 "+err.Error())
		return
	}

	err = SeviceProxy(config.ProxyConfig)
	if err != nil{
		fmt.Println("启动服务器失败 "+err.Error())
	}
}
