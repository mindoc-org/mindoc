package cache

import (
	"github.com/astaxie/beego/cache"
	"time"
	"encoding/gob"
	"fmt"
	"bytes"
	"errors"
	"github.com/astaxie/beego"
)

var bm cache.Cache

func Get(key string,e interface{}) error {

	val := bm.Get(key)

	if val == nil {
		return errors.New("cache does not exist")
	}
	if b,ok := val.([]byte); ok {
		buf := bytes.NewBuffer(b)

		decoder := gob.NewDecoder(buf)

		err := decoder.Decode(e)

		if err != nil {
			fmt.Println("反序列化对象失败 ->", err)
		}
		return err
	}
	return errors.New("value is not []byte")
}

func GetMulti(keys []string) []interface{} {
	return bm.GetMulti(keys)
}

func Put(key string, val interface{}, timeout time.Duration) error {

	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(val)
	if err != nil {
		beego.Error("序列化对象失败 ->",err)
		return err
	}

	return bm.Put(key, buf.String(), timeout)
}

func Delete(key string) error {
	return bm.Delete(key)
}
func Incr(key string) error {
	return bm.Incr(key)
}
func Decr(key string) error {
	return bm.Decr(key)
}
func IsExist(key string) bool {
	return bm.IsExist(key)
}
func ClearAll() error{
	return bm.ClearAll()
}

func StartAndGC(config string) error {
	return bm.StartAndGC(config)
}
//初始化缓存
func Init(c cache.Cache)  {
	bm = c
}