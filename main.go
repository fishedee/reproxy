package main

import (
	"flag"
	"fmt"

	. "./module"
)

func main() {
	configFilePath := flag.String("c", "config", "config file address")
	flag.Parse()

	config, err := GetConfigFromFile(*configFilePath)
	if err != nil {
		fmt.Printf("配置文件失败: '%s'\n", err.Error())
		return
	}

	err = InitUser(config.UserConfig)
	if err != nil {
		fmt.Printf("初始化用户失败: '%s'\n", err.Error())
		return
	}

	err = InitLogger(config.Log)
	if err != nil {
		fmt.Printf("启动日志失败: '%s'\n", err.Error())
		return
	}

	err = SeviceProxy(config.ProxyConfig)
	if err != nil {
		fmt.Printf("启动服务器失败: '%s'\n", err.Error())
	}
}
