"""Apache ERP bridge - receives ERP workorder, returns CAM"""
from http.server import BaseHTTPRequestHandler, HTTPServer
import json
from src.fem.direct import solve_erp_job
from src.cam.erp_to_cam import fem_to_gcode
import numpy as np

class ERPHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        length = int(self.headers.get('Content-Length',0))
        data = json.loads(self.rfile.read(length) or b'{}')
        nodes = np.array(data.get('nodes',[[0,0],[1,0],[0,1]]))
        elems = np.array(data.get('elems',[[0,1,2]]))
        u = solve_erp_job(nodes, elems)
        gcode = fem_to_gcode(u, nodes)
        self.send_response(200)
        self.send_header('Content-Type','text/plain')
        self.send_header('X-FEMoeneye-MIME','application/vnd.x-mosfetq-dos-image')
        self.end_headers()
        self.wfile.write(gcode.encode())

if __name__=="__main__":
    print("FEMoeneye ERP CAM on :8080 - Apache proxy to this")
    HTTPServer(("0.0.0.0",8080), ERPHandler).serve_forever()
