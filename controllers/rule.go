package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/zjjEcho/spo-api/models"
)

type RuleController struct {
	beego.Controller
}

type Resp struct {
	Status  string      `json:"status"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func (r *RuleController) Get() {
	page, _ := r.GetInt("page")
	size, _ := r.GetInt("size")
	searchItem := r.GetString("searchItem")
	resp := new(Resp)

	ruleList, err := models.GetRules(page, size, searchItem)
	if err != nil {
		beego.Error(err)
		resp.Status = "failed"
		resp.Message = fmt.Sprintf("%v", err)
	} else {
		resp.Status = "success"
		resp.Data = ruleList
	}
	r.Data["json"] = resp
	r.ServeJSON()
}

func (r *RuleController) Put() {
	resp := new(Resp)

	d := new(models.Rule)
	err := json.Unmarshal(r.Ctx.Input.RequestBody, d)
	if err == nil {
		err = models.SetRule(d)
		if err != nil {
			resp.Status = "failed"
			resp.Message = fmt.Sprintf("%v", err)
		} else {
			resp.Status = "success"
		}
	} else {
		resp.Status = "failed"
		resp.Message = fmt.Sprintf("%v", err)
	}

	r.Data["json"] = resp
	r.ServeJSON()
}

func (r *RuleController) Post() {
	resp := new(Resp)

	rule := r.GetString("ruleName")
	status, _ := r.GetBool("status")
	err := models.ChangeStatus(rule, status)
	if err != nil {
		resp.Status = "failed"
		resp.Message = fmt.Sprintf("%v", err)
	} else {
		resp.Status = "success"
	}

	r.Data["json"] = resp
	r.ServeJSON()
}

func (r *RuleController) Delete() {
	resp := new(Resp)

	rules := r.Input()["rules[]"]
	err := models.DeleteRule(rules)
	if err != nil {
		resp.Status = "failed"
		resp.Message = fmt.Sprintf("%v", err)
	} else {
		resp.Status = "success"
	}

	r.Data["json"] = resp
	r.ServeJSON()
}

func (r *RuleController) Count() {
	// m := map[string]string{"status": "success", "count": ""}
	searchItem := r.GetString("searchItem")
	resp := new(Resp)
	c, err := models.CountRules(searchItem)
	if err != nil {
		resp.Status = "failed"
		resp.Message = fmt.Sprintf("%v", err)
		beego.Error(err)
	} else {
		resp.Status = "success"
		resp.Data = map[string]int{"count": c}
		// m["count"] = strconv.Itoa(c)
	}
	r.Data["json"] = resp
	r.ServeJSON()
}
