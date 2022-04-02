package messenger

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Chat interface {
	AddChat(ctx context.Context, newChatID string) (string, error) //return chat id
}

type ChatMongoDB struct {
	collection *mongo.Collection
}

func NewChatMongoDb(collection *mongo.Collection) *ChatMongoDB {
	return &ChatMongoDB{
		collection: collection,
	}
}

type ChatMongo struct {
	ID       primitive.ObjectID `bson:"_id"`
	ChatName string             `bson:"name"`
}

func (db *ChatMongoDB) AddChat(ctx context.Context, newChatID string) (string, error) {
	//chatID := primitive.NewObjectID()
	chatID, err := primitive.ObjectIDFromHex(newChatID)
	if err != nil {
		log.Println(err)
	}
	chat := ChatMongo{
		ID:       chatID,
		ChatName: "",
	}
	_, err = db.collection.InsertOne(ctx, chat)
	return chatID.Hex(), err
}
