package handler

import (
	"net/http"
)

type ProxyHandler interface{
	Do(request *http.Request)(*http.Response,error)
}