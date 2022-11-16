package main

import (
	"fmt"
	"strings"
)

type MessageHandler interface {
	Handler(user *User, context string)
}

type MessageHandlerStrategy struct {
	MessageHandlerMap map[string]*MessageHandler
}

type WhoMessageHandler struct {
}

func (this *WhoMessageHandler) Handler(user *User, context string) {
	server := user.server
	server.mapLock.Lock()

	for _, user := range server.OnlineMap {
		onlineMsg := fmt.Sprintf("[%s]%s:online\n", user.Addr, user.Name)
		user.SendMsg(onlineMsg)
	}
	server.mapLock.Unlock()
}

type BroadCastMessageHandler struct {
}

func (this *BroadCastMessageHandler) Handler(user *User, context string) {
	user.server.BroadCast(user, context)
}

// MessageHandlerStrategyFactory 策略工厂
type MessageHandlerStrategyFactory struct {
	strategys map[string]MessageHandler
}

func NewStrategyFactory() *MessageHandlerStrategyFactory {
	factory := new(MessageHandlerStrategyFactory)
	//初始化 内部的策略
	var strategys = make(map[string]MessageHandler, 2)
	quickSort := new(WhoMessageHandler)
	bubbleSort := new(BroadCastMessageHandler)
	strategys["who"] = quickSort
	strategys["bc"] = bubbleSort
	factory.strategys = strategys
	return factory
}

// HandlerMessageStrategy 策略工厂提供该方法，客户端通过策略处理消息。
func (factory *MessageHandlerStrategyFactory) HandlerMessageStrategy(name string, user *User, context string) {
	if len(name) == 0 {
		return
	}
	name = strings.ToLower(name)
	if v, ok := factory.strategys[name]; ok {
		v.Handler(user, context)
		return
	}
}
