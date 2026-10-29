package adb_test

// common_test.go holds product-local copies of root-package HCL test
// constants referenced by the migrated adb tests. AdbCommonTestCase is a pure
// HCL data block (no Go logic), safe to copy 1:1; it is unexported in the root
// suite and not reachable from an external _test package.

const AdbCommonTestCase = `
data "alicloud_adb_zones" "default" {}
data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_adb_zones.default.ids.0
}

locals {
  vswitch_id = data.alicloud_vswitches.default.ids.0
}
`
