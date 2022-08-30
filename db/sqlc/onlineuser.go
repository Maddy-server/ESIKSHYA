package db

import "github.com/gorilla/websocket"

type OnlineUserMaps interface {
	Add(userID int64, conn *websocket.Conn)
	Get(userID int64) (*websocket.Conn, error)
	Disconnect(userID int64) error
	IsOnline(userID int64) bool
}

var _ OnlineUserMaps = (*OnlineUserMap)(nil)
