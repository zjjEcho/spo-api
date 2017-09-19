package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/zjjEcho/spo-api/models"
)

type ReplyController struct {
	beego.Controller
}

func (r *ReplyController) Get() {
	resp := new(Resp)
	list, err := models.GetReplys()
	if err != nil {
		beego.Error(err)
		resp.Status = "failed"
		resp.Message = fmt.Sprintf("%v", err)
	} else {
		resp.Status = "success"
		resp.Data = list
	}

	r.Data["json"] = resp
	r.ServeJSON()
}

func (r *ReplyController) Put() {
	resp := new(Resp)
	d := new(models.Reply)
	err := json.Unmarshal(r.Ctx.Input.RequestBody, d)
	if err != nil {
		beego.Error(err)
		resp.Status = "failed"
		resp.Message = fmt.Sprintf("%v", err)
	} else {
		beego.Debug(d)
		err = models.SetReply(d)
		if err != nil {
			resp.Status = "failed"
			resp.Message = fmt.Sprintf("%v", err)
		} else {
			resp.Status = "success"
		}
	}

	r.Data["json"] = resp
	r.ServeJSON()
}

func (r *ReplyController) Delete() {
	resp := new(Resp)

	name := r.GetString("name")
	err := models.DeleteReply(name)
	if err != nil {
		resp.Status = "failed"
		resp.Message = fmt.Sprintf("%v", err)
	} else {
		resp.Status = "success"
	}

	r.Data["json"] = resp
	r.ServeJSON()
}
