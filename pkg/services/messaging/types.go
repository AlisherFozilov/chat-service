package messaging

type IdJSON struct {
	Id int64
}

const (
	Text           = "text"
	Media          = "media"
	QueryUsersList = "users_list"
)

type Message struct {
	Type string `json:"type"`
}

type CommonMessage struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	SenderID   int64  `json:"sender_id"`
	ReceiverID int64  `json:"receiver_id"`
	Data       string `json:"text"`
	TimeStamp  int64  `json:"time_stamp"`
}

type QueryAllUsersListMessage struct {
	Type string `json:"type"`
}

// TODO: use it
type QuerySetUUIDtoSentMessage struct {
	Type string `json:"type"`
	Uuid string `json:"uuid"`
}
