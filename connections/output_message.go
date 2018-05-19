package connections

type OutputMessageType uint8

const (
	OutputMessageTypeGameEvent OutputMessageType = iota
	OutputMessageTypeGroupNotice
	OutputMessageTypeConnectionNotice
)

var outputMessageTypeLabels = map[OutputMessageType]string{
	OutputMessageTypeGameEvent:        "game_event",
	OutputMessageTypeGroupNotice:      "group_notice",
	OutputMessageTypeConnectionNotice: "connection_notice",
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
