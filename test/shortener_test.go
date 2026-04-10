package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shortener/handler"
	"shortener/storage"
	"testing"
)

func setupHanlder(t *testing.T) *handler.ShortenerHandler { 
	dbStore,err := storage.NewSQLiteStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to connection database.err: %v",err)
	}
	t.Cleanup(func(){
		dbStore.Close()
	})
	return handler.NewShortenerHandler(dbStore)
}

func createShortURLResponse(t *testing.T, h *handler.ShortenerHandler, longURL string) handler.CreateShortURLResponse {
	reqBody := handler.CreateShortURLRequest{LongURL: longURL}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("conversion JSON error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.CreateShortURL(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var resBody handler.CreateShortURLResponse
	if err := json.NewDecoder(resp.Body).Decode(&resBody); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	return resBody
}

func TestCreateShortURL(t *testing.T) {
	h := setupHanlder(t)
	
	longURL := "https://github.com/coulc"
	resBody := createShortURLResponse(t,h,longURL)

	if resBody.LongURL != longURL {
		t.Fatalf("expected long url %s, got %s",longURL,resBody.LongURL)
		return 
	}
}

func TestRedirect(t *testing.T) {
	h := setupHanlder(t)

	longURL := "http://github.com/coulc"
	resp := createShortURLResponse(t, h, longURL)  
	
	req := httptest.NewRequest(http.MethodGet, "/"+resp.ShortCode, nil)
	req.SetPathValue("shortCode", resp.ShortCode)
	
	w := httptest.NewRecorder()
	h.Redirect(w, req)

	redirectResp := w.Result()
	defer redirectResp.Body.Close()

	if redirectResp.StatusCode != http.StatusFound {
		t.Fatalf("expected status %d, got %d",http.StatusFound,redirectResp.StatusCode)
		return 
	}

	location := redirectResp.Header.Get("Location")
	if location != longURL {
		t.Fatalf("expected redirect to %s, got %s",longURL,location)
		return
	}
}

func TestCreateDuplicateURL(t *testing.T) {
	h := setupHanlder(t)
	longURL := "https://github.com/coulc"
	resBody1 := createShortURLResponse(t,h,longURL)
	t.Logf("First creation - ShortCode: %s", resBody1.ShortCode) 
	
	reqBody2 := handler.CreateShortURLRequest{LongURL: longURL}
	body,err := json.Marshal(reqBody2)	
	if err != nil {
		t.Fatalf("conversion JSON error: %v",err)
	}

	req := httptest.NewRequest(http.MethodPost,"/",bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.CreateShortURL(w,req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected status %d for duplicate URL, got %d", http.StatusConflict, resp.StatusCode)
	}
}

func TestDeleteURL(t *testing.T) {
	h := setupHanlder(t)
	
	longURL := "https://github.com/coulc"
	resBody := createShortURLResponse(t,h,longURL)
	
	req := httptest.NewRequest(http.MethodDelete,"/" + resBody.ShortCode,nil)
	req.SetPathValue("shortCode", resBody.ShortCode)
	w := httptest.NewRecorder()
	h.DeleteShortURL(w,req)
	resp := w.Result()
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d for delete URL, got %d", http.StatusNoContent, resp.StatusCode)
	}
}

