package util

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type RouteHandler struct{
	Host string
	TimeoutWarn time.Duration
	TimeoutError time.Duration
	Cache *Cache
	CacheExpireTime time.Duration
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
	
	client := &http.Client{
		Timeout:this.TimeoutError,
	}

    resp, err := client.Do(request)
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
			Logger.Err(
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
			Logger.Err(
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

func SeviceProxy(port int,location []Location)(error){
	for _,singleLocation := range location{
		if singleLocation.TimeoutWarn == 0{
			singleLocation.TimeoutWarn = 5*1000
		}
		if singleLocation.TimeoutError == 0 {
			singleLocation.TimeoutError = 30*1000
		}
		http.Handle(
			singleLocation.Url,
			&RouteHandler{
				Host:singleLocation.Proxy,
				TimeoutWarn:time.Duration(singleLocation.TimeoutWarn) * time.Millisecond,
				TimeoutError:time.Duration(singleLocation.TimeoutError) * time.Millisecond,
				Cache:NewCache(singleLocation.CacheSize),
				CacheExpireTime:time.Duration(singleLocation.CacheTime) * time.Millisecond,
			},
		)
	}
	Logger.Info("Start Proxy Server Listen On :"+strconv.Itoa(port))
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil{
		return err
	}
	return nil
}