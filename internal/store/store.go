package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Server struct{ID string `json:"id"`;Hostname string `json:"hostname"`;IP string `json:"ip"`;OS string `json:"os"`;Provider string `json:"provider"`;Region string `json:"region"`;Tags string `json:"tags"`;Status string `json:"status"`;Notes string `json:"notes"`;CreatedAt string `json:"created_at"`}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"homestead.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS servers(id TEXT PRIMARY KEY,hostname TEXT NOT NULL,ip TEXT DEFAULT '',os TEXT DEFAULT '',provider TEXT DEFAULT '',region TEXT DEFAULT '',tags TEXT DEFAULT '',status TEXT DEFAULT 'active',notes TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(s *Server)error{s.ID=genID();s.CreatedAt=now();if s.Status==""{s.Status="active"};_,err:=d.db.Exec(`INSERT INTO servers VALUES(?,?,?,?,?,?,?,?,?,?)`,s.ID,s.Hostname,s.IP,s.OS,s.Provider,s.Region,s.Tags,s.Status,s.Notes,s.CreatedAt);return err}
func(d *DB)Get(id string)*Server{var s Server;if d.db.QueryRow(`SELECT id,hostname,ip,os,provider,region,tags,status,notes,created_at FROM servers WHERE id=?`,id).Scan(&s.ID,&s.Hostname,&s.IP,&s.OS,&s.Provider,&s.Region,&s.Tags,&s.Status,&s.Notes,&s.CreatedAt)!=nil{return nil};return &s}
func(d *DB)List()[]Server{rows,_:=d.db.Query(`SELECT id,hostname,ip,os,provider,region,tags,status,notes,created_at FROM servers ORDER BY hostname`);if rows==nil{return nil};defer rows.Close()
var o []Server;for rows.Next(){var s Server;rows.Scan(&s.ID,&s.Hostname,&s.IP,&s.OS,&s.Provider,&s.Region,&s.Tags,&s.Status,&s.Notes,&s.CreatedAt);o=append(o,s)};return o}
func(d *DB)Update(s *Server)error{_,err:=d.db.Exec(`UPDATE servers SET hostname=?,ip=?,os=?,provider=?,region=?,tags=?,status=?,notes=? WHERE id=?`,s.Hostname,s.IP,s.OS,s.Provider,s.Region,s.Tags,s.Status,s.Notes,s.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM servers WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM servers`).Scan(&n);return n}
func(d *DB)Stats()map[string]any{var total,active int;d.db.QueryRow(`SELECT COUNT(*) FROM servers`).Scan(&total);d.db.QueryRow(`SELECT COUNT(*) FROM servers WHERE status='active'`).Scan(&active);return map[string]any{"total":total,"active":active}}
func(d *DB)Search(q string,filters map[string]string)[]Server{return d.List()}
