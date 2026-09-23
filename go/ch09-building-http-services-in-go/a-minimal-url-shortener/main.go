package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type URLEntry struct {
	ID          int64     `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	Clicks      int       `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func readJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func internalError(w http.ResponseWriter, err error) {
	slog.Error("internal error", "error", err)
	writeJSON(w, http.StatusInternalServerError, ErrorResponse{
		Error: "internal server error",
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start))
	})
}

// generateCode returns an 8-character hex short code, e.g. "a1b2c3d4".
func generateCode() (string, error) {
	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

type server struct {
	db *sql.DB
}

func (s *server) handleShorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := readJSON(r, &req); err != nil || req.URL == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	code, err := generateCode()
	if err != nil {
		internalError(w, err)
		return
	}

	_, err = s.db.Exec(
		"INSERT INTO urls (short_code, original_url) VALUES (?, ?)",
		code, req.URL,
	)
	if err != nil {
		internalError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, ShortenResponse{
		ShortCode: code,
		ShortURL:  "http://localhost:8080/" + code,
	})
}

func (s *server) handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	var originalURL string
	err := s.db.QueryRow(
		"SELECT original_url FROM urls WHERE short_code = ?", code,
	).Scan(&originalURL)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "short code not found"})
		return
	}
	if err != nil {
		internalError(w, err)
		return
	}

	// Fire-and-forget: don't make the redirect wait on the click count update.
	go s.incrementClicks(code)

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func (s *server) incrementClicks(code string) {
	if _, err := s.db.Exec("UPDATE urls SET clicks = clicks + 1 WHERE short_code = ?", code); err != nil {
		slog.Error("failed to increment clicks", "short_code", code, "error", err)
	}
}

func (s *server) handleStats(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	var e URLEntry
	err := s.db.QueryRow(
		"SELECT id, short_code, original_url, clicks, created_at FROM urls WHERE short_code = ?",
		code,
	).Scan(&e.ID, &e.ShortCode, &e.OriginalURL, &e.Clicks, &e.CreatedAt)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "short code not found"})
		return
	}
	if err != nil {
		internalError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, e)
}

func (s *server) handleList(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(
		"SELECT id, short_code, original_url, clicks, created_at FROM urls ORDER BY created_at DESC",
	)
	if err != nil {
		internalError(w, err)
		return
	}
	defer rows.Close()

	entries := []URLEntry{}
	for rows.Next() {
		var e URLEntry
		if err := rows.Scan(&e.ID, &e.ShortCode, &e.OriginalURL, &e.Clicks, &e.CreatedAt); err != nil {
			internalError(w, err)
			return
		}
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		internalError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, entries)
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "urlshort.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_code TEXT UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
		clicks INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	s := &server{db: db}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/shorten", s.handleShorten)
	mux.HandleFunc("GET /api/stats/{code}", s.handleStats)
	mux.HandleFunc("GET /api/urls", s.handleList)
	mux.HandleFunc("GET /{code}", s.handleRedirect)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("server starting", "addr", server.Addr)
	log.Fatal(server.ListenAndServe())
}
