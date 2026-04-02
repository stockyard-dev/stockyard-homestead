package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-homestead/internal/server";"github.com/stockyard-dev/stockyard-homestead/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="8880"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./homestead-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("homestead: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Homestead — Self-hosted server inventory\n  ─────────────────────────────────\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Data:       %s\n  ─────────────────────────────────\n\n",port,port,dataDir)
log.Printf("homestead: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
