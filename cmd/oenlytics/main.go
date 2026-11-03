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
Timestamp time.Time              `json:"ts"`
IP        string                 `json:"ip"`
UA        string                 `json:"ua"`
Page      string                 `json:"page"`
Event     string                 `json:"event"`
Params    map[string]interface{} `json:"params"`
Source    string                 `json:"source"` // ga4, oen, erp, ximg
}

var (
events []Event
mu     sync.Mutex
)

func trackHandler(w http.ResponseWriter, r *http.Request) {
// Support both: own oenLytics and GA4 MP proxy
var e Event
if r.Method == "POST" {
json.NewDecoder(r.Body).Decode(&e)
} else {
e.Page = r.URL.Query().Get("page")
e.Event = r.URL.Query().Get("event")
if e.Event == "" { e.Event = "page_view" }
}
e.Timestamp = time.Now()
e.IP = r.RemoteAddr
e.UA = r.UserAgent()
if e.Source == "" { e.Source = "oen" }
if e.Params == nil { e.Params = map[string]interface{}{} }

mu.Lock()
events = append(events, e)
mu.Unlock()

// CORS for JS
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Content-Type", "application/json")
w.Write([]byte(`{"status":"ok","analytics":"oenLytics.x","iana":"1460172"}`))

fmt.Printf("[oenLytics] %s | %s | %s | %s\n", e.Timestamp.Format("15:04:05"), e.Event, e.Page, e.Source)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
mu.Lock()
defer mu.Unlock()
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
json.NewEncoder(w).Encode(map[string]interface{}{
"total_events": len(events),
"engine": "oenLytics.x",
"ga4_compatible": true,
"self_hosted": true,
"mime": "application/vnd.x-mosfetq-dos-image",
"events": events[len(events)-20:],
})
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/javascript")
w.Write([]byte(`
// oenLytics.js - Own Analytics - Google Analytics Support
(function(){
  var ENDPOINT = location.origin + ":8082/track";
  if(location.port=="8081") ENDPOINT="http://localhost:8082/track";
  function send(event, params){
    var data={page:location.pathname+location.search, event:event, params:params||{}, source:"oen", ts:Date.now()};
    navigator.sendBeacon(ENDPOINT, JSON.stringify(data));
    // Also send to GA4 if gtag exists - support mode
    if(window.gtag){ gtag('event', event, params); }
  }
  window.oenLytics = { track: send, pageview: function(){ send('page_view',{title:document.title}) } };
  window.oenLytics.pageview();
  console.log("[oenLytics.x] Own analytics active - GA support mode - IANA 1460172");
})();
`))
}

func main() {
http.HandleFunc("/track", trackHandler)
http.HandleFunc("/collect", trackHandler) // GA4 /g/collect compatibility
http.HandleFunc("/stats", statsHandler)
http.HandleFunc("/oenlytics.js", jsHandler)
fmt.Println("oenLytics.x - Own Analytics on :8082")
fmt.Println("  /oenlytics.js -> JS snippet")
fmt.Println("  /track -> POST events")
fmt.Println("  /collect -> GA4 compatible")
fmt.Println("  /stats -> dashboard")
log.Fatal(http.ListenAndServe(":8082", nil))
}
