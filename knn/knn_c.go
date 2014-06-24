package knn

/*

#cgo CFLAGS: -O3 -std=c99
#cgo LDFLAGS: -lm
#include <string.h>
#include <stdio.h>
#include <math.h>

void kcompute(double *train, double *test, double *out, int attr_sz, int train_sz) {
    for (int i = 0; i < train_sz; i++) {
        double *ref = train + (attr_sz * i);
//        fprintf(stderr, "A Ref: %d, %.2f\n", attr_sz * i, *ref);
        double *dist = out + i;
        double d = 0;
        for (int j = 0; j < attr_sz; j++) {
            double a = *(train + (attr_sz * i) + j);
            double b = *(test + j);
            d += (a - b) * (a - b);
//            fprintf(stderr, "A: %.2f B: %.2f TEMP: %.2f %d\n", a, b, d, attr_sz * i + j);
        }
//        fprintf(stderr, "I: %d RESULT: %.2f\n\n", i, sqrt(d));
        *dist = sqrt(d);
    }
}

*/
import "C"

import (
	"github.com/sjwhitworth/golearn/base"
	"unsafe"
)

func computeDistances(train base.PondStorageRef, testing []float64, out []float64) {
	trainingRefPt := unsafe.Pointer(&train.Storage[0])
	cTrainingRef := (*C.double)(trainingRefPt)

	testingRefPt := unsafe.Pointer(&testing[0])
	cTestingRef := (*C.double)(testingRefPt)

	outPt := unsafe.Pointer(&out[0])
	cOutPt := (*C.double)(outPt)

	testingLen := C.int((len(testing)))
	trainLen := C.int(train.Rows)

	C.kcompute(cTrainingRef, cTestingRef, cOutPt, testingLen, trainLen)
}
