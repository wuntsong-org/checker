package checker

import (
	errors "github.com/wuntsong-org/wterrors"
	"reflect"
)

type Checker struct {
	intChecker    map[string]func(*reflect.StructField, int64) errors.WTError
	stringChecker map[string]func(*reflect.StructField, string) errors.WTError
	sliceChecker  map[string]func(*reflect.StructField, any) errors.WTError
	mapChecker    map[string]func(*reflect.StructField, any) errors.WTError
	strictSlice   bool
	strictMap     bool
}

func NewChecker() *Checker {
	return &Checker{
		intChecker:    make(map[string]func(*reflect.StructField, int64) errors.WTError),
		stringChecker: make(map[string]func(*reflect.StructField, string) errors.WTError),
		sliceChecker:  make(map[string]func(*reflect.StructField, interface{}) errors.WTError),
		mapChecker:    make(map[string]func(*reflect.StructField, interface{}) errors.WTError),
		strictSlice:   false,
		strictMap:     false,
	}
}

func (c *Checker) AddIntChecker(name string, checker func(*reflect.StructField, int64) errors.WTError) {
	c.intChecker[name] = checker
}

func (c *Checker) RemoveIntChecker(name string) {
	delete(c.intChecker, name)
}

func (c *Checker) GetIntChecker(name string) (func(*reflect.StructField, int64) errors.WTError, bool) {
	checker, exists := c.intChecker[name]
	return checker, exists
}

func (c *Checker) AddStringChecker(name string, checker func(*reflect.StructField, string) errors.WTError) {
	c.stringChecker[name] = checker
}

func (c *Checker) RemoveStringChecker(name string) {
	delete(c.stringChecker, name)
}

func (c *Checker) GetStringChecker(name string) (func(*reflect.StructField, string) errors.WTError, bool) {
	checker, exists := c.stringChecker[name]
	return checker, exists
}

func (c *Checker) AddSliceChecker(name string, checker func(*reflect.StructField, interface{}) errors.WTError) {
	c.sliceChecker[name] = checker
}

func (c *Checker) RemoveSliceChecker(name string) {
	delete(c.sliceChecker, name)
}

func (c *Checker) GetSliceChecker(name string) (func(*reflect.StructField, interface{}) errors.WTError, bool) {
	checker, exists := c.sliceChecker[name]
	return checker, exists
}

func (c *Checker) AddMapChecker(name string, checker func(*reflect.StructField, interface{}) errors.WTError) {
	c.mapChecker[name] = checker
}

func (c *Checker) RemoveMapChecker(name string) {
	delete(c.mapChecker, name)
}

func (c *Checker) GetMapChecker(name string) (func(*reflect.StructField, interface{}) errors.WTError, bool) {
	checker, exists := c.mapChecker[name]
	return checker, exists
}

func (c *Checker) SetStrictSlice(value bool) {
	c.strictSlice = value
}

func (c *Checker) GetStrictSlice() bool {
	return c.strictSlice
}

func (c *Checker) SetStrictMap(value bool) {
	c.strictMap = value
}

func (c *Checker) GetStrictMap() bool {
	return c.strictMap
}

func (c *Checker) Check(s any) errors.WTError {
	return c.checkInterfaceOrPointer(&s, nil)
}
