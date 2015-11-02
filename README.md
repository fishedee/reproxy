# [reproxy](https://github.com/fishedee/reverse-proxy)
high performance reverse proxy written by go

### Feature

* fast like nginx
* support fastcgi And http protocol
* cache data in memory identify by url
* slow request warning
* abnormal status warning

### Install

git clone github.com/fishedee/reproxy
make
./reproxy 

### Config
```
{
	"listen":"unix:reverse-proxy.sock",
	"log":{
		"filename":"log/access.log",
		"level":"info"
	},
	"server":[
		{
			"name":"php",
			"type":"fastcgi",
			"address":"unix:/var/run/fastcgi/php5-fpm.sock",
			"document_root":"/var/www/BakeWeb",
			"document_index":"/server/index.php"
		},
		{
			"name":"go",
			"type":"http",
			"address":":8001"
		}
	],
	"location":[
		{
			"url":"/activity/get",
			"server":"php",
			"cache_time":"1s",
			"cache_size":"16m"
		},
		{
			"url":"/activity/getComponent",
			"server":"php",
			"cache_time":"1s",
			"cache_size":"16m"
		},
		{
			"url":"/",
			"server":"php"
		}
	]
}
```
