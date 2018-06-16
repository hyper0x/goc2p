package lib

/*
#include <stdio.h>
#include <stdlib.h>

void myprint(char* s) {
        printf("%s", s);
}
*/
import "C"
import "unsafe"

func Print(s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.myprint(cs)
}
