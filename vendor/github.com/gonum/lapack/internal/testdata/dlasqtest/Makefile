FC = gfortran
FFLAGS = -O2

targets = testdlasq1 testdlasq2 testdlasq3 testdlasq4
objects = dlamch.o lsame.o xerbla.o ieeeck.o iparmq.o ilaenv.o dlasrt.o \
          dlaisnan.o disnan.o dlasq6.o dlasq5.o dlasq4.o dlasq3.o dlasq2.o

default : $(targets)
.PHONY : default

testdlasq1 : testdlasq1.o $(objects) dcopy.o dlas2.o dlascl.o dlasq1.o

testdlasq2 : testdlasq2.o $(objects)

testdlasq3 : testdlasq3.o $(objects)

testdlasq4 : testdlasq4.o $(objects)

% : %.o
	$(FC) $(FFLAGS) $^ -o $@

%.o : %.f90
	$(FC) $(FFLAGS) -c -o $@ $<

clean :
	rm -f *.o $(targets) *.txt
.PHONY : clean