package crudProvider

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
	"surface_attack/consts"
)

type Implamintation struct {
	conn redis.Conn
}

func (i *Implamintation) InitConn(Addr string) error {

	var err error
	i.conn, err = redis.Dial("tcp", Addr)
	if err != nil {
		return err
	}
	log.Println("API connected to REDIS succesfully on address: ", Addr)
	return nil
}

func (i *Implamintation) EndConn() {
	i.conn.Close()
	log.Println("API disconnected from REDIS succesfully")
}

func (i *Implamintation) Get(key string) (string, error) {
	reply, e := i.conn.Do(consts.REDIS_GET, key)
	if e != nil {
		return "", e
	}

	if reply == nil && key == consts.MV_COUNT_KEY_NAME {
		reply, e = i.handleBootReply(key)
		if e != nil {
			return "", e
		}
	}

	if reply == nil && key == consts.REQ_COUNT_KEY_NAME {
		reply, e = i.handleBootReply(key)
		if e != nil {
			return "", e
		}
	}
	if reply == nil {
		return "", errors.New("key: " + key + " not found")
	}

	return string(reply.([]byte)), nil
}

func (i *Implamintation) handleBootReply(key string) ([]byte, error) {
	_, e := i.conn.Do(consts.REDIS_SET, key, "0")
	if e != nil {
		return nil, e
	}
	reply, e := json.Marshal(0)
	if e != nil {
		return nil, e
	}
	return reply, nil
}

func (i *Implamintation) GetInt(key string) (int, error) {
	reply, e := i.conn.Do(consts.REDIS_GET, key)
	if e != nil {
		return 0, e
	}

	if reply == nil {
		_, e = i.conn.Do(consts.REDIS_SET, key, "0")
		if e != nil {
			return 0, e
		}
		reply, e = json.Marshal(0)
		if e != nil {
			return 0, e
		}
	}

	replyStr := string(reply.([]byte))
	replyInt, e := strconv.Atoi(replyStr)
	if e != nil {
		return 0, e
	}

	return replyInt, nil
}

func (i *Implamintation) Set(key, value string) error {
	_, e := i.conn.Do(consts.REDIS_SET, key, value)
	if e != nil {
		return e
	}

	return nil
}

func (i *Implamintation) Increment(key string) error {

	_, e := i.conn.Do(consts.REDIS_INC, key)
	if e != nil {
		if e.Error() == "WRONGTYPE Operation against a key holding the wrong kind of value" {
			_, e = i.conn.Do(consts.REDIS_SET, key, 0)
			if e != nil {
				return e
			} else {
				_, e = i.conn.Do(consts.REDIS_INC, key)
				if e != nil {
					return e
				}
			}
		} else {
			return e
		}
	}
	return nil
}
