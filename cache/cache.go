package cache

import (
	"github.com/astaxie/beego/cache"
	"time"
)

var bm cache.Cache

func Get(key string) interface{} {

	return bm.Get(key)
}

func GetMulti(keys []string) []interface{} {

	return bm.GetMulti(keys)
}

func Put(key string, val interface{}, timeout time.Duration) error {

	return bm.Put(key, val, timeout)
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