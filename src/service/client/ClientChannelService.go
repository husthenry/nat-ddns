package client

import (
	"entity"
	"sync"
	"log"
)

type clientChannelService struct {

}

var channelMap = make(map[string]entity.Channel)

var cscs *clientChannelService
var conce sync.Once

func GetCcsInstance() *clientChannelService {
	conce.Do(func() {
		cscs = &clientChannelService {}
	})
	return cscs
}

func (cscs *clientChannelService) IsContainsChannel(key string) bool {
	_, ok := channelMap[key]
	return ok
}

func (cscs *clientChannelService) AddChannel(channel entity.Channel) bool {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	if cscs.IsContainsChannel(channel.Key) {
		panic("channel is exists!!!")
		return false
	}
	channelMap[channel.Key]=channel
	log.Println("ServerChannelService add key:", channel.Key)

	return true
}

func (cscs *clientChannelService) GetChannel(key string) entity.Channel {
	return channelMap[key]
}

func (cscs *clientChannelService) RemoveChannel(key string) bool {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	delete(channelMap, key)

	return true
}
