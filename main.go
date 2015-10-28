package main

import (
	"fmt"
	. "github.com/fishedee/reverse-proxy/util"
)

func main(){
	config,err := GetConfigFromFile("config")
	if err != nil{
		fmt.Println("配置文件失败 "+err.Error())
		return
	}

	err = InitLogger(config.LogFile)
	if err != nil{
		fmt.Println("启动日志失败 "+err.Error())
		return
	}

	err = SeviceProxy(config.Listen,config.Location)
	if err != nil{
		fmt.Println("启动服务器失败 "+err.Error())
	}
}
