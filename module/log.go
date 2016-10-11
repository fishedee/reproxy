package module

import (
	"encoding/json"
	"errors"

	"github.com/astaxie/beego/logs"
)

var Logger *logs.BeeLogger
var RateLogger *logs.BeeLogger

type LogConfig struct {
	Maxlines int    `json:"maxlines"`
	Filename string `json:"filename"`
	Maxdays  int64  `json:"maxdays"`
	Rotate   bool   `json:"rotate"` // 回转
	Level    string `json:"level"`
}

type logConfigInner struct {
	Maxlines int    `json:"maxlines"`
	Filename string `json:"filename"`
	Daily    bool   `json:"daily"`
	Maxdays  int64  `json:"maxdays"`
	Rotate   bool   `json:"rotate"`
	Level    int    `json:"level"`
}

// 初始化日志模块
func newLogger(config LogConfig) (*logs.BeeLogger, error) {
	var logger *logs.BeeLogger

	var configInner logConfigInner
	if config.Filename == "" {
		return logger, errors.New("文件名不能为空哈")
	}
	configInner.Filename = config.Filename

	if config.Maxlines == 0 {
		config.Maxlines = 1000000
	}
	configInner.Maxlines = config.Maxlines

	configInner.Daily = true

	if config.Maxdays == 0 {
		config.Maxdays = 7
	}
	configInner.Maxdays = config.Maxdays

	logLevel := map[string]int{
		"error": logs.LevelError,
		"warn":  logs.LevelWarning,
		"info":  logs.LevelInformational,
	}
	var logLevelInt int
	if config.Level == "" {
		logLevelInt = logs.LevelInformational
	} else {
		var ok bool
		logLevelInt, ok = logLevel[config.Level]
		if ok == false {
			return logger, errors.New("日志等级设置不正确")
		}
	}
	configInner.Level = logLevelInt
	configInner.Rotate = true

	// 普通日志
	logger = logs.NewLogger(10000)
	logConfigString, err := json.Marshal(&configInner)
	if err != nil {
		return logger, err
	}

	err = logger.SetLogger("file", string(logConfigString))
	if err != nil {
		return logger, err
	}

	return logger, nil
}

func InitLogger(config LogConfig) error {
	var err error
	Logger, err = newLogger(config)
	if err != nil {
		return err
	}
	return nil
}

func InitRateLogger(config LogConfig) error {
	var err error
	RateLogger, err = newLogger(config)
	if err != nil {
		return err
	}
	return nil
}
