package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"simpleMesseger/internal/messenger"
	"time"
)

type Handler struct {
	chatsManager *messenger.ChatsManager
}

func NewHandler(cm *messenger.ChatsManager) *Handler {
	return &Handler{
		chatsManager: cm,
	}
}

func (h *Handler) Register(r *chi.Mux) {
	r.Get("/", h.addChatHandler)
	r.Get("/chat/{id}", h.getChatHandler)
	r.Post("/chat/{id}/add", h.addMessageHandler)
}

func (h *Handler) addChatHandler(w http.ResponseWriter, r *http.Request) {
	id := h.chatsManager.AddChat()
	http.Redirect(w, r, "http://192.168.1.33:8080/chat/"+id.String(), http.StatusFound)
}

func (h *Handler) getChatHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	chat, err := h.chatsManager.GetChat(uuid.MustParse(id))
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
	id := chi.URLParam(r, "id")
	d := json.NewDecoder(r.Body)
	var msg messenger.Message
	err := d.Decode(&msg)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	chat, err := h.chatsManager.GetChat(uuid.MustParse(id))
	if err != nil {
		log.Printf("Chat not found %s", id)
		http.Error(w, "Chat not found "+id, http.StatusNotFound)
		return
	}
	msg.CreatedAt = time.Now()
	chat.AddMessage(msg)
}
