package utils

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	RedisCliet *RedisInst
)

type RedisInst struct {
	Pool *redis.Pool
}

func NewClient(server, password string) {
	RedisCliet = new(RedisInst)
	RedisCliet.Pool = &redis.Pool{
		MaxIdle:     3,
		MaxActive:   5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// func (r *RedisInst) Hset(key, field, value string) error {
// 	conn := r.Pool.Get()
// 	defer conn.Close()

// 	_, err := conn.Do("HSET", key, field, value)
// 	if err != nil {
// 		fmt.Printf("%v\n", err)
// 		return err
// 	}
// 	return nil
// }

func (r *RedisInst) Hmset(kvs ...interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("HMSET", kvs...)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	return nil
}

func (r *RedisInst) Hlen(key string) (int64, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	res, err := redis.Int64(conn.Do("HLEN", key))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return res, nil
}

func (r *RedisInst) Del(keys ...string) error {
	k := make([]interface{}, len(keys))
	for i, v := range keys {
		k[i] = v
	}

	conn := r.Pool.Get()
	defer conn.Close()

	_, err := redis.Int64(conn.Do("DEL", k...))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (r *RedisInst) Hdel(kvs ...interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("HDEL", kvs...)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	return nil
}

func (r *RedisInst) Hget(key, field string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("HGET", key, field))
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *RedisInst) Hgetall(key string) ([]string, map[string]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	res1, err := redis.Strings(conn.Do("HKEYS", key))
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, nil, err
	}

	res2, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		fmt.Printf("%v\n", err)
		return res1, nil, err
	}
	return res1, res2, nil
}

func (r *RedisInst) Keys(key string) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	res, err := redis.Strings(conn.Do("KEYS", key))
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	return res, nil
}

func (r *RedisInst) FTadd(kvs ...interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("FT.ADD", kvs...))
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisInst) FTcreate(kvs ...interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("FT.CREATE", kvs...))
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisInst) FTdel(i, k string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := redis.Int64(conn.Do("FT.DEL", i, k))
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisInst) FTsearch(kvs ...interface{}) (int64, []string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	kvs = append(kvs, "NOCONTENT", "LIMIT", 0, 1000)
	res, err := conn.Do("FT.SEARCH", kvs...)
	if err != nil {
		return 0, nil, err
	}

	c := res.([]interface{})[0].(int64)
	if c == 0 {
		return 0, nil, nil
	}

	res1, err := redis.Strings(res.([]interface{})[1:], nil)
	if err != nil {
		return c, nil, err
	}
	return c, res1, nil
}
