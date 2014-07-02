# Platform detection
ifeq ($(OS),Windows_NT)
    ARTEFACT="linear.dll"
else
    ARTEFACT="liblinear.so"
endif

liblinear: tron.o linear.o blas.a	
	$(CXX) -shared -o $(ARTEFACT) tron.o linear.o blas/blas.a	

tron.o: tron.cpp
	$(CXX) -fpic -c tron.cpp -o tron.o

linear.o: linear.cpp
	$(CXX) -fpic -c linear.cpp -o linear.o

blas.a:	
	make -C blas

clean:
	$(RM) *.o
	make -C blas clean
	
