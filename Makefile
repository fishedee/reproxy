.PHONY:debug install build stop start_dev_linux start_dev_mac start_idc
dev_linux:build stop start_dev_linux
	echo "finish"
dev_mac:build stop start_dev_mac
	echo "finish"
idc:build stop_idc start_idc
	echo "finish"
build:
	go build main.go
	mv main reproxy
stop:
	rm -rf reproxy.sock
	-sudo pkill reproxy
	-ps aux | grep -v "grep" | grep reproxy
start_dev_linux:
	sudo -u www-data nohup ./reproxy -c config.dev.linux > nohup.out 2>&1 &
	sleep 1
	-ps aux | grep -v "grep" | grep reproxy
start_dev_mac:
	sudo -u nobody ./reproxy -c config.dev.mac > nohup.out 2>&1 &
	sleep 1
	-ps aux | grep -v "grep" | grep reproxy
start_idc:
	sudo supervisorctl restart reproxy
	sleep 1
	-ps aux | grep -v "grep" | grep reproxy
stop_idc:
	sudo supervisorctl stop reproxy
	-pkill -9 reproxy
