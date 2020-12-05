package lib

/*
#include <stdlib.h>
*/
import "C"

func Random() int {
	return int(C.rand())
}

func Seed(i int) {
	C.srand(C.uint(i))
}
