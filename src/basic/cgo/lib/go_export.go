package lib

/*
#include <stdio.h>
extern void CFunction1();
*/
import "C"

import "fmt"

//export GoFunction1
func GoFunction1() {
	fmt.Println("GoFunction1() is called.")
}

func CallCFunc() {
	C.CFunction1()
}
