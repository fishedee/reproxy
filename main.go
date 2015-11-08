package main

import (
	"fmt"
	. "github.com/fishedee/reproxy/module"
)

func main(){
	config,err := GetConfigFromFile("config")
	if err != nil{
		fmt.Println("配置文件失败 "+err.Error())
		return
	}

	err = InitUser(config.UserConfig)
	if err != nil{
		fmt.Println("初始化用户失败 "+err.Error())
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
