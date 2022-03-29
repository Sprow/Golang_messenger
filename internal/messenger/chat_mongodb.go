package messenger

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func (db *ChatMongoDB) AddChat(ctx context.Context) (primitive.ObjectID, error) {
	chatID := primitive.NewObjectID()
	chat := ChatMongo{
		ID:       chatID,
		ChatName: "",
	}
	_, err := db.collection.InsertOne(ctx, chat)
	return chatID, err
}
