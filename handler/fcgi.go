package handler

import (
	"github.com/mholt/caddy/middleware/fastcgi"
	"strings"
	"net/http"
	"time"
	"bytes"
	"io/ioutil"
	"io"
	"strconv"
)

type FastCgiHandler struct{
	protocol string
	address string
	documentRoot string
	documentIndex string
	params map[string]string
}

func NewFastCgiHandler(protocol string, address string, documentRoot string,documentIndex string,params map[string]string,timeoutError time.Duration )(ProxyHandler,error){
	return &FastCgiHandler{
		protocol:protocol,
		address:address,
		params:params,
		documentRoot:documentRoot,
		documentIndex:documentIndex,
	},nil
}

func (this *FastCgiHandler)Do(request *http.Request)(*http.Response,error){
	fcgi,err := fastcgi.Dial(this.protocol,this.address)
	if err != nil{
		return nil,err
	}

	//设置header数据
	header := map[string]string{}
	header["QUERY_STRING"] = request.URL.RawQuery
	header["REQUEST_METHOD"] = request.Method
	header["SCRIPT_FILENAME"] =  strings.TrimLeft(this.documentRoot,"/") + "/" + strings.TrimLeft(this.documentIndex,"/")
	header["SCRIPT_NAME"] = strings.TrimLeft(this.documentIndex,"/")
	header["REQUEST_URI"] = request.URL.RequestURI()
	header["DOCUMENT_URI"] = strings.TrimLeft(this.documentRoot,"/")
	header["DOCUMENT_ROOT"] = strings.TrimLeft(this.documentIndex,"/")
	header["SERVER_PROTOCOL"] = "HTTP/1.1"
	header["GATEWAY_INTERFACE"] = "CGI/1.1"
	header["SERVER_SOFTWARE"] = "reverse-proxy/1.0"
	header["REMOTE_PORT"] = "123"
	header["SERVER_ADDR"] = "127.0.0.1"
	header["SERVER_PORT"] = "8001"
	header["SERVER_NAME"] = ""
	header["HTTPS"] = "0"
	if request.Header["X-Real-Ip"] != nil &&
		len(request.Header["X-Real-Ip"]) != 0 {
		header["REMOTE_ADDR"] = request.Header["X-Real-Ip"][0]
	}else{
		header["REMOTE_ADDR"] = request.RemoteAddr
	}
	for key,value := range this.params{
		header[key] = value
	}
	header["HTTP_HOST"] = request.Host;
	for key,value := range request.Header{
		header["HTTP_"+strings.ToUpper(key)] = value[0]
	}

	//设置body素据
	var bodyReader io.Reader
	body := request.Body
	if body != nil{
		defer body.Close()
		if request.Header["Content-Type"] != nil &&
			len(request.Header["Content-Type"]) != 0 &&
			request.Header["Content-Length"] != nil &&
			len(request.Header["Content-Length"]) != 0 {
			header["CONTENT_TYPE"] = request.Header["Content-Type"][0]
			header["CONTENT_LENGTH"] = request.Header["Content-Length"][0]
			bodyReader = body
		}else{
			data,err := ioutil.ReadAll(body)
			if err != nil{
				return nil,err
			}
			header["CONTENT_TYPE"] = http.DetectContentType(data)
			header["CONTENT_LENGTH"] = strconv.Itoa(len(data))
			bodyReader = bytes.NewReader(data)
		}
	}
	
	return fcgi.Request(header,bodyReader)
}

