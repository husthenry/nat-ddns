package client

import (
	"entity"
	"sync"
	"log"
)

type clientChannelService struct {
	channelMap map[string]entity.Channel
	lock sync.Mutex
}

var cscs *clientChannelService
var conce sync.Once

func GetCcsInstance() *clientChannelService {
	conce.Do(func() {
		cscs = &clientChannelService {
			channelMap:make(map[string]entity.Channel),
		}
	})
	return cscs
}

func (cscs *clientChannelService) IsContainsChannel(key string) bool {
	_, ok := cscs.channelMap[key]
	return ok
}

func (cscs *clientChannelService) AddChannel(channel entity.Channel) bool {
	cscs.lock.Lock()
	defer cscs.lock.Unlock()

	if cscs.IsContainsChannel(channel.Key) {
		panic("channel is exists!!!")
		return false
	}
	cscs.channelMap[channel.Key]=channel
	log.Println("ServerChannelService add key:", channel.Key)

	return true
}

func (cscs *clientChannelService) GetChannel(key string) entity.Channel {
	cscs.lock.Lock()
	defer cscs.lock.Unlock()

	return cscs.channelMap[key]
}

func (cscs *clientChannelService) RemoveChannel(key string) bool {
	cscs.lock.Lock()
	defer cscs.lock.Unlock()

	delete(cscs.channelMap, key)

	return true
}

func (cscs *clientChannelService) IsContainsSubChannel(parentKey, subKey string) bool {
	if !cscs.IsContainsChannel(parentKey) {
		log.Println("no parent channel:", parentKey, "for sub channel:", subKey)
		return false
	}

	parentChan := cscs.channelMap[parentKey]
	subMap := parentChan.SubChan

	_, ok := subMap[subKey]
	return ok
}

func (cscs *clientChannelService) RemoveSubChannel(parentKey, subKey string){
	cscs.lock.Lock()
	defer cscs.lock.Unlock()

	if !cscs.IsContainsChannel(parentKey) {
		log.Println("remove sub channel:", subKey, " err, no parent channel:", parentKey)
		return
	}

	subMap := cscs.channelMap[parentKey].SubChan

	delete(subMap, subKey)
}

func (cscs *clientChannelService) GetSubChannel(parentKey, subKey string) entity.Channel{
	if !cscs.IsContainsChannel(parentKey) {
		log.Println("get sub channel:", subKey, " err, no parent channel:", parentKey)
		return entity.Channel{}
	}

	subMap := cscs.channelMap[parentKey].SubChan
	if nil == subMap {
		log.Println("sub channel is empty!!!")
		return entity.Channel{}
	}
	return subMap[subKey]
}

func (cscs *clientChannelService) AddSubChannel(channel entity.Channel) bool {
	cscs.lock.Lock()
	defer cscs.lock.Unlock()

	if !cscs.IsContainsChannel(channel.Key) {
		log.Println("add sub channel:", channel.Uri, " err, no parent channel:", channel.Key)
		return false
	}

	parentChan := cscs.channelMap[channel.Key]

	subMap := parentChan.SubChan
	_, ok := subMap[channel.Uri]
	if ok  {
		log.Println("the sub channel is exists! sub channel:", channel.Uri)
		return false
	}

	subMap[channel.Uri] = channel

	return true
}
