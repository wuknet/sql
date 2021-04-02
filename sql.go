package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	//"sync"

	//爬虫库

	"./fun"
	"./page"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"

	//browser "github.com/EDDYCJY/fake-useragent"
	"github.com/Unknwon/goconfig"
	"github.com/jakecoffman/cron"
	//"github.com/alexedwards/scs"
	//"github.com/astaxie/beego"
	//_ "github.com/mattn/go-sqlite3"
	//"github.com/go-session/session"
)

//var db *sql.DB
var (
	AppPath string                  //当前路径
	logger  = service.ConsoleLogger //注册服务
	//wg      sync.WaitGroup //定义一个同步等待的组

	addr     string //配置文件
	isSSL    string
	certFile string
	keyFile  string

	c *cron.Cron
)

//////////////////////////////////////////////////////
type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}
func (p *program) run() {
	go sql_wuknet_net()
}

func sql_wuknet_net() {
	router := httprouter.New()
	router.GET("/*url", page.Default)
	router.POST("/login", page.Login)
	router.POST("/run", page.Run)

	if isSSL == "0" {
		log.Println(http.ListenAndServe(addr, router))
	} else {
		log.Println(http.ListenAndServeTLS(addr, certFile, keyFile, router))
	}
}

func init() { //初始化函数
	AppPath = fun.GetAppPath() //读取当前目录路径
	cfg, err := goconfig.LoadConfigFile(AppPath + "/conf.ini")
	if err != nil {
		fun.Log(0, "读取配置文件错误")
		os.Exit(1)
	}
	addr, _ = cfg.GetValue("domain", "addr")
	isSSL, _ = cfg.GetValue("ssl", "isSSL")
	certFile, _ = cfg.GetValue("ssl", "certFile")
	keyFile, _ = cfg.GetValue("ssl", "keyFile")

	//TimingTask() //定时任务
	//分词词典加载调用
	//fun.Indexes() // 中文分词载入词典

	fmt.Println("GOPATH路径：" + os.Getenv("GOPATH"))
	fmt.Println("程序路径：" + AppPath)
}

func main() {

	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			//fmt.Println(err)
			fun.Log(0, "main:"+fmt.Sprintf("%s", err))
		}
	}()

	svcConfig := &service.Config{
		Name:        "sqladmin",         //服务名称
		DisplayName: "sql.wuknet.net",   //服务显示名称
		Description: "由wuknet.net提供的服务", //服务描述
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		//log.Fatal(err)
		fun.Log(0, "wuknet.main: "+fmt.Sprintf("%s", err))
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			s.Install()
			logger.Info("服务安装成功!")
			s.Start()
			logger.Info("服务启动成功!")
			break
		case "reinstall":
			s.Stop()
			logger.Info("服务关闭成功!")
			s.Uninstall()
			logger.Info("服务卸载成功!")
			s.Install()
			logger.Info("服务安装成功!")
			s.Start()
			logger.Info("服务启动成功!")
			break
		case "start":
			s.Start()
			logger.Info("服务启动成功!")
			break
		case "restart":
			s.Stop()
			logger.Info("服务关闭成功!")
			s.Start()
			logger.Info("服务启动成功!")
			break
		case "stop":
			s.Stop()
			logger.Info("服务关闭成功!")
			break
		case "remove":
			s.Stop()
			logger.Info("服务关闭成功!")
			s.Uninstall()
			logger.Info("服务卸载成功!")
			break
		case "version":
			fmt.Println("wuknet sql V1.0.0")
			break
		}
		return
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
