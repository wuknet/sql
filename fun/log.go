package fun

import (
	//"io"
	"fmt"
	"log"
	"os"
)

var logger *log.Logger

func Log(Levle int, Message string) {
	//Message = Gettime() + "  " + Message + "\r\n"
	Message = Message + "\r\n"
	switch Levle {
	case 0:
		{
			SaveLog(Message)
		}
	}
}

func SaveLog(Message string) {
	// 按照所需读写权限创建文件
	AppPath := GetAppPath() //读取当前目录路径
	f, err := os.OpenFile(AppPath+"/run.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger = log.New(f, "", log.LstdFlags|log.Lshortfile) // 日志文件格式:log包含时间及文件行数
	//logger.SetFlags(log.LstdFlags | log.Lshortfile) // 设置日志格式
	fmt.Println(Message) //输出日志到命令行终端
	//log.Println(Message) //输出日志到命令行终端
	logger.Println(Message) //将日志写入文件
}
