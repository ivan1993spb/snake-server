package game

import "encoding/json"

type eventType uint8

const (
	EVENT_TYPE_ERROR = iota
	EVENT_TYPE_OBJECT_CREATING
	EVENT_TYPE_OBJECT_DELETING
	EVENT_TYPE_OBJECT_UPDATING
)

// Implementing json.Marshaler interface
func (et eventType) MarshalJSON() ([]byte, error) {
	switch et {
	case EVENT_TYPE_ERROR:
		return []byte(`"error"`), nil
	case EVENT_TYPE_OBJECT_CREATING:
		return []byte(`"creating"`), nil
	case EVENT_TYPE_OBJECT_DELETING:
		return []byte(`"deleting"`), nil
	case EVENT_TYPE_OBJECT_UPDATING:
		return []byte(`"updating"`), nil
	}

	// Don't return error for undefined event type
	return []byte(`"undefined"`), nil
}

type event struct {
	Type eventType      `json:"type"`
	Data json.Marshaler `json:"data"`
}
