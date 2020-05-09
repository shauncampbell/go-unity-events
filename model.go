package events

type EventType string
type EventUnit string

type UnityEvent struct {
	DeviceId string `json:"d"`
	EventType EventType `json:"e"`
	Timestamp int64 `json:"t"`
	Units EventUnit `json:"u"`
	Value string `json:"v"`
}

type Config struct {
	ConsumerName string
	RabbitHost string
	RabbitPort int
	RabbitUser string
	RabbitPass string
	RabbitQueue string
}