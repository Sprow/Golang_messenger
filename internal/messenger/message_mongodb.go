package messenger

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Messages interface {
	AddMessage(ctx context.Context, msg Message) error
	GetAllChatMessages(ctx context.Context, chatID string) ([]Message, error)
}

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

type Message struct {
	ID        string
	ChatID    string
	Username  string
	Text      string
	CreatedAt time.Time
}

func messageToMongoMessage(msg Message) MessageMongo {
	msgMongo := MessageMongo{
		ID:        primitive.ObjectID{},
		ChatID:    primitive.ObjectID{},
		Username:  msg.Username,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}
	primitiveID, err := primitive.ObjectIDFromHex(msg.ID)
	if err != nil {
		log.Println(err)
	}
	msgMongo.ID = primitiveID

	primitiveChatID, err := primitive.ObjectIDFromHex(msg.ChatID)
	if err != nil {
		log.Println(err)
	}
	msgMongo.ChatID = primitiveChatID

	return msgMongo
}

func (db *MessageMongoDB) AddMessage(ctx context.Context, msg Message) error {
	msgMongo := messageToMongoMessage(msg)

	_, err := db.collection.InsertOne(ctx, msgMongo)
	return err
}

func (db *MessageMongoDB) GetAllChatMessages(ctx context.Context, chatID string) ([]Message, error) {
	var chatMessages []Message
	id, err := primitive.ObjectIDFromHex(chatID) // convert string to Primitive.ObjectID

	filter := bson.D{{"chat_id", id}}
	cursor, err := db.collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
		return chatMessages, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result Message
		err = cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
			return chatMessages, err
		}
		chatMessages = append(chatMessages, result)
	}
	return chatMessages, nil
}
