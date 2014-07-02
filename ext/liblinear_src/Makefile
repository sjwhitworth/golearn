liblinear: tron.o linear.o blas.a	
	$(CXX) -shared -o linear.dll tron.o linear.o blas/blas.a	

tron.o: tron.cpp
	$(CXX) -fpic -c tron.cpp -o tron.o

linear.o: linear.cpp
	$(CXX) -fpic -c linear.cpp -o linear.o

blas.a:	
	make -C blas

clean:
	$(RM) *.o
	make -C blas clean
	
