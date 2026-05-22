package chat

type MessagePayload struct {
	Text string `json:"text"`
}

type WSMessage struct {
	Type      string         `json:"type"`
	Payload   MessagePayload `json:"payload"`
	Timestamp string         `json:"timestamp"`
}
