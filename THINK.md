未优化前

8001 /login/islogin 49.41

8002 /login/islogin 46.51

加上nginx与proxy长连接

8001 /login/islogin 49.47

8002 /login/islogin 46.67

加入fastcgi直连

8001 /login/islogin 49.22

8002 /login/islogin 48.47

修复fastcgi没有关闭短连接的问题

8001 /login/islogin 40.32

8002 /login/islogin 39.55

总结

fastcig直连的优化比较明显，基本上达到了nginx连接fastcgi的速度，在单线程的go环境下能达到这样效率，基本上是ok了。
