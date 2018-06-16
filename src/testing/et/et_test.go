package et

import (
	"fmt"
	"testing"
)

func ExampleHello() {
	for i := 0; i < 3; i++ {
		fmt.Println("Hello, Golang~")
	}

	// Output: Hello, Golang~
	// Hello, Golang~
	// Hello, Golang~
}

func TestOne(t *testing.T) {
	t.Log("Hi~")
}
