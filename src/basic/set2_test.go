package basic

import (
	"fmt"
	"math/rand"
	"runtime/debug"
	"sort"
	"testing"
)

func TestSimpleSet(t *testing.T) {
	debugTag := false
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	n := 10
	r := 3
	seeds := make([]int, r)
	for k, _ := range seeds {
		seeds[k] = k
	}
	t.Logf("Seeds: %v\n", seeds)
	t.Logf("Seeds Length: %v\n", len(seeds))
	matrix := make([][]int, n)
	for i, _ := range matrix {
		indexs := make([]int, r)
		for j, _ := range indexs {
			m := rand.Intn(3)
			indexs[j] = seeds[m]
		}
		if debugTag {
			t.Logf("%v (i=%d)\n", indexs, i)
		}
		matrix[i] = indexs
	}
	t.Logf("Matrix: %v\n", matrix)
	t.Logf("Matrix Length: %v\n", len(matrix))
	matrixSet := SimpleSet{KeyGenerator: generateKey, Comparator: compare}
	for _, v := range matrix {
		matrixSet.Add(interface{}(v))
	}
	t.Logf("Matrix Set: %v\n", matrixSet)
	t.Logf("Matrix Set Length: %v\n", matrixSet.Len())
	matrixSlice := matrixSet.Slice()
	t.Logf("Matrix Sorted Slice: %v\n", matrixSlice)
	matrixIterator := matrixSet.Iterator()
	var pe interface{}
	for j := 0; j < matrixSet.Len(); j++ {
		e, has := matrixIterator()
		if !has {
			break
		}
		if pe != nil {
			if compare(pe, e) > 0 {
				t.Errorf("Error: %v should be LE %v. (j=%d)\n", pe, e, j)
				t.FailNow()
			}
			if debugTag {
				t.Logf("%v <= %v. (j=%d)\n", pe, e, j)
			}
		}
		pe = e
	}
	randElement := matrixSlice[rand.Intn(len(matrixSlice))]
	t.Logf("Rand Elements: %v\n", randElement)
	if !matrixSet.Contain(randElement) {
		t.Errorf("Error: The element '%v' shoud be in marix set '%v'.\n", randElement, matrixSet)
	}
	t.Logf("The matrix contains element '%v'.\n", randElement)
	if !matrixSet.Remove(randElement) {
		t.Errorf("Error: Remove element '%v' shoud be successful.\n", randElement)
	}
	t.Logf("The element '%v' is removed.\n", randElement)
	if matrixSet.Contain(randElement) {
		t.Errorf("Error: The removed element '%v' shoud not be in marix set '%v'.\n", randElement, matrixSet)
	}
	t.Logf("The matrix not contains element '%v'.\n", randElement)
	if matrixSet.Remove(randElement) {
		t.Errorf("Error: Remove removed element '%v' shoud not failing.\n", randElement)
	}
	t.Logf("Can not remove removed element '%v'.\n", randElement)
	if !matrixSet.Clear() {
		t.Errorf("Error: Clear matrix should be successful.\n", randElement)
	}
	t.Logf("The matrix is cleared.\n")
	if matrixSet.Len() > 0 {
		t.Errorf("Error: The length of matrix should be 0.\n", randElement)
	}
	t.Logf("The length of matrix is 0.\n")
}

func generateKey(x interface{}) string {
	xa := interface{}(x).([]int)
	xac := make([]int, len(xa))
	copy(xac, xa)
	sort.Ints(xac)
	return fmt.Sprintf("%v", xac)
}

func compare(i, j interface{}) int {
	ia := interface{}(i).([]int)
	ja := interface{}(j).([]int)
	sort.Ints(ia)
	sort.Ints(ja)
	il := len(ia)
	jl := len(ja)
	result := 0
	if il < jl {
		result = -1
	} else if il > jl {
		result = 1
	} else {
		for i, iv := range ia {
			jv := ja[i]
			if iv != jv {
				if iv < jv {
					result = -1
				} else if iv > jv {
					result = 1
				}
				break
			}
		}
	}
	return result
}
