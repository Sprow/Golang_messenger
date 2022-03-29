package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"simpleMesseger/internal/messenger"
	"time"
)

type Handler struct {
	//chatsManager *messenger.ChatsManager
	chatMongoDB *messenger.ChatMongoDB
	msgMongoDB  *messenger.MessageMongoDB
}

func NewHandler(chatMongoDb *messenger.ChatMongoDB, msgMongoDB *messenger.MessageMongoDB) *Handler {
	return &Handler{
		//chatsManager: cm,
		chatMongoDB: chatMongoDb,
		msgMongoDB:  msgMongoDB,
	}
}

func (h *Handler) Register(r *chi.Mux) {
	r.Get("/", h.addChatHandler)
	r.Get("/chat/{id}", h.getChatHandler)
	r.Post("/chat/{id}/add", h.addMessageHandler)
}

func (h *Handler) addChatHandler(w http.ResponseWriter, r *http.Request) {
	//d := json.NewDecoder(r.Body)
	//var chatName string
	//err := d.Decode(&chatName)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	id, err := h.chatMongoDB.AddChat(ctx)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "http://192.168.1.33:8080/chat/"+id.Hex(), http.StatusFound)
}

//func (h *Handler) addChatHandler(w http.ResponseWriter, r *http.Request) {
//	id := h.chatsManager.AddChat()
//	http.Redirect(w, r, "http://192.168.1.33:8080/chat/"+id.String(), http.StatusFound)
//}

func (h *Handler) getChatHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	chatID, err := primitive.ObjectIDFromHex(id) // convert string to Primitive.ObjectID
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	chat, err := h.msgMongoDB.GetAllChatMessages(ctx, chatID)
	if err != nil {
		log.Printf("Chat not found %s", id)
		http.Error(w, "Chat not found "+id, http.StatusNotFound)
		return
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(chat)
	if err != nil {
		log.Println(err)
	}
}

//func (h *Handler) getChatHandler(w http.ResponseWriter, r *http.Request) {
//	id := chi.URLParam(r, "id")
//	chat, err := h.chatsManager.GetChat(uuid.MustParse(id))
//	if err != nil {
//		log.Printf("Chat not found %s", id)
//		http.Error(w, "Chat not found "+id, http.StatusNotFound)
//		return
//	}
//	encoder := json.NewEncoder(w)
//	err = encoder.Encode(chat)
//	if err != nil {
//		log.Println(err)
//	}
//}

func (h *Handler) addMessageHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")                  //string
	chatID, err := primitive.ObjectIDFromHex(id) // convert string to Primitive.ObjectID
	if err != nil {
		panic(err)
	}
	d := json.NewDecoder(r.Body)
	var msg messenger.MessageMongo
	err = d.Decode(&msg)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	msg.ChatID = chatID
	msg.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	err = h.msgMongoDB.AddMessage(ctx, msg)

	//chat, err := h.chatsManager.GetChat(uuid.MustParse(id))
	//if err != nil {
	//	log.Printf("Chat not found %s", id)
	//	http.Error(w, "Chat not found "+id, http.StatusNotFound)
	//	return
	//}
	//msg.CreatedAt = time.Now()
	//chat.AddMessage(msg)
}

//func (h *Handler) addMessageHandler(w http.ResponseWriter, r *http.Request) {
//	id := chi.URLParam(r, "id")
//	d := json.NewDecoder(r.Body)
//	var msg messenger.Message
//	err := d.Decode(&msg)
//	if err != nil {
//		log.Println(err)
//		http.Error(w, "Bad request", http.StatusBadRequest)
//		return
//	}
//	chat, err := h.chatsManager.GetChat(uuid.MustParse(id))
//	if err != nil {
//		log.Printf("Chat not found %s", id)
//		http.Error(w, "Chat not found "+id, http.StatusNotFound)
//		return
//	}
//	msg.CreatedAt = time.Now()
//	chat.AddMessage(msg)
//}
