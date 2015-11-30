package main

import (
	"fmt"
	"flag"
	. "./module"
)

func main(){
	configFilePath := flag.String("c","config","config file address")
	flag.Parse()

	config,err := GetConfigFromFile(*configFilePath)
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
