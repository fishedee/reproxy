package module

import (
	"io/ioutil"
	"net/http"
	"time"
	"errors"
	. "github.com/fishedee/reverse-proxy/handler"
)

type ProxyConfig struct{
	Listen string `json:"listen"`
	Server []ProxyServerConfig `json:"server"`
	Location []ProxyLocationConfig `json:"location"`
}

type ProxyLocationConfig struct{
	Url string `json:"url"`
	Server string `json:"server"`
	TimeoutWarn string `json:"timeout_warn,omitempty"`
	TimeoutError string `json:"timeout_error,omitempty"`
	CacheTime string `json:"cache_time,omitempty"`
	CacheSize string `json:"cache_size,omitempty"`
}

type RouteHandler struct{
	TimeoutWarn time.Duration
	Cache *Cache
	CacheExpireTime time.Duration
	Client ProxyHandler
}

func (this *RouteHandler) HandleHttpRequest(request *http.Request)(*CacheResponse,error){
	request.RequestURI = ""
	queryParam := request.URL.Query()
	queryParam.Del("t")
	queryParam.Del("_")
	request.URL.RawQuery = queryParam.Encode()

	url := request.URL.String()
	method := request.Method
	cacheResponse := this.Cache.Get(method,url)
	if cacheResponse != nil{
		return cacheResponse,nil
	}

	hasLock := this.Cache.AcquireLock(method,url)
	if hasLock == false{
		cacheResponse = this.Cache.Get(method,url+"_old")
		if cacheResponse != nil{
			return cacheResponse,nil
		}
	}else{
		defer this.Cache.ReleaseLock(method,url)
	}

    resp, err := this.Client.Do(request)
    if err != nil{
    	return nil,err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil{
    	return nil,err
    }

    cacheResponse = &CacheResponse{
    	Header:resp.Header,
    	StatusCode:resp.StatusCode,
    	Body:body,
    }

    if hasLock{
    	this.Cache.Set(method,url,cacheResponse,this.CacheExpireTime)
   		this.Cache.Set(method,url+"_old",cacheResponse,0)
    }

    return cacheResponse,nil
}
func (this *RouteHandler) HandleHttp(writer http.ResponseWriter,request *http.Request)(int,error){
	resp,err := this.HandleHttpRequest(request)
	if err != nil{
		return 0,err
	}

    for k, v := range resp.Header {
        for _, vv := range v {
            writer.Header().Add(k, vv)
        }
    }
    writer.WriteHeader(resp.StatusCode)
    writer.Write(resp.Body)
   	return resp.StatusCode,nil
}

func (this *RouteHandler) HandleTimeoutAndHttp(logBeginner string,writer http.ResponseWriter,request *http.Request)(error){
	timer := time.AfterFunc(this.TimeoutWarn,func(){
		Logger.Warn(
			"%s 执行时间超长",
			logBeginner,
		)
	})
	defer timer.Stop()
	statusCode,err := this.HandleHttp(writer,request)
	if err != nil{
		return err
	}
	if statusCode != 200 &&
		statusCode != 304{
		Logger.Warn(
			"%s 返回码:%d",
			logBeginner,
			statusCode,
		)
	}
	return nil
}

func (this *RouteHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request){
	beginTime := time.Now().UnixNano()
	logBeginner := request.RemoteAddr + " -- [" +request.Method+" "+request.RequestURI+"]"
	err := this.HandleTimeoutAndHttp(logBeginner,writer,request)
	endTime := time.Now().UnixNano()
	
	if err != nil{
		Logger.Error(
			"%s %s",
			logBeginner,
			err.Error(),
		)
	}
	Logger.Info(
		"%s execution time: %f ms",
		logBeginner,
		float64(endTime-beginTime)/1000000,
	)
}

func SeviceProxy(config ProxyConfig)(error){
	serverMap := map[string]ProxyServerConfig{}
	for _,singleServer := range config.Server{
		serverMap[singleServer.Name] = singleServer
	}

	for _,singleLocation := range config.Location{
		url := singleLocation.Url
		server := singleLocation.Server

		timeoutWarn,err := GetConfigTime(singleLocation.TimeoutWarn)
		if err != nil{
			return err
		}
		if timeoutWarn == 0{
			timeoutWarn= 5*time.Second
		}

		timeoutError,err := GetConfigTime(singleLocation.TimeoutError)
		if err != nil{
			return err
		}
		if timeoutError == 0 {
			timeoutError = 30*time.Second
		}

		cacheExpireTime,err := GetConfigTime(singleLocation.CacheTime)
		if err != nil{
			return err
		}

		cacheSize,err := GetConfigSize(singleLocation.CacheSize)
		if err != nil{
			return err
		}

		singleServer,ok := serverMap[server]
		if ok == false{
			return errors.New("没有url找到对应的server "+url)
		}

		client,err := NewHandler(singleServer,timeoutError)
		if err != nil{
			return err
		}

		Logger.Info("Handle Url "+singleLocation.Url)
		http.Handle(
			singleLocation.Url,
			&RouteHandler{
				TimeoutWarn:timeoutWarn,
				Cache:NewCache(cacheSize),
				CacheExpireTime:cacheExpireTime,
				Client:client,
			},
		)
	}
	listener,err := GetConfigListener(config.Listen)
	if err != nil{
		return err
	}

	Logger.Info("Start Proxy Server Listen On "+config.Listen)
	err = http.Serve(listener, nil)
	if err != nil{
		return err
	}
	return nil
}