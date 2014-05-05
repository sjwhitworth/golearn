liblinear: tron.o linear.o blas.a
	g++ -shared -o linear.dll tron.o linear.o blas/blas.a
tron.o: tron.cpp
	g++ -c -fpic tron.cpp -o tron.o
linear.o: linear.cpp
	g++ -c -fpic linear.cpp -o linear.o
blas.a:	
	make -C blas
	