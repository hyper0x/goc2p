package map1

import (
	"math/rand"
	"reflect"
	"runtime/debug"
	"testing"
)

func testOrderedMap(
	t *testing.T,
	newOrderedMap func() OrderedMap,
	genKey func() interface{},
	genElem func() interface{},
	elemKind reflect.Kind) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: OrderedMap(type=%s): %s\n", elemKind, err)
		}
	}()
	t.Logf("Starting TestOrderedMap<elemType=%s>...", elemKind)

	// Basic
	omap := newOrderedMap()
	expectedLen := 0
	if omap.Len() != expectedLen {
		t.Errorf("ERROR: The length of OrderedMap(elemType=%s) value %d is not %d!\n",
			elemKind, omap.Len(), expectedLen)
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
		oldElem, ok := omap.Put(key, elem)
		if !ok {
			t.Errorf("ERROR: Put (%v, %v) to OrderedMap(elemType=%s) value %d is failing!\n",
				key, elem, elemKind, omap)
			t.FailNow()
		}
		if oldElem != nil {
			t.Errorf("ERROR: Already had a (%v, %v) in OrderedMap(elemType=%s) value %d!\n",
				key, elem, elemKind, omap)
			t.FailNow()
		}
		t.Logf("Put (%v, %v) to the OrderedMap(elemType=%s) value %v.",
			key, elem, elemKind, omap)
	}
	if omap.Len() != expectedLen {
		t.Errorf("ERROR: The length of OrderedMap(elemType=%s) value %d is not %d!\n",
			elemKind, omap.Len(), expectedLen)
		t.FailNow()
	}
	for key, elem := range testMap {
		contains := omap.Contains(key)
		if !contains {
			t.Errorf("ERROR: The OrderedMap(elemType=%s) value %v do not contains %v!",
				elemKind, omap, key)
			t.FailNow()
		}
		actualElem := omap.Get(key)
		if actualElem == nil {
			t.Errorf("ERROR: The OrderedMap(elemType=%s) value %v do not contains %v!",
				elemKind, omap, key)
			t.FailNow()
		}
		t.Logf("The OrderedMap(elemType=%s) value %v contains key %v.", elemKind, omap, key)
		if actualElem != elem {
			t.Errorf("ERROR: The element of OrderedMap(elemType=%s) value %v with key %v do not equals %v!\n",
				elemKind, actualElem, key, elem)
			t.FailNow()
		}
		t.Logf("The element of OrderedMap(elemType=%s) value %v to key %v is %v.",
			elemKind, omap, key, actualElem)
	}
	oldElem := omap.Remove(invalidKey)
	if oldElem == nil {
		t.Errorf("ERROR: Remove %v from OrderedMap(elemType=%s) value %d is failing!\n",
			invalidKey, elemKind, omap)
		t.FailNow()
	}
	t.Logf("Removed (%v, %v) from the OrderedMap(elemType=%s) value %v.",
		invalidKey, oldElem, elemKind, omap)
	delete(testMap, invalidKey)

	// Type
	actualElemType := omap.ElemType()
	if actualElemType == nil {
		t.Errorf("ERROR: The element type of OrderedMap(elemType=%s) value is nil!\n",
			elemKind)
		t.FailNow()
	}
	actualElemKind := actualElemType.Kind()
	if actualElemKind != elemKind {
		t.Errorf("ERROR: The element type of OrderedMap(elemType=%s) value %s is not %s!\n",
			elemKind, actualElemKind, elemKind)
		t.FailNow()
	}
	t.Logf("The element type of OrderedMap(elemType=%s) value %v is %s.",
		elemKind, omap, actualElemKind)
	actualKeyKind := omap.KeyType().Kind()
	keyKind := reflect.TypeOf(genKey()).Kind()
	if actualKeyKind != elemKind {
		t.Errorf("ERROR: The key type of OrderedMap(elemType=%s) value %s is not %s!\n",
			keyKind, actualKeyKind, keyKind)
		t.FailNow()
	}
	t.Logf("The key type of OrderedMap(elemType=%s) value %v is %s.",
		keyKind, omap, actualKeyKind)

	// Export
	keys := omap.Keys()
	elems := omap.Elems()
	pairs := omap.ToMap()
	for key, elem := range testMap {
		var hasKey bool
		for _, k := range keys {
			if k == key {
				hasKey = true
			}
		}
		if !hasKey {
			t.Errorf("ERROR: The keys of OrderedMap(elemType=%s) value %v do not contains %v!\n",
				elemKind, omap, key)
			t.FailNow()
		}
		var hasElem bool
		for _, e := range elems {
			if e == elem {
				hasElem = true
			}
		}
		if !hasElem {
			t.Errorf("ERROR: The elems of OrderedMap(elemType=%s) value %v do not contains %v!\n",
				elemKind, omap, elem)
			t.FailNow()
		}
		var hasPair bool
		for k, e := range pairs {
			if k == key && e == elem {
				hasPair = true
			}
		}
		if !hasPair {
			t.Errorf("ERROR: The elems of OrderedMap(elemType=%s) value %v do not contains (%v, %v)!\n",
				elemKind, omap, key, elem)
			t.FailNow()
		}
	}

	// Advance
	fKey := omap.FirstKey()
	if fKey != keys[0] {
		t.Errorf("ERROR: The first key of OrderedMap(elemType=%s) value %v is not equals %v!\n",
			elemKind, fKey, keys[0])
		t.FailNow()
	}
	t.Logf("The first key of OrderedMap(elemType=%s) value %v is %s.",
		elemKind, omap, fKey)
	lKey := omap.LastKey()
	if lKey != keys[len(keys)-1] {
		t.Errorf("ERROR: The last key of OrderedMap(elemType=%s) value %v is not equals %v!\n",
			elemKind, lKey, keys[len(keys)-1])
		t.FailNow()
	}
	t.Logf("The last key of OrderedMap(elemType=%s) value %v is %s.",
		elemKind, omap, lKey)
	endIndex := len(keys)/2 + 1
	toKey := keys[endIndex]
	headMap := omap.HeadMap(toKey)
	headKeys := headMap.Keys()
	for i := 0; i < endIndex; i++ {
		hKey := headKeys[i]
		tempKey := keys[i]
		if hKey != tempKey {
			t.Errorf("ERROR: The key of OrderedMap(elemType=%s) value %v with index %d is not equals %v!\n",
				elemKind, tempKey, i, hKey)
			t.FailNow()
		}
	}
	beginIndex := len(keys)/2 - 1
	endIndex = len(keys) - 1
	fromKey := keys[beginIndex]
	tailMap := omap.TailMap(fromKey)
	tailKeys := tailMap.Keys()
	for i := beginIndex; i < endIndex; i++ {
		tKey := tailKeys[i-beginIndex]
		tempKey := keys[i]
		if tKey != tempKey {
			t.Errorf("ERROR: The key of OrderedMap(elemType=%s) value %v with index %d is not equals %v!\n",
				elemKind, tempKey, i, tKey)
			t.FailNow()
		}
	}
	beginIndex = len(keys)/2 - 1
	endIndex = len(keys)/2 + 1
	fromKey = keys[beginIndex]
	toKey = keys[endIndex]
	subMap := omap.SubMap(fromKey, toKey)
	subKeys := subMap.Keys()
	for i := beginIndex; i < endIndex; i++ {
		sKey := subKeys[i-beginIndex]
		tempKey := keys[i]
		if sKey != tempKey {
			t.Errorf("ERROR: The key of OrderedMap(elemType=%s) value %v with index %d is not equals %v!\n",
				elemKind, tempKey, i, sKey)
			t.FailNow()
		}
	}

	// Clear
	omap.Clear()
	if omap.Len() != 0 {
		t.Errorf("ERROR: Clear OrderedMap(elemType=%s) value %d is failing!\n",
			elemKind, omap)
		t.FailNow()
	}
	t.Logf("The OrderedMap(elemType=%s) value %v has been cleared.", elemKind, omap)
}

func TestInt64Omap(t *testing.T) {
	keys := NewKeys(
		func(e1 interface{}, e2 interface{}) int8 {
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
		reflect.TypeOf(int64(1)))
	newOmap := func() OrderedMap {
		return NewOrderedMap(keys, reflect.TypeOf(int64(1)))
	}
	testOrderedMap(
		t,
		newOmap,
		func() interface{} { return rand.Int63n(1000) },
		func() interface{} { return rand.Int63n(1000) },
		reflect.Int64)
}

func TestFloat64Omap(t *testing.T) {
	keys := NewKeys(
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
	newOmap := func() OrderedMap {
		return NewOrderedMap(keys, reflect.TypeOf(float64(1)))
	}
	testOrderedMap(
		t,
		newOmap,
		func() interface{} { return rand.Float64() },
		func() interface{} { return rand.Float64() },
		reflect.Float64)
}

func TestStringOmap(t *testing.T) {
	keys := NewKeys(
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
	newOmap := func() OrderedMap {
		return NewOrderedMap(keys, reflect.TypeOf(string(1)))
	}
	testOrderedMap(
		t,
		newOmap,
		func() interface{} { return genRandString() },
		func() interface{} { return genRandString() },
		reflect.String)
}
