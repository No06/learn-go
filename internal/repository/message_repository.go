package repository

import "hinoob.net/learn-go/internal/model"

// CreateMessage saves a new message to the database.
func CreateMessage(message *model.Message) error {
	return DB.Create(message).Error
}

// GetMessagesBetweenUsers retrieves the message history between two users.
func GetMessagesBetweenUsers(user1ID, user2ID uint) ([]model.Message, error) {
	var messages []model.Message
	err := DB.Where("(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
		user1ID, user2ID, user2ID, user1ID).
		Order("sent_at asc").
		Find(&messages).Error
	return messages, err
}
