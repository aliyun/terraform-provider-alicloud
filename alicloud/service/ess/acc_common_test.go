package ess_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// ESS-specific acceptance-test harness.
//
// The generic resource/attribute reflection engine (resourceCheck /
// resourceAttr / resourceAttrCheck, the sentinel constants, and the HCL
// config builder) was duplicated verbatim from alicloud/acctest/common.go.
// The split ess_test package imports that package directly, so the copy is
// removed: ESS callers now use the shared acctest.* surface
// (ResourceCheckInit / ResourceAttrInit / ResourceAttrCheckInit /
// ResourceTestAccConfig / CheckResourceExists / CheckResourceDestroy /
// ResourceAttrMapUpdateSet and the NOSET/CHECKSET/REMOVEKEY/CLEARLIST/
// CLEARMAP/REGEXMATCH/ForceSleep sentinels) instead of package-local twins.
//
// What remains here is ESS-specific: the data-source runner adapter and the
// shared dependency HCL. alicloud/acctest.DataSourceTestCheck(WithPreCheck)
// hardcodes SDKProviderFactories, while the migrated ESS data-source tests
// run on the ProtoV5 provider factories; until acctest offers a ProtoV5
// data-source runner, this thin adapter reproduces the step-building with
// acctest.ProtoV5ProviderFactories() + acctest.PreCheck and delegates the
// attribute-check machinery to acctest.ResourceAttrInit / ResourceAttrMapUpdateSet.

type dataSourceAttr struct {
	// IDRefreshName, like "data.alicloud_dns_records.record"
	resourceId string

	// get existMap function
	existMapFunc func(rand int) map[string]string

	// get fakeMap function
	fakeMapFunc func(rand int) map[string]string
}

// get exist and empty ResourceAttrMapUpdate function
func (dsa *dataSourceAttr) checkDataSourceAttr(rand int) (exist, empty acctest.ResourceAttrMapUpdate) {
	exist = acctest.ResourceAttrInit(dsa.resourceId, dsa.existMapFunc(rand)).ResourceAttrMapUpdateSet()
	empty = acctest.ResourceAttrInit(dsa.resourceId, dsa.fakeMapFunc(rand)).ResourceAttrMapUpdateSet()
	return
}

// according to configs generate step list and execute the test
// lintignore: AT001
func (dsa *dataSourceAttr) dataSourceTestCheck(t *testing.T, rand int, configs ...dataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps:                    steps,
	})
}

// according to configs generate step list and execute the test with preCheck
// lintignore: AT001
func (dsa *dataSourceAttr) dataSourceTestCheckWithPreCheck(t *testing.T, rand int, preCheck func(), configs ...dataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 preCheck,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps:                    steps,
	})
}

// per schema attribute test config
type dataSourceTestAccConfig struct {
	// be equal to testCase config string,but the result has only one record
	existConfig string

	// if the dataSourceAttr.existMapFunc returned map value not match we want, existChangMap can alter checkMap for existConfig
	existChangMap map[string]string

	// be equal to testCase config string,but the result is empty
	fakeConfig string

	// if the dataSourceAttr.fakeMapFunc returned map value not match we want, fakeChangMap can alter checkMap for fakeConfig
	fakeChangMap map[string]string
}

// build test cases for each attribute
func (conf *dataSourceTestAccConfig) buildDataSourceSteps(t *testing.T, info *dataSourceAttr, rand int) []resource.TestStep {
	testAccCheckExist, testAccCheckEmpty := info.checkDataSourceAttr(rand)
	var steps []resource.TestStep
	if conf.existConfig != "" {
		step := resource.TestStep{
			Config: conf.existConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckExist(conf.existChangMap),
			),
		}
		steps = append(steps, step)
	}
	if conf.fakeConfig != "" {
		step := resource.TestStep{
			Config: conf.fakeConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckEmpty(conf.fakeChangMap),
			),
		}
		steps = append(steps, step)
	}
	return steps
}

// EcsInstanceCommonTestCase is shared dependency HCL (vpc/vswitch/security group + zone data) copied verbatim from alicloud/service_alicloud_common_test.go; several ESS resource test dependence configs embed it.
const EcsInstanceCommonTestCase = `
data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu"
  most_recent = true
  owners      = "system"
}
resource "alicloud_vpc" "default" {
  vpc_name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name              = "${var.name}"
}
resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}
resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}
`
