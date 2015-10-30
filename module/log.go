package module

import (
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"errors"
)

var Logger *logs.BeeLogger

type LogConfig struct{
	Maxlines int `json:"maxlines"`
	Filename string `json:"filename"`
	Maxdays int64  `json:"maxdays"`
	Rotate bool  `json:"rotate"`
	Level string  `json:"level"`
}

type logConfigInner struct{
	Maxlines int `json:"maxlines"`
	Filename string `json:"filename"`
	Daily bool  `json:"daily"`
	Maxdays int64  `json:"maxdays"`
	Rotate bool  `json:"rotate"`
	Level int  `json:"level"`
}

func InitLogger(config LogConfig)(error){
	var configInner logConfigInner
	if config.Filename == ""{
		return errors.New("文件名不能为空哈")
	}
	configInner.Filename = config.Filename

	if config.Maxlines == 0{
		config.Maxlines = 1000000
	}
	configInner.Maxlines = config.Maxlines

	configInner.Daily = true

	if config.Maxdays == 0{
		config.Maxdays = 7
	}
	configInner.Maxdays = config.Maxdays

	logLevel := map[string]int{
		"error":logs.LevelError,
		"warn":logs.LevelWarning,
		"info":logs.LevelInformational,
	}
	var logLevelInt int
	if config.Level == ""{
		logLevelInt = logs.LevelInformational
	}else{
		var ok bool
		logLevelInt,ok = logLevel[config.Level]
		if ok == false{
			return errors.New("日志等级设置不正确")
		}
	}
	configInner.Level = logLevelInt
	configInner.Rotate = true

	Logger = logs.NewLogger(10000)
	logConfigString,err := json.Marshal(&configInner)
	if err != nil{
		return err
	}

	err = Logger.SetLogger("file", string(logConfigString))
	if err != nil{
		return err
	}
	return nil
}