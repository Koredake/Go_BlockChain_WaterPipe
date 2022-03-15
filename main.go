package main

import (
	"fmt"
	"os"
	"todo/cfg"
	"todo/db"
	"todo/handler"
)

func main() {
	// 检查输入参数
	if len(os.Args) != 2 {
		fmt.Printf("usage: %v config-file\n", os.Args[0])
		return
	}

	// 载入配置文件
	c, err := cfg.LoadConfig(os.Args[1])
	if err != nil {
		fmt.Printf("载入配置文件错误:%v\n", err)
		return
	}
	//连接本地私链节点
	handler.InitNodes(c.GethPort, c.PrivateKey)
	//handler.SendEth()
	// 初始化数据库
	err = db.InitMysql(c.DbUserName, c.DbPassword, c.DbIp, c.DbPort, c.DbName)
	if err != nil {
		fmt.Printf("初始化数据库错误:%v\n", err)
		return
	}
	// 初始化web服务器，监听的地址格式为 0.0.0.0:8080
	fmt.Printf("host:%v port:%v webDir:%v", c.Host, c.Port, c.WebDir)
	err = handler.Start(fmt.Sprintf("%s:%s", c.Host, c.Port), c.WebDir)
	if err != nil {
		fmt.Printf("web服务错误:%v\n", err)
		return
	}
}
