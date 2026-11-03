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
Timestamp time.Time `json:"ts"`
Page string `json:"page"`
Event string `json:"event"`
Params map[string]interface{} `json:"params"`
Source string `json:"source"`
}
var (events []Event; mu sync.Mutex)

func trackHandler(w http.ResponseWriter, r *http.Request){
var e Event
json.NewDecoder(r.Body).Decode(&e)
if e.Page=="" { e.Page=r.URL.Query().Get("page") }
if e.Event=="" { e.Event=r.URL.Query().Get("event"); if e.Event=="" { e.Event="page_view" } }
e.Timestamp=time.Now()
if e.Source=="" { e.Source="oen" }
if e.Params==nil { e.Params=map[string]interface{}{} }
mu.Lock(); events=append(events,e); mu.Unlock()
w.Header().Set("Access-Control-Allow-Origin","*")
w.Header().Set("Content-Type","application/json")
w.Write([]byte(`{"status":"ok","analytics":"oenLytics.x eal7 matrix"}`))
fmt.Printf("[eal7] %s | %s | %s\n", e.Event, e.Page, e.Source)
}

func statsHandler(w http.ResponseWriter, r *http.Request){
mu.Lock(); defer mu.Unlock()
w.Header().Set("Access-Control-Allow-Origin","*")
w.Header().Set("Content-Type","application/json")
start := 0
if len(events) > 20 { start = len(events)-20 }
var last []Event
if len(events)>0 { last = events[start:] }
json.NewEncoder(w).Encode(map[string]interface{}{
"total":len(events), "engine":"oenLytics.x eal7", "events":last,
})
}

func matrixHandler(w http.ResponseWriter, r *http.Request){
mu.Lock(); defer mu.Unlock()
w.Header().Set("Access-Control-Allow-Origin","*")
w.Header().Set("Content-Type","application/json")
// Build tracking matrix: pages x events
pages := map[string]map[string]int{}
for _, e := range events {
if pages[e.Page]==nil { pages[e.Page]=map[string]int{} }
pages[e.Page][e.Event]++
}
json.NewEncoder(w).Encode(map[string]interface{}{
"matrix":pages, "total":len(events), "type":"eal7 tracking matrix",
})
}

func jsHandler(w http.ResponseWriter, r *http.Request){
w.Header().Set("Content-Type","application/javascript")
w.Write([]byte(`
(function(){
  var EP="http://localhost:8082/track";
  function send(ev,params){
    var d={page:location.pathname,event:ev,params:params||{},source:"oen"};
    navigator.sendBeacon(EP, JSON.stringify(d));
    if(window.gtag) gtag('event',ev,params);
  }
  window.oenLytics={track:send, pageview:function(){send('page_view',{title:document.title})}};
  window.oenLytics.pageview();
})();
`))
}

func main(){
http.HandleFunc("/track", trackHandler)
http.HandleFunc("/collect", trackHandler)
http.HandleFunc("/stats", statsHandler)
http.HandleFunc("/matrix", matrixHandler)
http.HandleFunc("/oenlytics.js", jsHandler)
fmt.Println("oenLytics.x eal7 matrix on :8082")
log.Fatal(http.ListenAndServe(":8082", nil))
}
