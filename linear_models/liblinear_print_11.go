// +build go1.1
// +build !go1.2
// +build !go1.3

package linear_models

import "C"

//export libLinearPrintFunc
func libLinearPrintFunc(str *C.char) {
	// Stubbed
}

func libLinearHookPrintFunc() {
	// Stubbed
}
