package realtime

import "golang.org/x/net/websocket"

/**
 * Интерфейс Задачника
 */
type AbstractModel interface {
	Create() bool
	Update() bool
	Delete() bool
	GetId(id int64) bool
}

type AbstractServiceEngine interface {
	Start(ws *websocket.Conn)
	NewService(ws *websocket.Conn)
	DeleteService(ws *websocket.Conn)
}
