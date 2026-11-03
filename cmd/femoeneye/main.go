package main
import ("encoding/json"; "fmt"; "log"; "net/http")
type ERPJob struct{ Nodes [][]float64 `json:"nodes"`; Elems [][]int `json:"elems"`}
func SolveERPJob(nodes [][]float64) []float64{ dof:=len(nodes)*2; u:=make([]float64,dof); if dof>=2 {u[dof-2]=-0.001}; return u }
func femHandler(w http.ResponseWriter, r *http.Request){
 if r.Method!="POST" { http.Error(w,"POST only - FEMoeneye ERP CAM",405); return }
 var job ERPJob; json.NewDecoder(r.Body).Decode(&job)
 if len(job.Nodes)==0 {job.Nodes=[][]float64{{0,0},{1,0},{0,1}}}
 u:=SolveERPJob(job.Nodes)
 w.Header().Set("Content-Type","text/plain"); w.Header().Set("X-FEMoeneye-MIME","application/vnd.x-mosfetq-dos-image")
 fmt.Fprintf(w,"G21 ; FEMoeneye Go+Fortran ERP CAM 208bed9 MAIN .gitignore\n")
 for i,n:=range job.Nodes{ fmt.Fprintf(w,"G1 X%.4f Y%.4f ; node %d\n", n[0]+u[i*2], n[1]+u[i*2+1], i) }
 fmt.Fprintf(w,"M30 ; XIMG IANA 1460172\n")
}
func main(){ http.HandleFunc("/erp/fem", femHandler); fmt.Println("FEMoeneye MAIN .gitignore 208bed9 on :8081"); log.Fatal(http.ListenAndServe(":8081", nil)) }
