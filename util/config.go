package util

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"strings"
	"errors"
	"strconv"
)

type Location struct{
	Url string
	Proxy string
	TimeoutWarn int
	TimeoutError int
	CacheTime int
	CacheSize int
}

type Config struct{
	Listen int
	LogFile string
	Location []Location
}

type JsonConfig struct{
	Listen int `json:"listen"`
	LogFile string `json:"log_file"`
	Location []JsonLocation `json:"location"`
}

type JsonLocation struct{
	Url string `json:"url"`
	Proxy string `json:"proxy"`
	TimeoutWarn string `json:"timeout_warn,omitempty"`
	TimeoutError string `json:"timeout_error,omitempty"`
	CacheTime string `json:"cache_time,omitempty"`
	CacheSize string `json:"cache_size,omitempty"`
}

func analyseTime(in string)(int,error){
	if len(in) == 0{
		return 0,nil
	}
	timeType := string(in[len(in)-1])
	timeTypeDuration := 0 
	if timeType == "s"{
		timeTypeDuration = 1000
	}else if timeType == "m"{
		timeTypeDuration = 1000 * 60
	}else if timeType == "h"{
		timeTypeDuration = 1000 * 60 * 60
	}else{
		return 0,errors.New("时间设置不正确"+in)
	}
	result,err :=strconv.Atoi(strings.Trim(in,timeType))
	if err != nil{
		return 0,err
	}
	return result*timeTypeDuration,nil
}


func analyseSize(in string)(int,error){
	if len(in) == 0{
		return 0,nil
	}
	sizeType := string(in[len(in)-1])
	sizeTypeDuration := 0 
	if sizeType == "b"{
		sizeTypeDuration = 1
	}else if sizeType == "k"{
		sizeTypeDuration = 1024
	}else if sizeType == "m"{
		sizeTypeDuration = 1024 * 1024
	}else{
		return 0,errors.New("容量设置不正确"+in)
	}
	result,err :=strconv.Atoi(strings.Trim(in,sizeType))
	if err != nil{
		return 0,err
	}
	return result*sizeTypeDuration,nil
}

func analyseConfig(config *JsonConfig)(*Config,error){
	result := &Config{}
	result.Listen = config.Listen
	result.LogFile = config.LogFile
	for _,singleLocation := range config.Location{

		singleUrl := singleLocation.Url

		singleProxy := singleLocation.Proxy

		singleTimeoutWarn,err := analyseTime(singleLocation.TimeoutWarn)
		if err != nil{
			return nil,err
		}

		singleTimeoutError,err := analyseTime(singleLocation.TimeoutError)
		if err != nil{
			return nil,err
		}

		singleCacheTime,err := analyseTime(singleLocation.CacheTime)
		if err != nil{
			return nil,err
		}

		singleCacheSize,err := analyseSize(singleLocation.CacheSize)
		if err != nil{
			return nil,err
		}

		singleResultLocation := &Location{
			Url:singleUrl,
			Proxy:singleProxy,
			TimeoutWarn:singleTimeoutWarn,
			TimeoutError:singleTimeoutError,
			CacheTime:singleCacheTime,
			CacheSize:singleCacheSize,
		}
		result.Location = append(result.Location,*singleResultLocation)
		
	}
	return result,nil
}

func GetConfigFromFile(fileName string)(*Config,error){
	file, err := os.Open(fileName)
    if err != nil {
    	return nil,err
    }
    defer file.Close()

    configFile,err := ioutil.ReadAll(file)
    if err != nil{
    	return nil,err
    }

    var result *JsonConfig;
	err = json.Unmarshal(configFile,&result);
	if err != nil{
		return nil,err
	}

	result2,err := analyseConfig(result)
	if err != nil{
		return nil,err
	}

	return result2,err
}