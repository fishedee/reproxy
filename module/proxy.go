package module

import (
	"io/ioutil"
	"net/http"
	"time"
	. "github.com/fishedee/reverse-proxy/handler"
)

type ProxyConfig struct{
	Listen string `json:"listen"`
	Server []ProxyServerConfig `json:"server"`
	Location []ProxyLocationConfig `json:"location"`
}

type ProxyServerConfig struct{
	Name string `json:"name"`
	Type string `json:"type"`
	Address string `json:"address"`
}

type ProxyLocationConfig struct{
	Url string `json:"url"`
	Proxy string `json:"proxy"`
	TimeoutWarn string `json:"timeout_warn,omitempty"`
	TimeoutError string `json:"timeout_error,omitempty"`
	CacheTime string `json:"cache_time,omitempty"`
	CacheSize string `json:"cache_size,omitempty"`
}

type RouteHandler struct{
	Host string
	TimeoutWarn time.Duration
	Cache *Cache
	CacheExpireTime time.Duration
	Client ProxyHandler
}

func (this *RouteHandler) HandleHttpRequest(request *http.Request)(*CacheResponse,error){
	request.URL.Scheme = "http"
	request.URL.Host = this.Host
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
func (this *RouteHandler) HandleHttp(result chan error ,writer http.ResponseWriter,request *http.Request){
	resp,err := this.HandleHttpRequest(request)
	if err != nil{
		result <- err
		return
	}

    for k, v := range resp.Header {
        for _, vv := range v {
            writer.Header().Add(k, vv)
        }
    }
    writer.WriteHeader(resp.StatusCode)
    writer.Write(resp.Body)
    result <- nil
}

func (this *RouteHandler) HandleTimeoutAndHttp(logBeginner string,writer http.ResponseWriter,request *http.Request){
	resultChan := make(chan error)
	go this.HandleHttp(resultChan,writer,request)
	select {
	case result := <- resultChan:
		if result != nil{
			Logger.Error(
				logBeginner,
				result.Error(),
			)
		}
	case <- time.After(this.TimeoutWarn):
		Logger.Warn(
			logBeginner,
			"执行时间超长 ",
		)
		result := <-resultChan
		if result != nil{
			Logger.Error(
				logBeginner,
				result.Error(),
			)
		}
	}
	close(resultChan)
}

func (this *RouteHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request){
	beginTime := time.Now().UnixNano()
	logBeginner := request.RemoteAddr + " -- [" +request.Method+" "+request.RequestURI+"] "
	this.HandleTimeoutAndHttp(logBeginner,writer,request)
	endTime := time.Now().UnixNano()
	
	Logger.Info(
		logBeginner,
		"execution time: ",
		float64(endTime-beginTime)/1000000,
		"ms",
	)
}

func SeviceProxy(config ProxyConfig)(error){
	for _,singleLocation := range config.Location{
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

		proxy := singleLocation.Proxy

		Logger.Info("Handle Url "+singleLocation.Url)
		http.Handle(
			singleLocation.Url,
			&RouteHandler{
				Host:proxy,
				TimeoutWarn:timeoutWarn,
				Cache:NewCache(cacheSize),
				CacheExpireTime:cacheExpireTime,
				Client:NewFastCgiHandler(timeoutError),
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