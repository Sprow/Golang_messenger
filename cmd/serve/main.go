package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"simpleMesseger/cmd/serve/handler"
	"simpleMesseger/internal/messenger"
)

func main() {
	chatManager := messenger.NewChatsManager()
	h := handler.NewHandler(chatManager)
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	h.Register(router)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}
