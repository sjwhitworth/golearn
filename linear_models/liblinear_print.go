// +build go1.2

package linear_models

/*
#cgo LDFLAGS: -llinear
#cgo CFLAGS:
#include <linear.h>

typedef void (*print_func)(char *);
void golearn_liblinear_print_func_cgo(char *);
*/
import "C"

import (
	"fmt"
	"unsafe"
)

//export libLinearPrintFunc
func libLinearPrintFunc(str *C.char) {
	fmt.Println(C.GoString(str))
}

func libLinearHookPrintFunc() {
	C.set_print_string_function((C.print_func)(unsafe.Pointer(C.golearn_liblinear_print_func_cgo)))
}
