package app

import (
	"log"
	"net/http"
)

func (s *Server) handleMain() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("OK"))
	}
}

func (s *Server) handleConnection() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Print("handle connection")
		conn, err := s.upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Print("can't upgrade connection: ", err)
			return
		}

		go func() {
			log.Print("serve user")
			err := s.messagingSvc.ServeUser(conn)
			if err != nil {
				log.Print(err)
			}
		}()
	}
}
