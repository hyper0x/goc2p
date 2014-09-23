package main

import (
	cgolib "basic/cgo/lib"
	"fmt"
)

func main() {
	input := float32(2.33)
	output, err := cgolib.Sqrt(input)
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
	fmt.Printf("The square root of %f is %f.\n", input, output)

	cgolib.Print("ABC\n")

	cgolib.CallCFunc()
}
