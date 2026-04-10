package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"shortener/storage"
	"shortener/utils"
	"time"
)

type ShortenerHandler struct {
	Storage storage.Storage
}

func NewShortenerHandler(storage storage.Storage) *ShortenerHandler {
	return &ShortenerHandler{
		Storage: storage,
	}
}

type CreateShortURLRequest struct {
	LongURL string `json:"long_url"`
}

type CreateShortURLResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
	LongURL   string `json:"long_url"`
}

func (h *ShortenerHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var req CreateShortURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request.", http.StatusBadRequest)
		slog.Error("CreateShortURL Error.","err",err)
		return
	}

	if !validateURL(req.LongURL) {
		http.Error(w,"Invalid  url",http.StatusBadRequest)
		slog.Error("Invalid url","url",req.LongURL)
		return 
	}

	shortCode := utils.GenerateShortCode(req.LongURL)

	mapping := &storage.URLMapping{
		ShortCode:  shortCode,
		LongURL:    req.LongURL,
		CreatedAt:   time.Now().Unix(),
		VisitCount: 0,
	}

	if err := h.Storage.Save(mapping); err != nil {
		if errors.Is(err,storage.URLExist) {
			http.Error(w,storage.URLExist.Error(),http.StatusConflict)
			slog.Warn("Repeat adding URL.","err",err,"url",req.LongURL)	
			return 
		}
		http.Error(w,"Internal error",http.StatusInternalServerError)
		slog.Error("Save mapping error","err",err)	
		return 
	}

	resp := CreateShortURLResponse{
		ShortCode: mapping.ShortCode,
		ShortURL:  "http://" + r.Host + "/" + mapping.ShortCode,
		LongURL:   req.LongURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
	slog.Info("Successfully added url","url",resp.LongURL)
}

func (h *ShortenerHandler) DeleteShortURL(w http.ResponseWriter,r *http.Request) {
	shortCode := r.PathValue("shortCode")

	if shortCode == "" {
		http.Error(w, "Missing short code.", http.StatusBadRequest)
		return
	}

	if err := h.Storage.Delete(shortCode);err != nil {
		if errors.Is(err,storage.URLNotFound) {
			http.Error(w,storage.URLNotFound.Error(),http.StatusNotFound)
			slog.Warn("Attempt to delete non-existent URLs.","shortCode",shortCode)
			return 
		}
		http.Error(w,"Internal error",http.StatusInternalServerError)
		slog.Error("Delete url error","err",err,"shortCode",shortCode)	
		return 
	}

	w.WriteHeader(http.StatusNoContent)
	slog.Info("Successfully deleted url","shortCode",shortCode)
}


func (h *ShortenerHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")

	if shortCode == "" {
			http.Error(w, "Missing short code.", http.StatusBadRequest)
			return
	}

	mapping, err := h.Storage.Get(shortCode)
	if err != nil || mapping == nil {
		http.NotFound(w, r)
		slog.Error("Short code not found", "shortCode", shortCode)
		return
	}

	if err := h.Storage.IncrementVisit(shortCode); err != nil {
		slog.Error("Failed to Increment visit count.","short_code",shortCode,"err",err)
	}
	
	http.Redirect(w, r, mapping.LongURL, http.StatusFound)
}

func validateURL(raw string) bool {
	u,err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	if u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

func Index(w http.ResponseWriter,r *http.Request) {
	w.Write([]byte("hello world!"))
	slog.Info("hello world!")
}
