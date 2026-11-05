package rkvstore

// Local copies of root-package (package alicloud) types and constants that
// cannot be reached via qualified leaf-package calls: the root package is not
// importable from a product sub-package. Per stock-method-disposition, the
// cross-product shared root copy stays in the root package; this private copy
// preserves the required behavior for the migrated R-kvstore types.

// Status mirrors alicloud.Status (common.go), used by the kvstore wait helpers.
type Status string

var (
	Normal  = Status("Normal")
	Deleted = Status("Deleted")
)

const (
	DefaultTimeout     = 120
	DefaultLongTimeout = 1000
)
