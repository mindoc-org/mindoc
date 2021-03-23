package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/client/cache"
)

var bm cache.Cache

var nilctx = context.TODO()

func Get(key string, e interface{}) error {

	val, err := bm.Get(nilctx, key)

	if err != nil {
		return errors.New("get cache error:"+ err.Error())
	}

	if val == nil {
		return errors.New("cache does not exist")
	}
	if b, ok := val.([]byte); ok {
		buf := bytes.NewBuffer(b)

		decoder := gob.NewDecoder(buf)

		err := decoder.Decode(e)

		if err != nil {
			logs.Error("反序列化对象失败 ->", err)
		}
		return err
	} else if s, ok := val.(string); ok && s != "" {

		buf := bytes.NewBufferString(s)

		decoder := gob.NewDecoder(buf)

		err := decoder.Decode(e)

		if err != nil {
			logs.Error("反序列化对象失败 ->", err)
		}
		return err
	}
	return errors.New("value is not []byte or string")
}

func Put(key string, val interface{}, timeout time.Duration) error {

	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(val)
	if err != nil {
		logs.Error("序列化对象失败 ->", err)
		return err
	}

	return bm.Put(nilctx, key, buf.String(), timeout)
}

func Delete(key string) error {
	return bm.Delete(nilctx, key)
}
func Incr(key string) error {
	return bm.Incr(nilctx, key)
}
func Decr(key string) error {
	return bm.Decr(nilctx, key)
}
func IsExist(key string) (bool, error) {
	return bm.IsExist(nilctx, key)
}
func ClearAll() error {
	return bm.ClearAll(nilctx)
}

func StartAndGC(config string) error {
	return bm.StartAndGC(config)
}

//Init will initialize cache
func Init(c cache.Cache) {
	bm = c
}
