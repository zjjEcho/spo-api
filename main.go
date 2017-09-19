package main

import (
	_ "github.com/zjjEcho/spo-api/routers"

	"github.com/astaxie/beego"
	"github.com/zjjEcho/spo-api/utils"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
	}

	beego.SetLogger("file", `{"filename":"logs/test.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":5}`)

	utils.NewClient("127.0.0.1:6379", "")

	beego.Run()
}
