package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Server struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Hostname string `json:"hostname"`
	IP string `json:"ip"`
	OS string `json:"os"`
	Role string `json:"role"`
	Location string `json:"location"`
	Status string `json:"status"`
	Notes string `json:"notes"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"homestead.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS servers(id TEXT PRIMARY KEY,name TEXT NOT NULL,hostname TEXT DEFAULT '',ip TEXT DEFAULT '',os TEXT DEFAULT '',role TEXT DEFAULT '',location TEXT DEFAULT '',status TEXT DEFAULT 'active',notes TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Server)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO servers(id,name,hostname,ip,os,role,location,status,notes,created_at)VALUES(?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Hostname,e.IP,e.OS,e.Role,e.Location,e.Status,e.Notes,e.CreatedAt);return err}
func(d *DB)Get(id string)*Server{var e Server;if d.db.QueryRow(`SELECT id,name,hostname,ip,os,role,location,status,notes,created_at FROM servers WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Hostname,&e.IP,&e.OS,&e.Role,&e.Location,&e.Status,&e.Notes,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Server{rows,_:=d.db.Query(`SELECT id,name,hostname,ip,os,role,location,status,notes,created_at FROM servers ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Server;for rows.Next(){var e Server;rows.Scan(&e.ID,&e.Name,&e.Hostname,&e.IP,&e.OS,&e.Role,&e.Location,&e.Status,&e.Notes,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM servers WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM servers`).Scan(&n);return n}
