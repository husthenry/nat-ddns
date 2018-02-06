package server

import (
	"fmt"
	"sync"
	"util/cache"
)

//todo:接下来密钥管理中要抽象成一个实体，带密码管理，增加安全性

type serverClientKeyService struct {
}

var sscks *serverClientKeyService
var konce sync.Once

func GetScksInstance() *serverClientKeyService {
	konce.Do(func() {
		sscks = &serverClientKeyService {}
	})
	return sscks
}


func (scks *serverClientKeyService) IsContainsKey(clientKey string) bool {
	return cache.IsContains(clientKey)
}

func (scks *serverClientKeyService) AddKey(clientKey string) bool {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	if cache.IsContains(clientKey) {
		fmt.Println("client key is exists!!!")
		return false
	}

	cache.Add(clientKey, clientKey)
	return true
}

func (scks *serverClientKeyService) RemoveKey(clientKey string) bool {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	cache.Remove(clientKey)
	return true
}
