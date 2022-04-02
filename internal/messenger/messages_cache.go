package messenger

import (
	"context"
	"log"
)

type MessageCache struct {
	messages Messages
	msgCache map[string][]Message
}

func NewMessagesCache(messages Messages) *MessageCache {
	return &MessageCache{
		messages: messages,
		msgCache: map[string][]Message{},
	}
}

func (mc *MessageCache) AddMessage(ctx context.Context, msg Message) error {
	err := mc.messages.AddMessage(ctx, msg)
	if err != nil {
		return err
	}
	delete(mc.msgCache, msg.ChatID)
	return nil
}

func (mc *MessageCache) GetAllChatMessages(ctx context.Context, chatID string) ([]Message, error) {
	allMsg, ok := mc.msgCache[chatID]
	if ok {
		log.Println("get all messages from cache")
		return allMsg, nil
	}

	messages, err := mc.messages.GetAllChatMessages(ctx, chatID)
	if err != nil {
		return []Message{}, err
	}
	mc.msgCache[chatID] = messages
	log.Println("get all messages from DB")
	return messages, nil
}
