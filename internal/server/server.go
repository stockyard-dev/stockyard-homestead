package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/stockyard-dev/stockyard-homestead/internal/store"
)

type Server struct { db *store.DB; mux *http.ServeMux; port int; limits Limits }

func New(db *store.DB, port int, limits Limits) *Server {
	s := &Server{db: db, mux: http.NewServeMux(), port: port, limits: limits}
	s.mux.HandleFunc("POST /api/bookmarks", s.hCreateBM)
	s.mux.HandleFunc("GET /api/bookmarks", s.hListBM)
	s.mux.HandleFunc("DELETE /api/bookmarks/{id}", s.hDelBM)
	s.mux.HandleFunc("POST /api/notes", s.hCreateNote)
	s.mux.HandleFunc("GET /api/notes", s.hListNotes)
	s.mux.HandleFunc("PUT /api/notes/{id}", s.hUpdateNote)
	s.mux.HandleFunc("DELETE /api/notes/{id}", s.hDelNote)
	s.mux.HandleFunc("POST /api/feeds", s.hCreateFeed)
	s.mux.HandleFunc("GET /api/feeds", s.hListFeeds)
	s.mux.HandleFunc("DELETE /api/feeds/{id}", s.hDelFeed)
	s.mux.HandleFunc("GET /api/status", func(w http.ResponseWriter, r *http.Request) { wj(w, 200, s.db.Stats()) })
	s.mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) { wj(w, 200, map[string]string{"status": "ok"}) })
	s.mux.HandleFunc("GET /ui", s.handleUI)
	s.mux.HandleFunc("GET /api/version", func(w http.ResponseWriter, r *http.Request) { wj(w, 200, map[string]any{"product": "stockyard-homestead", "version": "0.1.0"}) })
	return s
}

func (s *Server) Start() error {
	log.Printf("[homestead] listening on :%d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}

func (s *Server) hCreateBM(w http.ResponseWriter, r *http.Request) {
	var req struct { Title string `json:"title"`; URL string `json:"url"`; Category string `json:"category"` }
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Title == "" || req.URL == "" {
		wj(w, 400, map[string]string{"error": "title and url required"}); return
	}
	bm, err := s.db.CreateBookmark(req.Title, req.URL, req.Category)
	if err != nil { wj(w, 500, map[string]string{"error": err.Error()}); return }
	wj(w, 201, map[string]any{"bookmark": bm})
}
func (s *Server) hListBM(w http.ResponseWriter, r *http.Request) {
	bm, _ := s.db.ListBookmarks(); if bm == nil { bm = []store.Bookmark{} }
	wj(w, 200, map[string]any{"bookmarks": bm, "count": len(bm)})
}
func (s *Server) hDelBM(w http.ResponseWriter, r *http.Request) { s.db.DeleteBookmark(r.PathValue("id")); wj(w, 200, map[string]string{"status": "deleted"}) }

func (s *Server) hCreateNote(w http.ResponseWriter, r *http.Request) {
	var req struct { Title string `json:"title"`; Content string `json:"content"` }
	json.NewDecoder(r.Body).Decode(&req)
	n, _ := s.db.CreateNote(req.Title, req.Content)
	wj(w, 201, map[string]any{"note": n})
}
func (s *Server) hListNotes(w http.ResponseWriter, r *http.Request) {
	n, _ := s.db.ListNotes(); if n == nil { n = []store.Note{} }
	wj(w, 200, map[string]any{"notes": n, "count": len(n)})
}
func (s *Server) hUpdateNote(w http.ResponseWriter, r *http.Request) {
	var req struct { Title *string `json:"title"`; Content *string `json:"content"` }
	json.NewDecoder(r.Body).Decode(&req)
	s.db.UpdateNote(r.PathValue("id"), req.Title, req.Content)
	wj(w, 200, map[string]string{"status": "updated"})
}
func (s *Server) hDelNote(w http.ResponseWriter, r *http.Request) { s.db.DeleteNote(r.PathValue("id")); wj(w, 200, map[string]string{"status": "deleted"}) }

func (s *Server) hCreateFeed(w http.ResponseWriter, r *http.Request) {
	var req struct { Title string `json:"title"`; URL string `json:"url"` }
	json.NewDecoder(r.Body).Decode(&req)
	f, _ := s.db.CreateFeed(req.Title, req.URL)
	wj(w, 201, map[string]any{"feed": f})
}
func (s *Server) hListFeeds(w http.ResponseWriter, r *http.Request) {
	f, _ := s.db.ListFeeds(); if f == nil { f = []store.Feed{} }
	wj(w, 200, map[string]any{"feeds": f, "count": len(f)})
}
func (s *Server) hDelFeed(w http.ResponseWriter, r *http.Request) { s.db.DeleteFeed(r.PathValue("id")); wj(w, 200, map[string]string{"status": "deleted"}) }

func wj(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json"); w.WriteHeader(code); json.NewEncoder(w).Encode(v)
}
