package gopool

import (
	"sync"
	"errors"
	"fmt"
)
var (
	ErrHandlerIsExist = errors.New("指定的键已存在")
	ErrWorkerChanClosed = errors.New("队列已关闭")
)
type ChannelHandler func()

type entry struct {
	handler ChannelHandler
	key string
}

type ChannelPool struct {
	maxWorkerNum int
	maxPoolNum int
	wait *sync.WaitGroup
	cache *sync.Map
	worker chan *entry
	limit chan bool
	isClosed bool
	once *sync.Once
}

func NewChannelPool(maxWorkerNum, maxPoolNum int) (*ChannelPool) {
	if maxWorkerNum <= 0 {
		maxWorkerNum = 1
	}
	if maxPoolNum <= 0 {
		maxWorkerNum = 100
	}
	return &ChannelPool{
		maxWorkerNum: maxWorkerNum,
		maxPoolNum: maxPoolNum,
		wait: &sync.WaitGroup{},
		cache: &sync.Map{},
		worker: make(chan  *entry, maxWorkerNum),
		limit: make(chan bool, maxWorkerNum),
		isClosed: false,
		once: &sync.Once{},
	}
}

func (pool *ChannelPool) LoadOrStore(key string,value ChannelHandler) error  {
	if pool.isClosed {
		return ErrWorkerChanClosed
	}
	if _,loaded := pool.cache.LoadOrStore(key,false); loaded {
		return ErrHandlerIsExist
	}else{
		pool.worker <- &entry{handler:value,key:key}
		return  nil
	}
}

func (pool *ChannelPool) Start() {
	pool.once.Do(func() {
		go func() {
			for i :=0; i < pool.maxWorkerNum; i ++ {
				pool.limit <- true
			}
			for {
				actual, isClosed := <-pool.worker
				//当队列被关闭，则跳出循环
				if actual == nil && !isClosed {
					fmt.Println("工作队列已关闭")
					break
				}
				limit := <-pool.limit

				if limit {
					pool.wait.Add(1)
					go func(actual *entry) {
						defer func() {
							pool.cache.Delete(actual.key)
							pool.limit <- true
							pool.wait.Done()
						}()

						actual.handler()

					}(actual)
				}
			}
		}()
	})
}

func (pool *ChannelPool) Wait() {
	close(pool.worker)

	pool.wait.Wait()
}

