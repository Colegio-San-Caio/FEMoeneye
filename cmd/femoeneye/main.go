package main
import ("encoding/json"; "fmt"; "log"; "net/http")
type ERPJob struct{ Nodes [][]float64 `json:"nodes"`; Elems [][]int `json:"elems"`}
func SolveERPJob(nodes [][]float64) []float64{ dof:=len(nodes)*2; u:=make([]float64,dof); if dof>=2 {u[dof-2]=-0.001}; return u }
func femHandler(w http.ResponseWriter, r *http.Request){
 var job ERPJob; json.NewDecoder(r.Body).Decode(&job)
 if len(job.Nodes)==0 {job.Nodes=[][]float64{{0,0},{1,0},{0,1}}}
 u:=SolveERPJob(job.Nodes)
 w.Header().Set("Content-Type","text/plain"); w.Header().Set("X-FEMoeneye-MIME","application/vnd.x-mosfetq-dos-image"); w.Header().Set("X-FEMoeneye-Lang","go+fortran")
 fmt.Fprintf(w,"G21 ; FEMoeneye Go+Fortran ERP CAM\n")
 for i,n:=range job.Nodes{ fmt.Fprintf(w,"G1 X%.4f Y%.4f ; node %d\n", n[0]+u[i*2], n[1]+u[i*2+1], i) }
 fmt.Fprintf(w,"M30 ; XIMG\n")
}
func main(){ http.HandleFunc("/erp/fem", femHandler); fmt.Println("FEMoeneye Go+Fortran on :8080"); log.Fatal(http.ListenAndServe(":8080", nil)) }
