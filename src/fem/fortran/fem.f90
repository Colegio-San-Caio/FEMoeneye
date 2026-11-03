module fem_direct
contains
subroutine stiffness_triangle(E, nu, t, coords, ke)
 real(8), intent(in) :: E, nu, t, coords(6)
 real(8), intent(out) :: ke(6,6)
 real(8) :: A, B(3,6), D(3,3)
 A=0.5d0*abs((coords(3)-coords(1))*(coords(6)-coords(2))-(coords(5)-coords(1))*(coords(4)-coords(2)))
 B=0.0d0; B(1,1)=coords(4)-coords(6); B(1,3)=coords(6)-coords(2); B(1,5)=coords(2)-coords(4)
 B(2,2)=coords(5)-coords(3); B(2,4)=coords(1)-coords(5); B(2,6)=coords(3)-coords(1)
 B(3,1)=B(2,2); B(3,2)=B(1,1); B(3,3)=B(2,4); B(3,4)=B(1,3); B(3,5)=B(2,6); B(3,6)=B(1,5)
 B=B/(2.0d0*A); D=0.0d0; D(1,1)=1.0d0; D(1,2)=nu; D(2,1)=nu; D(2,2)=1.0d0; D(3,3)=(1.0d0-nu)/2.0d0; D=D*E/(1.0d0-nu*nu)
 ke=t*A*matmul(transpose(B),matmul(D,B))
end subroutine
end module
program femoeneye_main
 use fem_direct; implicit none; real(8)::coords(6)=[0.0d0,0.0d0,1.0d0,0.0d0,0.0d0,1.0d0]; real(8)::ke(6,6)
 call stiffness_triangle(210.0d9,0.3d0,0.01d0,coords,ke)
 print *, "FEMoeneye Fortran ERP CAM - XIMG MIME application/vnd.x-mosfetq-dos-image"
end program
