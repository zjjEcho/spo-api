package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/zjjEcho/spo-api/utils"
)

type Rule struct {
	RuleName     string   `json:"ruleName"`
	WhiteList    []string `json:"whiteList"`
	BlackList    []string `json:"blackList"`
	DegradeRule  string   `json:"degradeRule"`
	LimitRule    []string `json:"limitRule"`
	AccessRule   []string `json:"accessRule"`
	BlackReply   string   `json:"blackReply"`
	LimitReply   string   `json:"limitReply"`
	AccessReply  string   `json:"accessReply"`
	DegradeReply string   `json:"degradeReply"`
	Active       bool     `json:"active"`
	UpdateTime   string   `json:"updateTime"`
}

type RuleList []*Rule

func GetRules(page, size int, searchItem string) (RuleList, error) {
	beego.Debug(searchItem)
	// searchItem := "rule"

	sl := []string{}
	if searchItem != "" {
		search := searchItem + "*"
		if len(searchItem) < 3 {
			search = searchItem
		}
		_, rl, err := utils.RedisCliet.FTsearch("ruleindex", search)
		if err != nil {
			beego.Error(err)
		} else {
			sl = rl
		}
	}
	beego.Debug(sl)

	ol, m, err := utils.RedisCliet.Hgetall("config:item:status")
	if err != nil {
		return nil, err
	}

	l := []string{}
	if len(sl) != 0 {
		for _, v := range ol {
			for _, av := range sl {
				if v == av {
					l = append(l, v)
					continue
				}
			}
		}
	} else {
		l = ol
	}
	beego.Debug(l)

	min := (page - 1) * size
	max := len(l)
	if page*size > len(l) || len(l) < size {
		max = len(l)
	} else {
		max = page * size
	}
	lp := l[min:max]

	ruleList := RuleList{}
	for _, v := range lp {
		rule := new(Rule)
		rule.RuleName = v
		if m[v] == "active" {
			rule.Active = true
		}

		_, items, err := utils.RedisCliet.Hgetall("config:item:" + v)
		if err != nil {
			beego.Error(err)
			continue
		}
		for ki, vi := range items {
			if vi == "" {
				continue
			}
			switch ki {
			case "white-iplist":
				rule.WhiteList = strings.Split(vi, ",")
			case "black-iplist":
				rule.BlackList = strings.Split(vi, ",")
			case "degrade-rule":
				rule.DegradeRule = vi
			case "limit-rule":
				rule.LimitRule = strings.Split(vi, ",")
			case "access-rule":
				rule.AccessRule = strings.Split(vi, ",")
			case "black-reply":
				rule.BlackReply = vi
			case "limit-reply":
				rule.LimitReply = vi
			case "access-reply":
				rule.AccessReply = vi
			case "degrade-reply":
				rule.DegradeReply = vi
			case "update-time":
				rule.UpdateTime = vi
			default:
			}
		}
		// rule.UpdateTime = time.Now().Format("2006-01-02 15:04:05.000")
		ruleList = append(ruleList, rule)
	}

	return ruleList, nil
}

func SetRule(r *Rule) error {
	a := "standby"
	if r.Active {
		a = "active"
	}
	err := utils.RedisCliet.Hmset("config:item:status", r.RuleName, a)
	if err != nil {
		return err
	}

	err = utils.RedisCliet.Del("config:item:" + r.RuleName)
	if err != nil {
		return err
	}

	s1 := []interface{}{}
	s1 = append(s1, "config:item:"+r.RuleName)
	s2 := []interface{}{}
	s2 = append(s2, "ruleindex")
	s2 = append(s2, r.RuleName)
	s2 = append(s2, 1.0)
	s2 = append(s2, "REPLACE")
	s2 = append(s2, "FIELDS")
	s2 = append(s2, "name")
	s2 = append(s2, r.RuleName)

	s := []interface{}{}
	s = append(s, "white-iplist")
	s = append(s, strings.Join(r.WhiteList, ","))
	s = append(s, "black-iplist")
	s = append(s, strings.Join(r.BlackList, ","))
	s = append(s, "degrade-rule")
	s = append(s, r.DegradeRule)
	s = append(s, "limit-rule")
	s = append(s, strings.Join(r.LimitRule, ","))
	s = append(s, "access-rule")
	s = append(s, strings.Join(r.AccessRule, ","))

	s = append(s, "black-reply")
	s = append(s, r.BlackReply)
	s = append(s, "limit-reply")
	s = append(s, r.LimitReply)
	s = append(s, "access-reply")
	s = append(s, r.AccessReply)
	s = append(s, "degrade-reply")
	s = append(s, r.DegradeReply)

	s = append(s, "update-time")
	s = append(s, time.Now().Format("2006-01-02 15:04:05.000"))

	s1 = append(s1, s...)
	s2 = append(s2, s...)
	beego.Debug(s2)

	err = utils.RedisCliet.Hmset(s1...)
	if err != nil {
		return err
	}

	// sync redisearch
	err = utils.RedisCliet.FTadd(s2...)
	if err != nil {
		if fmt.Sprintf("%v", err) == "Unknown index name" {
			err = utils.RedisCliet.FTcreate("ruleindex", "SCHEMA", "name", "TEXT", "white-iplist",
				"TEXT", "black-iplist", "TEXT", "degrade-rule", "TEXT", "limit-rule", "TEXT",
				"access-rule", "TEXT", "black-reply", "TEXT", "limit-reply", "TEXT",
				"degrade-reply", "TEXT", "access-reply", "TEXT", "update-time", "TEXT")
			if err == nil {
				utils.RedisCliet.FTadd(s2...)
			} else {
				beego.Error(err)
			}
		} else {
			beego.Error(err)
		}
	}

	// sync lua config
	err = SyncRule()
	if err != nil {
		return err
	}
	return nil
}

func ChangeRelation(rList map[string]string) {
	l, _, err := utils.RedisCliet.Hgetall("config:reply:relation")
	if err != nil {
		beego.Error(err)
		return
	}
	s := []interface{}{}
	s = append(s, "config:reply:relation")
	for _, v := range l {
		s = append(s, v)
		s = append(s, nil)
	}
	err = utils.RedisCliet.Hmset(s...)
	if err != nil {
		beego.Error(err)
	}

	for k, v := range rList {
		err = utils.RedisCliet.Hmset("config:reply:relation", k, v)
		if err != nil {
			beego.Error(err)
			continue
		}
	}
}

func DeleteRule(rules []string) error {
	s := []interface{}{}
	s = append(s, "config:item:status")
	for _, v := range rules {
		s = append(s, v)
	}
	beego.Debug(s)
	err := utils.RedisCliet.Hdel(s...)
	if err != nil {
		return err
	}

	k := make([]string, len(rules))
	for i, v := range rules {
		k[i] = "config:item:" + v
		err = utils.RedisCliet.FTdel("ruleindex", v)
		if err != nil {
			beego.Error(err)
		}
	}
	beego.Debug(k)
	err = utils.RedisCliet.Del(k...)
	if err != nil {
		return err
	}

	err = SyncRule()
	if err != nil {
		return err
	}
	return nil
}

func ChangeStatus(rule string, status bool) error {
	a := "standby"
	if status {
		a = "active"
	}

	err := utils.RedisCliet.Hmset("config:item:status", rule, a)
	if err != nil {
		return err
	}

	err = SyncRule()
	if err != nil {
		return err
	}

	return nil
}

func CountRules(searchItem string) (int, error) {
	beego.Debug(searchItem)
	// searchItem := "rule"

	if searchItem != "" {
		search := searchItem + "*"
		if len(searchItem) < 3 {
			search = searchItem
		}
		c, _, err := utils.RedisCliet.FTsearch("ruleindex", search)
		if err != nil {
			beego.Error(err)
		} else {
			beego.Debug(c)
			return int(c), nil
		}
	}

	c, err := utils.RedisCliet.Hlen("config:item:status")
	if err != nil {
		return 0, err
	}
	return int(c), nil
}

func SyncRule() error {
	r, m, err := utils.RedisCliet.Hgetall("config:item:status")
	if err != nil {
		return err
	}

	l := []string{}
	for k, v := range m {
		if v == "active" {
			l = append(l, k)
		}
	}

	kv := []interface{}{}
	kv = append(kv, "nginx:lua:config")

	for _, v := range l {
		_, m, err := utils.RedisCliet.Hgetall("config:item:" + v)
		if err != nil {
			beego.Error(err)
			continue
		}

		for mk, mv := range m {
			if mk == "update-time" || mv == "" || mv == "0" {
				continue
			}

			if strings.Contains(mk, "-reply") {
				if mv != "" {
					// mid := strings.TrimRight(mk, "-reply")
					mid := strings.Split(mk, "-")[0]
					_, rm, err := utils.RedisCliet.Hgetall("config:reply:" + mv)
					if err != nil {
						beego.Error(err)
					} else {
						for rk, rv := range rm {
							if rv != "" {
								kv = append(kv, v+"-"+mid+"-"+rk)
								kv = append(kv, rv)
							}
						}
					}
				}
				continue
			}

			kv = append(kv, v+"-"+mk)
			kv = append(kv, mv)
		}
	}

	// beego.Debug(kv)
	err = utils.RedisCliet.Del("nginx:lua:config")
	if err != nil {
		return err
	}

	if len(kv) > 1 {
		err = utils.RedisCliet.Hmset(kv...)
		if err != nil {
			return err
		}
	}

	// sync relation
	replyList := make(map[string]string)
	for _, v := range r {
		_, m, err := utils.RedisCliet.Hgetall("config:item:" + v)
		if err != nil {
			beego.Error(err)
			continue
		}

		for kk, vv := range m {
			if strings.Contains(kk, "-reply") && vv != "" {
				if _, ok := replyList[vv]; ok {
					if !strings.Contains(replyList[vv], v) {
						replyList[vv] = replyList[vv] + "," + v
					}
				} else {
					replyList[vv] = v
				}
			}
		}
	}

	beego.Debug(replyList)
	ChangeRelation(replyList)

	return nil
}
