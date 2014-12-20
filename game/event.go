package game

// import (
// 	"encoding/json"
// )

// type EventType uint8

// const (
// 	EVENT_TYPE_ERROR = iota
// 	EVENT_TYPE_OBJECT_CREATING
// 	EVENT_TYPE_OBJECT_DELETING
// 	EVENT_TYPE_OBJECT_UPDATING
// )

// // Implementing json.Marshaler interface
// func (et EventType) MarshalJSON() ([]byte, error) {
// 	switch et {
// 	case EVENT_TYPE_ERROR:
// 		return []byte(`"error"`), nil
// 	case EVENT_TYPE_OBJECT_CREATING:
// 		return []byte(`"creating"`), nil
// 	case EVENT_TYPE_OBJECT_DELETING:
// 		return []byte(`"deleting"`), nil
// 	case EVENT_TYPE_OBJECT_UPDATING:
// 		return []byte(`"updating"`), nil
// 	}

// 	// Don't return error for undefined event type
// 	return []byte(`"undefined"`), nil
// }

// type Event struct {
// 	Type EventType       `json:"type"`
// 	Data json.RawMessage `json:"data"`
// }
