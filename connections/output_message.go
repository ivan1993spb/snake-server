package connections

type OutputMessageType uint8

const (
	OutputMessageTypeGame OutputMessageType = iota
	OutputMessageTypePlayer
	OutputMessageTypeBroadcast
)

var outputMessageTypeLabels = map[OutputMessageType]string{
	OutputMessageTypeGame:      "game",
	OutputMessageTypePlayer:    "player",
	OutputMessageTypeBroadcast: "broadcast",
}

func (t OutputMessageType) String() string {
	if label, ok := outputMessageTypeLabels[t]; ok {
		return label
	}
	return "unknown"
}

var outputMessageTypeJSONs = map[OutputMessageType][]byte{
	OutputMessageTypeGame:      []byte(`"game"`),
	OutputMessageTypePlayer:    []byte(`"player"`),
	OutputMessageTypeBroadcast: []byte(`"broadcast"`),
}

func (t OutputMessageType) MarshalJSON() ([]byte, error) {
	if json, ok := outputMessageTypeJSONs[t]; ok {
		return json, nil
	}
	return []byte(`"unknown"`), nil
}

//go:generate ffjson $GOFILE

// ffjson: nodecoder
type OutputMessage struct {
	Type    OutputMessageType `json:"type"`
	Payload interface{}       `json:"payload"`
}
