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
type Event struct {
TS time.Time `json:"ts"`
Page string `json:"page"`
Event string `json:"event"`
Flops float64 `json:"flops"` // quadruple FLOP count
Source string `json:"source"`
}
var (events []Event; mu sync.Mutex; totalFlops float64)

func track(w http.ResponseWriter, r *http.Request){
w.Header().Set("Access-Control-Allow-Origin","*")
var e Event
json.NewDecoder(r.Body).Decode(&e)
if e.Page=="" { e.Page=r.URL.Query().Get("page") }
if e.Event=="" { e.Event=r.URL.Query().Get("event"); if e.Event=="" { e.Event="page_view" } }
e.TS=time.Now()
if e.Source=="" { e.Source="oen" }
// FLOP QUADRUPLE: 4x precision for FEM matrix
// page_view=4 FLOPS, fem_run= 4*4=16, ximg=4*4*4=64
switch e.Event {
case "page_view": e.Flops=4
case "fem_run": e.Flops=16
case "ximg_view": e.Flops=64
default: e.Flops=4
}
mu.Lock()
events=append(events,e)
totalFlops+=e.Flops
mu.Unlock()
w.Header().Set("Content-Type","application/json")
json.NewEncoder(w).Encode(map[string]interface{}{"status":"ok","port":"8038","flop_quadruple":e.Flops,"total_flops":totalFlops})
fmt.Printf("[eal7 QUAD] %s | %s | %.0f FLOPS | total %.0f\n", e.Event, e.Page, e.Flops, totalFlops)
}
func matrix(w http.ResponseWriter, r *http.Request){
w.Header().Set("Access-Control-Allow-Origin","*")
mu.Lock(); defer mu.Unlock()
m := map[string]map[string]interface{}{}
for _, e := range events {
if m[e.Page]==nil { m[e.Page]=map[string]interface{}{} }
c := 0
if v, ok := m[e.Page][e.Event]; ok { c = v.(int) }
m[e.Page][e.Event]=c+1
}
json.NewEncoder(w).Encode(map[string]interface{}{"matrix":m,"total":len(events),"total_flops_quadruple":totalFlops,"port":"8038 flop quadruple"})
}
func js(w http.ResponseWriter, r *http.Request){
w.Header().Set("Access-Control-Allow-Origin","*")
w.Header().Set("Content-Type","application/javascript")
w.Write([]byte(`(function(){var EP="http://localhost:8038/track";function s(ev,p){var d={page:location.pathname,event:ev,params:p||{},source:"oen"};try{navigator.sendBeacon(EP,JSON.stringify(d))}catch(e){fetch(EP,{method:"POST",body:JSON.stringify(d),mode:"no-cors"})}}window.oenLytics={track:s,pageview:function(){s('page_view')}};s('page_view')})();`))
}
func main(){
port:=os.Getenv("PORT")
if port=="" { port="8038" } // 38 or 80 as flop quadruple = 8038
http.HandleFunc("/track", track)
http.HandleFunc("/collect", track)
http.HandleFunc("/matrix", matrix)
http.HandleFunc("/oenlytics.js", js)
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
w.Header().Set("Access-Control-Allow-Origin","*")
if r.URL.Path=="/"{ w.Header().Set("Content-Type","text/html"); w.Write([]byte(`<html><script src="/oenlytics.js"></script><body><h1>eal7 QUAD 8038</h1><div id=m></div><script>fetch('/matrix').then(r=>r.json()).then(d=>document.getElementById('m').innerHTML='<pre>'+JSON.stringify(d,null,2)+'</pre>')</script></body></html>`)); return }
http.NotFound(w,r)
})
fmt.Printf("oenLytics.x eal7.x FLOP QUADRUPLE on :%s - Port 38 or 80 as 8038\n", port)
log.Fatal(http.ListenAndServe(":"+port, nil))
}
