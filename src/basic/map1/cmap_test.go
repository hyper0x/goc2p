package map1

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime/debug"
	"testing"
)

func testConcurrentMap(
	t *testing.T,
	newConcurrentMap func() ConcurrentMap,
	genKey func() interface{},
	genElem func() interface{},
	keyKind reflect.Kind,
	elemKind reflect.Kind) {
	mapType := fmt.Sprintf("ConcurrentMap<keyType=%s, elemType=%s>", keyKind, elemKind)
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s: %s\n", mapType, err)
		}
	}()
	t.Logf("Starting Test%s...", mapType)

	// Basic
	cmap := newConcurrentMap()
	expectedLen := 0
	if cmap.Len() != expectedLen {
		t.Errorf("ERROR: The length of %s value %d is not %d!\n",
			mapType, cmap.Len(), expectedLen)
		t.FailNow()
	}
	expectedLen = 5
	testMap := make(map[interface{}]interface{}, expectedLen)
	var invalidKey interface{}
	for i := 0; i < expectedLen; i++ {
		key := genKey()
		testMap[key] = genElem()
		if invalidKey == nil {
			invalidKey = key
		}
	}
	for key, elem := range testMap {
		oldElem, ok := cmap.Put(key, elem)
		if !ok {
			t.Errorf("ERROR: Put (%v, %v) to %s value %d is failing!\n",
				key, elem, mapType, cmap)
			t.FailNow()
		}
		if oldElem != nil {
			t.Errorf("ERROR: Already had a (%v, %v) in %s value %d!\n",
				key, elem, mapType, cmap)
			t.FailNow()
		}
		t.Logf("Put (%v, %v) to the %s value %v.",
			key, elem, mapType, cmap)
	}
	if cmap.Len() != expectedLen {
		t.Errorf("ERROR: The length of %s value %d is not %d!\n",
			mapType, cmap.Len(), expectedLen)
		t.FailNow()
	}
	for key, elem := range testMap {
		contains := cmap.Contains(key)
		if !contains {
			t.Errorf("ERROR: The %s value %v do not contains %v!",
				mapType, cmap, key)
			t.FailNow()
		}
		actualElem := cmap.Get(key)
		if actualElem == nil {
			t.Errorf("ERROR: The %s value %v do not contains %v!",
				mapType, cmap, key)
			t.FailNow()
		}
		t.Logf("The %s value %v contains key %v.", mapType, cmap, key)
		if actualElem != elem {
			t.Errorf("ERROR: The element of %s value %v with key %v do not equals %v!\n",
				mapType, actualElem, key, elem)
			t.FailNow()
		}
		t.Logf("The element of %s value %v to key %v is %v.",
			mapType, cmap, key, actualElem)
	}
	oldElem := cmap.Remove(invalidKey)
	if oldElem == nil {
		t.Errorf("ERROR: Remove %v from %s value %d is failing!\n",
			invalidKey, mapType, cmap)
		t.FailNow()
	}
	t.Logf("Removed (%v, %v) from the %s value %v.",
		invalidKey, oldElem, mapType, cmap)
	delete(testMap, invalidKey)

	// Type
	actualElemType := cmap.ElemType()
	if actualElemType == nil {
		t.Errorf("ERROR: The element type of %s value is nil!\n",
			mapType)
		t.FailNow()
	}
	actualElemKind := actualElemType.Kind()
	if actualElemKind != elemKind {
		t.Errorf("ERROR: The element type of %s value %s is not %s!\n",
			mapType, actualElemKind, elemKind)
		t.FailNow()
	}
	t.Logf("The element type of %s value %v is %s.",
		mapType, cmap, actualElemKind)
	actualKeyKind := cmap.KeyType().Kind()
	if actualKeyKind != elemKind {
		t.Errorf("ERROR: The key type of %s value %s is not %s!\n",
			mapType, actualKeyKind, keyKind)
		t.FailNow()
	}
	t.Logf("The key type of %s value %v is %s.",
		mapType, cmap, actualKeyKind)

	// Export
	keys := cmap.Keys()
	elems := cmap.Elems()
	pairs := cmap.ToMap()
	for key, elem := range testMap {
		var hasKey bool
		for _, k := range keys {
			if k == key {
				hasKey = true
			}
		}
		if !hasKey {
			t.Errorf("ERROR: The keys of %s value %v do not contains %v!\n",
				mapType, cmap, key)
			t.FailNow()
		}
		var hasElem bool
		for _, e := range elems {
			if e == elem {
				hasElem = true
			}
		}
		if !hasElem {
			t.Errorf("ERROR: The elems of %s value %v do not contains %v!\n",
				mapType, cmap, elem)
			t.FailNow()
		}
		var hasPair bool
		for k, e := range pairs {
			if k == key && e == elem {
				hasPair = true
			}
		}
		if !hasPair {
			t.Errorf("ERROR: The elems of %s value %v do not contains (%v, %v)!\n",
				mapType, cmap, key, elem)
			t.FailNow()
		}
	}

	// Clear
	cmap.Clear()
	if cmap.Len() != 0 {
		t.Errorf("ERROR: Clear %s value %d is failing!\n",
			mapType, cmap)
		t.FailNow()
	}
	t.Logf("The %s value %v has been cleared.", mapType, cmap)
}

func TestInt64Cmap(t *testing.T) {
	newCmap := func() ConcurrentMap {
		keyType := reflect.TypeOf(int64(2))
		elemType := keyType
		return NewConcurrentMap(keyType, elemType)
	}
	testConcurrentMap(
		t,
		newCmap,
		func() interface{} { return rand.Int63n(1000) },
		func() interface{} { return rand.Int63n(1000) },
		reflect.Int64,
		reflect.Int64)
}

func TestFloat64Cmap(t *testing.T) {
	newCmap := func() ConcurrentMap {
		keyType := reflect.TypeOf(float64(2))
		elemType := keyType
		return NewConcurrentMap(keyType, elemType)
	}
	testConcurrentMap(
		t,
		newCmap,
		func() interface{} { return rand.Float64() },
		func() interface{} { return rand.Float64() },
		reflect.Float64,
		reflect.Float64)
}

func TestStringCmap(t *testing.T) {
	newCmap := func() ConcurrentMap {
		keyType := reflect.TypeOf(string(2))
		elemType := keyType
		return NewConcurrentMap(keyType, elemType)
	}
	testConcurrentMap(
		t,
		newCmap,
		func() interface{} { return genRandString() },
		func() interface{} { return genRandString() },
		reflect.String,
		reflect.String)
}

func BenchmarkConcurrentMap(b *testing.B) {
	keyType := reflect.TypeOf(int32(2))
	elemType := keyType
	cmap := NewConcurrentMap(keyType, elemType)
	var key, elem int32
	fmt.Printf("N=%d.\n", b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		seed := int32(i)
		key = seed
		elem = seed << 10
		b.StartTimer()
		cmap.Put(key, elem)
		_ = cmap.Get(key)
		b.StopTimer()
		b.SetBytes(8)
		b.StartTimer()
	}
	ml := cmap.Len()
	b.StopTimer()
	mapType := fmt.Sprintf("ConcurrentMap<%s, %s>",
		keyType.Kind().String(), elemType.Kind().String())
	b.Logf("The length of %s value is %d.\n", mapType, ml)
	b.StartTimer()
}

func BenchmarkMap(b *testing.B) {
	keyType := reflect.TypeOf(int32(2))
	elemType := keyType
	imap := make(map[interface{}]interface{})
	var key, elem int32
	fmt.Printf("N=%d.\n", b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		seed := int32(i)
		key = seed
		elem = seed << 10
		b.StartTimer()
		imap[key] = elem
		b.StopTimer()
		_ = imap[key]
		b.StopTimer()
		b.SetBytes(8)
		b.StartTimer()
	}
	ml := len(imap)
	b.StopTimer()
	mapType := fmt.Sprintf("Map<%s, %s>",
		keyType.Kind().String(), elemType.Kind().String())
	b.Logf("The length of %s value is %d.\n", mapType, ml)
	b.StartTimer()
}
