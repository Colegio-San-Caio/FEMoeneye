"""FEMoeneye direct FEM solver - 2D plane stress"""
import numpy as np

def stiffness_triangle(E, nu, t, coords):
    """Direct stiffness for CST triangle"""
    # Simplified - for ERP CAM proof
    x1,y1,x2,y2,x3,y3 = coords.flatten()
    A = 0.5*abs((x2-x1)*(y3-y1)-(x3-x1)*(y2-y1))
    B = np.array([[y2-y3,0,y3-y1,0,y1-y2,0],
                  [0,x3-x2,0,x1-x3,0,x2-x1],
                  [x3-x2,y2-y3,x1-x3,y3-y1,x2-x1,y1-y2]])/(2*A)
    D = E/(1-nu**2)*np.array([[1,nu,0],[nu,1,0],[0,0,(1-nu)/2]])
    k = t*A*B.T@D@B
    return k

def solve_erp_job(nodes, elems, E=210e9, nu=0.3):
    """nodes: (n,2), elems: (m,3) -> displacements"""
    dof = len(nodes)*2
    K = np.zeros((dof,dof))
    for el in elems:
        coords = nodes[el]
        ke = stiffness_triangle(E,nu,0.01,coords)
        for i in range(3):
            for j in range(3):
                K[el[i]*2:el[i]*2+2, el[j]*2:el[j]*2+2] += ke[i*2:i*2+2, j*2:j*2+2]
    # ERP boundary: fix first node
    K[0:2,:]=0; K[:,0:2]=0; K[0,0]=1; K[1,1]=1
    F = np.zeros(dof); F[-2]= -1000 # load from ERP work order
    u = np.linalg.solve(K,F)
    return u

if __name__=="__main__":
    nodes = np.array([[0,0],[1,0],[0,1]])
    elems = np.array([[0,1,2]])
    print(solve_erp_job(nodes,elems))
