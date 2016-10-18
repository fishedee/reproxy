package module

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	. "../handler"

	"github.com/coocood/freecache"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/config"
	"github.com/didip/tollbooth/libstring"
)

var rateCache *freecache.Cache

type ProxyConfig struct {
	Listen   string                `json:"listen"`
	Rate     ProxyRate             `json:"rate"`
	Server   []ProxyServerConfig   `json:"server"`
	Location []ProxyLocationConfig `json:"location"`
}

type ProxyRate struct {
	Max       int64     `json:"max"`
	Time      string    `json:"time"`
	Log       LogConfig `json:"log"`        // 记录频繁IP
	CacheSize string    `json:"cache_size"` // 缓存大小
}

type ProxyLocationConfig struct {
	Url          string `json:"url"`
	Server       string `json:"server"`
	TimeoutWarn  string `json:"timeout_warn,omitempty"`
	TimeoutError string `json:"timeout_error,omitempty"`
	CacheTime    string `json:"cache_time,omitempty"`
	CacheSize    string `json:"cache_size,omitempty"`
}

type RouteHandler struct {
	TimeoutWarn time.Duration
	Cache       *Cache
	Client      ProxyHandler
}

func (this *RouteHandler) SetCache(request *http.Request, cacheResponse **CacheResponse, hasLock bool) {
	this.Cache.Set(request, *cacheResponse, hasLock)
}

// 发出请求
func (this *RouteHandler) HandleHttpRequest(request *http.Request) (*CacheResponse, error) {
	cacheResponse, hasLock := this.Cache.Get(request)
	defer this.SetCache(request, &cacheResponse, hasLock)

	if cacheResponse != nil {
		return cacheResponse, nil
	}

	request.RequestURI = ""
	resp, err := this.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//FIXME 特殊逻辑
	if request.URL.Path == "/appstatic/getFile" {
		Logger.Info(
			"%v => Code:[%v],Header:[%v]",
			request.URL,
			resp.StatusCode,
			resp.Header,
		)
		_, isExist := resp.Header["Etag"]
		if resp.StatusCode == 200 && isExist == false {
			Logger.Error(
				"%v => Has No Etag!!!",
				request.URL,
			)
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cacheResponse = &CacheResponse{
		Header:     resp.Header,
		StatusCode: resp.StatusCode,
		Body:       body,
	}
	return cacheResponse, nil
}

// 处理返回
func (this *RouteHandler) HandleHttp(writer http.ResponseWriter, request *http.Request) (int, error) {
	resp, err := this.HandleHttpRequest(request)
	if err != nil {
		return 0, err
	}

	for k, v := range resp.Header {
		for _, vv := range v {
			writer.Header().Add(k, vv)
		}
	}
	writer.WriteHeader(resp.StatusCode)
	statusCode := resp.StatusCode
	if statusCode != 200 &&
		statusCode != 304 {
		Logger.Warn(
			"响应: %+v",
			resp,
		)
	}
	writer.Write(resp.Body)
	return resp.StatusCode, nil
}

// 处理超时
func (this *RouteHandler) HandleTimeoutAndHttp(logBeginner string, writer http.ResponseWriter, request *http.Request) error {
	// TimeoutWarn时间到后，请求还没结束，则调用参数函数
	timer := time.AfterFunc(this.TimeoutWarn, func() {
		Logger.Warn(
			"%s 执行时间超长",
			logBeginner,
		)
	})
	defer timer.Stop()

	// 检查返回码
	statusCode, err := this.HandleHttp(writer, request)
	if err != nil {
		return err
	}
	if statusCode != 200 &&
		statusCode != 304 {
		Logger.Warn(
			"%s 返回码:%d",
			logBeginner,
			statusCode,
		)
	}
	return nil
}

// 请求
func (this *RouteHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	beginTime := time.Now().UnixNano()
	logBeginner := request.RemoteAddr + " -- [" + request.Method + " " + request.RequestURI + "]"
	err := this.HandleTimeoutAndHttp(logBeginner, writer, request)
	endTime := time.Now().UnixNano()

	if err != nil {
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

// 启动代理
func SeviceProxy(config ProxyConfig) error {
	// 服务器配置映射
	serverMap := map[string]ProxyServerConfig{}
	for _, singleServer := range config.Server {
		serverMap[singleServer.Name] = singleServer
	}

	// 频率限制时间
	timeDuration, err := GetConfigTime(config.Rate.Time)
	if err != nil {
		return err
	}

	// 启动日志
	err = InitRateLogger(config.Rate.Log)
	if err != nil {
		return err
	}

	// Ip黑名单文件
	rateIpFilename := config.Rate.Log.Filename

	// 缓存大小
	rateCacheSize, err := GetConfigSize(config.Rate.CacheSize)
	if err != nil {
		return err
	}

	// 缓存初始化
	rateCache, _ = initRateCache(rateCacheSize, rateIpFilename)

	// 路由分发
	for _, singleLocation := range config.Location {
		url := singleLocation.Url
		server := singleLocation.Server

		timeoutWarn, err := GetConfigTime(singleLocation.TimeoutWarn)
		if err != nil {
			return err
		}
		if timeoutWarn == 0 {
			timeoutWarn = 5 * time.Second
		}

		timeoutError, err := GetConfigTime(singleLocation.TimeoutError)
		if err != nil {
			return err
		}
		if timeoutError == 0 {
			timeoutError = 30 * time.Second
		}

		cacheExpireTime, err := GetConfigTime(singleLocation.CacheTime)
		if err != nil {
			return err
		}

		cacheSize, err := GetConfigSize(singleLocation.CacheSize)
		if err != nil {
			return err
		}

		singleServer, ok := serverMap[server]
		if ok == false {
			return errors.New("没有url找到对应的server " + url)
		}

		client, err := NewHandler(singleServer, timeoutError)
		if err != nil {
			return err
		}

		Logger.Info("Handle Url " + singleLocation.Url)

		// 初始化请求频率限制中间件
		limiter := tollbooth.NewLimiter(config.Rate.Max, timeDuration)
		limiter.IPLookups = []string{"X-Real-IP"} // 指定获取IP的header字段
		handler := RateLimitHandler(
			limiter,
			&RouteHandler{
				TimeoutWarn: timeoutWarn,
				Cache:       NewCache(cacheSize, cacheExpireTime),
				Client:      client,
			},
		)

		// 配置路由
		http.Handle(singleLocation.Url, handler)
	}

	// 代理监听端口
	listener, err := GetConfigListener(config.Listen)
	if err != nil {
		return err
	}

	// 开启代理
	Logger.Info("Start Proxy Server Listen On " + config.Listen)
	err = http.Serve(listener, nil)
	if err != nil {
		return err
	}
	return nil
}

// Ip请求频率校验
func RateLimitHandler(limiter *config.Limiter, next http.Handler) http.Handler {
	middle := func(w http.ResponseWriter, r *http.Request) {
		tollbooth.SetResponseHeaders(limiter, w)

		// 检查Ip
		remoteIP := libstring.RemoteIP(limiter.IPLookups, r)
		if remoteIP != "127.0.0.1" {
			_, err := rateCache.Get([]byte(remoteIP))
			if err == nil {
				w.WriteHeader(200)
				w.Write([]byte(`{"code":10002,"msg":"","data":"","remindPoint":{"count":0,"data":[]}}`))
				return
			}

			// 频率校验
			httpError := tollbooth.LimitByRequest(limiter, r)
			if httpError != nil {
				// 记录超频Ip
				rateCache.Set([]byte(remoteIP), []byte("true"), 0)

				// 同步到日志
				RateLogger.Info(remoteIP)

				// 返回429
				w.Header().Add("Content-Type", limiter.MessageContentType)
				w.WriteHeader(httpError.StatusCode)
				w.Write([]byte(httpError.Message))
				return
			}
		}

		// 正常访问
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(middle)
}

func initRateCache(rateCacheSize int, fileName string) (*freecache.Cache, error) {
	cache := freecache.NewCache(rateCacheSize)

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return cache, err
	}
	lines := strings.Split(string(data), "\n")
	for _, single := range lines {
		lineInfo := strings.Split(single, " ")
		if len(lineInfo) == 0 {
			continue
		}
		cache.Set([]byte(lineInfo[len(lineInfo)-1]), []byte("true"), 0)
	}

	return cache, nil
}
