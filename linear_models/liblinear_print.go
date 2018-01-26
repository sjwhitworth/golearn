// +build go1.2

package linear_models

/*
#include "linear.h"

typedef void (*print_func)(char *);
void golearn_liblinear_print_func_cgo(char *);
*/
import "C"

import (
	"github.com/amclay/golearn/base"
	"unsafe"
)

//export libLinearPrintFunc
func libLinearPrintFunc(str *C.char) {
	base.Logger.Println(C.GoString(str))
}

func libLinearHookPrintFunc() {
	C.set_print_string_function((C.print_func)(unsafe.Pointer(C.golearn_liblinear_print_func_cgo)))
}
