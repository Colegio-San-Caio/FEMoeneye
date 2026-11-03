package main

import (
"encoding/json"
"fmt"
"log"
"net/http"
"os"
"sync"
"time"
)
type Event struct{TS time.Time `json:"ts"`; Page string `json:"page"`; Event string `json:"event"`; Flops float64 `json:"flops"`}
var (events []Event; mu sync.Mutex; totalFlops float64; DB="eal7.sqlite.json")

func loadDB(){b,err:=os.ReadFile(DB); if err!=nil{return}; var d struct{Events []Event; Total float64}; if json.Unmarshal(b,&d)==nil{events=d.Events; totalFlops=d.Total; fmt.Printf("Loaded DB: %d events, %.0f flops\n",len(events),totalFlops)}}
func saveDB(){d:=struct{Events []Event; Total float64}{events,totalFlops}; b,_:=json.Marshal(d); os.WriteFile(DB,b,0644)}

func track(w http.ResponseWriter, r *http.Request){
w.Header().Set("Access-Control-Allow-Origin","*")
var e Event; json.NewDecoder(r.Body).Decode(&e)
if e.Page==""{e.Page=r.URL.Query().Get("page")}; if e.Event==""{e.Event=r.URL.Query().Get("event"); if e.Event==""{e.Event="page_view"}}
e.TS=time.Now()
switch e.Event{case "page_view":e.Flops=4; case "fem_run":e.Flops=16; case "ximg_view":e.Flops=64; default:e.Flops=4}
mu.Lock(); events=append(events,e); totalFlops+=e.Flops; saveDB(); mu.Unlock()
json.NewEncoder(w).Encode(map[string]interface{}{"status":"ok","total_flops":totalFlops})
fmt.Printf("[eal7 QUAD DB] %s | %s | %.0f | total %.0f\n",e.Event,e.Page,e.Flops,totalFlops)
}
func matrix(w http.ResponseWriter, r *http.Request){
w.Header().Set("Access-Control-Allow-Origin","*")
mu.Lock(); defer mu.Unlock()
m:=map[string]map[string]int{}; for _,e:=range events{if m[e.Page]==nil{m[e.Page]=map[string]int{}}; m[e.Page][e.Event]++}
json.NewEncoder(w).Encode(map[string]interface{}{"matrix":m,"total":len(events),"total_flops_quadruple":totalFlops,"db":DB,"persistent":true,"port":"8038 flop quadruple"})
}
func js(w http.ResponseWriter, r *http.Request){w.Header().Set("Access-Control-Allow-Origin","*"); w.Header().Set("Content-Type","application/javascript"); w.Write([]byte(`(function(){var EP="http://localhost:8038/track";function s(ev){var d={page:location.pathname,event:ev};try{navigator.sendBeacon(EP,JSON.stringify(d))}catch(e){fetch(EP,{method:"POST",body:JSON.stringify(d)})}}window.oenLytics={track:s};s('page_view')})();`))}

func main(){
loadDB()
http.HandleFunc("/track",track); http.HandleFunc("/matrix",matrix); http.HandleFunc("/oenlytics.js",js)
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
w.Header().Set("Access-Control-Allow-Origin","*")
if _,err:=os.Stat("eal7.html"); err==nil{http.ServeFile(w,r,"eal7.html"); return}
w.Write([]byte(`<html><script src="/oenlytics.js"></script><body><h1>eal7 QUAD 8038 PERSISTENT</h1><pre id=m></pre><script>setInterval(()=>fetch('/matrix').then(r=>r.json()).then(d=>m.textContent=JSON.stringify(d,null,2)),1000)</script></body></html>`))
})
port:=os.Getenv("PORT"); if port==""{port="8038"}
fmt.Printf("oenLytics.x eal7.x FLOP QUADRUPLE PERSISTENT on :%s DB=%s\n",port,DB)
log.Fatal(http.ListenAndServe(":"+port,nil))
}
