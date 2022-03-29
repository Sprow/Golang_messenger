package messenger

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type MessageMongoDB struct {
	collection *mongo.Collection
}

func NewMessageMongoDB(collection *mongo.Collection) *MessageMongoDB {
	return &MessageMongoDB{
		collection: collection,
	}
}

type MessageMongo struct {
	ID        primitive.ObjectID `bson:"_id"`
	ChatID    primitive.ObjectID `bson:"chat_id"`
	Username  string             `bson:"username"`
	Text      string             `bson:"text"`
	CreatedAt time.Time          `bson:"created_at"`
}

func (db *MessageMongoDB) AddMessage(ctx context.Context, msg MessageMongo) error {
	msg.ID = primitive.NewObjectID()
	_, err := db.collection.InsertOne(ctx, msg)
	return err
}

func (db *MessageMongoDB) GetAllChatMessages(ctx context.Context, chatID primitive.ObjectID) ([]MessageMongo, error) {
	var chatMessages []MessageMongo
	filter := bson.D{{"chat_id", chatID}}
	cursor, err := db.collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
		return chatMessages, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result MessageMongo
		err = cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
			return chatMessages, err
		}
		chatMessages = append(chatMessages, result)
	}
	return chatMessages, nil
}
