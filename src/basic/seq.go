package basic

import (
	"reflect"
	"sort"
)

type Sortable interface {
	sort.Interface
	Sort()
}

type GenericSeq interface {
	Sortable
	Append(e interface{}) bool
	Set(index int, e interface{}) bool
	Delete(index int) (interface{}, bool)
	ElemValue(index int) interface{}
	ElemType() reflect.Type
	Value() interface{}
}

type StringSeq struct {
	value []string
}

func (self *StringSeq) Len() int {
	return len(self.value)
}

func (self *StringSeq) Less(i, j int) bool {
	return self.value[i] < self.value[j]
}

func (self *StringSeq) Swap(i, j int) {
	self.value[i], self.value[j] = self.value[j], self.value[i]
}

func (self *StringSeq) Sort() {
	sort.Sort(self)
}

func (self *StringSeq) Append(e interface{}) bool {
	s, ok := e.(string)
	if !ok {
		return false
	}
	self.value = append(self.value, s)
	return true
}

func (self *StringSeq) Set(index int, e interface{}) bool {
	if index >= self.Len() {
		return false
	}
	s, ok := e.(string)
	if !ok {
		return false
	}
	self.value[index] = s
	return true
}

func (self *StringSeq) Delete(index int) (interface{}, bool) {
	length := self.Len()
	if index >= length {
		return nil, false
	}
	s := self.value[index]
	if index < (length - 1) {
		copy(self.value[index:], self.value[index+1:])
	}
	invalidIndex := length - 1
	self.value[invalidIndex] = ""
	self.value = self.value[:invalidIndex]
	return s, true
}

func (self StringSeq) ElemValue(index int) interface{} {
	if index >= self.Len() {
		return nil
	}
	return self.value[index]
}

func (self *StringSeq) ElemType() reflect.Type {
	return reflect.TypeOf(self.value).Elem()
}

func (self StringSeq) Value() interface{} {
	return self.value
}

type Sequence struct {
	GenericSeq
	sorted   bool
	elemType reflect.Type
}

func (self *Sequence) Sort() {
	self.GenericSeq.Sort()
	self.sorted = true
}

func (self *Sequence) Append(e interface{}) bool {
	result := self.GenericSeq.Append(e)
	if result && self.sorted {
		self.sorted = false
	}
	return result
}

func (self *Sequence) Set(index int, e interface{}) bool {
	result := self.GenericSeq.Set(index, e)
	if result && self.sorted {
		self.sorted = false
	}
	return result
}

func (self *Sequence) ElemType() reflect.Type {
	if self.elemType == nil && self.GenericSeq != nil {
		self.elemType = self.GenericSeq.ElemType()
	}
	return self.elemType
}

func (self *Sequence) Init() (ok bool) {
	if self.GenericSeq != nil {
		self.elemType = self.GenericSeq.ElemType()
		ok = true
	}
	return ok
}

func (self *Sequence) Sorted() bool {
	return self.sorted
}
