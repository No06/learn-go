package model

import (
	"gorm.io/gorm"
	"time"
)

// MessageType defines the type of content in a message.
type MessageType string

const (
	MsgTypeText  MessageType = "text"
	MsgTypeImage MessageType = "image"
	MsgTypeVideo MessageType = "video"
	MsgTypeVoice MessageType = "voice"
)

// Message represents a single chat message in the system.
type Message struct {
	gorm.Model
	SenderID    uint   `gorm:"not null"`
	RecipientID uint   `gorm:"not null"`
	Content     string `gorm:"type:text;not null"` // For text, this is the message. For files, this is the URL.
	MessageType MessageType `gorm:"type:varchar(10);not null"`
	SentAt      time.Time `gorm:"autoCreateTime"`

	Sender    User `gorm:"foreignKey:SenderID"`
	Recipient User `gorm:"foreignKey:RecipientID"`
}
