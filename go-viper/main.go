package main

import (
	"fmt"
	"go-viper/internal/pkg/setting"
)

func main() {

	v := setting.Viper
	cfg := setting.InitConfig()

	setting.InitLog()

	//打印环境变量
	fmt.Println(v.Get("MYSQL_HOST"))
	fmt.Println(v.Get("log-level"))
	fmt.Println(v.Get("ENV"))
	//
	fmt.Println(cfg.Mysql.Master.Username)
	fmt.Println(cfg.Mysql.Master.Host)
	fmt.Println(cfg.Server.Port)
	fmt.Println(v.Get("server.port"))
	fmt.Println(v.Get("log.logFile"))
	fmt.Println(v.Get("log.level"))
	setting.Logger.Info("测试")

}
