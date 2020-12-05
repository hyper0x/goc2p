package map1

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
)

type CompareFunction func(interface{}, interface{}) int8

type Keys interface {
	sort.Interface
	Add(k interface{}) bool
	Remove(k interface{}) bool
	Clear()
	Get(index int) interface{}
	GetAll() []interface{}
	Search(k interface{}) (index int, contains bool)
	CompareFunc() CompareFunction
	ElemType() reflect.Type
}

// compareFunc的结果值：
//   小于0: 第一个参数小于第二个参数
//   等于0: 第一个参数等于第二个参数
//   大于1: 第一个参数大于第二个参数
type myKeys struct {
	container   []interface{}
	compareFunc CompareFunction
	elemType    reflect.Type
}

func (keys *myKeys) Len() int {
	return len(keys.container)
}

func (keys *myKeys) Less(i, j int) bool {
	return keys.compareFunc(keys.container[i], keys.container[j]) == -1
}

func (keys *myKeys) Swap(i, j int) {
	keys.container[i], keys.container[j] = keys.container[j], keys.container[i]
}

func (keys *myKeys) isAcceptableElem(k interface{}) bool {
	if k == nil {
		return false
	}
	if reflect.TypeOf(k) != keys.elemType {
		return false
	}
	return true
}

func (keys *myKeys) Add(k interface{}) bool {
	ok := keys.isAcceptableElem(k)
	if !ok {
		return false
	}
	keys.container = append(keys.container, k)
	sort.Sort(keys)
	return true
}

func (keys *myKeys) Remove(k interface{}) bool {
	index, contains := keys.Search(k)
	if !contains {
		return false
	}
	keys.container = append(keys.container[0:index], keys.container[index+1:]...)
	return true
}

func (keys *myKeys) Clear() {
	keys.container = make([]interface{}, 0)
}

func (keys *myKeys) Get(index int) interface{} {
	if index >= keys.Len() {
		return nil
	}
	return keys.container[index]
}

func (keys *myKeys) GetAll() []interface{} {
	initialLen := len(keys.container)
	snapshot := make([]interface{}, initialLen)
	actualLen := 0
	for _, key := range keys.container {
		if actualLen < initialLen {
			snapshot[actualLen] = key
		} else {
			snapshot = append(snapshot, key)
		}
		actualLen++
	}
	if actualLen < initialLen {
		snapshot = snapshot[:actualLen]
	}
	return snapshot
}

func (keys *myKeys) Search(k interface{}) (index int, contains bool) {
	ok := keys.isAcceptableElem(k)
	if !ok {
		return
	}
	index = sort.Search(
		keys.Len(),
		func(i int) bool { return keys.compareFunc(keys.container[i], k) >= 0 })
	if index < keys.Len() && keys.container[index] == k {
		contains = true
	}
	return
}

func (keys *myKeys) ElemType() reflect.Type {
	return keys.elemType
}

func (keys *myKeys) CompareFunc() CompareFunction {
	return keys.compareFunc
}

func (keys *myKeys) String() string {
	var buf bytes.Buffer
	buf.WriteString("Keys<")
	buf.WriteString(keys.elemType.Kind().String())
	buf.WriteString(">{")
	first := true
	buf.WriteString("[")
	for _, key := range keys.container {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", key))
	}
	buf.WriteString("]")
	buf.WriteString("}")
	return buf.String()
}

func NewKeys(
	compareFunc func(interface{}, interface{}) int8,
	elemType reflect.Type) Keys {
	return &myKeys{
		container:   make([]interface{}, 0),
		compareFunc: compareFunc,
		elemType:    elemType,
	}
}
