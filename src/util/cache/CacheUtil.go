package cache

import "sync"


//通用缓存
var cacheMap = make(map[string]interface{})

func Add(k string,t interface{}) {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	cacheMap[k]=t
}

func Get(k string) interface{} {
	return cacheMap[k]
}

func IsContains(k string) bool  {
	_, ok := cacheMap[k]
	return ok
}

func Remove(k string)  {
	delete(cacheMap, k)
}