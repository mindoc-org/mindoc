package cache

import (
	"context"
	"time"
)

type NullCache struct {

}


func (bm *NullCache) Get(ctx context.Context, key string) (interface{}, error) {
	return nil, nil
}

func (bm *NullCache)GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	return nil, nil
}

func (bm *NullCache)Put(ctx context.Context,key string, val interface{}, timeout time.Duration) error {
	return nil
}
func (bm *NullCache)Delete(ctx context.Context,key string) error {
	return nil
}
func (bm *NullCache)Incr(ctx context.Context,key string) error {
	return nil
}
func (bm *NullCache)Decr(ctx context.Context,key string) error {
	return nil
}
func (bm *NullCache)IsExist(ctx context.Context,key string) (bool, error) {
	return false, nil
}
func (bm *NullCache)ClearAll(ctx context.Context) error{
	return nil
}

func (bm *NullCache)StartAndGC(config string) error {
	return nil
}
