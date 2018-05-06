package connections

type OutputMessageType uint8

const (
	OutputMessageTypeGameEvent OutputMessageType = iota
)

type OutputMessage struct {
	Type    OutputMessageType
	Payload interface{}
}
