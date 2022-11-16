package main

import (
	"fmt"
	"strings"
)

type Handler interface {
	Handler(user *User, context string)
}

type HandlerStrategy struct {
	MessageHandlerMap map[string]*Handler
}

type WhoHandler struct {
}

func (this *WhoHandler) Handler(user *User, context string) {
	server := user.server
	server.mapLock.Lock()

	for _, u := range server.OnlineMap {
		onlineMsg := fmt.Sprintf("[%s]%s:online\n", u.Addr, u.Name)
		if user.Addr == u.Addr {
			onlineMsg = "*" + onlineMsg
		} else {
			onlineMsg = " " + onlineMsg
		}
		user.SendMsg(onlineMsg)
	}
	server.mapLock.Unlock()
}

type BroadCastHandler struct {
}

func (this *BroadCastHandler) Handler(user *User, context string) {
	user.server.BroadCast(user, context)
}

type RenameHandler struct {
}

func (this *RenameHandler) Handler(user *User, context string) {
	server := user.server
	if user.Name == context {
		user.SendMsg("No modification required")
	} else if _, ok := server.OnlineMap[context]; ok {
		user.SendMsg("User name already exists")
	} else {
		server.mapLock.Lock()
		delete(server.OnlineMap, user.Name)
		user.Name = context
		server.OnlineMap[user.Name] = user
		server.mapLock.Unlock()
		user.SendMsg(fmt.Sprintf("Your user name is changed to : %s", user.Name))
	}
}

// MessageHandlerStrategyFactory 策略工厂
type MessageHandlerStrategyFactory struct {
	strategys map[string]Handler
}

func NewStrategyFactory() *MessageHandlerStrategyFactory {
	factory := new(MessageHandlerStrategyFactory)
	//初始化 内部的策略
	var strategys = make(map[string]Handler, 2)
	whoHandler := new(WhoHandler)
	broadCastHandler := new(BroadCastHandler)
	renameHandler := new(RenameHandler)
	strategys["who"] = whoHandler
	strategys["bc"] = broadCastHandler
	strategys["rename"] = renameHandler
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
