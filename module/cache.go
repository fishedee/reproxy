package module

import (
	"encoding/gob"
	"net/http"
	"time"
	"bytes"
	"sync"
	"strings"
	"github.com/coocood/freecache"
)

type Cache struct{
	cache *freecache.Cache
	mutex *sync.Mutex
	expireTime time.Duration
	data map[string]bool
}

type CacheResponse struct{
	Header map[string][]string
	StatusCode int
	Body []byte
}

func NewCache(size int,expireTime time.Duration)(*Cache){
	if size == 0 {
		return &Cache{}
	}else{
		return &Cache{
			cache:freecache.NewCache(size),
			mutex:&sync.Mutex{},
			data:map[string]bool{},
			expireTime:expireTime,
		}
	}	
}

func (this *Cache)AcquireLock(url string)(bool){
	this.mutex.Lock()
	defer this.mutex.Unlock()

	result,ok := this.data[url]
	if ok == false || result == false{
		this.data[url] = true
		return true
	}else{
		return false
	}
}

func (this *Cache)ReleaseLock(url string){
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.data[url] = false
}

func (this *Cache)GetInner(url string)(*CacheResponse){
	data,err := this.cache.Get([]byte(url))
	if err != nil{
		return nil
	}

	var result CacheResponse
    buf := bytes.NewBuffer(data)  
    enc := gob.NewDecoder(buf)  
    err = enc.Decode(&result)  
    if err != nil {  
    	Logger.Error("decode error!"+err.Error())
        return nil
    }
	return &result
}


func (this *Cache)SetInner(url string,response *CacheResponse,expireTime time.Duration){
	buf := bytes.NewBuffer(nil)
    enc := gob.NewEncoder(buf)  
    err := enc.Encode(response)
    if err != nil {
    	Logger.Error("encode error!"+err.Error())
        return
    }

    this.cache.Set([]byte(url),buf.Bytes(),int(expireTime.Seconds()))
    return
}

func (this *Cache)GetCacheUrlString(request *http.Request)(string){
	requestUrl := *request.URL
	requestQueryParam := requestUrl.Query()
	requestQueryParam.Del("t")
	requestQueryParam.Del("_")
	requestUrl.RawQuery = requestQueryParam.Encode()
	requestUrlString := requestUrl.String()
	return requestUrlString
}

func (this *Cache)Get(request *http.Request)(*CacheResponse,bool){
	if request.Method != "GET" || this.cache == nil{
		return nil,false;
	}

	cacheUrl := this.GetCacheUrlString(request)
	response := this.GetInner(cacheUrl)
	if response != nil{
		return response,false
	}

	hasLock := this.AcquireLock(cacheUrl)
	if hasLock == false{
		response = this.GetInner(cacheUrl+"_old")
	}
	return response,hasLock
}

func (this *Cache)Set(request *http.Request,response *CacheResponse,hasLock bool){
	if request.Method != "GET" || this.cache == nil || hasLock == false{
		return 
	}

	var newResponse CacheResponse
	newResponse.StatusCode = response.StatusCode
	newResponse.Body = response.Body
	newResponse.Header = map[string][]string{}
	for key,value := range response.Header{
		lowerKey := strings.ToLower(key)
		if lowerKey == "set-cookie"{
			continue
		}
		newResponse.Header[key] = value
	}

	cacheUrl := this.GetCacheUrlString(request)
	this.SetInner(cacheUrl,&newResponse,this.expireTime)
	this.SetInner(cacheUrl+"_old",&newResponse,0)
}