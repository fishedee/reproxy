package handler

import (
	"net/http"
	"time"
)

type HttpHandler struct{
	client *http.Client
}

func NewHttpHandler(host string,timeoutError time.Duration )(ProxyHandler){
	return &HttpHandler{
		client:&http.Client{Timeout:timeoutError},
	}
}

func (this *HttpHandler)Do(request *http.Request)(*http.Response,error){
	resp, err := this.client.Do(request)
    if err != nil{
    	return nil,err
    }
    return resp,nil
}

