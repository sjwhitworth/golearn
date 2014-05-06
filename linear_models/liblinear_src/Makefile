liblinear: tron.o linear.o blas.a
	$(CXX) -shared -o linear.dll tron.o linear.o blas/blas.a
	cp linear.dll liblinear.so
tron.o: tron.cpp
	$(CXX) -c -fpic tron.cpp -o tron.o
linear.o: linear.cpp
	$(CXX) -c -fpic linear.cpp -o linear.o
blas.a:	
	make -C blas
clean:
	rm -f *.o
	make -C blas clean
	
install:
	cp linear.dll ..
	cp liblinear.so ..
