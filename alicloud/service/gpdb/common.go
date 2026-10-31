// Package gpdb contains the gpdb service resources and data sources migrated
// from the root alicloud package. The root package cannot be imported from a
// service package, so the general status type and the trivial separators/host
// constants that have no equivalent export in a leaf package are mirrored here
// verbatim from alicloud/common.go.
package gpdb

// Status mirrors the root package status type used by the WaitFor helpers.
type Status string

const (
	// Deleted is the target status used by WaitFor* helpers when a resource is
	// expected to be gone after destruction.
	Deleted = Status("Deleted")

	// DefaultTimeoutMedium is the default medium timeout (in seconds) for gpdb
	// wait helpers. Mirrors alicloud/common.go.
	DefaultTimeoutMedium = 500

	// Trivial separator/host constants mirroring alicloud/common.go that have no
	// equivalent export in a leaf package.
	COLON_SEPARATED = ":"
	LOCAL_HOST_IP   = "127.0.0.1"
)
