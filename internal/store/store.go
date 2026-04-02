package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Property struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Address string `json:"address"`
	Type string `json:"type"`
	Bedrooms int `json:"bedrooms"`
	Rent int `json:"rent"`
	Status string `json:"status"`
	TenantName string `json:"tenant_name"`
	Notes string `json:"notes"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"homestead.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS properties(id TEXT PRIMARY KEY,name TEXT NOT NULL,address TEXT DEFAULT '',type TEXT DEFAULT 'residential',bedrooms INTEGER DEFAULT 0,rent INTEGER DEFAULT 0,status TEXT DEFAULT 'available',tenant_name TEXT DEFAULT '',notes TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Property)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO properties(id,name,address,type,bedrooms,rent,status,tenant_name,notes,created_at)VALUES(?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Address,e.Type,e.Bedrooms,e.Rent,e.Status,e.TenantName,e.Notes,e.CreatedAt);return err}
func(d *DB)Get(id string)*Property{var e Property;if d.db.QueryRow(`SELECT id,name,address,type,bedrooms,rent,status,tenant_name,notes,created_at FROM properties WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Address,&e.Type,&e.Bedrooms,&e.Rent,&e.Status,&e.TenantName,&e.Notes,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Property{rows,_:=d.db.Query(`SELECT id,name,address,type,bedrooms,rent,status,tenant_name,notes,created_at FROM properties ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Property;for rows.Next(){var e Property;rows.Scan(&e.ID,&e.Name,&e.Address,&e.Type,&e.Bedrooms,&e.Rent,&e.Status,&e.TenantName,&e.Notes,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Property)error{_,err:=d.db.Exec(`UPDATE properties SET name=?,address=?,type=?,bedrooms=?,rent=?,status=?,tenant_name=?,notes=? WHERE id=?`,e.Name,e.Address,e.Type,e.Bedrooms,e.Rent,e.Status,e.TenantName,e.Notes,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM properties WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM properties`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Property{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["type"];ok&&v!=""{where+=" AND type=?";args=append(args,v)}
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,address,type,bedrooms,rent,status,tenant_name,notes,created_at FROM properties WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Property;for rows.Next(){var e Property;rows.Scan(&e.ID,&e.Name,&e.Address,&e.Type,&e.Bedrooms,&e.Rent,&e.Status,&e.TenantName,&e.Notes,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM properties GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}
