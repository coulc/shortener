package main

import (
	"net/http"
	"shortener/handler"
	"shortener/storage"
)


func dbInit() (*handler.ShortenerHandler,func()) {
	dbStore,err := storage.NewSQLiteStorage("db/urls.db")
	if err != nil {
		panic(err)
	}
	shortenerHandler := handler.NewShortenerHandler(dbStore)

	return shortenerHandler,func() {
		dbStore.Close()
	}
}

func routesInit(h *handler.ShortenerHandler) {
	http.HandleFunc("POST /shorten",h.CreateShortURL)
// 	http.HandleFunc("GET /shorten/{shortCode}",h.Redirect)
	http.HandleFunc("GET /{shortCode}",h.Redirect)
	http.HandleFunc("DELETE /{shortCode}",h.DeleteShortURL)
	http.HandleFunc("GET /{$}",handler.Index)
}
