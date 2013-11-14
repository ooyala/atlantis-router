package zk

import (
	"testing"
)

func TestArrayDiff(t *testing.T) {
	arr0 := []string{}
	arr1 := []string{"mary", "had", "a"}
	arr2 := []string{"mary", "had", "a", "little", "lamb"}

	if len(ArrayDiff(arr0, arr1)) != 0 {
		t.Errorf("should be empty")
	}

	if len(ArrayDiff(arr1, arr0)) != 3 {
		t.Errorf("should be {mary, had, a}")
	}

	if len(ArrayDiff(arr2, arr1)) != 2 {
		t.Errorf("should be {little, lamb}")
	}
}
