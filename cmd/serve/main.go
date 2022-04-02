package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"simpleMesseger/cmd/serve/handler"
	"simpleMesseger/internal/messenger"
	"time"
)

func mongoDB() (chat messenger.Chat, msg messenger.Messages) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println(err)
	}
	collectionChat := client.Database("messenger").Collection("chat")
	chatMongoDb := messenger.NewChatMongoDb(collectionChat)

	collectionChatMessages := client.Database("messenger").Collection("chat_messages")
	messageMongoDB := messenger.NewMessageMongoDB(collectionChatMessages)

	return chatMongoDb, messageMongoDB
}

func postgreSql() (chat messenger.Chat, msg messenger.Messages) {
	err := godotenv.Load() // load .env
	if err != nil {
		log.Println(err)
		return
	}
	cfg := messenger.Config{
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
	}
	postgreSqlDB, err := messenger.Initialize(cfg) // instialize posetreSql db
	if err != nil {
		log.Fatalf("failed to initialize database: %s", err.Error())
	}
	chatPostgreSql := messenger.NewChatPostgreSql(postgreSqlDB)
	messagePostgreSql := messenger.NewMessagePostgreSql(postgreSqlDB)
	return chatPostgreSql, messagePostgreSql
}

func main() {
	err := godotenv.Load() // load .env
	if err != nil {
		log.Println(err)
		return
	}
	var chat messenger.Chat
	var message messenger.Messages
	
	if os.Getenv("CHOOSE_DATABASE") == "MONGO_DB" {
		chat, message = mongoDB()
		log.Println("connect to MongoDB")
	} else if os.Getenv("CHOOSE_DATABASE") == "POSTGRE_SQL" {
		chat, message = postgreSql()
		log.Println("connect to PostgreSql")
	} else {
		panic("cant connect database")
	}

	message = messenger.NewMessagesCache(message) // add cache

	h := handler.NewHandler(chat, message)
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	h.Register(router)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}
