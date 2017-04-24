package main

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"ruleAnaly/models"

	"github.com/astaxie/beego"
	"github.com/coreos/etcd/clientv3"
	"github.com/hprose/hprose-golang/rpc"
)

func init() {
	logDir, _ := filepath.Abs("../logs/ruleAnaly")
	//if the directory does not exist to create
	if _, err := os.Stat(logDir); err != nil {
		os.Mkdir(logDir, os.ModePerm)
	}
	//initialize the log configuration
	beego.SetLogger("file", `{
		"filename":"`+logDir+`/app.log",
		"level":7,
		"maxlines":0,
		"maxsize":0,
		"daily":true,
		"maxdays":10 
	}`)
}
func main() {
	//connection service
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   beego.AppConfig.Strings("etcdAddrs"),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		beego.Critical("Connection service registry failed, ", err)
		return
	}
	defer client.Close()

	port := beego.AppConfig.String("tcpport")
	//registration service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = client.Put(ctx, "ruleAnaly", "tcp://:"+port)
	cancel()
	if err != nil {
		beego.Critical(err)
		return
	}
	//start service
	srv := rpc.NewTCPServer("tcp://:" + port)
	srv.AddFunction("Run", models.Run)

	beego.Info("ruleAnaly starting with :" + port)

	err = srv.Start()
	if err != nil {
		beego.Critical(err)
		return
	}
	beego.Run()
}
