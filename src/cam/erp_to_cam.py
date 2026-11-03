"""ERP -> CAM: FEM result to XIMG + G-code"""
def fem_to_gcode(displacements, nodes):
    g = ["G21 ; ERP FEMoeneye","G90"]
    for n, (x,y) in enumerate(nodes):
        dx, dy = displacements[n*2], displacements[n*2+1]
        g.append(f"G1 X{x+dx:.4f} Y{y+dy:.4f} ; node {n} stress")
    g.append("M30 ; end FEMoeneye CAM")
    return "\n".join(g)

def fem_to_ximg(displacements):
    # Write XIMG header for BOA viewer - IANA 1460172
    header = b"XIMG"
    body = displacements.tobytes()
    return header + body
