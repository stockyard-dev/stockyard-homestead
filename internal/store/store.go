package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

// Host is a single tracked server / VM in the inventory. Renamed from
// the original "Server" type to avoid clashing with the HTTP Server type.
// Status is one of: active, inactive, decommissioned, maintenance.
type Host struct {
	ID        string `json:"id"`
	Hostname  string `json:"hostname"`
	IP        string `json:"ip"`
	OS        string `json:"os"`
	Provider  string `json:"provider"`
	Region    string `json:"region"`
	Tags      string `json:"tags"`
	Status    string `json:"status"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "homestead.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS servers(
		id TEXT PRIMARY KEY,
		hostname TEXT NOT NULL,
		ip TEXT DEFAULT '',
		os TEXT DEFAULT '',
		provider TEXT DEFAULT '',
		region TEXT DEFAULT '',
		tags TEXT DEFAULT '',
		status TEXT DEFAULT 'active',
		notes TEXT DEFAULT '',
		created_at TEXT DEFAULT(datetime('now'))
	)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_hosts_status ON servers(status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_hosts_provider ON servers(provider)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_hosts_region ON servers(region)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
		resource TEXT NOT NULL,
		record_id TEXT NOT NULL,
		data TEXT NOT NULL DEFAULT '{}',
		PRIMARY KEY(resource, record_id)
	)`)
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string   { return time.Now().UTC().Format(time.RFC3339) }

func (d *DB) Create(s *Host) error {
	s.ID = genID()
	s.CreatedAt = now()
	if s.Status == "" {
		s.Status = "active"
	}
	_, err := d.db.Exec(
		`INSERT INTO servers(id, hostname, ip, os, provider, region, tags, status, notes, created_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.ID, s.Hostname, s.IP, s.OS, s.Provider, s.Region, s.Tags, s.Status, s.Notes, s.CreatedAt,
	)
	return err
}

func (d *DB) Get(id string) *Host {
	var s Host
	err := d.db.QueryRow(
		`SELECT id, hostname, ip, os, provider, region, tags, status, notes, created_at
		 FROM servers WHERE id=?`,
		id,
	).Scan(&s.ID, &s.Hostname, &s.IP, &s.OS, &s.Provider, &s.Region, &s.Tags, &s.Status, &s.Notes, &s.CreatedAt)
	if err != nil {
		return nil
	}
	return &s
}

func (d *DB) List() []Host {
	rows, _ := d.db.Query(
		`SELECT id, hostname, ip, os, provider, region, tags, status, notes, created_at
		 FROM servers ORDER BY hostname ASC`,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Host
	for rows.Next() {
		var s Host
		rows.Scan(&s.ID, &s.Hostname, &s.IP, &s.OS, &s.Provider, &s.Region, &s.Tags, &s.Status, &s.Notes, &s.CreatedAt)
		o = append(o, s)
	}
	return o
}

func (d *DB) Update(s *Host) error {
	_, err := d.db.Exec(
		`UPDATE servers SET hostname=?, ip=?, os=?, provider=?, region=?, tags=?, status=?, notes=?
		 WHERE id=?`,
		s.Hostname, s.IP, s.OS, s.Provider, s.Region, s.Tags, s.Status, s.Notes, s.ID,
	)
	return err
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM servers WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM servers`).Scan(&n)
	return n
}

// Search is a real implementation now. The original was a stub that
// just returned List() and ignored all filters.
func (d *DB) Search(q string, filters map[string]string) []Host {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (hostname LIKE ? OR ip LIKE ? OR notes LIKE ? OR tags LIKE ?)"
		s := "%" + q + "%"
		args = append(args, s, s, s, s)
	}
	if v, ok := filters["status"]; ok && v != "" {
		where += " AND status=?"
		args = append(args, v)
	}
	if v, ok := filters["provider"]; ok && v != "" {
		where += " AND provider=?"
		args = append(args, v)
	}
	if v, ok := filters["region"]; ok && v != "" {
		where += " AND region=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(
		`SELECT id, hostname, ip, os, provider, region, tags, status, notes, created_at
		 FROM servers WHERE `+where+`
		 ORDER BY hostname ASC`,
		args...,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Host
	for rows.Next() {
		var s Host
		rows.Scan(&s.ID, &s.Hostname, &s.IP, &s.OS, &s.Provider, &s.Region, &s.Tags, &s.Status, &s.Notes, &s.CreatedAt)
		o = append(o, s)
	}
	return o
}

// Stats returns total hosts plus by_status, by_provider, and by_region
// breakdowns. The original only returned total and active counts.
func (d *DB) Stats() map[string]any {
	m := map[string]any{
		"total":       d.Count(),
		"active":      0,
		"by_status":   map[string]int{},
		"by_provider": map[string]int{},
		"by_region":   map[string]int{},
	}

	var active int
	d.db.QueryRow(`SELECT COUNT(*) FROM servers WHERE status='active'`).Scan(&active)
	m["active"] = active

	if rows, _ := d.db.Query(`SELECT status, COUNT(*) FROM servers GROUP BY status`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_status"] = by
	}

	if rows, _ := d.db.Query(`SELECT provider, COUNT(*) FROM servers WHERE provider != '' GROUP BY provider`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_provider"] = by
	}

	if rows, _ := d.db.Query(`SELECT region, COUNT(*) FROM servers WHERE region != '' GROUP BY region`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_region"] = by
	}

	return m
}

// ─── Extras ───────────────────────────────────────────────────────

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
