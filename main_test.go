package main

import (
	// "errors"
	// "fmt"
	// . "github.com/fishedee/assert"
	// "github.com/fishedee/web"
	"net/http"
	"runtime"
	"sync"
	"testing"
	// "time"
	// "os"
	"fmt"
	"log"
	"os"
)

func TestMain(t *testing.T) {
	Concurrent(19, 19, func() {
		resp, err := http.Get("http://localhost/") // 需要测试的目标网址
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		if resp.StatusCode == http.StatusOK {
			fmt.Println(resp.StatusCode)
		}
		defer resp.Body.Close()

		buf := make([]byte, 1024)
		f, err1 := os.OpenFile("path.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm) //可读写，追加的方式打开（或创建文件）
		if err1 != nil {
			panic(err1)
		}
		defer f.Close()

		for {
			n, _ := resp.Body.Read(buf)
			if 0 == n {
				break
			}
			f.WriteString(string(buf[:n]))
		}
	})

	runtime.Gosched()
	runtime.Gosched()
	runtime.Gosched()
	runtime.Gosched()
	runtime.Gosched()
	runtime.Gosched()

}

func Concurrent(number int, concurrency int, handler func()) {
	if number <= 0 {
		panic("benchmark numer is invalid")
	}
	if concurrency <= 0 {
		panic("benchmark concurrency is invalid")
	}
	singleConcurrency := number / concurrency
	if singleConcurrency <= 0 ||
		number%concurrency != 0 {
		panic("benchmark numer/concurrency is invalid")
	}

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runtime.LockOSThread()
			for i := 0; i < singleConcurrency; i++ {
				handler()
			}
		}()
	}
	wg.Wait()
}
