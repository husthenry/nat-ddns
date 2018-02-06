package server

import (
	"entity"
	"sync"
	"log"
)

type serverChannelService struct {
	channelMap map[string]entity.Channel
	lock sync.Mutex
}

var sscs *serverChannelService
var conce sync.Once

func GetScsInstance() *serverChannelService {
	conce.Do(func() {
		sscs = &serverChannelService {
			channelMap:make(map[string]entity.Channel),
		}
	})
	return sscs
}

func (scs *serverChannelService) IsContainsChannel(key string) bool {
	_, ok := scs.channelMap[key]
	return ok
}

func (scs *serverChannelService) AddChannel(channel entity.Channel) bool {
	scs.lock.Lock()
	defer scs.lock.Unlock()

	if scs.IsContainsChannel(channel.Key) {
		panic("channel is exists!!!")
		return false
	}
	scs.channelMap[channel.Key]=channel
	log.Println("ServerChannelService add key:", channel.Key)

	return true
}

func (scs *serverChannelService) GetChannel(key string) entity.Channel {
	scs.lock.Lock()
	defer scs.lock.Unlock()

	return scs.channelMap[key]
}

func (scs *serverChannelService) RemoveChannel(key string) bool {
	scs.lock.Lock()
	defer scs.lock.Unlock()

	delete(scs.channelMap, key)

	return true
}

func (scs *serverChannelService) IsContainsSubChannel(parentKey, subKey string) bool {

	if !scs.IsContainsChannel(parentKey) {
		log.Println("no parent channel:", parentKey, "for sub channel:", subKey)
		return false
	}

	parentChan := scs.channelMap[parentKey]
	subMap := parentChan.SubChan

	_, ok := subMap[subKey]
	return ok
}

func (scs *serverChannelService) RemoveSubChannel(parentKey, subKey string){
	scs.lock.Lock()
	defer scs.lock.Unlock()

	if !scs.IsContainsChannel(parentKey) {
		log.Println("remove sub channel:", subKey, " err, no parent channel:", parentKey)
		return
	}

	subMap := scs.channelMap[parentKey].SubChan

	delete(subMap, subKey)
}

func (scs *serverChannelService) GetSubChannel(parentKey, subKey string) entity.Channel{
	scs.lock.Lock()
	defer scs.lock.Unlock()

	if !scs.IsContainsChannel(parentKey) {
		log.Println("get sub channel:", subKey, " err, no parent channel:", parentKey)
		return entity.Channel{}
	}

	subMap := scs.channelMap[parentKey].SubChan
	if nil == subMap {
		log.Println("sub channel is empty!!!")
		return entity.Channel{}
	}
	return subMap[subKey]
}

func (scs *serverChannelService) AddSubChannel(channel entity.Channel) bool {
	scs.lock.Lock()
	defer scs.lock.Unlock()

	if !scs.IsContainsChannel(channel.Key) {
		log.Println("add sub channel:", channel.Uri, " err, no parent channel:", channel.Key)
		return false
	}

	parentChan := scs.channelMap[channel.Key]

	subMap := parentChan.SubChan
	//_, ok := subMap[channel.Uri]
	//if ok  {
	//	log.Println("the sub channel is exists! sub channel:", channel.Uri)
	//	return false
	//}

	subMap[channel.Uri] = channel

	return true
}