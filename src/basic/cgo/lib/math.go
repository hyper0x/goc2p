package lib

/*
#cgo LDFLAGS: -lm
#include <math.h>
*/
import "C"

func Sqrt(p float32) (float32, error) {
	result, err := C.sqrt(C.double(p))
	return float32(result), err
}
