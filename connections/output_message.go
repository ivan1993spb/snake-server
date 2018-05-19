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

type OutputMessage struct {
	Type    OutputMessageType
	Payload interface{}
}
