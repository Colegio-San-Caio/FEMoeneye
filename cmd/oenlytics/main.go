package main

import (
"encoding/json"
"fmt"
"log"
"net/http"
"os"
"strings"
"sync"
"time"
)
type Event struct{TS time.Time `json:"ts"`; Page, Event, Nickname, Code, Role string; Flops float64; Points, Port int; BOA string}
var (events []Event; mu sync.Mutex; totalFlops float64; DB="eal7.sqlite.json")
func loadDB(){b,err:=os.ReadFile(DB); if err!=nil{return}; var d struct{Events []Event; Total float64}; if json.Unmarshal(b,&d)==nil{events=d.Events; totalFlops=d.Total; fmt.Printf("Loaded DB: %d events %.0f flops\n",len(events),totalFlops)}}
func saveDB(){d:=struct{Events []Event; Total float64}{events,totalFlops}; b,_:=json.Marshal(d); os.WriteFile(DB,b,0644)}

func parseKDP(s string) Event{
e:=Event{TS:time.Now(), Port:8038, BOA:"application/vnd.x-mosfetq-dos-image", Flops:64}
for _,line:=range strings.Split(s,"\n"){
kv:=strings.SplitN(line,":",2); if len(kv)!=2{continue}
k:=strings.TrimSpace(kv[0]); v:=strings.TrimSpace(kv[1])
switch k{
case "Nickname": e.Nickname=v
case "Code": e.Code=v
case "Role": e.Role=v
case "Points": fmt.Sscanf(v,"%d",&e.Points)
case "Page": e.Page=v
case "Event": e.Event=v
case "Flops": fmt.Sscanf(v,"%f",&e.Flops)
}
}
if e.Event==""{e.Event="ximg_view"}; if e.Page==""{e.Page="main.ximg"}
return e
}

func track(w http.ResponseWriter, r *http.Request){
w.Header().Set("Access-Control-Allow-Origin","*")
body, _ := os.ReadFile("kdp_profile.txt")
var e Event
if len(body)>0 { e=parseKDP(string(body)) } else { json.NewDecoder(r.Body).Decode(&e) }
if e.Flops==0 { switch e.Event{case "page_view":e.Flops=4; case "fem_run":e.Flops=16; case "ximg_view":e.Flops=64; default:e.Flops=4} }
mu.Lock(); events=append(events,e); totalFlops+=e.Flops; saveDB(); mu.Unlock()
json.NewEncoder(w).Encode(map[string]interface{}{"status":"ok","kdp":e,"total_flops":totalFlops})
fmt.Printf("[ADAPTED MAIN] KDP %s Code %s | %s | %.0f flops | total %.0f\n",e.Nickname,e.Code,e.Event,e.Flops,totalFlops)
}
func matrix(w http.ResponseWriter, r *http.Request){
w.Header().Set("Access-Control-Allow-Origin","*")
mu.Lock(); defer mu.Unlock()
m:=map[string]map[string]int{}; for _,ev:=range events{if m[ev.Page]==nil{m[ev.Page]=map[string]int{}}; m[ev.Page][ev.Event]++}
json.NewEncoder(w).Encode(map[string]interface{}{"matrix":m,"total":len(events),"total_flops_quadruple":totalFlops,"kdp_profiles":len(events),"port":"8038 flop quadruple - adapted to all main files"})
}
func main(){
loadDB()
http.HandleFunc("/track",track); http.HandleFunc("/matrix",matrix)
http.HandleFunc("/oenlytics.js",func(w http.ResponseWriter, r *http.Request){w.Header().Set("Access-Control-Allow-Origin","*"); w.Header().Set("Content-Type","application/javascript"); w.Write([]byte(`(function(){var EP="/track";function s(ev){fetch(EP,{method:"POST",body:JSON.stringify({page:location.pathname,event:ev})}) }window.oenLytics={track:s};s('page_view')})();`))})
http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){w.Header().Set("Access-Control-Allow-Origin","*"); if _,err:=os.Stat("eal7.html");err==nil{http.ServeFile(w,r,"eal7.html"); return}; http.ServeFile(w,r,"eal7.html")})
port:=os.Getenv("PORT"); if port==""{port="8038"}
fmt.Printf("ADAPTED MAIN - KDP + eal7 QUAD on :%s\n",port)
log.Fatal(http.ListenAndServe(":"+port,nil))
}
