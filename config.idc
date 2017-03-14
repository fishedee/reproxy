{
	"listen":":8002",
	"user":"www-data",
	"log":{
		"filename":"log/access.log",
		"level":"info"
	},
	"rate":{
		"comment":"请求频率限制：time时间内只能请求max次",
		"max":40,
		"time":"1s",
		"log":{
			"filename":"log/rateips.log",
			"level":"info"
		},
		"cache_size":"100m"
	},
	"server":[
		{
			"name":"php",
			"type":"fastcgi",
			"address":"10.163.126.155:9000",
			"document_root":"/var/www/BakeWeb",
			"document_index":"/server/index.php",
			"params":{
				"CI_ENV":"production"
			}
		},
		{
			"name":"go",
			"type":"http",
			"comment":"idc1",
			"address":"10.163.126.155:9001"
		},
		{
			"name":"live",
			"type":"http",
			"comment":"idc3",
			"address":"10.144.144.249:9001"
		}
	],
	"location":[
	{
		"url":"/live/",
		"server":"live"
	},
	{
		"url":"/course/",
		"server":"live"
	},
	{
		"url":"/paycallback/",
		"server":"live"
	},
	{
		"url":"/activity/get",
		"server":"go",
		"cache_time":"1s",
		"cache_size":"16m"
	},
	{
		"url":"/activity/getComponent",
		"server":"go",
		"cache_time":"1s",
		"cache_size":"16m"
	},
	{
		"url":"/backstage/topicbanner/",
		"server":"php"
	},
	{
		"url":"/backstage/topic/",
		"server":"php"
	},
	{
		"url":"/backstage/post/",
		"server":"php"
	},
	{
		"url":"/backstage/wxspilder/",
		"server":"php"
	},
	{
		"url":"/backstage/crawl/",
		"server":"php"
	},
	{
		"url":"/",
		"server":"go"
	}
	]
}
