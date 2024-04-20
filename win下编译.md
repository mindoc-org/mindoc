win下
修改conf里
注释掉mysql
取消sqlite注释
拷贝一个sqlite数据库到database文件夹里，更名为mindoc.db
cmd进入mindoc文件夹
因为用sqlite需要cgo，go env -w CGO_ENABLED=1
https://blog.csdn.net/qq_43625649/article/details/134488353

D:\gowork>cd src

D:\gowork\src>cd github.com

D:\gowork\src\github.com>cd mindoc

D:\gowork\src\github.com\mindoc>bee run
______
| ___ \
| |_/ /  ___   ___
| ___ \ / _ \ / _ \
| |_/ /|  __/|  __/
\____/  \___| \___| v2.0.4
2024/04/19 21:21:29 WARN     ▶ 0001 Running application outside of GOPATH
2024/04/19 21:21:29 INFO     ▶ 0002 Using 'mindoc' as 'appname'
2024/04/19 21:21:29 INFO     ▶ 0003 Initializing watcher...
github.com/mattn/go-sqlite3
github.com/mindoc-org/mindoc
2024/04/19 21:22:04 SUCCESS  ▶ 0004 Built Successfully!
2024/04/19 21:22:04 INFO     ▶ 0005 Restarting 'mindoc.exe'...
2024/04/19 21:22:04 SUCCESS  ▶ 0006 './mindoc.exe' is running...
2024/04/19 21:22:05.230 [I] [command.go:38]  正在初始化数据库配置.
2024/04/19 21:22:05.260 [I] [command.go:115]  数据库初始化完成.
MinDoc version =>
build time =>
start directory => D:\gowork\src\github.com\mindoc\mindoc.exe

2024/04/19 21:22:05.453 [I] [server.go:281]  http server Running on http://:8181

