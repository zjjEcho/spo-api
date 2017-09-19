package controllers

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/zjjEcho/spo-api/models"
)

type ServerController struct {
	beego.Controller
}

func (s *ServerController) Get() {
	serverList := models.ServerList{
		&models.Server{
			Id:         0,
			Name:       "test",
			Status:     1,
			IP:         []string{"192.168.11.11", "127.0.0.1"},
			HttpHosts:  []string{"flash.dianwoda.com", "api.dianwoda.com"},
			Rules:      []string{"rule1", "rule2"},
			UpdateTime: time.Now().Format("2006-01-02 15:04:05.000"),
		},
		&models.Server{
			Id:         1,
			Name:       "test1",
			Status:     2,
			IP:         []string{"192.168.11.12"},
			HttpHosts:  []string{},
			Rules:      []string{"rule1", "rule2"},
			UpdateTime: time.Now().Format("2006-01-02 15:04:05.000"),
		},
		&models.Server{
			Id:         2,
			Name:       "test2",
			Status:     3,
			IP:         []string{"192.168.11.13"},
			HttpHosts:  []string{"flash3.dianwoda.com", "api3.dianwoda.com"},
			Rules:      []string{"rule8", "rule1"},
			UpdateTime: time.Now().Format("2006-01-02 15:04:05.000"),
		},
	}

	s.Data["json"] = serverList
	s.ServeJSON()
}

func (s *ServerController) Count() {
	m := map[string]string{"status": "success", "num": "3"}
	s.Data["json"] = m
	s.ServeJSON()
}

func (s *ServerController) Post() {
	t, err := models.GetSaltServers("compute01.openstack.office", "poc002.novalocal")
	if err != nil {
		beego.Error(err)
	}
	beego.Debug(t)

	m := map[string]string{"status": "success"}
	s.Data["json"] = m
	s.ServeJSON()
}

// CREATE TABLE IF NOT EXISTS server( id INT UNSIGNED AUTO_INCREMENT, server_name VARCHAR(100), status TINYINT, ip_list VARCHAR(300), update_time TIMESTAMP, PRIMARY KEY ( id ) )ENGINE=InnoDB DEFAULT CHARSET=utf8;
