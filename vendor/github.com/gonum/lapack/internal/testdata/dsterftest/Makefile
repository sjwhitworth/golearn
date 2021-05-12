FC = gfortran
FFLAGS = -O2

targets = testdsterf
objects = disnan.o dlamch.o dlanst.o dlaisnan.o dlassq.o dlapy2.o lsame.o dlae2.o dlascl.o dlasrt.o dsterf.o xerbla.o

default : $(targets)
.PHONY : default

testdsterf : testdsterf.o $(objects)

% : %.o
	$(FC) $(FFLAGS) $^ -o $@

%.o : %.f90
	$(FC) $(FFLAGS) -c -o $@ $<

clean :
	rm -f *.o $(targets) *.txt
.PHONY : clean