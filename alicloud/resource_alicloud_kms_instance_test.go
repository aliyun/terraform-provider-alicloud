package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudKmsInstance_basic4048(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsInstanceMap4048)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsInstanceBasicDependence4048)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KmsInstanceSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num":         "1",
					"key_num":         "1000",
					"secret_num":      "0",
					"spec":            "1000",
					"product_version": "3",
					"vpc_id":          "${alicloud_vpc.default.id}",
					"log":             "0",
					"log_storage":     "0",
					"zone_ids": []string{
						"cn-hangzhou-k", "cn-hangzhou-j"},
					"vswitch_ids": []string{
						"${alicloud_vswitch.vswitch.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num":         "1",
						"key_num":         "1000",
						"secret_num":      "0",
						"spec":            "1000",
						"product_version": "3",
						"vpc_id":          CHECKSET,
						"zone_ids.#":      "2",
						"vswitch_ids.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_vpcs": []map[string]interface{}{
						{
							"vpc_id":       "${alicloud_vswitch.shareVswitch.vpc_id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.shareVswitch.id}",
							"vpc_owner_id": "${data.alicloud_account.current.id}",
						},
						{
							"vpc_id":       "${alicloud_vswitch.share-vswitch2.vpc_id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.share-vswitch2.id}",
							"vpc_owner_id": "${data.alicloud_account.current.id}",
						},
						{
							"vpc_id":       "${alicloud_vswitch.share-vsw3.vpc_id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.share-vsw3.id}",
							"vpc_owner_id": "${data.alicloud_account.current.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bind_vpcs.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_num": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_num": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_num": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_num": "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_num": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_num": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log":         "1",
					"log_storage": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log":         "1",
						"log_storage": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_vpcs": []map[string]interface{}{
						{
							"vpc_id":       "${alicloud_vswitch.share-vswitch2.vpc_id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.share-vswitch2.id}",
							"vpc_owner_id": "${data.alicloud_account.current.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bind_vpcs.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_vpcs": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bind_vpcs.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num":         "5",
					"key_num":         "2000",
					"secret_num":      "2000",
					"spec":            "2000",
					"renew_status":    "ManualRenewal",
					"product_version": "3",
					"renew_period":    "3",
					"vpc_id":          "${alicloud_vpc.default.id}",
					"zone_ids": []string{
						"cn-hangzhou-k", "cn-hangzhou-j"},
					"vswitch_ids": []string{
						"${alicloud_vswitch.vswitch-j.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num":         "5",
						"key_num":         "2000",
						"secret_num":      "2000",
						"spec":            "2000",
						"renew_status":    "ManualRenewal",
						"product_version": "3",
						"renew_period":    "3",
						"vpc_id":          CHECKSET,
						"zone_ids.#":      "2",
						"vswitch_ids.#":   "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_version", "renew_period", "renew_status", "period"},
			},
		},
	})
}

var AlicloudKmsInstanceMap4048 = map[string]string{
	"status":                   CHECKSET,
	"create_time":              CHECKSET,
	"end_date":                 CHECKSET,
	"instance_name":            CHECKSET,
	"ca_certificate_chain_pem": CHECKSET,
}

func AlicloudKmsInstanceBasicDependence4048(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "current" {}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "vswitch-j" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = "cn-hangzhou-j"
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_vpc" "shareVPC" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}3"
}

resource "alicloud_vswitch" "shareVswitch" {
  vpc_id     = alicloud_vpc.shareVPC.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}5"
}

resource "alicloud_vswitch" "share-vswitch2" {
  vpc_id     = alicloud_vpc.share-VPC2.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC3" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}7"
}

resource "alicloud_vswitch" "share-vsw3" {
  vpc_id     = alicloud_vpc.share-VPC3.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}


`, name)
}

func AlicloudKmsInstanceBasicDependence4048_intl(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "current" {}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = "ap-southeast-1a"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "vswitch-j" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = "ap-southeast-1b"
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_vpc" "shareVPC" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}3"
}

resource "alicloud_vswitch" "shareVswitch" {
  vpc_id     = alicloud_vpc.shareVPC.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}5"
}

resource "alicloud_vswitch" "share-vswitch2" {
  vpc_id     = alicloud_vpc.share-VPC2.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC3" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}7"
}

resource "alicloud_vswitch" "share-vsw3" {
  vpc_id     = alicloud_vpc.share-VPC3.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}


`, name)
}

// Case 4048  twin
func TestAccAliCloudKmsInstance_basic4048_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsInstanceMap4048)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsInstanceBasicDependence4048)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.KmsInstanceSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":   name,
					"vpc_num":         "7",
					"key_num":         "2000",
					"secret_num":      "1000",
					"spec":            "2000",
					"renew_status":    "ManualRenewal",
					"product_version": "3",
					"renew_period":    "3",
					"log":             "1",
					"log_storage":     "1000",
					"period":          "2",
					"vpc_id":          "${alicloud_vpc.default.id}",
					"zone_ids": []string{
						"cn-hangzhou-k", "cn-hangzhou-j"},
					"vswitch_ids": []string{
						"${alicloud_vswitch.vswitch.id}"},
					"bind_vpcs": []map[string]interface{}{
						{
							"vpc_id":       "${alicloud_vpc.shareVPC.id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.shareVswitch.id}",
							"vpc_owner_id": "${data.alicloud_account.current.id}",
						},
						{
							"vpc_id":       "${alicloud_vswitch.share-vsw3.vpc_id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.share-vsw3.id}",
							"vpc_owner_id": "${data.alicloud_account.current.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":   name,
						"vpc_num":         "7",
						"key_num":         "2000",
						"secret_num":      "1000",
						"spec":            "2000",
						"renew_status":    "ManualRenewal",
						"product_version": "3",
						"renew_period":    "3",
						"vpc_id":          CHECKSET,
						"zone_ids.#":      "2",
						"vswitch_ids.#":   "1",
						"bind_vpcs.#":     "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_vpcs": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bind_vpcs.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_version", "renew_period", "renew_status", "period"},
			},
		},
	})
}

func TestAccAliCloudKmsInstance_basic4048_postpaid(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsInstanceMap4048)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsInstanceBasicDependence4048)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.KmsInstanceSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":    "PayAsYouGo",
					"product_version": "3",
					"vpc_id":          "${alicloud_vpc.default.id}",
					"zone_ids": []string{
						"cn-hangzhou-k", "cn-hangzhou-j"},
					"vswitch_ids": []string{
						"${alicloud_vswitch.vswitch-j.id}"},
					"bind_vpcs": []map[string]interface{}{
						{
							"vpc_id":       "${alicloud_vpc.shareVPC.id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.shareVswitch.id}",
							"vpc_owner_id": "${data.alicloud_account.current.id}",
						},
						{
							"vpc_id":       "${alicloud_vswitch.share-vsw3.vpc_id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.share-vsw3.id}",
							"vpc_owner_id": "${data.alicloud_account.current.id}",
						},
					},
					"force_delete_without_backup": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_version": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_vpcs": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bind_vpcs.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_vpcs": []map[string]interface{}{
						{
							"vpc_id":       "${alicloud_vswitch.shareVswitch.vpc_id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.shareVswitch.id}",
							"vpc_owner_id": "${data.alicloud_account.current.id}",
						},
						{
							"vpc_id":       "${alicloud_vswitch.share-vswitch2.vpc_id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.share-vswitch2.id}",
							"vpc_owner_id": "${data.alicloud_account.current.id}",
						},
						{
							"vpc_id":       "${alicloud_vswitch.share-vsw3.vpc_id}",
							"region_id":    "cn-hangzhou",
							"vswitch_id":   "${alicloud_vswitch.share-vsw3.id}",
							"vpc_owner_id": "${data.alicloud_account.current.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bind_vpcs.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_vpcs": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bind_vpcs.#": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_version", "renew_period", "renew_status", "period", "force_delete_without_backup"},
			},
		},
	})
}

func TestAccAliCloudKmsInstance_basic4048_postpaid_intl(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsInstanceMap4048)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsInstanceBasicDependence5405)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":    "PayAsYouGo",
					"product_version": "3",
					"vpc_id":          "${alicloud_vswitch.vswitch.vpc_id}",
					"zone_ids": []string{
						"${alicloud_vswitch.vswitch.zone_id}", "${alicloud_vswitch.vswitch-j.zone_id}"},
					"vswitch_ids": []string{
						"${alicloud_vswitch.vswitch.id}"},
					"force_delete_without_backup": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_version": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bind_vpcs": []map[string]interface{}{
						{
							"vpc_id":       "${alicloud_vswitch.shareVswitch.vpc_id}",
							"region_id":    defaultRegionToTest,
							"vswitch_id":   "${alicloud_vswitch.shareVswitch.id}",
							"vpc_owner_id": "${data.alicloud_account.current.id}",
						},
						{
							"vpc_id":       "${alicloud_vswitch.share-vswitch2.vpc_id}",
							"region_id":    defaultRegionToTest,
							"vswitch_id":   "${alicloud_vswitch.share-vswitch2.id}",
							"vpc_owner_id": "${data.alicloud_account.current.id}",
						},
						{
							"vpc_id":       "${alicloud_vswitch.share-vsw3.vpc_id}",
							"region_id":    defaultRegionToTest,
							"vswitch_id":   "${alicloud_vswitch.share-vsw3.id}",
							"vpc_owner_id": "${data.alicloud_account.current.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bind_vpcs.#": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_version", "renew_period", "renew_status", "period", "force_delete_without_backup"},
			},
		},
	})
}

func TestAccAliCloudKmsInstance_basic4048_intl(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsInstanceMap4048)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsInstanceBasicDependence4048_intl)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KmsInstanceIntlSupportRegions)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num":         "2",
					"key_num":         "1000",
					"secret_num":      "1000",
					"spec":            "1000",
					"renew_status":    "ManualRenewal",
					"product_version": "3",
					"renew_period":    "3",
					"vpc_id":          "${alicloud_vpc.default.id}",
					"zone_ids": []string{
						"ap-southeast-1a", "ap-southeast-1b"},
					"vswitch_ids": []string{
						"${alicloud_vswitch.vswitch.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num":         "2",
						"key_num":         "1000",
						"secret_num":      "1000",
						"spec":            "1000",
						"renew_status":    "ManualRenewal",
						"product_version": "3",
						"renew_period":    "3",
						"vpc_id":          CHECKSET,
						"zone_ids.#":      "2",
						"vswitch_ids.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_num": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_num": "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_num": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_num": "2000",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_version", "renew_period", "renew_status", "period"},
			},
		},
	})
}

// Test Kms Instance. >>> Resource test cases, automatically generated.
// Case 国际站小规格日志——国际账号 5405
func TestAccAliCloudKmsInstance_basic5405(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsInstanceMap5405)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsInstanceBasicDependence5405)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num":    "5",
					"key_num":    "100",
					"secret_num": "0",
					"spec":       "200",
					"vpc_id":     "${alicloud_vswitch.vswitch.vpc_id}",
					"zone_ids": []string{
						"${alicloud_vswitch.vswitch.zone_id}", "${alicloud_vswitch.vswitch-j.zone_id}"},
					"vswitch_ids": []string{
						"${alicloud_vswitch.vswitch.id}"},
					"period": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num":       "5",
						"key_num":       "100",
						"secret_num":    "0",
						"spec":          "200",
						"vpc_id":        CHECKSET,
						"zone_ids.#":    "2",
						"vswitch_ids.#": "1",
						"period":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_num": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_num": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_num": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_num": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num":         "5",
					"key_num":         "100",
					"secret_num":      "0",
					"spec":            "200",
					"renew_status":    "AutoRenewal",
					"product_version": "5",
					"vpc_id":          "${alicloud_vswitch.vswitch.vpc_id}",
					"zone_ids": []string{
						"${alicloud_vswitch.vswitch.zone_id}", "${alicloud_vswitch.vswitch-j.zone_id}"},
					"vswitch_ids": []string{
						"${alicloud_vswitch.vswitch.id}"},
					"period":       "1",
					"renew_period": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num":         "5",
						"key_num":         "100",
						"secret_num":      "0",
						"spec":            "200",
						"renew_status":    "AutoRenewal",
						"product_version": "5",
						"vpc_id":          CHECKSET,
						"zone_ids.#":      "2",
						"vswitch_ids.#":   "1",
						"period":          "1",
						"renew_period":    "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "product_version", "renew_period", "renew_status"},
			},
		},
	})
}

var AlicloudKmsInstanceMap5405 = map[string]string{
	"ca_certificate_chain_pem": CHECKSET,
	"log_storage":              "0",
	"status":                   CHECKSET,
	"log":                      "0",
	"create_time":              CHECKSET,
	"instance_name":            CHECKSET,
}

func AlicloudKmsInstanceBasicDependence5405(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "current" {}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc-amp-instance-test" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name

}

resource "alicloud_vswitch" "vswitch" {
  vpc_id     = alicloud_vpc.vpc-amp-instance-test.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "vswitch-j" {
  vpc_id     = alicloud_vpc.vpc-amp-instance-test.id
  zone_id    = data.alicloud_zones.default.zones.2.id
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_vpc" "shareVPC" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}3"
}

resource "alicloud_vswitch" "shareVswitch" {
  vpc_id     = alicloud_vpc.shareVPC.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}5"
}

resource "alicloud_vswitch" "share-vswitch2" {
  vpc_id     = alicloud_vpc.share-VPC2.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC3" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}7"
}

resource "alicloud_vswitch" "share-vsw3" {
  vpc_id     = alicloud_vpc.share-VPC3.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}
`, name)
}

// Case 国际站小规格日志——国际账号 5405  twin
func TestAccAliCloudKmsInstance_basic5405_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsInstanceMap5405)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsInstanceBasicDependence5405)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num":         "5",
					"key_num":         "100",
					"secret_num":      "0",
					"spec":            "200",
					"renew_status":    "AutoRenewal",
					"product_version": "5",
					"vpc_id":          "${alicloud_vswitch.vswitch.vpc_id}",
					"zone_ids": []string{
						"${alicloud_vswitch.vswitch.zone_id}", "${alicloud_vswitch.vswitch-j.zone_id}"},
					"vswitch_ids": []string{
						"${alicloud_vswitch.vswitch.id}"},
					"period":       "1",
					"renew_period": "1",
					"log":          "1",
					"log_storage":  "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num":         "5",
						"key_num":         "100",
						"secret_num":      "0",
						"spec":            "200",
						"renew_status":    "AutoRenewal",
						"product_version": "5",
						"vpc_id":          CHECKSET,
						"zone_ids.#":      "2",
						"vswitch_ids.#":   "1",
						"period":          "1",
						"renew_period":    "1",
						"log":             "1",
						"log_storage":     "1000",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "product_version", "renew_period", "renew_status"},
			},
		},
	})
}

// Test Kms Instance. <<< Resource test cases, automatically generated.
