package handler

import "sync"

var (
	serverHandler *ServerHandler
	once          sync.Once
)

type ServerHandler struct {
	*AlbumHandler
	*UserHandler
}

func NewHandler() *ServerHandler {
	once.Do(func() {
		serverHandler = &ServerHandler{}
	})
	return serverHandler
}

func (h *ServerHandler) Register(i interface{}) *ServerHandler {
	switch interfaceType := i.(type) {
	case *AlbumHandler:
		serverHandler.AlbumHandler = interfaceType
	case *UserHandler:
		serverHandler.UserHandler = interfaceType
	}
	return serverHandler
}
