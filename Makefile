.PHONY:debug install build stop start restart
debug:build stop start
	echo "finish"
install:
	go get github.com/coocood/freecache
build:
	go build main.go
	mv main reverse-proxy
stop:
	rm -rf reverse-proxy.sock
	-pkill reverse-proxy
	-ps aux | grep -v "grep" | grep reverse-proxy
start:
	nohup ./reverse-proxy > nohup.out 2>&1 &
	sleep 2s
	sudo chown www-data:www-data reverse-proxy.sock
	-ps aux | grep -v "grep" | grep reverse-proxy
restart:build stop start
	echo "finish"
