package alb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudAlbLoadBalancersDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_load_balancer.default.id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_load_balancer.default.id}_fake"]`,
		}),
	}
	addressTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_alb_load_balancer.default.id}"]`,
			"address_type": `"Internet"`,
		}),
		FakeConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"address_type": `"Intranet"`,
		}),
	}
	loadBalancerBussinessstatus := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":                           `["${alicloud_alb_load_balancer.default.id}"]`,
			"load_balancer_business_status": `"Normal"`,
		}),
		FakeConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":                           `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"load_balancer_business_status": `"Abnormal"`,
		}),
	}
	vpcIdsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_alb_load_balancer.default.id}"]`,
			"vpc_ids": `["${alicloud_alb_load_balancer.default.vpc_id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"vpc_ids": `["${alicloud_alb_load_balancer.default.vpc_id}"]`,
		}),
	}
	vpcIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_load_balancer.default.id}"]`,
			"vpc_id": `"${alicloud_alb_load_balancer.default.vpc_id}"`,
		}),
		FakeConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"vpc_id": `"${alicloud_alb_load_balancer.default.vpc_id}"`,
		}),
	}
	loadBalancerIdsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"load_balancer_ids": `["${alicloud_alb_load_balancer.default.id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"load_balancer_ids": `["${alicloud_alb_load_balancer.default.id}_fake"]`,
		}),
	}
	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_load_balancer.default.id}"]`,
			"resource_group_id": `"${alicloud_alb_load_balancer.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"resource_group_id": `"${alicloud_alb_load_balancer.default.resource_group_id}_fake"`,
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_alb_load_balancer.default.id}"]`,
			"tags": `{Created = "TF"}`,
		}),
		FakeConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"tags": `{Created = "TF_fake"}`,
		}),
	}
	loadBalancerNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"load_balancer_name": `"${alicloud_alb_load_balancer.default.load_balancer_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"load_balancer_name": `"${alicloud_alb_load_balancer.default.load_balancer_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_load_balancer.default.id}"]`,
			"status": `"Active"`,
		}),
		FakeConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_load_balancer.default.id}"]`,
			"status": `"Configuring"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_load_balancer.default.load_balancer_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_load_balancer.default.load_balancer_name}_fake"`,
		}),
	}
	zoneIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_alb_load_balancer.default.id}"]`,
			"zone_id": `"${data.alicloud_alb_zones.default.zones.0.id}"`,
		}),
		FakeConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"zone_id": `"${data.alicloud_alb_zones.default.zones.0.id}"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":                            `["${alicloud_alb_load_balancer.default.id}"]`,
			"name_regex":                     `"${alicloud_alb_load_balancer.default.load_balancer_name}"`,
			"vpc_id":                         `"${alicloud_alb_load_balancer.default.vpc_id}"`,
			"zone_id":                        `"${data.alicloud_alb_zones.default.zones.0.id}"`,
			"load_balancer_name":             `"${alicloud_alb_load_balancer.default.load_balancer_name}"`,
			"status":                         `"Active"`,
			"tags":                           `{Created = "TF"}`,
			"resource_group_id":              `"${alicloud_alb_load_balancer.default.resource_group_id}"`,
			"load_balancer_ids":              `["${alicloud_alb_load_balancer.default.id}"]`,
			"load_balancer_business_status":  `"Normal"`,
			"load_balancer_bussiness_status": `"Normal"`,
			"address_type":                   `"Internet"`,
		}),
		FakeConfig: testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":                            `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"name_regex":                     `"${alicloud_alb_load_balancer.default.load_balancer_name}_fake"`,
			"vpc_id":                         `"${alicloud_alb_load_balancer.default.vpc_id}"`,
			"zone_id":                        `"${data.alicloud_alb_zones.default.zones.0.id}"`,
			"load_balancer_name":             `"${alicloud_alb_load_balancer.default.load_balancer_name}_fake"`,
			"status":                         `"Configuring"`,
			"tags":                           `{Created = "TF_fake"}`,
			"resource_group_id":              `"${alicloud_alb_load_balancer.default.resource_group_id}_fake"`,
			"load_balancer_ids":              `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"load_balancer_business_status":  `"Abnormal"`,
			"load_balancer_bussiness_status": `"Abnormal"`,
			"address_type":                   `"Intranet"`,
		}),
	}
	var existAliCloudAlbLoadBalancersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                                 "1",
			"names.#":                                               "1",
			"balancers.#":                                           "1",
			"balancers.0.load_balancer_name":                        acctest.CHECKSET,
			"balancers.0.address_type":                              acctest.CHECKSET,
			"balancers.0.address_allocated_mode":                    acctest.CHECKSET,
			"balancers.0.load_balancer_edition":                     acctest.CHECKSET,
			"balancers.0.zone_mappings.#":                           acctest.CHECKSET,
			"balancers.0.zone_mappings.0.vswitch_id":                acctest.CHECKSET,
			"balancers.0.zone_mappings.0.zone_id":                   acctest.CHECKSET,
			"balancers.0.zone_mappings.0.status":                    acctest.CHECKSET,
			"balancers.0.zone_mappings.0.load_balancer_addresses.#": acctest.CHECKSET,
			"balancers.0.zone_mappings.0.load_balancer_addresses.0.allocation_id":              acctest.CHECKSET,
			"balancers.0.zone_mappings.0.load_balancer_addresses.0.eip_type":                   acctest.CHECKSET,
			"balancers.0.zone_mappings.0.load_balancer_addresses.0.address":                    acctest.CHECKSET,
			"balancers.0.zone_mappings.0.load_balancer_addresses.0.intranet_address":           acctest.CHECKSET,
			"balancers.0.zone_mappings.0.load_balancer_addresses.0.intranet_address_hc_status": acctest.CHECKSET,
			"balancers.0.zone_mappings.0.load_balancer_addresses.0.ipv6_address":               acctest.CHECKSET,
			"balancers.0.zone_mappings.0.load_balancer_addresses.0.ipv6_address_hc_status":     acctest.CHECKSET,
			"balancers.0.zone_mappings.0.load_balancer_addresses.0.ipv4_local_addresses.#":     acctest.CHECKSET,
			"balancers.0.zone_mappings.0.load_balancer_addresses.0.ipv6_local_addresses.#":     acctest.CHECKSET,
		}
	}
	var fakeAliCloudAlbLoadBalancersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"balancers.#": "0",
		}
	}
	var AliCloudAlbLoadBalancersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_alb_load_balancers.default",
		ExistMapFunc: existAliCloudAlbLoadBalancersDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudAlbLoadBalancersDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}
	AliCloudAlbLoadBalancersCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, addressTypeConf, loadBalancerBussinessstatus, vpcIdsConf, vpcIdConf, loadBalancerIdsConf, statusConf, resourceGroupIdConf, tagsConf, loadBalancerNameConf, nameRegexConf, zoneIdConf, allConf)
}
func testAccCheckAliCloudAlbLoadBalancersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {	
		default = "tf-testaccalb%d"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_alb_zones" "default" {
	}

	resource "alicloud_vpc" "default" {
  		vpc_name    = var.name
  		cidr_block  = "192.168.0.0/16"
  		enable_ipv6 = "true"
	}

	resource "alicloud_eip" "zone_a" {
  		bandwidth            = "10"
  		internet_charge_type = "PayByTraffic"
	}

	resource "alicloud_vswitch" "zone_a" {
  		vswitch_name         = var.name
  		vpc_id               = alicloud_vpc.default.id
  		cidr_block           = "192.168.0.0/18"
  		zone_id              = data.alicloud_alb_zones.default.zones.0.id
  		ipv6_cidr_block_mask = "6"
	}

	resource "alicloud_vswitch" "zone_b" {
  		vswitch_name         = var.name
  		vpc_id               = alicloud_vpc.default.id
  		cidr_block           = "192.168.128.0/18"
  		zone_id              = data.alicloud_alb_zones.default.zones.1.id
  		ipv6_cidr_block_mask = "8"
	}

	resource "alicloud_vpc_ipv6_gateway" "default" {
  		ipv6_gateway_name = var.name
  		vpc_id            = alicloud_vpc.default.id
	}

	resource "alicloud_common_bandwidth_package" "default" {
  		bandwidth            = 1000
  		internet_charge_type = "PayByBandwidth"
	}

	resource "alicloud_log_project" "default" {
  		project_name = var.name
	}

	resource "alicloud_log_store" "default" {
  		project_name  = alicloud_log_project.default.project_name
  		logstore_name = var.name
	}

	resource "alicloud_alb_load_balancer" "default" {
  		load_balancer_edition       = "Basic"
  		address_type                = "Internet"
  		vpc_id                      = alicloud_vpc_ipv6_gateway.default.vpc_id
  		address_allocated_mode      = "Fixed"
  		address_ip_version          = "DualStack"
  		ipv6_address_type           = "Internet"
  		bandwidth_package_id        = alicloud_common_bandwidth_package.default.id
  		resource_group_id           = data.alicloud_resource_manager_resource_groups.default.groups.1.id
  		load_balancer_name          = var.name
  		deletion_protection_enabled = false
  		load_balancer_billing_config {
    		pay_type = "PayAsYouGo"
  		}
  		modification_protection_config {
    		status = "NonProtection"
  		}
  		access_log_config {
    		log_project = alicloud_log_store.default.project_name
    		log_store   = alicloud_log_store.default.logstore_name
  		}
  		zone_mappings {
    		vswitch_id       = alicloud_vswitch.zone_a.id
    		zone_id          = alicloud_vswitch.zone_a.zone_id
    		eip_type         = "Common"
    		allocation_id    = alicloud_eip.zone_a.id
    		intranet_address = "192.168.10.1"
  		}
  		zone_mappings {
    		vswitch_id = alicloud_vswitch.zone_b.id
    		zone_id    = alicloud_vswitch.zone_b.zone_id
  		}
  		tags = {
    		Created = "TF"
  		}
	}

	data "alicloud_alb_load_balancers" "default" {	
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
