wuknet在线SQL管理器是一个免费开源的数据库管理工具，采用GoLang开发， 目前仅支持mysql。 该工具主要为了解决客户端管理数据库的跨平台问题，在服务器上安装后，可以在手机， 平板等设备上通过浏览器来管理数据库。
如果你不想在自己的服务器上安装，也可以通过下面的地址直接管理你的数据库。<br />
https://sql.wuknet.net<br /><br />
Github源码下载地址：<br />
https://github.com/wuknet/sql<br />
下载包根目录的SQL文件为已经编译好的Linux程序。如果你要在其它操作系统上安装，请安装Golang自行编译。<br /><br />

安装方法：<br />
./sql install #安装程序<br />
./sql reinstall #重新安装<br />
./sql start #启动<br />
./sql stop #停止<br />
./sql remove #删除安装<br />
./sql version #查看版本<br />
注意：安装后默认在后台运行，如果要开机自动启动请另行配置。
