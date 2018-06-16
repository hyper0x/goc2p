package map1

import (
	"reflect"
)

// 泛化的Map的接口类型
type GenericMap interface {
	// 获取给定键值对应的元素值。若没有对应元素值则返回nil。
	Get(key interface{}) interface{}
	// 添加键值对，并返回与给定键值对应的旧的元素值。若没有旧元素值则返回(nil, true)。
	Put(key interface{}, elem interface{}) (interface{}, bool)
	// 删除与给定键值对应的键值对，并返回旧的元素值。若没有旧元素值则返回nil。
	Remove(key interface{}) interface{}
	// 清除所有的键值对。
	Clear()
	// 获取键值对的数量。
	Len() int
	// 判断是否包含给定的键值。
	Contains(key interface{}) bool
	// 获取已排序的键值所组成的切片值。
	Keys() []interface{}
	// 获取已排序的元素值所组成的切片值。
	Elems() []interface{}
	// 获取已包含的键值对所组成的字典值。
	ToMap() map[interface{}]interface{}
	// 获取键的类型。
	KeyType() reflect.Type
	// 获取元素的类型。
	ElemType() reflect.Type
}
