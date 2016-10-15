package module

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ProxyConfig
	UserConfig
	Log LogConfig `json:"log"`
}

func GetConfigTime(in string) (time.Duration, error) {
	in = strings.TrimSpace(in)
	if len(in) == 0 {
		return 0, nil
	}
	var timeTypeDuration time.Duration
	timeType := string(in[len(in)-1])
	timeTypeDuration = 0
	if timeType == "s" {
		timeTypeDuration = time.Second
	} else if timeType == "m" {
		timeTypeDuration = time.Minute
	} else if timeType == "h" {
		timeTypeDuration = time.Hour
	} else {
		return 0, errors.New("时间设置不正确" + in)
	}
	result, err := strconv.Atoi(strings.Trim(in, timeType))
	if err != nil {
		return 0, err
	}
	return time.Duration(result) * timeTypeDuration, nil
}

func GetConfigSize(in string) (int, error) {
	in = strings.TrimSpace(in)
	if len(in) == 0 {
		return 0, nil
	}
	sizeType := string(in[len(in)-1])
	sizeTypeDuration := 0
	if sizeType == "b" {
		sizeTypeDuration = 1
	} else if sizeType == "k" {
		sizeTypeDuration = 1024
	} else if sizeType == "m" {
		sizeTypeDuration = 1024 * 1024
	} else {
		return 0, errors.New("容量设置不正确" + in)
	}
	result, err := strconv.Atoi(strings.Trim(in, sizeType))
	if err != nil {
		return 0, err
	}
	return result * sizeTypeDuration, nil
}

func GetConfigNetInfo(address string) (string, string, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return "", "", errors.New("输入地址不能为空")
	}
	addrInfo := strings.Split(address, ":")
	if len(addrInfo) == 1 {
		return "tcp", address, nil
	} else if len(addrInfo) == 2 {
		if addrInfo[0] == "unix" {
			return "unix", addrInfo[1], nil
		} else {
			return "tcp", address, nil
		}
	} else {
		return "", "", errors.New("不合法的地址信息" + address)
	}
}

func GetConfigListener(address string) (net.Listener, error) {
	protocol, addr, err := GetConfigNetInfo(address)
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen(protocol, addr)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func GetConfigConnect(address string) (net.Conn, error) {
	protocol, addr, err := GetConfigNetInfo(address)
	if err != nil {
		return nil, err
	}

	dialer, err := net.Dial(protocol, addr)
	if err != nil {
		return nil, err
	}
	return dialer, nil
}

func GetConfigFromFile(fileName string) (*Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	configFile, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var result *Config
	err = json.Unmarshal(configFile, &result)
	if err != nil {
		return nil, err
	}

	return result, err
}
