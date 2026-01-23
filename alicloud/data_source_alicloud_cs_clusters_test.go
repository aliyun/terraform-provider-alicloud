package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAckClusterDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAckClusterSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_cs_managed_kubernetes.default.id}"]`,
			"enable_details": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudAckClusterSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_cs_managed_kubernetes.default.id}_fake"]`,
			"enable_details": `"true"`,
		}),
	}

	ClusterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAckClusterSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_cs_managed_kubernetes.default.id}"]`,
			"cluster_id":     `"${alicloud_cs_managed_kubernetes.default.id}"`,
			"enable_details": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudAckClusterSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_cs_managed_kubernetes.default.id}_fake"]`,
			"cluster_id":     `"${alicloud_cs_managed_kubernetes.default.id}_fake"`,
			"enable_details": `"true"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAckClusterSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_cs_managed_kubernetes.default.id}"]`,
			"cluster_id":     `"${alicloud_cs_managed_kubernetes.default.id}"`,
			"enable_details": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudAckClusterSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_cs_managed_kubernetes.default.id}_fake"]`,
			"cluster_id":     `"${alicloud_cs_managed_kubernetes.default.id}_fake"`,
			"enable_details": `"true"`,
		}),
	}

	AckClusterCheckInfo.dataSourceTestCheck(t, rand, idsConf, ClusterIdConf, allConf)
}

var existAckClusterMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clusters.#":                      "1",
		"clusters.0.resource_group_id":    CHECKSET,
		"clusters.0.ip_stack":             CHECKSET,
		"clusters.0.tags.%":               CHECKSET,
		"clusters.0.proxy_mode":           CHECKSET,
		"clusters.0.state":                CHECKSET,
		"clusters.0.deletion_protection":  CHECKSET,
		"clusters.0.vpc_id":               CHECKSET,
		"clusters.0.operation_policy.#":   CHECKSET,
		"clusters.0.maintenance_window.#": CHECKSET,
		"clusters.0.pod_cidr":             CHECKSET,
		"clusters.0.cluster_domain":       CHECKSET,
		"clusters.0.current_version":      CHECKSET,
		"clusters.0.profile":              CHECKSET,
		"clusters.0.vswitch_ids.#":        CHECKSET,
		"clusters.0.service_cidr":         CHECKSET,
		"clusters.0.cluster_name":         CHECKSET,
		"clusters.0.cluster_id":           CHECKSET,
		"clusters.0.security_group_id":    CHECKSET,
		"clusters.0.cluster_type":         CHECKSET,
		"clusters.0.cluster_spec":         CHECKSET,
		"clusters.0.region_id":            CHECKSET,
	}
}

var fakeAckClusterMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clusters.#": "0",
	}
}

var AckClusterCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cs_clusters.default",
	existMapFunc: existAckClusterMapFunc,
	fakeMapFunc:  fakeAckClusterMapFunc,
}

func testAccCheckAlicloudAckClusterSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccAckCluster%d"
}
variable "zone_1" {
  default = "cn-hangzhou-k"
}

variable "zone_2" {
  default = "cn-hangzhou-g"
}

variable "vsw1_cidr" {
  default = "10.1.0.0/24"
}

variable "vsw4_cidr" {
  default = "10.1.3.0/24"
}

variable "rg_name_1" {
  default = "tf-test-resource-group-1"
}

variable "vsw2_cidr" {
  default = "10.1.1.0/24"
}

variable "rg_name_2" {
  default = "tf-test-resource-group-2"
}

variable "container_cidr" {
  default = "172.17.3.0/24"
}

variable "user_data" {
  default = "I18vYmluL3No"
}

variable "service_cidr" {
  default = "172.17.2.0/24"
}

variable "vsw3_cidr" {
  default = "10.1.2.0/24"
}

variable "kubernetes_version" {
  default = "1.32.1-aliyun.1"
}

variable "user_data_1" {
  default = "IyEvYmluL3NoIGVjaG8gIkhlbGxvIFdvcmxkLiBUaGUgdGltZSBpcyBudWcgJChkYXRlIC1SKSkhfCB0ZWUgL3Jvb3QvdXNlcmRhdGFfdGVzdC50eHQ="
}

variable "zone_3" {
  default = "cn-hangzhou-i"
}

variable "zone_4" {
  default = "cn-hangzhou-j"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultqe0KHK" {
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_security_group" "defaultKHUbRj" {
  vpc_id              = alicloud_vpc.defaultqe0KHK.id
  security_group_name = "tf-test-security-group"
  security_group_type = "normal"
}

resource "alicloud_security_group" "defaultKYDOFD" {
  security_group_name = "tf-test-security-group-2"
  vpc_id              = alicloud_vpc.defaultqe0KHK.id
  security_group_type = "normal"
}

resource "alicloud_vswitch" "defaultVTblQn" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  cidr_block = var.vsw1_cidr
  zone_id    = var.zone_1
}

resource "alicloud_vswitch" "defaultziRRat" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_2
  cidr_block = var.vsw2_cidr
}

resource "alicloud_vswitch" "defaultT8D8ss" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_3
  cidr_block = var.vsw3_cidr
}

resource "alicloud_vswitch" "defaultFsk7cj" {
  vpc_id     = alicloud_vpc.defaultqe0KHK.id
  zone_id    = var.zone_4
  cidr_block = var.vsw4_cidr
}

resource "alicloud_cs_managed_kubernetes" "default" {
  pod_cidr          = var.container_cidr
  vswitch_ids       = ["${alicloud_vswitch.defaultT8D8ss.id}", "${alicloud_vswitch.defaultziRRat.id}"]
  service_cidr      = var.service_cidr
  security_group_id = alicloud_security_group.defaultKHUbRj.id
  cluster_spec      = "ack.pro.small"
}

data "alicloud_cs_clusters" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
