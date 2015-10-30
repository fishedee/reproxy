package handler

import (
	"github.com/mholt/caddy/middleware/fastcgi"
	"strings"
	"net/http"
	"time"
)

type FastCgiHandler struct{
}

func NewFastCgiHandler(timeoutError time.Duration )(ProxyHandler){
	return &FastCgiHandler{}
}

func (this *FastCgiHandler)Do(request *http.Request)(*http.Response,error){
	fcgi,err := fastcgi.Dial("unix",request.URL.Host)
	if err != nil{
		return nil,err
	}

	body := request.Body
	if body != nil{
		defer body.Close()
	}

	header := map[string]string{}
	header["QUERY_STRING"] = request.URL.RawQuery
	header["REQUEST_METHOD"] = request.Method
	//header["CONTENT_TYPE"] = request.Header["Content-Type"][0]
	//header["CONTENT_LENGTH"] = request.Header["Content-Length"][0]
	header["SCRIPT_FILENAME"] = "var/www/BakeWeb/server/index.php"
	header["SCRIPT_NAME"] = "server/index.php"
	header["REQUEST_URI"] = request.URL.RequestURI()
	header["DOCUMENT_URI"] = "var/www/BakeWeb"
	header["DOCUMENT_ROOT"] = "server/index.php"
	header["SERVER_PROTOCOL"] = "HTTP/1.1"
	header["GATEWAY_INTERFACE"] = "CGI/1.1"
	header["SERVER_SOFTWARE"] = "nginx/1.4.6"
	header["REMOTE_ADDR"] = request.RemoteAddr
	header["REMOTE_PORT"] = "123"
	header["SERVER_ADDR"] = "127.0.0.1"
	header["SERVER_PORT"] = "8001"
	header["SERVER_NAME"] = ""
	header["HTTPS"] = "0"
	for key,value := range request.Header{
		header["HTTP_"+strings.ToUpper(key)] = value[0]
	}
	return fcgi.Request(header,body)
}

