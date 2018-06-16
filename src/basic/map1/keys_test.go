package map1

import (
	"bytes"
	"math/rand"
	"reflect"
	"runtime/debug"
	"sort"
	"testing"
	"time"
)

func testKeys(
	t *testing.T,
	newKeys func() Keys,
	genKey func() interface{},
	elemKind reflect.Kind) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: Keys(%s): %s\n", elemKind, err)
		}
	}()
	t.Logf("Starting TestKeys<%s>...", elemKind)
	keys := newKeys()
	expectedLen := 0
	if keys.Len() != expectedLen {
		t.Errorf("ERROR: The length of Keys(%s) value %d is not %d!\n",
			elemKind, keys.Len(), expectedLen)
		t.FailNow()
	}
	expectedLen = 5
	testKeys := make([]interface{}, expectedLen)
	for i := 0; i < expectedLen; i++ {
		testKeys[i] = genKey()
	}
	for _, key := range testKeys {
		result := keys.Add(key)
		if !result {
			t.Errorf("ERROR: Add %v to Keys(%s) value %d is failing!\n",
				key, elemKind, keys)
			t.FailNow()
		}
		t.Logf("Added %v to the Keys(%s) value %v.", key, elemKind, keys)
	}
	if keys.Len() != expectedLen {
		t.Errorf("ERROR: The length of Keys(%s) value %d is not %d!\n",
			elemKind, keys.Len(), expectedLen)
		t.FailNow()
	}
	for _, key := range testKeys {
		index, contains := keys.Search(key)
		if !contains {
			t.Errorf("ERROR: The Keys(%s) value %v do not contains %v!",
				elemKind, keys, key)
			t.FailNow()
		}
		t.Logf("The Keys(%s) value %v contains key %v.", elemKind, keys, key)
		actualElem := keys.Get(index)
		if actualElem != key {
			t.Errorf("ERROR: The element of Keys(%s) value %v with index %d do not equals %v!\n",
				elemKind, actualElem, index, key)
			t.FailNow()
		}
		t.Logf("The element of Keys(%s) value %v with index %d is %v.",
			elemKind, keys, index, actualElem)
	}
	invalidElem := testKeys[len(testKeys)/2]
	result := keys.Remove(invalidElem)
	if !result {
		t.Errorf("ERROR: Remove %v from Keys(%s) value %d is failing!\n",
			invalidElem, elemKind, keys)
		t.FailNow()
	}
	t.Logf("Removed %v from the Keys(%s) value %v.", invalidElem, elemKind, keys)
	if !sort.IsSorted(keys) {
		t.Errorf("ERROR: The Keys(%s) value %v is not sorted yet?!\n",
			elemKind, keys)
		t.FailNow()
	}
	t.Logf("The Keys(%s) value %v is sorted.", elemKind, keys)
	actualElemType := keys.ElemType()
	if actualElemType == nil {
		t.Errorf("ERROR: The element type of Keys(%s) value is nil!\n",
			elemKind)
		t.FailNow()
	}
	actualElemKind := actualElemType.Kind()
	if actualElemKind != elemKind {
		t.Errorf("ERROR: The element type of Keys(%s) value %s is not %s!\n",
			elemKind, actualElemKind, elemKind)
		t.FailNow()
	}
	t.Logf("The element type of Keys(%s) value %v is %s.", elemKind, keys, actualElemKind)
	currCompFunc := keys.CompareFunc()
	if currCompFunc == nil {
		t.Errorf("ERROR: The compare function of Keys(%s) value is nil!\n",
			elemKind)
		t.FailNow()
	}
	keys.Clear()
	if keys.Len() != 0 {
		t.Errorf("ERROR: Clear Keys(%s) value %d is failing!\n",
			elemKind, keys)
		t.FailNow()
	}
	t.Logf("The Keys(%s) value %v have been cleared.", elemKind, keys)
}

func TestInt64Keys(t *testing.T) {
	testKeys(t,
		func() Keys {
			//return NewKeys(
			//func(e1 interface{}, e2 interface{}) int8 {
			//	k1 := e1.(int64)
			//	k2 := e2.(int64)
			//	if k1 < k2 {
			//		return -1
			//	} else if k1 > k2 {
			//		return 1
			//	} else {
			//		return 0
			//	}
			//},
			//reflect.TypeOf(int64(1)))
			int64Keys := &myKeys{
				container: make([]interface{}, 0),
				compareFunc: func(e1 interface{}, e2 interface{}) int8 {
					k1 := e1.(int64)
					k2 := e2.(int64)
					if k1 < k2 {
						return -1
					} else if k1 > k2 {
						return 1
					} else {
						return 0
					}
				},
				elemType: reflect.TypeOf(int64(1))}
			return int64Keys
		},
		func() interface{} { return rand.Int63n(1000) },
		reflect.Int64)
}

func TestFloat64Keys(t *testing.T) {
	testKeys(t,
		func() Keys {
			return NewKeys(
				func(e1 interface{}, e2 interface{}) int8 {
					k1 := e1.(float64)
					k2 := e2.(float64)
					if k1 < k2 {
						return -1
					} else if k1 > k2 {
						return 1
					} else {
						return 0
					}
				},
				reflect.TypeOf(float64(1)))
		},
		func() interface{} { return rand.Float64() },
		reflect.Float64)
}

func TestStringKeys(t *testing.T) {
	testKeys(t,
		func() Keys {
			return NewKeys(
				func(e1 interface{}, e2 interface{}) int8 {
					k1 := e1.(string)
					k2 := e2.(string)
					if k1 < k2 {
						return -1
					} else if k1 > k2 {
						return 1
					} else {
						return 0
					}
				},
				reflect.TypeOf(string(1)))
		},
		func() interface{} { return genRandString() },
		reflect.String)
}

func genRandString() string {
	var buff bytes.Buffer
	var prev string
	var curr string
	for i := 0; buff.Len() < 3; i++ {
		curr = string(genRandAZAscii())
		if curr == prev {
			continue
		}
		prev = curr
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
