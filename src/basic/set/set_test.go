package set

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime/debug"
	"strings"
	"testing"
	"time"
)

func testSetLenAndContains(t *testing.T, newSet func() Set, typeName string) {
	t.Logf("Starting Test%sLenAndContains...", typeName)
	set, expectedElemMap := genRandSet(newSet)
	t.Logf("Got a %s value: %v.", typeName, set)
	expectedLen := len(expectedElemMap)
	if set.Len() != expectedLen {
		t.Errorf("ERROR: The length of %s value %d is not %d!\n",
			set.Len(), typeName, expectedLen)
		t.FailNow()
	}
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
	for k := range expectedElemMap {
		if !set.Contains(k) {
			t.Errorf("ERROR: The %s value %v do not contains %v!",
				set, typeName, k)
			t.FailNow()
		}
	}
}

func testSetAdd(t *testing.T, newSet func() Set, typeName string) {
	t.Logf("Starting Test%sAdd...", typeName)
	set := newSet()
	var randElem interface{}
	var result bool
	expectedElemMap := make(map[interface{}]bool)
	for i := 0; i < 5; i++ {
		randElem = genRandElement()
		t.Logf("Add %v to the %s value %v.\n", randElem, typeName, set)
		result = set.Add(randElem)
		if expectedElemMap[randElem] && result {
			t.Errorf("ERROR: The element adding (%v => %v) is successful but should be failing!\n",
				randElem, set)
			t.FailNow()
		}
		if !expectedElemMap[randElem] && !result {
			t.Errorf("ERROR: The element adding (%v => %v) is failing!\n",
				randElem, set)
			t.FailNow()
		}
		expectedElemMap[randElem] = true
	}
	t.Logf("The %s value: %v.", typeName, set)
	expectedLen := len(expectedElemMap)
	if set.Len() != expectedLen {
		t.Errorf("ERROR: The length of %s value %d is not %d!\n",
			set.Len(), typeName, expectedLen)
		t.FailNow()
	}
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
	for k := range expectedElemMap {
		if !set.Contains(k) {
			t.Errorf("ERROR: The %s value %v do not contains %v!",
				set, typeName, k)
			t.FailNow()
		}
	}
}

func testSetRemove(t *testing.T, newSet func() Set, typeName string) {
	t.Logf("Starting Test%sRemove...", typeName)
	set, expectedElemMap := genRandSet(newSet)
	t.Logf("Got a %s value: %v.", typeName, set)
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
	var number int
	for k, _ := range expectedElemMap {
		if number%2 == 0 {
			t.Logf("Remove %v from the HashSet value %v.\n", k, set)
			set.Remove(k)
			if set.Contains(k) {
				t.Errorf("ERROR: The element removing (%v => %v) is failing!\n",
					k, set)
				t.FailNow()
			}
			delete(expectedElemMap, k)
		}
		number++
	}
	expectedLen := len(expectedElemMap)
	if set.Len() != expectedLen {
		t.Errorf("ERROR: The length of HashSet value %d is not %d!\n", set.Len(), expectedLen)
		t.FailNow()
	}
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
	for _, v := range set.Elements() {
		if !expectedElemMap[v] {
			t.Errorf("ERROR: The HashSet value %v contains %v but should not contains!", set, v)
			t.FailNow()
		}
	}
}

func testSetClear(t *testing.T, newSet func() Set, typeName string) {
	t.Logf("Starting Test%sClear...", typeName)
	set, _ := genRandSet(newSet)
	t.Logf("Got a %s value: %v.", typeName, set)
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
	t.Logf("Clear the HashSet value %v.\n", set)
	set.Clear()
	expectedLen := 0
	if set.Len() != expectedLen {
		t.Errorf("ERROR: The length of HashSet value %d is not %d!\n", set.Len(), expectedLen)
		t.FailNow()
	}
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
}

func testSetElements(t *testing.T, newSet func() Set, typeName string) {
	t.Logf("Starting Test%sElements...", typeName)
	set, expectedElemMap := genRandSet(newSet)
	t.Logf("Got a %s value: %v.", typeName, set)
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
	elems := set.Elements()
	t.Logf("The elements of %s value is %v.\n", typeName, elems)
	expectedLen := len(expectedElemMap)
	if len(elems) != expectedLen {
		t.Errorf("ERROR: The length of HashSet value %d is not %d!\n", len(elems), expectedLen)
		t.FailNow()
	}
	t.Logf("The length of elements is %d.\n", len(elems))
	for _, v := range elems {
		if !expectedElemMap[v] {
			t.Errorf("ERROR: The elements %v contains %v but should not contains!", set, v)
			t.FailNow()
		}
	}
}

func testSetSame(t *testing.T, newSet func() Set, typeName string) {
	t.Logf("Starting Test%sSame...", typeName)
	set, _ := genRandSet(newSet)
	t.Logf("Got a %s value: %v.", typeName, set)
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
	set2 := newSet()
	t.Logf("Clone the HashSet value %v...\n", set)
	for _, v := range set.Elements() {
		set2.Add(v)
	}
	result := set2.Same(set)
	if !result {
		t.Errorf("ERROR: Two sets are not same!")
	}
	t.Logf("Two sets are same.")
}

func testSetString(t *testing.T, newSet func() Set, typeName string) {
	t.Logf("Starting Test%sString...", typeName)
	set, _ := genRandSet(newSet)
	t.Logf("Got a %s value: %v.", typeName, set)
	setStr := set.String()
	t.Logf("The string of %s value is %s.\n", typeName, setStr)
	var elemStr string
	for _, v := range set.Elements() {
		elemStr = fmt.Sprintf("%v", v)
		if !strings.Contains(setStr, elemStr) {
			t.Errorf("ERROR: The string of %s value %s do not contains %s!",
				typeName, setStr, elemStr)
			t.FailNow()
		}
	}
}

// ----- Set 公用函数测试 -----

func TestIsSuperset(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	t.Log("Starting TestIsSuperset...")
	set, _ := genRandSet(func() Set { return NewSimpleSet() })
	set2 := NewSimpleSet()
	for _, v := range set.Elements() {
		set2.Add(v)
	}
	for extraElem := genRandElement(); ; {
		if set2.Add(extraElem) {
			break
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
	if !IsSuperset(set2, set) {
		t.Errorf("ERROR: The HashSet value %v is not a superset of %v!\n", set2, set)
		t.FailNow()
	} else {
		t.Logf("The HashSet value %v is a superset of %v.\n", set2, set)
	}
	for extraElem := genRandElement(); ; {
		if set.Add(extraElem) {
			break
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
	if IsSuperset(set2, set) {
		t.Errorf("ERROR: The HashSet value %v should not be a superset of %v!\n", set2, set)
		t.FailNow()
	} else {
		t.Logf("The HashSet value %v is not a superset of %v.\n", set2, set)
	}
}

func TestUnion(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	t.Log("Starting TestUnion...")
	set, _ := genRandSet(func() Set { return NewSimpleSet() })
	t.Logf("The set value: %v", set)
	set2, _ := genRandSet(func() Set { return NewSimpleSet() })
	uSet := Union(set, set2)
	t.Logf("The set value (2): %v", set2)
	for _, v := range set.Elements() {
		if !uSet.Contains(v) {
			t.Errorf("ERROR: The union set value %v do not contains %v!",
				uSet, v)
			t.FailNow()
		}
	}
	for _, v := range set2.Elements() {
		if !uSet.Contains(v) {
			t.Errorf("ERROR: The union set value %v do not contains %v!",
				uSet, v)
			t.FailNow()
		}
	}
	t.Logf("The set value %v is a unioned set of %v and %v", uSet, set, set2)
}

func TestIntersect(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	t.Log("Starting TestIntersect...")
	commonElem := genRandElement()
	set, _ := genRandSet(func() Set { return NewSimpleSet() })
	set.Add(commonElem)
	t.Logf("The set value: %v", set)
	set2, _ := genRandSet(func() Set { return NewSimpleSet() })
	set2.Add(commonElem)
	t.Logf("The set value (2): %v", set2)
	iSet := Intersect(set, set2)
	for _, v := range iSet.Elements() {
		if !set.Contains(v) {
			t.Errorf("ERROR: The set value %v do not contains %v!",
				set, v)
			t.FailNow()
		}
		if !set2.Contains(v) {
			t.Errorf("ERROR: The set value %v do not contains %v!",
				set2, v)
			t.FailNow()
		}
	}
	t.Logf("The set value %v is a intersected set of %v and %v", iSet, set, set2)
}

func TestDifference(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	t.Log("Starting TestDifference...")
	commonElem := genRandElement()
	set, _ := genRandSet(func() Set { return NewSimpleSet() })
	set.Add(commonElem)
	t.Logf("The set value: %v", set)
	set2, _ := genRandSet(func() Set { return NewSimpleSet() })
	set2.Add(commonElem)
	t.Logf("The set value (2): %v", set2)
	dSet := Difference(set, set2)
	for _, v := range dSet.Elements() {
		if !set.Contains(v) {
			t.Errorf("ERROR: The set value %v do not contains %v!",
				set, v)
			t.FailNow()
		}
		if set2.Contains(v) {
			t.Errorf("ERROR: The set value %v contains %v!",
				set2, v)
			t.FailNow()
		}
	}
	t.Logf("The set value %v is a differenced set of %v to %v", dSet, set, set2)
}

func TestSymmetricDifference(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	t.Log("Starting TestSymmetricDifference...")
	commonElem := genRandElement()
	set, _ := genRandSet(func() Set { return NewSimpleSet() })
	set.Add(commonElem)
	t.Logf("The set value: %v", set)
	set2, _ := genRandSet(func() Set { return NewSimpleSet() })
	set2.Add(commonElem)
	t.Logf("The set value (2): %v", set2)
	sdSet := SymmetricDifference(set, set2)
	for _, v := range sdSet.Elements() {
		if set.Contains(v) && set2.Contains(v) {
			t.Errorf("ERROR: The element %v can not be a common element of %v to %v!",
				v, set, set2)
			t.FailNow()
		}
	}
	t.Logf("The set value %v is a symmetric differenced set of %v to %v", sdSet, set, set2)
}

// ----- 随机测试对象生成函数 -----

func genRandSet(newSet func() Set) (set Set, elemMap map[interface{}]bool) {
	set = newSet()
	elemMap = make(map[interface{}]bool)
	var enough bool
	for !enough {
		e := genRandElement()
		set.Add(e)
		elemMap[e] = true
		if len(elemMap) >= 3 {
			enough = true
		}
	}
	return
}

func genRandElement() interface{} {
	seed := rand.Int63n(10000)
	switch seed {
	case 0:
		return genRandInt()
	case 1:
		return genRandString()
	case 2:
		return struct {
			num int64
			str string
		}{genRandInt(), genRandString()}
	default:
		const length = 2
		arr := new([length]interface{})
		for i := 0; i < length; i++ {
			if i%2 == 0 {
				arr[i] = genRandInt()
			} else {
				arr[i] = genRandString()
			}
		}
		return *arr
	}
}

func genRandString() string {
	var buff bytes.Buffer
	var prev string
	var curr string
	for i := 0; buff.Len() < 3; i++ {
		curr = string(genRandAZAscii())
		if curr == prev {
			continue
		} else {
			prev = curr
		}
		buff.WriteString(curr)
	}
	return buff.String()
}

func genRandAZAscii() int {
	min := 65 // A
	max := 90 // Z
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func genRandInt() int64 {
	return rand.Int63n(10000)
}
