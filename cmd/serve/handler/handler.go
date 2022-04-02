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
	chat messenger.Chat
	msg  messenger.Messages
}

func NewHandler(chat messenger.Chat, msg messenger.Messages) *Handler {
	return &Handler{
		chat: chat,
		msg:  msg,
	}
}

func (h *Handler) Register(r *chi.Mux) {
	r.Get("/", h.addChatHandler)
	r.Get("/chat/{id}", h.getChatHandler)
	r.Post("/chat/{id}/add", h.addMessageHandler)
}

func (h *Handler) addChatHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	newChatID := primitive.NewObjectID().Hex()

	id, err := h.chat.AddChat(ctx, newChatID)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "http://192.168.1.33:8080/chat/"+id, http.StatusFound)
}

func (h *Handler) getChatHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	chat, err := h.msg.GetAllChatMessages(ctx, id)
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

func (h *Handler) addMessageHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id") //string
	d := json.NewDecoder(r.Body)
	var msg messenger.Message
	err := d.Decode(&msg)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	msg.ChatID = chatID
	msg.CreatedAt = time.Now()
	msg.ID = primitive.NewObjectID().Hex()

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	err = h.msg.AddMessage(ctx, msg)
	if err != nil {
		log.Println(err)
	}
}
