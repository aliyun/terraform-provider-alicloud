package ecs_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	ecs "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
)

// Bare enum/alias identifiers used by ECS acceptance and unit tests that were
// in-package under the flat root layout (root alicloud/common.go Running,
// PostPaid; extension_ecs.go NoSpot; the helper.String alias exposed as
// String). In this split package those declarations live in package ecs; these
// package-level aliases restore bare-name visibility for ecs_test so the
// original test sources read verbatim, mirroring their root-layout
// availability. Values are the identical ecs declarations, not redefinitions.
var (
	Running  = ecs.Running
	PostPaid = ecs.PostPaid
	NoSpot   = ecs.NoSpot
	String   = ecs.String

	// ECS acceptance tests also reference these root-layout enums by bare name:
	// common.go DomesticSite / PayByBandwidth and extension_ecs.go
	// DiskTypeData / InstanceTypeResource. PayByBandwidth / DiskTypeData /
	// InstanceTypeResource are exported in package ecs (types.go /
	// extension_ecs.go) and are passed to ecs methods, so they keep the ecs
	// type here. DomesticSite is passed to acctest.PreCheckWithAccountSiteType,
	// which expects acctest.AccountSite; package ecs also re-declares its own
	// AccountSite (types.go), so the ecs.DomesticSite value is not assignable.
	// acctest defines its own DomesticSite of the matching acctest.AccountSite
	// type with the same "Domestic" value, so we alias that to preserve both
	// behavior and the bare-name call site.
	DomesticSite         = acctest.DomesticSite
	PayByBandwidth       = ecs.PayByBandwidth
	DiskTypeData         = ecs.DiskTypeData
	InstanceTypeResource = ecs.InstanceTypeResource
)
