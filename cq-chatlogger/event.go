package cqchat

type MetaEvent struct {
	MetaEventType string `json:"meta_event_type"`
}

type LifeCycleEvent struct {
	MetaPost
	MetaEvent
	MetaEventType string `json:"meta_event_type"`
}

type HeartbeatEvent struct {
	MetaPost
	MetaEvent
	MetaEventType string `json:"meta_event_type"`
	Interval      int    `json:"interval"`
	Status        map[string]interface{}
}

func ParseEventFromData(data []byte) (interface{}, error) {
	return nil, nil
	// todo
}
