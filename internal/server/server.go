package server
import ("encoding/json";"net/http";"github.com/stockyard-dev/stockyard-homestead/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux;limits Limits}
func New(db *store.DB,limits Limits)*Server{s:=&Server{db:db,mux:http.NewServeMux(),limits:limits}
s.mux.HandleFunc("GET /api/servers",s.list)
s.mux.HandleFunc("POST /api/servers",s.create)
s.mux.HandleFunc("GET /api/servers/{id}",s.get)
s.mux.HandleFunc("PUT /api/servers/{id}",s.update)
s.mux.HandleFunc("DELETE /api/servers/{id}",s.del)
s.mux.HandleFunc("GET /api/stats",s.stats)
s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /api/tier",func(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"tier":s.limits.Tier,"upgrade_url":"https://stockyard.dev/homestead/"})})
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root)
return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)list(w http.ResponseWriter,r *http.Request){servers:=s.db.List();if servers==nil{servers=[]store.Server{}};wj(w,200,map[string]any{"servers":servers})}
func(s *Server)create(w http.ResponseWriter,r *http.Request){if s.limits.MaxItems>0&&s.db.Count()>=s.limits.MaxItems{we(w,402,"Free tier limit reached");return};var srv store.Server;json.NewDecoder(r.Body).Decode(&srv);if srv.Hostname==""{we(w,400,"hostname required");return};s.db.Create(&srv);wj(w,201,s.db.Get(srv.ID))}
func(s *Server)get(w http.ResponseWriter,r *http.Request){srv:=s.db.Get(r.PathValue("id"));if srv==nil{we(w,404,"not found");return};wj(w,200,srv)}
func(s *Server)update(w http.ResponseWriter,r *http.Request){existing:=s.db.Get(r.PathValue("id"));if existing==nil{we(w,404,"not found");return};var patch store.Server;json.NewDecoder(r.Body).Decode(&patch);patch.ID=existing.ID;if patch.Hostname==""{patch.Hostname=existing.Hostname};s.db.Update(&patch);wj(w,200,s.db.Get(patch.ID))}
func(s *Server)del(w http.ResponseWriter,r *http.Request){s.db.Delete(r.PathValue("id"));wj(w,200,map[string]string{"status":"deleted"})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,s.db.Stats())}
func(s *Server)health(w http.ResponseWriter,r *http.Request){st:=s.db.Stats();wj(w,200,map[string]any{"service":"homestead","status":"ok","servers":st["total"]})}
