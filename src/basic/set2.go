package basic

import (
	"fmt"
	"sort"
)

type KeyGeneratorFunc func(x interface{}) string

type ComparisonFunc func(i, j interface{}) int

type SimpleSetIterator func() (interface{}, bool)

type SimpleSet struct {
	KeyGenerator KeyGeneratorFunc
	Comparator   ComparisonFunc
	keys         []string
	elementMap   map[string]interface{}
}

func (self *SimpleSet) Len() int {
	return len(self.keys)
}

func (self *SimpleSet) Less(i, j int) bool {
	ki := self.keys[i]
	kj := self.keys[j]
	var result bool
	if !self.Sortable() {
		result = ki < kj
	} else {
		ii := self.elementMap[ki]
		ij := self.elementMap[kj]
		sign := self.Comparator(ii, ij)
		if sign < 0 {
			result = true
		} else {
			result = false
		}
	}
	return result
}

func (self *SimpleSet) Swap(i, j int) {
	self.keys[i], self.keys[j] = self.keys[j], self.keys[i]
}

func (self *SimpleSet) initialize(rebuild bool) {
	if rebuild {
		if self.elementMap != nil {
			self.elementMap = nil
			self.elementMap = make(map[string]interface{})
		}
		if self.keys != nil {
			self.keys = nil
			self.keys = make([]string, 0)
		}
	} else {
		if self.elementMap == nil {
			self.elementMap = make(map[string]interface{})
		}
		if self.keys == nil {
			self.keys = make([]string, 0)
		}
	}
}

func (self *SimpleSet) generateKey(value interface{}) string {
	var k string
	if self.KeyGenerator != nil {
		k = self.KeyGenerator(value) + "}"
	} else {
		k = fmt.Sprintf("%v", value)
	}
	return "K{" + k + "}"
}

func (self *SimpleSet) Add(element interface{}) bool {
	self.initialize(false)
	if element == nil {
		return false
	}
	done := false
	key := self.generateKey(element)
	if self.elementMap[key] == nil {
		self.elementMap[key] = element
		self.keys = append(self.keys, key)
		done = true
	}
	return done
}

func (self *SimpleSet) Remove(element interface{}) bool {
	if len(self.elementMap) == 0 {
		return false
	}
	key := self.generateKey(element)
	done := false
	if self.elementMap[key] != nil {
		delete(self.elementMap, key)
		var keyX string
		sort.Strings(self.keys)
		keySize := len(self.keys)
		index := sort.Search(keySize, func(x int) bool {
			keyX = self.keys[x]
			return keyX >= key
		})
		if index >= 0 || index < keySize && keyX == key {
			copy(self.keys[index:], self.keys[index+1:])
			endIndex := keySize - 1
			self.keys[endIndex] = ""
			self.keys = self.keys[:endIndex]
			done = true
		}
	}
	return done
}

func (self *SimpleSet) Clear() bool {
	self.initialize(true)
	return true
}

func (self *SimpleSet) Contain(element interface{}) bool {
	if self.elementMap == nil || len(self.elementMap) == 0 {
		return false
	}
	key := self.generateKey(element)
	if self.elementMap[key] == nil {
		return false
	}
	return true
}

func (self *SimpleSet) Iterator() SimpleSetIterator {
	self.initialize(false)
	return func() SimpleSetIterator {
		index := 0
		snapshots := self.Slice()
		return func() (interface{}, bool) {
			if index >= 0 && index < len(snapshots) {
				element := snapshots[index]
				index++
				return element, true
			}
			return nil, false
		}
	}()
}

func (self *SimpleSet) Slice() []interface{} {
	if len(self.keys) == 0 {
		return make([]interface{}, 0)
	}
	snapshots := make([]interface{}, len(self.keys))
	if self.Sortable() {
		sort.Sort(self)
		for i, k := range self.keys {
			snapshots[i] = self.elementMap[k]
		}
	} else {
		count := 0
		for _, v := range self.elementMap {
			snapshots[count] = v
			count++
		}
	}
	return snapshots
}

func (self *SimpleSet) String() string {
	return fmt.Sprintf("%v", self.Slice())
}

func (self *SimpleSet) Sortable() bool {
	return self.Comparator != nil
}

func (self *SimpleSet) GetComparator() ComparisonFunc {
	return self.Comparator
}
