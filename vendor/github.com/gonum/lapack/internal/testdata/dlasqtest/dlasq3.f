*> \brief \b DLASQ3 checks for deflation, computes a shift and calls dqds. Used by sbdsqr.
*
*  =========== DOCUMENTATION ===========
*
* Online html documentation available at 
*            http://www.netlib.org/lapack/explore-html/ 
*
*> \htmlonly
*> Download DLASQ3 + dependencies 
*> <a href="http://www.netlib.org/cgi-bin/netlibfiles.tgz?format=tgz&filename=/lapack/lapack_routine/dlasq3.f"> 
*> [TGZ]</a> 
*> <a href="http://www.netlib.org/cgi-bin/netlibfiles.zip?format=zip&filename=/lapack/lapack_routine/dlasq3.f"> 
*> [ZIP]</a> 
*> <a href="http://www.netlib.org/cgi-bin/netlibfiles.txt?format=txt&filename=/lapack/lapack_routine/dlasq3.f"> 
*> [TXT]</a>
*> \endhtmlonly 
*
*  Definition:
*  ===========
*
*       SUBROUTINE DLASQ3( I0, N0, Z, PP, DMIN, SIGMA, DESIG, QMAX, NFAIL,
*                          ITER, NDIV, IEEE, TTYPE, DMIN1, DMIN2, DN, DN1,
*                          DN2, G, TAU )
* 
*       .. Scalar Arguments ..
*       LOGICAL            IEEE
*       INTEGER            I0, ITER, N0, NDIV, NFAIL, PP
*       DOUBLE PRECISION   DESIG, DMIN, DMIN1, DMIN2, DN, DN1, DN2, G,
*      $                   QMAX, SIGMA, TAU
*       ..
*       .. Array Arguments ..
*       DOUBLE PRECISION   Z( * )
*       ..
*  
*
*> \par Purpose:
*  =============
*>
*> \verbatim
*>
*> DLASQ3 checks for deflation, computes a shift (TAU) and calls dqds.
*> In case of failure it changes shifts, and tries again until output
*> is positive.
*> \endverbatim
*
*  Arguments:
*  ==========
*
*> \param[in] I0
*> \verbatim
*>          I0 is INTEGER
*>         First index.
*> \endverbatim
*>
*> \param[in,out] N0
*> \verbatim
*>          N0 is INTEGER
*>         Last index.
*> \endverbatim
*>
*> \param[in] Z
*> \verbatim
*>          Z is DOUBLE PRECISION array, dimension ( 4*N )
*>         Z holds the qd array.
*> \endverbatim
*>
*> \param[in,out] PP
*> \verbatim
*>          PP is INTEGER
*>         PP=0 for ping, PP=1 for pong.
*>         PP=2 indicates that flipping was applied to the Z array   
*>         and that the initial tests for deflation should not be 
*>         performed.
*> \endverbatim
*>
*> \param[out] DMIN
*> \verbatim
*>          DMIN is DOUBLE PRECISION
*>         Minimum value of d.
*> \endverbatim
*>
*> \param[out] SIGMA
*> \verbatim
*>          SIGMA is DOUBLE PRECISION
*>         Sum of shifts used in current segment.
*> \endverbatim
*>
*> \param[in,out] DESIG
*> \verbatim
*>          DESIG is DOUBLE PRECISION
*>         Lower order part of SIGMA
*> \endverbatim
*>
*> \param[in] QMAX
*> \verbatim
*>          QMAX is DOUBLE PRECISION
*>         Maximum value of q.
*> \endverbatim
*>
*> \param[out] NFAIL
*> \verbatim
*>          NFAIL is INTEGER
*>         Number of times shift was too big.
*> \endverbatim
*>
*> \param[out] ITER
*> \verbatim
*>          ITER is INTEGER
*>         Number of iterations.
*> \endverbatim
*>
*> \param[out] NDIV
*> \verbatim
*>          NDIV is INTEGER
*>         Number of divisions.
*> \endverbatim
*>
*> \param[in] IEEE
*> \verbatim
*>          IEEE is LOGICAL
*>         Flag for IEEE or non IEEE arithmetic (passed to DLASQ5).
*> \endverbatim
*>
*> \param[in,out] TTYPE
*> \verbatim
*>          TTYPE is INTEGER
*>         Shift type.
*> \endverbatim
*>
*> \param[in,out] DMIN1
*> \verbatim
*>          DMIN1 is DOUBLE PRECISION
*> \endverbatim
*>
*> \param[in,out] DMIN2
*> \verbatim
*>          DMIN2 is DOUBLE PRECISION
*> \endverbatim
*>
*> \param[in,out] DN
*> \verbatim
*>          DN is DOUBLE PRECISION
*> \endverbatim
*>
*> \param[in,out] DN1
*> \verbatim
*>          DN1 is DOUBLE PRECISION
*> \endverbatim
*>
*> \param[in,out] DN2
*> \verbatim
*>          DN2 is DOUBLE PRECISION
*> \endverbatim
*>
*> \param[in,out] G
*> \verbatim
*>          G is DOUBLE PRECISION
*> \endverbatim
*>
*> \param[in,out] TAU
*> \verbatim
*>          TAU is DOUBLE PRECISION
*>
*>         These are passed as arguments in order to save their values
*>         between calls to DLASQ3.
*> \endverbatim
*
*  Authors:
*  ========
*
*> \author Univ. of Tennessee 
*> \author Univ. of California Berkeley 
*> \author Univ. of Colorado Denver 
*> \author NAG Ltd. 
*
*> \date September 2012
*
*> \ingroup auxOTHERcomputational
*
*  =====================================================================
      SUBROUTINE DLASQ3( I0, N0, Z, PP, DMIN, SIGMA, DESIG, QMAX, NFAIL,
     $                   ITER, NDIV, IEEE, TTYPE, DMIN1, DMIN2, DN, DN1,
     $                   DN2, G, TAU )
*
*  -- LAPACK computational routine (version 3.4.2) --
*  -- LAPACK is a software package provided by Univ. of Tennessee,    --
*  -- Univ. of California Berkeley, Univ. of Colorado Denver and NAG Ltd..--
*     September 2012
*
*     .. Scalar Arguments ..
      LOGICAL            IEEE
      INTEGER            I0, ITER, N0, NDIV, NFAIL, PP
      DOUBLE PRECISION   DESIG, DMIN, DMIN1, DMIN2, DN, DN1, DN2, G,
     $                   QMAX, SIGMA, TAU
*     ..
*     .. Array Arguments ..
      DOUBLE PRECISION   Z( * )
*     ..
*
*  =====================================================================
*
*     .. Parameters ..
      DOUBLE PRECISION   CBIAS
      PARAMETER          ( CBIAS = 1.50D0 )
      DOUBLE PRECISION   ZERO, QURTR, HALF, ONE, TWO, HUNDRD
      PARAMETER          ( ZERO = 0.0D0, QURTR = 0.250D0, HALF = 0.5D0,
     $                     ONE = 1.0D0, TWO = 2.0D0, HUNDRD = 100.0D0 )
*     ..
*     .. Local Scalars ..
      INTEGER            IPN4, J4, N0IN, NN, TTYPE
      DOUBLE PRECISION   EPS, S, T, TEMP, TOL, TOL2
*     ..
*     .. External Subroutines ..
      EXTERNAL           DLASQ4, DLASQ5, DLASQ6
*     ..
*     .. External Function ..
      DOUBLE PRECISION   DLAMCH
      LOGICAL            DISNAN
      EXTERNAL           DISNAN, DLAMCH
*     ..
*     .. Intrinsic Functions ..
      INTRINSIC          ABS, MAX, MIN, SQRT
*     ..
*     .. Executable Statements ..
*

      N0IN = N0
      EPS = DLAMCH( 'Precision' )
      TOL = EPS*HUNDRD
      TOL2 = TOL**2
*
*     Check for deflation.
*
   10 CONTINUE
*
      IF( N0.LT.I0 )
     $   RETURN
      IF( N0.EQ.I0 )
     $   GO TO 20
      NN = 4*N0 + PP
      IF( N0.EQ.( I0+1 ) )
     $   GO TO 40
*
*     Check whether E(N0-1) is negligible, 1 eigenvalue.
*
      IF( Z( NN-5 ).GT.TOL2*( SIGMA+Z( NN-3 ) ) .AND.
     $    Z( NN-2*PP-4 ).GT.TOL2*Z( NN-7 ) )
     $   GO TO 30
*
   20 CONTINUE
*
      Z( 4*N0-3 ) = Z( 4*N0+PP-3 ) + SIGMA
      N0 = N0 - 1
      GO TO 10
*
*     Check  whether E(N0-2) is negligible, 2 eigenvalues.
*
   30 CONTINUE
*
      IF( Z( NN-9 ).GT.TOL2*SIGMA .AND.
     $    Z( NN-2*PP-8 ).GT.TOL2*Z( NN-11 ) )
     $   GO TO 50
*
   40 CONTINUE
*
      IF( Z( NN-3 ).GT.Z( NN-7 ) ) THEN
         S = Z( NN-3 )
         Z( NN-3 ) = Z( NN-7 )
         Z( NN-7 ) = S
      END IF
      T = HALF*( ( Z( NN-7 )-Z( NN-3 ) )+Z( NN-5 ) )
      IF( Z( NN-5 ).GT.Z( NN-3 )*TOL2.AND.T.NE.ZERO ) THEN
         S = Z( NN-3 )*( Z( NN-5 ) / T )
         IF( S.LE.T ) THEN
            S = Z( NN-3 )*( Z( NN-5 ) /
     $          ( T*( ONE+SQRT( ONE+S / T ) ) ) )
         ELSE
            S = Z( NN-3 )*( Z( NN-5 ) / ( T+SQRT( T )*SQRT( T+S ) ) )
         END IF
         T = Z( NN-7 ) + ( S+Z( NN-5 ) )
         Z( NN-3 ) = Z( NN-3 )*( Z( NN-7 ) / T )
         Z( NN-7 ) = T
      END IF
      Z( 4*N0-7 ) = Z( NN-7 ) + SIGMA
      Z( 4*N0-3 ) = Z( NN-3 ) + SIGMA
      N0 = N0 - 2
      GO TO 10
*
   50 CONTINUE
      IF( PP.EQ.2 ) 
     $   PP = 0
*
*     Reverse the qd-array, if warranted.
*

      IF( DMIN.LE.ZERO .OR. N0.LT.N0IN ) THEN
         IF( CBIAS*Z( 4*I0+PP-3 ).LT.Z( 4*N0+PP-3 ) ) THEN
            IPN4 = 4*( I0+N0 )
            DO 60 J4 = 4*I0, 2*( I0+N0-1 ), 4
               TEMP = Z( J4-3 )
               Z( J4-3 ) = Z( IPN4-J4-3 )
               Z( IPN4-J4-3 ) = TEMP
               TEMP = Z( J4-2 )
               Z( J4-2 ) = Z( IPN4-J4-2 )
               Z( IPN4-J4-2 ) = TEMP
               TEMP = Z( J4-1 )
               Z( J4-1 ) = Z( IPN4-J4-5 )
               Z( IPN4-J4-5 ) = TEMP
               TEMP = Z( J4 )
               Z( J4 ) = Z( IPN4-J4-4 )
               Z( IPN4-J4-4 ) = TEMP
   60       CONTINUE
            IF( N0-I0.LE.4 ) THEN
               Z( 4*N0+PP-1 ) = Z( 4*I0+PP-1 )
               Z( 4*N0-PP ) = Z( 4*I0-PP )
            END IF
            DMIN2 = MIN( DMIN2, Z( 4*N0+PP-1 ) )
            Z( 4*N0+PP-1 ) = MIN( Z( 4*N0+PP-1 ), Z( 4*I0+PP-1 ),
     $                            Z( 4*I0+PP+3 ) )
            Z( 4*N0-PP ) = MIN( Z( 4*N0-PP ), Z( 4*I0-PP ),
     $                          Z( 4*I0-PP+4 ) )
            QMAX = MAX( QMAX, Z( 4*I0+PP-3 ), Z( 4*I0+PP+1 ) )
            DMIN = -ZERO
         END IF
      END IF
*
*     Choose a shift.
*
      ! Print out DLASQ4 test cases
      write(4,*) "{"
      write(4,'(9999(g0))',advance="no") "z: []float64{"
      do i = 1, NN
        write (4,'(99999(e24.16,a))',advance="no") z(i), ","
      end do
      write (4,*) "},"
      write (4,*) "i0: ", I0, ","
      write (4,*) "n0: ", N0, ","
      write (4,*) "pp:   ", PP, ","
      write (4,*) "n0in: ", N0IN, ","
      write (4,*) "dmin: ", DMIN, ","
      write (4,*) "dmin1:", DMIN1, ","
      write (4,*) "dmin2:", DMIN2, ","
      write (4,*) "dn:   ", DN, ","
      write (4,*) "dn1:  ", DN1, ","
      write (4,*) "dn2:  ", DN2, ","
      write (4,*) "tau: ", TAU, ","
      write (4,*) "ttype: ", TTYPE, ","
      write (4,*) "g:    ", G, ","
      CALL DLASQ4( I0, N0, Z, PP, N0IN, DMIN, DMIN1, DMIN2, DN, DN1,
     $             DN2, TAU, TTYPE, G )

      write(4,'(9999(g0))',advance="no") "zOut: []float64{"
      do i = 1, NN
        write (4,'(99999(e24.16,a))',advance="no") z(i), ","
      end do
      write (4,*) "},"
      write (4,*) "tauOut: ", TAU, ","
      write (4,*) "ttypeOut: ", TTYPE, ","
      write (4,*) "gOut:   ", G, ","
      write (4,*) "},"

*
*     Call dqds until DMIN > 0.
*
   70 CONTINUE
*

      write(5,*) "{"
      write(5,'(9999(g0))',advance="no") "z: []float64{"
      do i = 1, NN
        write (5,'(99999(e24.16,a))',advance="no") z(i), ","
      end do
      write (5,*) "},"
      write (5,*) "i0: ", I0, ","
      write (5,*) "n0: ", N0, ","
      write (5,*) "pp:   ", PP, ","
      write (5,*) "tau: ", TAU, ","
      write (5,*) "sigma: ", SIGMA, ","
      write (5,*) "dmin: ", DMIN, ","
      write (5,*) "dmin1:", DMIN1, ","
      write (5,*) "dmin2:", DMIN2, ","
      write (5,*) "dn:   ", DN, ","
      write (5,*) "dnm1:  ", DN1, ","
      write (5,*) "dnm2:  ", DN2, ","


      CALL DLASQ5( I0, N0, Z, PP, TAU, SIGMA, DMIN, DMIN1, DMIN2, DN,
     $             DN1, DN2, IEEE, EPS )



      write (5,*) "i0Out: ", I0, ","
      write (5,*) "n0Out: ", N0, ","
      write (5,*) "ppOut:   ", PP, ","
      write (5,*) "tauOut: ", TAU, ","
      write (5,*) "sigmaOut: ", SIGMA, ","
      write (5,*) "dminOut: ", DMIN, ","
      write (5,*) "dmin1Out:", DMIN1, ","
      write (5,*) "dmin2Out:", DMIN2, ","
      write (5,*) "dnOut:   ", DN, ","
      write (5,*) "dnm1Out:  ", DN1, ","
      write (5,*) "dnm2Out:  ", DN2, ","
      write (5,*) "},"

*
      NDIV = NDIV + ( N0-I0+2 )

      ITER = ITER + 1
*
*     Check status.
*

      IF( DMIN.GE.ZERO .AND. DMIN1.GE.ZERO ) THEN
*
*        Success.
*
         GO TO 90
*
      ELSE IF( DMIN.LT.ZERO .AND. DMIN1.GT.ZERO .AND. 
     $         Z( 4*( N0-1 )-PP ).LT.TOL*( SIGMA+DN1 ) .AND.
     $         ABS( DN ).LT.TOL*SIGMA ) THEN

*
*        Convergence hidden by negative DN.
*
         Z( 4*( N0-1 )-PP+2 ) = ZERO
         DMIN = ZERO
         GO TO 90
      ELSE IF( DMIN.LT.ZERO ) THEN

*
*        TAU too big. Select new TAU and try again.
*
         NFAIL = NFAIL + 1
         IF( TTYPE.LT.-22 ) THEN
*
*           Failed twice. Play it safe.
*
            TAU = ZERO
         ELSE IF( DMIN1.GT.ZERO ) THEN
*
*           Late failure. Gives excellent shift.
*
            TAU = ( TAU+DMIN )*( ONE-TWO*EPS )
            TTYPE = TTYPE - 11
         ELSE
*
*           Early failure. Divide by 4.
*
            TAU = QURTR*TAU
            TTYPE = TTYPE - 12
         END IF
         GO TO 70
      ELSE IF( DISNAN( DMIN ) ) THEN
*
*        NaN.
*
         IF( TAU.EQ.ZERO ) THEN
            GO TO 80
         ELSE
            TAU = ZERO
            GO TO 70
         END IF
      ELSE
*            
*        Possible underflow. Play it safe.
*
         GO TO 80
      END IF
*
*     Risk of underflow.
*
   80 CONTINUE

      CALL DLASQ6( I0, N0, Z, PP, DMIN, DMIN1, DMIN2, DN, DN1, DN2 )


      NDIV = NDIV + ( N0-I0+2 )
      ITER = ITER + 1
      TAU = ZERO
*
   90 CONTINUE

      IF( TAU.LT.SIGMA ) THEN
         DESIG = DESIG + TAU
         T = SIGMA + DESIG
         DESIG = DESIG - ( T-SIGMA )
      ELSE
         T = SIGMA + TAU
         DESIG = SIGMA - ( T-TAU ) + DESIG
      END IF
      SIGMA = T
*
      RETURN
*
*     End of DLASQ3
*
      END
