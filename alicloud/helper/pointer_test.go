package helper

import (
	"testing"
)

// Tests for the pointer helpers moved from the root package
// (alicloud/common.go); test names are unchanged from
// alicloud/common_test.go.

func TestUnitCommonPointerFunctions(t *testing.T) {
	t.Run("StringPointer", func(t *testing.T) {
		str := "test"
		ptr := StringPointer(str)
		if ptr == nil || *ptr != str {
			t.Error("StringPointer failed")
		}
	})

	t.Run("BoolPointer", func(t *testing.T) {
		b := true
		ptr := BoolPointer(b)
		if ptr == nil || *ptr != b {
			t.Error("BoolPointer failed")
		}
	})

	t.Run("Int32Pointer", func(t *testing.T) {
		i := int32(42)
		ptr := Int32Pointer(i)
		if ptr == nil || *ptr != i {
			t.Error("Int32Pointer failed")
		}
	})

	t.Run("Int64Pointer", func(t *testing.T) {
		i := int64(100)
		ptr := Int64Pointer(i)
		if ptr == nil || *ptr != i {
			t.Error("Int64Pointer failed")
		}
	})
}
