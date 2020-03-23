package app

import (
	"github.com/AlisherFozilov/chat-service/pkg/services/messaging"
	"github.com/AlisherFozilov/mymux/pkg/exactmux"
	"github.com/gorilla/websocket"
	"net/http"
)

type Server struct {
	router       *exactmux.ExactMux
	upgrader     *websocket.Upgrader
	messagingSvc *messaging.Service
}

func NewServer(router *exactmux.ExactMux, upgrader *websocket.Upgrader, messagingSvc *messaging.Service) *Server {
	if router == nil {
		panic("router must be not nil")
	}
	if upgrader == nil {
		panic("upgrader must be not nil")
	}
	if messagingSvc == nil {
		panic("messagingSvc must be not nil")
	}

	return &Server{router: router, upgrader: upgrader, messagingSvc: messagingSvc}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) InitRoutes() {
	mux := s.router
	mux.GET("/", s.handleMain())
	mux.GET("/ws", s.handleConnection())
}

func (s *Server) Start() {
	s.InitRoutes()
}
