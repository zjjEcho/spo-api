package models

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/zjjEcho/spo-api/utils"
)

type Reply struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type ReplyList []*Reply

func GetReplys() (ReplyList, error) {
	// l, err := utils.RedisCliet.Keys("config:reply:*")
	// if err != nil {
	// 	return nil, err
	// }

	l, _, err := utils.RedisCliet.Hgetall("config:reply:relation")
	if err != nil {
		return nil, err
	}

	replyList := ReplyList{}
	for _, v := range l {
		_, m, err := utils.RedisCliet.Hgetall("config:reply:" + v)
		if err != nil {
			beego.Error(err)
			continue
		}
		reply := new(Reply)
		reply.Name = m["name"]
		reply.Code = m["code"]
		reply.Type = m["header"]
		reply.Content = m["content"]
		replyList = append(replyList, reply)
	}

	return replyList, nil
}

func SetReply(r *Reply) error {
	_, err := utils.RedisCliet.Hget("config:reply:relation", r.Name)
	if err != nil {
		if fmt.Sprintf("%v", err) == "redigo: nil returned" {
			err = utils.RedisCliet.Hmset("config:reply:relation", r.Name, nil)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	err = utils.RedisCliet.Del("config:reply:" + r.Name)
	if err != nil {
		return err
	}

	err = utils.RedisCliet.Hmset("config:reply:"+r.Name, "name", r.Name, "code", r.Code, "header", r.Type, "content", r.Content)
	if err != nil {
		return err
	}

	err = SyncRule()
	if err != nil {
		return err
	}

	return nil
}

func DeleteReply(n string) error {
	res, err := utils.RedisCliet.Hget("config:reply:relation", n)
	if err != nil && fmt.Sprintf("%v", err) != "redigo: nil returned" {
		return err
	}

	if res != "" {
		return errors.New(fmt.Sprintf("rule in [%v], can't remove", res))
	}

	err = utils.RedisCliet.Hdel("config:reply:relation", n)
	if err != nil {
		return err
	}

	err = utils.RedisCliet.Del("config:reply:" + n)
	if err != nil {
		return err
	}
	return nil
}
