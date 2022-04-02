package messenger

import "context"

type MessagePostgreSql struct {
	db PostgreSqlDB
}

func NewMessagePostgreSql(db PostgreSqlDB) *MessagePostgreSql {
	return &MessagePostgreSql{
		db: db,
	}
}

func (m *MessagePostgreSql) AddMessage(ctx context.Context, msg Message) error {
	_, err := m.db.Conn.ExecContext(ctx,
		"INSERT INTO messages (id, chat_id, username, text, created_at) VALUES ($1, $2, $3, $4, $5)",
		msg.ID, msg.ChatID, msg.Username, msg.Text, msg.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (m *MessagePostgreSql) GetAllChatMessages(ctx context.Context, chatID string) ([]Message, error) {
	var messages []Message
	rows, err := m.db.Conn.QueryContext(ctx, "SELECT * FROM messages WHERE chat_id=$1", chatID)
	if err != nil {
		return messages, err
	}
	defer rows.Close()
	for rows.Next() {
		var msg Message
		err = rows.Scan(&msg.ID, &msg.ChatID, &msg.Username, &msg.Text, &msg.CreatedAt)
		if err != nil {
			return messages, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
