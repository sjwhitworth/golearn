*> \brief \b DLAHR2 reduces the specified number of first columns of a general rectangular matrix A so that elements below the specified subdiagonal are zero, and returns auxiliary matrices which are needed to apply the transformation to the unreduced part of A.
*
*  =========== DOCUMENTATION ===========
*
* Online html documentation available at 
*            http://www.netlib.org/lapack/explore-html/ 
*
*> \htmlonly
*> Download DLAHR2 + dependencies 
*> <a href="http://www.netlib.org/cgi-bin/netlibfiles.tgz?format=tgz&filename=/lapack/lapack_routine/dlahr2.f"> 
*> [TGZ]</a> 
*> <a href="http://www.netlib.org/cgi-bin/netlibfiles.zip?format=zip&filename=/lapack/lapack_routine/dlahr2.f"> 
*> [ZIP]</a> 
*> <a href="http://www.netlib.org/cgi-bin/netlibfiles.txt?format=txt&filename=/lapack/lapack_routine/dlahr2.f"> 
*> [TXT]</a>
*> \endhtmlonly 
*
*  Definition:
*  ===========
*
*       SUBROUTINE DLAHR2( N, K, NB, A, LDA, TAU, T, LDT, Y, LDY )
* 
*       .. Scalar Arguments ..
*       INTEGER            K, LDA, LDT, LDY, N, NB
*       ..
*       .. Array Arguments ..
*       DOUBLE PRECISION  A( LDA, * ), T( LDT, NB ), TAU( NB ),
*      $                   Y( LDY, NB )
*       ..
*  
*
*> \par Purpose:
*  =============
*>
*> \verbatim
*>
*> DLAHR2 reduces the first NB columns of A real general n-BY-(n-k+1)
*> matrix A so that elements below the k-th subdiagonal are zero. The
*> reduction is performed by an orthogonal similarity transformation
*> Q**T * A * Q. The routine returns the matrices V and T which determine
*> Q as a block reflector I - V*T*V**T, and also the matrix Y = A * V * T.
*>
*> This is an auxiliary routine called by DGEHRD.
*> \endverbatim
*
*  Arguments:
*  ==========
*
*> \param[in] N
*> \verbatim
*>          N is INTEGER
*>          The order of the matrix A.
*> \endverbatim
*>
*> \param[in] K
*> \verbatim
*>          K is INTEGER
*>          The offset for the reduction. Elements below the k-th
*>          subdiagonal in the first NB columns are reduced to zero.
*>          K < N.
*> \endverbatim
*>
*> \param[in] NB
*> \verbatim
*>          NB is INTEGER
*>          The number of columns to be reduced.
*> \endverbatim
*>
*> \param[in,out] A
*> \verbatim
*>          A is DOUBLE PRECISION array, dimension (LDA,N-K+1)
*>          On entry, the n-by-(n-k+1) general matrix A.
*>          On exit, the elements on and above the k-th subdiagonal in
*>          the first NB columns are overwritten with the corresponding
*>          elements of the reduced matrix; the elements below the k-th
*>          subdiagonal, with the array TAU, represent the matrix Q as a
*>          product of elementary reflectors. The other columns of A are
*>          unchanged. See Further Details.
*> \endverbatim
*>
*> \param[in] LDA
*> \verbatim
*>          LDA is INTEGER
*>          The leading dimension of the array A.  LDA >= max(1,N).
*> \endverbatim
*>
*> \param[out] TAU
*> \verbatim
*>          TAU is DOUBLE PRECISION array, dimension (NB)
*>          The scalar factors of the elementary reflectors. See Further
*>          Details.
*> \endverbatim
*>
*> \param[out] T
*> \verbatim
*>          T is DOUBLE PRECISION array, dimension (LDT,NB)
*>          The upper triangular matrix T.
*> \endverbatim
*>
*> \param[in] LDT
*> \verbatim
*>          LDT is INTEGER
*>          The leading dimension of the array T.  LDT >= NB.
*> \endverbatim
*>
*> \param[out] Y
*> \verbatim
*>          Y is DOUBLE PRECISION array, dimension (LDY,NB)
*>          The n-by-nb matrix Y.
*> \endverbatim
*>
*> \param[in] LDY
*> \verbatim
*>          LDY is INTEGER
*>          The leading dimension of the array Y. LDY >= N.
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
*> \ingroup doubleOTHERauxiliary
*
*> \par Further Details:
*  =====================
*>
*> \verbatim
*>
*>  The matrix Q is represented as a product of nb elementary reflectors
*>
*>     Q = H(1) H(2) . . . H(nb).
*>
*>  Each H(i) has the form
*>
*>     H(i) = I - tau * v * v**T
*>
*>  where tau is a real scalar, and v is a real vector with
*>  v(1:i+k-1) = 0, v(i+k) = 1; v(i+k+1:n) is stored on exit in
*>  A(i+k+1:n,i), and tau in TAU(i).
*>
*>  The elements of the vectors v together form the (n-k+1)-by-nb matrix
*>  V which is needed, with T and Y, to apply the transformation to the
*>  unreduced part of the matrix, using an update of the form:
*>  A := (I - V*T*V**T) * (A - Y*V**T).
*>
*>  The contents of A on exit are illustrated by the following example
*>  with n = 7, k = 3 and nb = 2:
*>
*>     ( a   a   a   a   a )
*>     ( a   a   a   a   a )
*>     ( a   a   a   a   a )
*>     ( h   h   a   a   a )
*>     ( v1  h   a   a   a )
*>     ( v1  v2  a   a   a )
*>     ( v1  v2  a   a   a )
*>
*>  where a denotes an element of the original matrix A, h denotes a
*>  modified element of the upper Hessenberg matrix H, and vi denotes an
*>  element of the vector defining H(i).
*>
*>  This subroutine is a slight modification of LAPACK-3.0's DLAHRD
*>  incorporating improvements proposed by Quintana-Orti and Van de
*>  Gejin. Note that the entries of A(1:K,2:NB) differ from those
*>  returned by the original LAPACK-3.0's DLAHRD routine. (This
*>  subroutine is not backward compatible with LAPACK-3.0's DLAHRD.)
*> \endverbatim
*
*> \par References:
*  ================
*>
*>  Gregorio Quintana-Orti and Robert van de Geijn, "Improving the
*>  performance of reduction to Hessenberg form," ACM Transactions on
*>  Mathematical Software, 32(2):180-194, June 2006.
*>
*  =====================================================================
      SUBROUTINE DLAHR2( N, K, NB, A, LDA, TAU, T, LDT, Y, LDY )
*
*  -- LAPACK auxiliary routine (version 3.4.2) --
*  -- LAPACK is a software package provided by Univ. of Tennessee,    --
*  -- Univ. of California Berkeley, Univ. of Colorado Denver and NAG Ltd..--
*     September 2012
*
*     .. Scalar Arguments ..
      INTEGER            K, LDA, LDT, LDY, N, NB
*     ..
*     .. Array Arguments ..
      DOUBLE PRECISION  A( LDA, * ), T( LDT, NB ), TAU( NB ),
     $                   Y( LDY, NB )
*     ..
*
*  =====================================================================
*
*     .. Parameters ..
      DOUBLE PRECISION  ZERO, ONE
      PARAMETER          ( ZERO = 0.0D+0, 
     $                     ONE = 1.0D+0 )
*     ..
*     .. Local Scalars ..
      INTEGER            I
      DOUBLE PRECISION  EI
*     ..
*     .. External Subroutines ..
      EXTERNAL           DAXPY, DCOPY, DGEMM, DGEMV, DLACPY,
     $                   DLARFG, DSCAL, DTRMM, DTRMV
*     ..
*     .. Intrinsic Functions ..
      INTRINSIC          MIN
*     ..
*     .. Executable Statements ..
*
*     Quick return if possible
*
      IF( N.LE.1 )
     $   RETURN
*
      DO 10 I = 1, NB
         IF( I.GT.1 ) THEN
*
*           Update A(K+1:N,I)
*
*           Update I-th column of A - Y * V**T
*
            CALL DGEMV( 'NO TRANSPOSE', N-K, I-1, -ONE, Y(K+1,1), LDY,
     $                  A( K+I-1, 1 ), LDA, ONE, A( K+1, I ), 1 )
*
*           Apply I - V * T**T * V**T to this column (call it b) from the
*           left, using the last column of T as workspace
*
*           Let  V = ( V1 )   and   b = ( b1 )   (first I-1 rows)
*                    ( V2 )             ( b2 )
*
*           where V1 is unit lower triangular
*
*           w := V1**T * b1
*
            CALL DCOPY( I-1, A( K+1, I ), 1, T( 1, NB ), 1 )
            CALL DTRMV( 'Lower', 'Transpose', 'UNIT', 
     $                  I-1, A( K+1, 1 ),
     $                  LDA, T( 1, NB ), 1 )
*
*           w := w + V2**T * b2
*
            CALL DGEMV( 'Transpose', N-K-I+1, I-1, 
     $                  ONE, A( K+I, 1 ),
     $                  LDA, A( K+I, I ), 1, ONE, T( 1, NB ), 1 )
*
*           w := T**T * w
*
            CALL DTRMV( 'Upper', 'Transpose', 'NON-UNIT', 
     $                  I-1, T, LDT,
     $                  T( 1, NB ), 1 )
*
*           b2 := b2 - V2*w
*
            CALL DGEMV( 'NO TRANSPOSE', N-K-I+1, I-1, -ONE, 
     $                  A( K+I, 1 ),
     $                  LDA, T( 1, NB ), 1, ONE, A( K+I, I ), 1 )
*
*           b1 := b1 - V1*w
*
            CALL DTRMV( 'Lower', 'NO TRANSPOSE', 
     $                  'UNIT', I-1,
     $                  A( K+1, 1 ), LDA, T( 1, NB ), 1 )
            CALL DAXPY( I-1, -ONE, T( 1, NB ), 1, A( K+1, I ), 1 )
*
            A( K+I-1, I-1 ) = EI
         END IF
*
*        Generate the elementary reflector H(I) to annihilate
*        A(K+I+1:N,I)
*
         CALL DLARFG( N-K-I+1, A( K+I, I ), A( MIN( K+I+1, N ), I ), 1,
     $                TAU( I ) )
         EI = A( K+I, I )
         A( K+I, I ) = ONE
*
*        Compute  Y(K+1:N,I)
*
         CALL DGEMV( 'NO TRANSPOSE', N-K, N-K-I+1, 
     $               ONE, A( K+1, I+1 ),
     $               LDA, A( K+I, I ), 1, ZERO, Y( K+1, I ), 1 )
         CALL DGEMV( 'Transpose', N-K-I+1, I-1, 
     $               ONE, A( K+I, 1 ), LDA,
     $               A( K+I, I ), 1, ZERO, T( 1, I ), 1 )
         CALL DGEMV( 'NO TRANSPOSE', N-K, I-1, -ONE, 
     $               Y( K+1, 1 ), LDY,
     $               T( 1, I ), 1, ONE, Y( K+1, I ), 1 )
         CALL DSCAL( N-K, TAU( I ), Y( K+1, I ), 1 )
*
*        Compute T(1:I,I)
*
         CALL DSCAL( I-1, -TAU( I ), T( 1, I ), 1 )
         CALL DTRMV( 'Upper', 'No Transpose', 'NON-UNIT', 
     $               I-1, T, LDT,
     $               T( 1, I ), 1 )
         T( I, I ) = TAU( I )
*
   10 CONTINUE
      A( K+NB, NB ) = EI
*
*     Compute Y(1:K,1:NB)
*
      CALL DLACPY( 'ALL', K, NB, A( 1, 2 ), LDA, Y, LDY )
      CALL DTRMM( 'RIGHT', 'Lower', 'NO TRANSPOSE', 
     $            'UNIT', K, NB,
     $            ONE, A( K+1, 1 ), LDA, Y, LDY )
      IF( N.GT.K+NB )
     $   CALL DGEMM( 'NO TRANSPOSE', 'NO TRANSPOSE', K, 
     $               NB, N-K-NB, ONE,
     $               A( 1, 2+NB ), LDA, A( K+1+NB, 1 ), LDA, ONE, Y,
     $               LDY )
      CALL DTRMM( 'RIGHT', 'Upper', 'NO TRANSPOSE', 
     $            'NON-UNIT', K, NB,
     $            ONE, T, LDT, Y, LDY )
*
      RETURN
*
*     End of DLAHR2
*
      END
