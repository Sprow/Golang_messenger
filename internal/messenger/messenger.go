package messenger

import (
	"fmt"
	"github.com/google/uuid"
	"log"
)

type ChatsManager struct {
	listOfChats map[uuid.UUID]*Chat
}

func NewChatsManager() *ChatsManager {
	return &ChatsManager{
		listOfChats: make(map[uuid.UUID]*Chat),
	}
}

func (cm *ChatsManager) AddChat() uuid.UUID {
	id := uuid.New()
	cm.listOfChats[id] = NewChat()
	log.Println("new chat id >>>", id)
	return id
}

func (cm *ChatsManager) GetChat(id uuid.UUID) (*Chat, error) {
	chat, ok := cm.listOfChats[id]
	if ok {
		return chat, nil
	}
	return nil, fmt.Errorf("chat not found")
}

//func (cm *ChatsManager) AddMessage(id uuid.UUID, msg Message) {
//	cm.listOfChats[id].ListOfMessages = append(cm.listOfChats[id].ListOfMessages, msg)
//}
