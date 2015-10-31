package handler

import (
	"net/http"
	"time"
	"strings"
	"errors"
)

type ProxyServerConfig struct{
	Name string `json:"name"`
	Type string `json:"type"`
	Address string `json:"address"`
	DocumentRoot string `json:"document_root,omitempty"`
	DocumentIndex string `json:"document_index,omitempty"`
}

type ProxyHandler interface{
	Do(request *http.Request)(*http.Response,error)
}

func getConfigNetInfo(address string)(string,string,error){
	addrInfo := strings.Split(address,":")
	if len(addrInfo) == 1{
		return "tcp",address,nil
	}else if len(addrInfo) == 2{
		if addrInfo[0] == "unix"{
			return "unix",addrInfo[1],nil
		}else{
			return "tcp",address,nil
		}
	}else{
		return "","",errors.New("不合法的地址信息"+address)
	}
}

func NewHandler(singleServer ProxyServerConfig,timeoutError time.Duration)(ProxyHandler,error){
	singleServerProtocol,singleServerAddr,err := getConfigNetInfo(singleServer.Address)
	if err != nil{
		return nil,err
	}
	if singleServer.Type == "fastcgi"{
		return NewFastCgiHandler(
			singleServerProtocol,
			singleServerAddr,
			singleServer.DocumentRoot,
			singleServer.DocumentIndex,
			timeoutError,
		)
	}else if singleServer.Type == "" || singleServer.Type == "http"{
		return NewHttpHandler(
			singleServerProtocol,
			singleServerAddr,
			timeoutError,
		)
	}else{
		return nil,errors.New("不合法的server type: "+singleServer.Type)
	}

}