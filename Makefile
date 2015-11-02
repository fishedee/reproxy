.PHONY:debug install build stop start restart
debug:build stop start
	echo "finish"
install:
	go get github.com/coocood/freecache
build:
	go build main.go
	mv main reproxy
stop:
	rm -rf reproxy.sock
	-pkill reproxy
	-ps aux | grep -v "grep" | grep reproxy
start:
	nohup ./reproxy > nohup.out 2>&1 &
	sleep 2s
	sudo chown www-data:www-data reproxy.sock
	-ps aux | grep -v "grep" | grep reproxy
restart:build stop start
	echo "finish"
