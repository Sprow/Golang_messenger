package messenger

import "context"

type ChatPostgreSql struct {
	db PostgreSqlDB
}

func NewChatPostgreSql(db PostgreSqlDB) *ChatPostgreSql {
	return &ChatPostgreSql{
		db: db,
	}
}

func (c *ChatPostgreSql) AddChat(ctx context.Context, newChatID string) (string, error) {
	_, err := c.db.Conn.ExecContext(ctx, "INSERT INTO chat (id, chat_name) VALUES ($1, $2)", newChatID, "")
	if err != nil {
		return "", err
	}
	return newChatID, nil
}
