package checker

import "testing"

func TestChecker(t *testing.T) {
	type TestStruct struct {
		Int int64 `ci-max:"10" ci-min:"0"`
	}

	a := TestStruct{Int: 5}

	checker := NewChecker()
	err := checker.Check(&a)
	if err != nil {
		t.Error(err.Error())
	}
}
