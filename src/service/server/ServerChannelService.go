package server

import (
	"entity"
	"sync"
	"log"
)

type serverChannelService struct {
}

var channelMap = make(map[string]entity.Channel)

var sscs *serverChannelService
var conce sync.Once

func GetScsInstance() *serverChannelService {
	conce.Do(func() {
		sscs = &serverChannelService {}
	})
	return sscs
}

func (scs *serverChannelService) IsContainsChannel(key string) bool {
	_, ok := channelMap[key]
	return ok
}

func (scs *serverChannelService) AddChannel(channel entity.Channel) bool {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	if scs.IsContainsChannel(channel.Key) {
		panic("channel is exists!!!")
		return false
	}
	channelMap[channel.Key]=channel
	log.Println("ServerChannelService add key:", channel.Key)

	return true
}

func (scs *serverChannelService) GetChannel(key string) entity.Channel {
	return channelMap[key]
}

func (scs *serverChannelService) RemoveChannel(key string) bool {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	delete(channelMap, key)

	return true
}

func (scs *serverChannelService) IsContainsSubChannel(parentKey, subKey string) bool {
	if !scs.IsContainsChannel(parentKey) {
		log.Println("no parent channel:", parentKey, "for sub channel:", subKey)
		return false
	}

	parentChan := channelMap[parentKey]
	subMap := parentChan.SubChan

	_, ok := subMap[subKey]
	return ok
}

func (scs *serverChannelService) RemoveSubChannel(parentKey, subKey string){
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	if !scs.IsContainsChannel(parentKey) {
		log.Println("remove sub channel:", subKey, " err, no parent channel:", parentKey)
		return
	}

	subMap := channelMap[parentKey].SubChan

	delete(subMap, subKey)
}

func (scs *serverChannelService) GetSubChannel(parentKey, subKey string) entity.Channel{
	if !scs.IsContainsChannel(parentKey) {
		log.Println("get sub channel:", subKey, " err, no parent channel:", parentKey)
		return entity.Channel{}
	}

	subMap := channelMap[parentKey].SubChan
	if nil == subMap {
		log.Println("sub channel is empty!!!")
		return entity.Channel{}
	}
	return subMap[subKey]
}

func (scs *serverChannelService) AddSubChannel(channel entity.Channel) bool {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	if !scs.IsContainsChannel(channel.Key) {
		log.Println("add sub channel:", channel.Uri, " err, no parent channel:", channel.Key)
		return false
	}

	parentChan := channelMap[channel.Key]

	subMap := parentChan.SubChan
	_, ok := subMap[channel.Uri]
	if ok  {
		log.Println("the sub channel is exists! sub channel:", channel.Uri)
		return false
	}

	subMap[channel.Uri] = channel

	return true
}