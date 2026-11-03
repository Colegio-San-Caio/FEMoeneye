package main

import (
"encoding/json"
"fmt"
"log"
"net/http"
"sync"
"time"
)
type Event struct {
TS time.Time `json:"ts"`
Page string `json:"page"`
Event string `json:"event"`
Source string `json:"source"`
Params map[string]interface{} `json:"params"`
}
var (events []Event; mu sync.Mutex)

func track(w http.ResponseWriter, r *http.Request){
var e Event
json.NewDecoder(r.Body).Decode(&e)
if e.Page=="" { e.Page=r.URL.Query().Get("page") }
if e.Event=="" { e.Event=r.URL.Query().Get("event"); if e.Event=="" { e.Event="page_view" } }
e.TS=time.Now()
if e.Source=="" { e.Source="oen" }
mu.Lock(); events=append(events,e); mu.Unlock()
w.Header().Set("Access-Control-Allow-Origin","*")
w.Header().Set("Content-Type","application/json")
json.NewEncoder(w).Encode(map[string]interface{}{"status":"ok","matrix":"eal7","total":len(events)})
fmt.Printf("[eal7] %s | %s\n", e.Event, e.Page)
}
func stats(w http.ResponseWriter, r *http.Request){
mu.Lock(); defer mu.Unlock()
w.Header().Set("Access-Control-Allow-Origin","*")
w.Header().Set("Content-Type","application/json")
start:=0; if len(events)>50 { start=len(events)-50 }
var last []Event; if len(events)>0 { last=events[start:] }
json.NewEncoder(w).Encode(map[string]interface{}{"total":len(events),"events":last})
}
func matrix(w http.ResponseWriter, r *http.Request){
mu.Lock(); defer mu.Unlock()
w.Header().Set("Access-Control-Allow-Origin","*")
w.Header().Set("Content-Type","application/json")
pages:=map[string]map[string]int{}
for _, e := range events {
if pages[e.Page]==nil { pages[e.Page]=map[string]int{} }
pages[e.Page][e.Event]++
}
json.NewEncoder(w).Encode(map[string]interface{}{"matrix":pages,"total":len(events),"type":"eal7 tracking matrix - Pages x Events"})
}
func js(w http.ResponseWriter, r *http.Request){
w.Header().Set("Content-Type","application/javascript")
w.Header().Set("Access-Control-Allow-Origin","*")
w.Write([]byte(`(function(){var EP="http://localhost:8082/track";function s(ev,p){var d={page:location.pathname,event:ev,params:p||{},source:"oen"};try{navigator.sendBeacon(EP,JSON.stringify(d))}catch(e){fetch(EP,{method:"POST",body:JSON.stringify(d)})}if(window.gtag)gtag('event',ev,p)}window.oenLytics={track:s,pageview:function(){s('page_view',{title:document.title})}};s('page_view',{title:document.title})})();`))
}
func main(){
http.HandleFunc("/track", track)
http.HandleFunc("/collect", track)
http.HandleFunc("/stats", stats)
http.HandleFunc("/matrix", matrix)
http.HandleFunc("/oenlytics.js", js)
// serve eal7.html directly from Go too
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
if r.URL.Path=="/" { http.ServeFile(w,r,"eal7.html"); return }
http.NotFound(w,r)
})
fmt.Println("oenLytics.x eal7 matrix on :8082 - FIXED")
log.Fatal(http.ListenAndServe(":8082", nil))
}
