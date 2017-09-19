package models

import (
	"strings"

	"github.com/zjjEcho/spo-api/utils"
)

type Server struct {
	Id         int32    `json:"id"`
	Name       string   `json:"serverName"`
	Status     int32    `json:"status"`
	IP         []string `json:"ipAddress"`
	HttpHosts  []string `json:"httpHosts"`
	Rules      []string `json:"rules"`
	UpdateTime string   `json:"updateTime"`
}

type ServerList []*Server

type SaltServer map[string][]string

func GetSaltServers(servers ...string) (SaltServer, error) {
	salt := utils.NewSaltClient("192.168.11.23:8000", "saltapi", "saltapi")
	resp, err := salt.GetIP(strings.Join(servers, ","))
	if err != nil {
		return nil, err
	}
	if resp.Code == 401 {
		_, err := salt.Auth()
		if err != nil {
			return nil, err
		}
		resp, err = salt.GetIP(strings.Join(servers, ","))
		if err != nil {
			return nil, err
		}
	}

	retu := resp.Data.(map[string]interface{})
	slist := retu["return"].([]interface{})[0].(map[string]interface{})
	saltServers := make(SaltServer)
	for k, v := range slist {
		ipList := v.(map[string]interface{})["ipv4"].([]interface{})
		temp := []string{}
		for _, v := range ipList {
			if v.(string) != "127.0.0.1" {
				temp = append(temp, v.(string))
			}
		}
		saltServers[k] = temp
	}

	return saltServers, nil
}
