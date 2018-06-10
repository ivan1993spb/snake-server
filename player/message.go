package player

type MessageType uint8

const (
	MessageTypeSize MessageType = iota
	MessageTypeSnake
	MessageTypeNotice
	MessageTypeError
	MessageTypeCountdown
	MessageTypeObjects
)

var messageTypeJSONs = map[MessageType][]byte{
	MessageTypeSize:      []byte(`"size"`),
	MessageTypeSnake:     []byte(`"snake"`),
	MessageTypeNotice:    []byte(`"notice"`),
	MessageTypeError:     []byte(`"error"`),
	MessageTypeCountdown: []byte(`"countdown"`),
	MessageTypeObjects:   []byte(`"objects"`),
}

func (t MessageType) MarshalJSON() ([]byte, error) {
	if jsonBytes, ok := messageTypeJSONs[t]; ok {
		return jsonBytes, nil
	}
	return []byte(`"unknown"`), nil
}

var messageTypeLabels = map[MessageType]string{
	MessageTypeSize:      "size",
	MessageTypeSnake:     "snake",
	MessageTypeNotice:    "notice",
	MessageTypeError:     "error",
	MessageTypeCountdown: "countdown",
	MessageTypeObjects:   "objects",
}

func (t MessageType) String() string {
	if label, ok := messageTypeLabels[t]; ok {
		return label
	}
	return "unknown"
}

type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

type MessageSize struct {
	Width  uint8 `json:"width"`
	Height uint8 `json:"height"`
}

func NewMessageSize(w, h uint8) Message {
	return Message{
		Type: MessageTypeSize,
		Payload: MessageSize{
			Width:  w,
			Height: h,
		},
	}
}

type MessageSnake string

func NewMessageSnake(uuid string) Message {
	return Message{
		Type:    MessageTypeSnake,
		Payload: MessageSnake(uuid),
	}
}

type MessageNotice string

func NewMessageNotice(notice string) Message {
	return Message{
		Type:    MessageTypeNotice,
		Payload: MessageNotice(notice),
	}
}

type MessageError string

func NewMessageError(errorStr string) Message {
	return Message{
		Type:    MessageTypeError,
		Payload: MessageError(errorStr),
	}
}

type MessageCountdown uint

func NewMessageCountdown(countdown uint) Message {
	return Message{
		Type:    MessageTypeCountdown,
		Payload: MessageCountdown(countdown),
	}
}

type MessageObjects []interface{}

func NewMessageObjects(objects []interface{}) Message {
	return Message{
		Type:    MessageTypeObjects,
		Payload: MessageObjects(objects),
	}
}
