package handler

import (
	"net/http"
	"net"
	"time"
)

type HttpHandler struct{
	client *http.Client
}

func getFakeDial(protocol string,address string)(func (string,string)(conn net.Conn, err error) ){
	return func(a1 string,a2 string) (conn net.Conn, err error) {
		return net.Dial(protocol, address)
	}
}
func NewHttpHandler(protocol string,address string,timeoutError time.Duration )(ProxyHandler,error){
	return &HttpHandler{
		client:&http.Client{
			Timeout:timeoutError,
			Transport:&http.Transport{
			    Dial: getFakeDial(protocol,address),
			},
		},
	},nil
}

func (this *HttpHandler)Do(request *http.Request)(*http.Response,error){
	request.URL.Scheme = "http"
	request.URL.Host = "test"
	resp, err := this.client.Do(request)
    if err != nil{
    	return nil,err
    }
    return resp,nil
}

