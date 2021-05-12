This set of codes generates tests for the dlasq* routines.

The high level routines are testdlasq*.f90. The testdlasq*.f90 routines are intended in combination with the routines in gonum/lapack/testlapack. 
During execution, these high-level routines record the inputs and outputs
to the executed subroutines. For example, testdlasq3, when executed,
generates files gen4tests.txt and gen5tests.txt that record the inputs
and outputs to the evaluation of dlasq4 and dlasq5 respectively. The output 
format in gen*tests.txt is the struct literal that matches the respective test
in gonum/lapack/testlapack/dlasq*.go. Thus, these generated tests can be copied
into the testing routine to test the native implementation.

The testing routines in testlapack have code for generating inputs for these 
fortran routines. Typically, one would isolate the particular failing test,
and modify the testlapack routine to print it to the terminal, for example
in testlapack/dlasq3.go one might add

    printDlasq3FortranInput(test)
    os.Exit(1)

This prints variable initialization for dlasq3 routine to the terminal, which
can be copied and pasted into testdlasq3.f90. Please note that this process
is not completely automated. Some of the other initialization may need to
change, particularly the size allocation of the data array(s).