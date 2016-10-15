package module

import (
	"errors"
	"fmt"
	. "github.com/fishedee/assert"
	"testing"
	"time"
)

func TestGetConfigTime(t *testing.T) {
	testCase := []struct {
		origin string
		target time.Duration
		err    error
	}{
		{
			"10s",
			10 * time.Second,
			nil,
		},
		{
			"110s",
			110 * time.Second,
			nil,
		},
		{
			"10m",
			10 * time.Minute,
			nil,
		},
		{
			"110m",
			110 * time.Minute,
			nil,
		},
		{
			"10h",
			10 * time.Hour,
			nil,
		},
		{
			"110h",
			110 * time.Hour,
			nil,
		},
		{
			" 110h",
			110 * time.Hour,
			nil,
		}, //左右有空格，没有过滤
		{
			" 110h ",
			110 * time.Hour,
			nil,
		}, //左右有空格，没有过滤
		{
			"110h ",
			110 * time.Hour,
			nil,
		}, //左右有空格，没有过滤
		{
			"2592000s",
			30 * 24 * time.Hour,
			nil,
		},
		{
			"0s",
			0 * time.Second,
			nil,
		},
		/*
			    {
						"6.6s",
						6.6 * time.Second,
						nil,//编译报错 : constant 6.6 truncated to integer
					},
					{
						"-6.6s",
						-6.6 * time.Second,
						nil,//编译报错 : constant -6.6 truncated to integer
					},
		*/
		{
			"0",
			0 * time.Second,
			errors.New("时间设置不正确0"),
		},
		{
			"11",
			0 * time.Second,
			errors.New("时间设置不正确11"),
		},
		{
			"-10s",
			-10 * time.Second,
			nil,
		},
	}

	for singlekey, singleValue := range testCase {
		resultTime, err := GetConfigTime(singleValue.origin)
		AssertEqual(t, resultTime, singleValue.target, singlekey)
		AssertEqual(t, err, singleValue.err, singlekey)
	}

}

func TestGetConfigSize(t *testing.T) {
	testCase := []struct {
		origin string
		target int
		err    error
	}{
		{
			"11b",
			11,
			nil,
		},
		{
			"11k",
			11264,
			nil,
		},
		{
			"11m",
			11534336,
			nil,
		},
		{
			" 11m",
			11534336,
			nil,
		},
		{
			" 11m ",
			11534336,
			nil,
		},
		{
			"11m ",
			11534336,
			nil,
		},
		{
			"0m",
			0,
			nil,
		},
		{
			"-11m",
			-11534336,
			nil,
		},
		{
			"0",
			0,
			errors.New("容量设置不正确0"),
		},
		{
			"11",
			0,
			errors.New("容量设置不正确11"),
		},
	}

	for singlekey, singleValue := range testCase {
		resultSize, err := GetConfigSize(singleValue.origin)
		AssertEqual(t, resultSize, singleValue.target, singlekey)
		AssertEqual(t, err, singleValue.err, singlekey)
	}

}
func TestGetConfigNetInfo(t *testing.T) {
	testCase := []struct {
		origin  string
		target  string
		target2 string
		err     error
	}{
		{
			"tcp://12.123.142.22",
			"tcp",
			"tcp://12.123.142.22",
			nil,
		},
		{
			"unix://12.123.142.22",
			"unix",
			"//12.123.142.22",
			nil,
		},
		{
			":9001",
			"tcp",
			":9001",
			nil,
		},
		{
			" :9001",
			"tcp",
			":9001",
			nil,
		},
		{
			" :9001 ",
			"tcp",
			":9001",
			nil,
		},
		{
			":9001 ",
			"tcp",
			":9001",
			nil,
		},
		{
			"9001",
			"tcp",
			"9001",
			nil,
		},
		{
			"1",
			"tcp",
			"1",
			nil,
		},
		{
			"-1",
			"tcp",
			"-1",
			nil,
		},
		{
			"s",
			"tcp",
			"s",
			nil,
		},
		{
			"",
			"",
			"",
			errors.New("输入地址不能为空"),
		},
	}

	for singlekey, singleValue := range testCase {
		result, result2, err := GetConfigNetInfo(singleValue.origin)
		AssertEqual(t, result, singleValue.target, singlekey)
		AssertEqual(t, result2, singleValue.target2, singlekey)
		AssertEqual(t, err, singleValue.err, singlekey)
	}

}

func TestGetConfigListener(t *testing.T) {
	testCase := []struct {
		origin        string
		network       string
		stringNetwork string
		err           string
	}{
		{
			"tcp://12.123.142.22",
			"tcp",
			"tcp://12.123.142.22",
			"listen tcp: lookup tcp///12.123.142.22: nodename nor servname provided, or not known",
		},
		{
			":9001",
			"tcp",
			"[::]:9001",
			"",
		},
		{
			" :9001",
			"tcp",
			"[::]:9001",
			"",
		},
		{
			" :9001 ",
			"tcp",
			"[::]:9001",
			"",
		},
		{
			":9001 ",
			"tcp",
			"[::]:9001",
			"",
		},
	}

	for singlekey, singleValue := range testCase {

		result, err := GetConfigListener(singleValue.origin)

		if err == nil {
			AssertEqual(t, result.Addr().Network(), singleValue.network, singlekey)
			AssertEqual(t, result.Addr().String(), singleValue.stringNetwork, singlekey)
			result.Close()
		} else {
			AssertEqual(t, fmt.Sprintf("%+v", err.Error()), singleValue.err, singlekey)
		}

	}

}
