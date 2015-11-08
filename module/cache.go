package module

import (
	"encoding/gob"
	"time"
	"bytes"
	"sync"
	"github.com/coocood/freecache"
)

type Cache struct{
	cache *freecache.Cache
	mutex *sync.Mutex
	data map[string]bool
}

type CacheResponse struct{
	Header map[string][]string
	StatusCode int
	Body []byte
}

func NewCache(size int)(*Cache){
	if size == 0 {
		return &Cache{}
	}else{
		return &Cache{
			cache:freecache.NewCache(size),
			mutex:&sync.Mutex{},
			data:map[string]bool{},
		}
	}	
}

func (this *Cache)AcquireLock(method string,url string)(bool){
	if this.cache == nil || method != "GET"{
		return true
	}

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

func (this *Cache)ReleaseLock(method string,url string){
	if this.cache == nil || method != "GET"{
		return
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.data[url] = false
}

func (this *Cache)Get(method string,url string)(*CacheResponse){
	if this.cache == nil || method != "GET"{
		return nil
	}

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

    delete(result.Header,"Set-Cookie")
	return &result
}


func (this *Cache)Set(method string,url string,response *CacheResponse,expireTime time.Duration){
	if this.cache == nil || method != "GET"{
		return
	}

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