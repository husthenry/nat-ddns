package server

import (
	"fmt"
	"sync"
	"log"
	"entity"
)

type serverClientKeyService struct {
	sc entity.ServerConfig
	clientKeyMap map[string]string
	lock sync.Mutex
}

var sscks *serverClientKeyService
var konce sync.Once

func GetScksInstance() *serverClientKeyService {
	konce.Do(func() {
		sscks = &serverClientKeyService {
			clientKeyMap: make(map[string]string),
		}
	})
	return sscks
}

func (scks *serverClientKeyService) ServerClientKeyServiceInit(sc entity.ServerConfig)  {
	scks.sc = sc

	for i:=0; i<len(scks.sc.ClientKey); i++{
		item := scks.sc.ClientKey[i]
		for _, v := range item{
			scks.AddKey(v)
		}
	}
}

func (scks *serverClientKeyService) IsContainsKey(clientKey string) bool {
	_, ok := scks.clientKeyMap[clientKey]
	return ok
}

func (scks *serverClientKeyService) AddKey(clientKey string) bool {
	scks.lock.Lock()
	defer scks.lock.Unlock()

	if scks.IsContainsKey(clientKey) {
		fmt.Println("client key is exists!!!")
		return false
	}

	scks.clientKeyMap[clientKey] = clientKey

	log.Println("add client key:", clientKey, "success")
	return true
}

func (scks *serverClientKeyService) RemoveKey(clientKey string) bool {
	scks.lock.Lock()
	defer scks.lock.Unlock()

	delete(scks.clientKeyMap, clientKey)
	return true
}
