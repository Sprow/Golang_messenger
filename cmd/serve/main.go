package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"simpleMesseger/cmd/serve/handler"
	"simpleMesseger/internal/messenger"
	"time"
)

func main() {
	//chatManager := messenger.NewChatsManager()
	//h := handler.NewHandler(chatManager)
	//router := chi.NewRouter()
	//router.Use(middleware.Logger)
	//h.Register(router)

	// MongoDb
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

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	h := handler.NewHandler(chatMongoDb, messageMongoDB)
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	h.Register(router)

	// check connection
	//ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	//defer cancel()
	//err = client.Ping(ctx, readpref.Primary())
	//log.Println(err)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}
