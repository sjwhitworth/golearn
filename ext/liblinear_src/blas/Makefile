AR     = ar rcv
RANLIB = ranlib 

HEADERS = blas.h blasp.h
FILES = dnrm2.o daxpy.o ddot.o dscal.o 

CFLAGS = $(OPTFLAGS) 
FFLAGS = $(OPTFLAGS)

blas: $(FILES) $(HEADERS)
	$(AR) blas.a $(FILES)  
	$(RANLIB) blas.a

clean:
	- $(RM) -f *.o
	- $(RM) -f *.a
	- $(RM) -f *~

.c.o:
	$(CC) $(CFLAGS) -fpic -c $*.c


