{
	"listen":":8002",
	"user":"www-data",
	"log":{
		"filename":"log/access.log",
		"level":"info"
	},
	"rate":{
		"comment":"请求频率限制：time时间内只能请求max次",
		"max":20,
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
			"address":"10.144.144.249:9001"
		}
	],
	"location":[
		{
			"url":"/appstatic/",
			"server":"go"
		},
		{
			"url":"/backstage/goldenstatue/",
			"server":"go"
		},
		{
			"url":"/file/",
			"server":"go"
		},
		{
			"url":"/backstage/file/",
			"server":"go"
		},
		{
			"url":"/rate/",
			"server":"go"
		},
		{
			"url":"/backstage/rate/",
			"server":"go"
		},
		{
			"url":"/appmsg/",
			"server":"go"
		},
		{
			"url":"/backstage/appmsg/",
			"server":"go"
		},
		{
			"url":"/point/sendDayPoint",
			"server":"go"
		},
		{
			"url":"/weixin/",
			"server":"go"
		},
		{
			"url":"/backstage/weixin/",
			"server":"go"
		},
		{
			"url":"/invite/checkInvite",
			"server":"go"
		},
		{
			"url":"/invite/inviteClient",
			"server":"go"
		},
		{
			"url":"/backstage/invite/",
			"server":"go"
		},
		{
             "url":"/export/",
             "server":"go"
	    },
		{
             "url":"/backstage/export/",
             "server":"go"
	    },
		{
             "url":"/backstage/virtual/",
             "server":"go"
	    },
		{
			"url":"/appidfa/",
			"server":"go"
		},
		{
			"url":"/backstage/appidfa/",
			"server":"go"
		},
		{
			"url":"/search/",
			"server":"go"
		},
		{
			"url":"/keyword/",
			"server":"go"
		},
		{
			"url":"/backstage/search/",
			"server":"go"
		},
		{
			"url":"/feed/",
			"server":"go"
		},
		{
			"url":"/question/getNew",
			"server":"go"
		},
		{
			"url":"/question/getHot",
			"server":"go"
		},
		{
			"url":"/question/getUnAnswer",
			"server":"go"
		},
		{
			"url":"/client/get",
			"server":"go"
		},
		{
			"url":"/client/getDish",
			"server":"go"
		},
		{
			"url":"/client/getRecipe",
			"server":"go"
		},
		{
			"url":"/client/getQuestionAndAnswer",
			"server":"go"
		},
		{
			"url":"/backstage/category/",
			"server":"go"
		},
		{
			"url":"/index/",
			"server":"go"
		},
		{
			"url":"/backstage/link/",
			"server":"go"
		},
		{
			"url":"/link/",
			"server":"go"
		},
		{
			"url":"/backstage/complaint/",
			"server":"go"
		},
		{
			"url":"/complaint/",
			"server":"go"
		},
		{
			"url":"/backstage/user/",
			"server":"go"
		},
		{
			"url":"/backstage/notice/",
			"server":"go"
		},
		{
			"url":"/notice/",
			"server":"go"
		},
		{
			"url":"/backstage/comment/batchDel",
			"server":"go"
		},
		{
			"url":"/backstage/question/batchDelQuestion",
			"server":"go"
		},
		{
			"url":"/backstage/question/batchDelAnswer",
			"server":"go"
		},
		{
			"url":"/bind/",
			"server":"go"
		},
		{
			"url":"/backstage/bind/",
			"server":"go"
		},
		{
			"url":"/backstage/block/",
			"server":"go"
		},
		{
			"url":"/login/",
			"server":"go"
		},
		{
			"url":"/login/wxlogin",
			"server":"php"
		},
		{
			"url":"/login/wxlogincallback",
			"server":"php"
		},
		{
			"url":"/account/",
			"server":"go"
		},
		{
			"url":"/backstage/letter/",
			"server":"go"
		},
		{
			"url":"/letter/",
			"server":"go"
		},
		{
			"url":"/backstage/groupbuy/",
			"server":"go"
		},
		{
			"url":"/groupbuy/",
			"server":"go"
		},
		{
			"url":"/backstage/feedback/",
			"server":"go"
		},
		{
			"url":"/feedback/",
			"server":"go"
		},
		{
			"url":"/backstage/client/testlogin",
			"server":"go"
		},
		{
			"url":"/backstage/pointsign/",
			"server":"go"
		},
		{
			"url":"/sign/",
			"server":"go"
		},
		{
			"url":"/backstage/tribe/",
			"server":"go"
		},
		{
			"url":"/tribe/",
			"server":"go"
		},
		{
			"url":"/backstage/shield/",
			"server":"go"
		},
		{
			"url":"/backstage/appversion/",
			"server":"go"
		},
		{
			"url":"/backstage/appstart/",
			"server":"go"
		},
		{
			"url":"/appstart/",
			"server":"go"
		},
		{
			"url":"/item/getTaoBaoKeFavorites",
			"server":"go"
		},
		{
			"url":"/item/getByFavoritesId",
			"server":"go"
		},
		{
			"url":"/backstage/remind/",
			"server":"go"
		},
		{
			"url":"/backstage/clientaction/",
			"server":"go"
		},
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
