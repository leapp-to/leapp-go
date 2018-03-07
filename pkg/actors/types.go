package actors

import (
	"encoding/json"
	"time"
)

// Host specifies the host the message comes from
type Host struct {
	Context  string `json:"context"`
	Hostname string `json:"hostname"`
}

// DataSource specifies the origin of the message
type DataSource struct {
	Host
	Actor string `json:"actor"`
	Phase string `json:"phase"`
}

// MessageData contains the data of the message and its hash
type MessageData struct {
	Hash string `json:"hash"`
	Data string `json:"data"`
}

// Message is a channel message
type Message struct {
	DataSource
	ID      int64       `json:"id"`
	Stamp   time.Time   `json:"stamp"`
	Type    string      `json:"type"`
	Channel string      `json:"channel"`
	Message MessageData `json:"message"`
}

// LogMessage a log message sent
type LogMessage struct {
	DataSource
	Stamp time.Time `json:"stamp"`
	Log   struct {
		Level   string `json:"level"`
		Message string `json:"message"`
	} `json:"log"`
}

// Audit message
type Audit struct {
	DataSource
	Event string    `json:"event"`
	Stamp time.Time `json:"stamp"`

	Message *Message `json:"message,omitempty"`
	Data    *string  `json:"data,omitempty"`
}

func convertMessage(msg *Message) *Audit {
	return &Audit{
		DataSource: msg.DataSource,
		Message:    msg,
		Event:      "new-message",
		Stamp:      msg.Stamp,
	}
}

func convertLogMessage(log LogMessage) *Audit {
	data, _ := json.Marshal(&log.Log)
	dataString := new(string)
	*dataString = string(data)
	return &Audit{
		DataSource: DataSource{
			Host: Host{
				Hostname: log.Hostname,
				Context:  log.Context,
			},
			Actor: log.Actor,
			Phase: log.Phase,
		},
		Event: "log-message",
		Stamp: log.Stamp,
		Data:  dataString,
	}
}
