package vpc_test

// External-package test shims. package vpc_test cannot construct the
// unexported-client structs VpcService / VpcServiceV2 by literal, cannot
// see the unexported funcs buildRouteEntryResourceId / isRouteEntryWaitTimeout,
// and references the VpcService/VpcServiceV2 type names in reflect / func
// signatures. These aliases and constructor vars re-export the production
// exported symbols under the unqualified names the migrated tests use, so
// the test call sites stay byte-identical.

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpc"
)

type VpcService = vpc.VpcService
type VpcServiceV2 = vpc.VpcServiceV2

var NewVpcService = vpc.NewVpcService
var NewVpcServiceV2 = vpc.NewVpcServiceV2

var buildRouteEntryResourceId = vpc.BuildRouteEntryResourceId
var isRouteEntryWaitTimeout = vpc.IsRouteEntryWaitTimeout

// newInstanceDiff is the root alicloud/common.go unexported helper. Its body
// was migrated verbatim to alicloud/acctest.NewInstanceDiff (see the
// "Migrated to alicloud/acctest.NewInstanceDiff" note in common.go and the
// reciprocal "Migrated from alicloud.newInstanceDiff" note in
// acctest/common.go). acctest.Provider is alicloud.Provider() (acctest.go),
// the same provider instance the vpc_test Provider() shim returns, so
// ResourcesMap lookups for the migrated VPC resource names resolve identically.
// The SDKv2 terraform InstanceState/InstanceDiff parameter types match the
// root helper exactly. Aliasing the exported canonical helper keeps the
// migrated call sites byte-identical.
var newInstanceDiff = acctest.NewInstanceDiff

// COLON_SEPARATED is the root alicloud/common.go unexported string const (":").
// The migrated route_table_attachment test splits a composite ID on it
// (strings.Split(rs.Primary.ID, COLON_SEPARATED)); re-declaring the const
// value verbatim keeps the call site byte-identical.
const COLON_SEPARATED = ":"

// Two migrated unit tests call the exported factory funcs by their historical
// lowercase names (pre-existing typos that compiled only because the factory
// and the call site shared package vpc). They are not CRUD registrations, so
// they cannot route through the factory-field wrappers; aliasing the exported
// factory under the lowercase name fixes the references without editing the
// call sites and keeps these test files free of a direct vpc import.
var resourceAlicloudVpcDhcpOptionsSetAttachement = vpc.ResourceAlicloudVpcDhcpOptionsSetAttachement
var resourceAliCloudVpcVpc = vpc.ResourceAliCloudVpcVpc

// VPC status and timeout constants live in package vpc
// (alicloud/service/vpc/status.go); the migrated tests referenced them
// unqualified as the root baseline did (WaitForRouteEntry(id, Available,
// DefaultTimeout), WaitForNetworkAclAttachment(..., Deleted, DefaultTimeout),
// case string(Inactive), string(Idle)). Re-export the subset the tests use
// under the same names so the call sites stay byte-identical.
var (
	Available      = vpc.Available
	Active         = vpc.Active
	Deleted        = vpc.Deleted
	Idle           = vpc.Idle
	Inactive       = vpc.Inactive
	DefaultTimeout = vpc.DefaultTimeout
)
