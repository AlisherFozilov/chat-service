package messaging

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type Service struct {
	usersMu        sync.RWMutex
	connWriteMu    sync.Mutex
	onlineUsers    map[int64]*websocket.Conn
	dbSvcConnector *ConnectorService
	saverCh        chan CommonMessage
}
// Mutex vs Channel
func NewService(dbSvcConnector *ConnectorService) *Service {
	return &Service{
		dbSvcConnector: dbSvcConnector,
		saverCh:        make(chan CommonMessage, 1024),
		onlineUsers:    map[int64]*websocket.Conn{},
	}
}

func (s *Service) Start() {
	go s.saverToDbService()
}

func (s *Service) saverToDbService() {
	const numberOfMessages = 1
	cash := make([]CommonMessage, 0, numberOfMessages)
	for {
		for i := 0; i < numberOfMessages; i++ {
			message := <-s.saverCh
			cash = append(cash, message)
		}
		s.dbSvcConnector.Send(cash)

		cash = cash[:0]
	}
}

func (s *Service) ServeUser(conn *websocket.Conn) (err error) {

	_, message, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	//TODO: take id from JWT token
	idMarshalled := IdJSON{}

	err = json.Unmarshal(message, &idMarshalled)
	if err != nil {
		return err
	}
	id := idMarshalled.Id

	s.usersMu.RLock()
	if someoneWithThisID, exists := s.onlineUsers[id]; exists {
		_ = someoneWithThisID
		panic("situation is impossible! Inner program incoordination")
	}
	s.usersMu.RUnlock()

	s.usersMu.Lock()
	s.onlineUsers[id] = conn
	s.usersMu.Unlock()

	defer func() {
		s.usersMu.Lock()
		delete(s.onlineUsers, id)
		s.usersMu.Unlock()

		errdefer := conn.Close()
		if errdefer != nil {
			log.Print("can't close connection: ", errdefer)
		}
	}()

	msg := Message{}
	commonMessage := CommonMessage{}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		err = json.Unmarshal(message, &msg)
		if err != nil {
			return err
		}

		switch msg.Type {
		case Text:
			err := s.sendCommonMessage(message, id, &commonMessage)
			if err != nil {
				return err
			}
		case Media:
			err := s.sendCommonMessage(message, id, &commonMessage)
			if err != nil {
				return err
			}
		case QueryUsersList:
			err := s.sendUsersList(id)
			if err != nil {
				return err
			}
		default:
			return errors.New("wrong message type")
		}
	}
}

func (s *Service) sendCommonMessage(message []byte, senderID int64, msg *CommonMessage) (err error) {
	err = json.Unmarshal(message, msg)
	if err != nil {
		return err
	}

	uuidV4 := uuid.New().String()
	msg.ID = uuidV4
	msg.SenderID = senderID
	msg.TimeStamp = time.Now().Unix()

	s.saverCh <- *msg

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = s.sendBytesToUserById(msgJSON, msg.ReceiverID)
	if err != nil {
		return err
	}

	return
}

func (s *Service) sendBytesToUserById(message []byte, receiverID int64) (err error) {

	s.usersMu.RLock()
	conn, ok := s.onlineUsers[receiverID]
	s.usersMu.RUnlock()

	if !ok {
		return nil //user is offline, not an error
	}
	s.connWriteMu.Lock()
	err = conn.WriteMessage(websocket.TextMessage, message)
	s.connWriteMu.Unlock()

	if err != nil {
		return err
	}
	return
}

func (s *Service) sendUsersList(id int64) (err error) {
	onlineUsersId := make([]int64, 0, len(s.onlineUsers))
	for id := range s.onlineUsers {
		onlineUsersId = append(onlineUsersId, id)
	}
	data, err := json.Marshal(onlineUsersId)
	if err != nil {
		return err
	}
	err = s.sendBytesToUserById(data, id)
	if err != nil {
		return err
	}
	return nil
}
