package helper

// Pointer constructors moved from the root package (alicloud/common.go);
// package alicloud re-exports them as var aliases so the flat layout
// keeps compiling unchanged.

func String(v string) *string {
	return &v
}

func StringPointer(s string) *string {
	return &s
}

func BoolPointer(b bool) *bool {
	return &b
}

func Int32Pointer(i int32) *int32 {
	return &i
}

func Int64Pointer(i int64) *int64 {
	return &i
}
