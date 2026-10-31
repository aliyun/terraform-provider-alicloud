package vpc

import (
	"fmt"
	"strings"
)

// Root-package-local CEN helpers duplicated into the VPC service package. The
// migrated alicloud_cen_instance_grant resource (and VpcService v1) reference
// these unqualified while in the root alicloud package; the root package cannot
// be imported from a service package, so the subset they use is kept here
// verbatim, preserving identical values and logic.

const COLON_SEPARATED = ":"

// DefaultCenTimeout mirrors service_alicloud_cen.go.
const DefaultCenTimeout = 60

const ChildInstanceTypeVpc = "VPC"
const ChildInstanceTypeVbr = "VBR"
const ChildInstanceTypeCcn = "CCN"

// GetCenChildInstanceType mirrors common.go; used by the CEN instance grant
// resource to resolve the child instance kind from its ID prefix.
func GetCenChildInstanceType(id string) (c string, e error) {
	if strings.HasPrefix(id, "vpc") {
		return ChildInstanceTypeVpc, nil
	} else if strings.HasPrefix(id, "vbr") {
		return ChildInstanceTypeVbr, nil
	} else if strings.HasPrefix(id, "ccn") {
		return ChildInstanceTypeCcn, nil
	} else {
		return c, fmt.Errorf("CEN child instance ID invalid. Now, it only supports VPC or VBR or CCN instance.")
	}
}
