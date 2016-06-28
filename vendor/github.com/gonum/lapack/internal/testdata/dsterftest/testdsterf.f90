program testdsterf
implicit none
integer, parameter :: n = 4
real(kind=8), dimension(n) :: d
real(kind=8), dimension(n-1) :: e
integer :: info,i

d(1:4) = (/1D+00, 3D+00, 4D+00, 6D+00/)
e(1:3) = (/2D+00, 4D+00, 5D+00/)

call dsterf(n,d,e,info)
DO i = 1, n
    print *, d(i)
end do
end