package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudKmsInstanceDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudKmsInstanceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_kms_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudKmsInstanceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_kms_instance.default.id}_fake"]`,
		}),
	}

	KmsInstanceCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

var existKmsInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":             "1",
		"instances.0.instance_id": CHECKSET,
	}
}

var fakeKmsInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#": "0",
	}
}

var KmsInstanceCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_kms_instances.default",
	existMapFunc: existKmsInstanceMapFunc,
	fakeMapFunc:  fakeKmsInstanceMapFunc,
}

func testAccCheckAlicloudKmsInstanceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccKmsInstance%d"
}
resource "alicloud_vpc" "vpc-amp-instance-test" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id     = alicloud_vpc.vpc-amp-instance-test.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "vswitch-j" {
  vpc_id     = alicloud_vpc.vpc-amp-instance-test.id
  zone_id    = "cn-hangzhou-j"
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_vpc" "shareVPC" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = format("%%s3", var.name)
}

resource "alicloud_vswitch" "shareVswitch" {
  vpc_id     = alicloud_vpc.shareVPC.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = format("%%s5", var.name)
}

resource "alicloud_vswitch" "share-vswitch2" {
  vpc_id     = alicloud_vpc.share-VPC2.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC3" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = format("%%s7", var.name)
}

resource "alicloud_vswitch" "share-vsw3" {
  vpc_id     = alicloud_vpc.share-VPC3.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

data "alicloud_account" "current" {}

resource "alicloud_kms_instance" "default" {
  vpc_num         = "7"
  key_num         = "1000"
  secret_num      = "0"
  spec            = "1000"
  renew_status    = "ManualRenewal"
  product_version = "3"
  renew_period    = "3"
  vpc_id          = alicloud_vswitch.vswitch.vpc_id
  zone_ids        = ["cn-hangzhou-k", "cn-hangzhou-j"]
  vswitch_ids     = ["${alicloud_vswitch.vswitch.id}"]
  bind_vpcs {
    vpc_id       = alicloud_vswitch.shareVswitch.vpc_id
    region_id    = "cn-hangzhou"
    vswitch_id   = alicloud_vswitch.shareVswitch.id
    vpc_owner_id = data.alicloud_account.current.id
  }
  bind_vpcs {
    vpc_id       = alicloud_vswitch.share-vswitch2.vpc_id
    region_id    = "cn-hangzhou"
    vswitch_id   = alicloud_vswitch.share-vswitch2.id
    vpc_owner_id = data.alicloud_account.current.id
  }
  bind_vpcs {
    vpc_id       = alicloud_vswitch.share-vsw3.vpc_id
    region_id    = "cn-hangzhou"
    vswitch_id   = alicloud_vswitch.share-vsw3.id
    vpc_owner_id = data.alicloud_account.current.id
  }
  log          = "0"
  period       = "1"
  log_storage  = "0"
  payment_type = "Subscription"
}

data "alicloud_kms_instances" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
