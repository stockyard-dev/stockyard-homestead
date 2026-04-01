package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"
	_ "modernc.org/sqlite"
)

type DB struct{ conn *sql.DB }

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	conn, err := sql.Open("sqlite", filepath.Join(dataDir, "homestead.db"))
	if err != nil { return nil, err }
	conn.Exec("PRAGMA journal_mode=WAL")
	conn.Exec("PRAGMA busy_timeout=5000")
	conn.SetMaxOpenConns(4)
	db := &DB{conn: conn}
	return db, db.migrate()
}

func (db *DB) Close() error { return db.conn.Close() }

func (db *DB) migrate() error {
	_, err := db.conn.Exec(`
CREATE TABLE IF NOT EXISTS bookmarks (
    id TEXT PRIMARY KEY, title TEXT NOT NULL, url TEXT NOT NULL,
    category TEXT DEFAULT '', icon TEXT DEFAULT '', sort_order INTEGER DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now'))
);
CREATE TABLE IF NOT EXISTS notes (
    id TEXT PRIMARY KEY, title TEXT DEFAULT '', content TEXT DEFAULT '',
    pinned INTEGER DEFAULT 0, created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);
CREATE TABLE IF NOT EXISTS feeds (
    id TEXT PRIMARY KEY, title TEXT NOT NULL, url TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now'))
);
`)
	return err
}

type Bookmark struct {
	ID string `json:"id"`; Title string `json:"title"`; URL string `json:"url"`
	Category string `json:"category"`; SortOrder int `json:"sort_order"`; CreatedAt string `json:"created_at"`
}

func (db *DB) CreateBookmark(title, url, category string) (*Bookmark, error) {
	id := "bm_" + genID(6)
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := db.conn.Exec("INSERT INTO bookmarks (id,title,url,category,created_at) VALUES (?,?,?,?,?)", id, title, url, category, now)
	if err != nil { return nil, err }
	return &Bookmark{ID: id, Title: title, URL: url, Category: category, CreatedAt: now}, nil
}

func (db *DB) ListBookmarks() ([]Bookmark, error) {
	rows, err := db.conn.Query("SELECT id,title,url,category,sort_order,created_at FROM bookmarks ORDER BY category, sort_order, title")
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Bookmark
	for rows.Next() {
		var b Bookmark
		rows.Scan(&b.ID, &b.Title, &b.URL, &b.Category, &b.SortOrder, &b.CreatedAt)
		out = append(out, b)
	}
	return out, rows.Err()
}

func (db *DB) DeleteBookmark(id string) { db.conn.Exec("DELETE FROM bookmarks WHERE id=?", id) }

type Note struct {
	ID string `json:"id"`; Title string `json:"title"`; Content string `json:"content"`
	Pinned bool `json:"pinned"`; CreatedAt string `json:"created_at"`; UpdatedAt string `json:"updated_at"`
}

func (db *DB) CreateNote(title, content string) (*Note, error) {
	id := "note_" + genID(6)
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := db.conn.Exec("INSERT INTO notes (id,title,content,created_at,updated_at) VALUES (?,?,?,?,?)", id, title, content, now, now)
	if err != nil { return nil, err }
	return &Note{ID: id, Title: title, Content: content, CreatedAt: now, UpdatedAt: now}, nil
}

func (db *DB) ListNotes() ([]Note, error) {
	rows, err := db.conn.Query("SELECT id,title,content,pinned,created_at,updated_at FROM notes ORDER BY pinned DESC, updated_at DESC")
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Note
	for rows.Next() {
		var n Note; var p int
		rows.Scan(&n.ID, &n.Title, &n.Content, &p, &n.CreatedAt, &n.UpdatedAt)
		n.Pinned = p == 1
		out = append(out, n)
	}
	return out, rows.Err()
}

func (db *DB) UpdateNote(id string, title, content *string) {
	now := time.Now().UTC().Format(time.RFC3339)
	if title != nil { db.conn.Exec("UPDATE notes SET title=?, updated_at=? WHERE id=?", *title, now, id) }
	if content != nil { db.conn.Exec("UPDATE notes SET content=?, updated_at=? WHERE id=?", *content, now, id) }
}

func (db *DB) DeleteNote(id string) { db.conn.Exec("DELETE FROM notes WHERE id=?", id) }

type Feed struct {
	ID string `json:"id"`; Title string `json:"title"`; URL string `json:"url"`; CreatedAt string `json:"created_at"`
}

func (db *DB) CreateFeed(title, url string) (*Feed, error) {
	id := "feed_" + genID(6)
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := db.conn.Exec("INSERT INTO feeds (id,title,url,created_at) VALUES (?,?,?,?)", id, title, url, now)
	if err != nil { return nil, err }
	return &Feed{ID: id, Title: title, URL: url, CreatedAt: now}, nil
}

func (db *DB) ListFeeds() ([]Feed, error) {
	rows, err := db.conn.Query("SELECT id,title,url,created_at FROM feeds ORDER BY title")
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Feed
	for rows.Next() {
		var f Feed
		rows.Scan(&f.ID, &f.Title, &f.URL, &f.CreatedAt)
		out = append(out, f)
	}
	return out, rows.Err()
}

func (db *DB) DeleteFeed(id string) { db.conn.Exec("DELETE FROM feeds WHERE id=?", id) }

func (db *DB) Stats() map[string]any {
	var bm, notes, feeds int
	db.conn.QueryRow("SELECT COUNT(*) FROM bookmarks").Scan(&bm)
	db.conn.QueryRow("SELECT COUNT(*) FROM notes").Scan(&notes)
	db.conn.QueryRow("SELECT COUNT(*) FROM feeds").Scan(&feeds)
	return map[string]any{"bookmarks": bm, "notes": notes, "feeds": feeds}
}

func genID(n int) string { b := make([]byte, n); rand.Read(b); return hex.EncodeToString(b) }
