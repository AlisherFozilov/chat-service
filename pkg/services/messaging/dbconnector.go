package messaging

import (
	"context"
	user "github.com/AlisherFozilov/db-storage/pkg/api"
	"log"
)

type RemoteURL string

type ConnectorService struct {
	client user.StorageClient
	url    string
}

func NewConnectorService(url RemoteURL, client user.StorageClient) *ConnectorService {
	return &ConnectorService{url: string(url), client: client}
}

const (
	textType  = 1
	mediaType = 2
)

func (s *ConnectorService) Send(messagesCash []CommonMessage) {
	for _, message := range messagesCash {
		var msgType int64
		switch message.Type {
		case Text:
			msgType = textType
		case Media:
			msgType = mediaType
		default:
			panic("situation is impossible! Inner program incoordination")
		}
		_, err := s.client.SaveMessages(context.Background(), &user.Messages{
			Id:         message.ID,
			Type:       msgType,
			SenderId:   message.SenderID,
			ReceiverId: message.ReceiverID,
			Data:       message.Data,
			Timestamp:  message.TimeStamp,
			Removed:    false,
		})
		if err != nil {
			log.Print(err)
		}
	}
}
