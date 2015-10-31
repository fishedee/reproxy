.PHONY:debug install watch release
debug:
	-pkill reverse-proxy
	go run main.go
install:
	go get github.com/coocood/freecache
watch:
	bee run
release:
	go build main.go
	mv main reverse-proxy
	-pkill reverse-proxy
	nohup ./reverse-proxy > nohup.out 2>&1 &
	tail -f nohup.out
