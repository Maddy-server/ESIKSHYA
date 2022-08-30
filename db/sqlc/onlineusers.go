package db

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type OnlineUserMap struct {
	mutex       sync.Mutex
	onlineUsers map[int64]*websocket.Conn
}

func newOnlineUserMap() OnlineUserMap {
	return OnlineUserMap{
		mutex:       sync.Mutex{},
		onlineUsers: make(map[int64]*websocket.Conn),
	}
}

func (ou *OnlineUserMap) Add(userID int64, conn *websocket.Conn) {
	ou.mutex.Lock()
	defer ou.mutex.Unlock()
	c, ok := ou.onlineUsers[userID]
	if ok {
		c.Close()
	}
	ou.onlineUsers[userID] = conn
	conn.SetCloseHandler(func(code int, text string) error {
		ou.mutex.Lock()
		delete(ou.onlineUsers, userID)
		ou.mutex.Unlock()
		return nil
	})
}

func (ou *OnlineUserMap) Get(userID int64) (*websocket.Conn, error) {
	if ou.onlineUsers == nil {
		logrus.Error("onlineUsers not initialised")
		return nil, errors.New("onlineUsers not initialised")
	}
	c, ok := ou.onlineUsers[userID]
	if !ok {
		return nil, errors.New("connection not found")
	}
	return c, nil
}
func (ou *OnlineUserMap) GetUsers(userID int64) (*websocket.Conn, error) {
	if ou.onlineUsers == nil {
		logrus.Error("onlineUsers not initialised")
		return nil, errors.New("onlineUsers not initialised")
	}
	c, ok := ou.onlineUsers[userID]
	if !ok {
		return nil, errors.New("connection not found")
	}
	return c, nil
}

func (ou *OnlineUserMap) Disconnect(userID int64) error {
	if ou.onlineUsers == nil {
		logrus.Error("onlineUsers not initialised")
		return errors.New("onlineUsers not initialised")
	}
	delete(ou.onlineUsers, userID)
	return nil
}

func (ou *OnlineUserMap) IsOnline(userID int64) bool {
	if ou.onlineUsers == nil {
		logrus.Error("onlineUsers not initialised")
		return false
	}
	_, isOnline := ou.onlineUsers[userID]
	return isOnline
}
