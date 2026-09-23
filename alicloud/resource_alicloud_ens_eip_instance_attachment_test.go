package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens EipInstanceAttachment. >>> Resource test cases, automatically generated.
// Case LB绑定EIP_20240426 6608
func TestAccAliCloudEnsEipInstanceAttachment_basic6608(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_eip_instance_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsEipInstanceAttachmentMap6608)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsEipInstanceAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%senseipinstanceattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsEipInstanceAttachmentBasicDependence6608)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":   "${alicloud_ens_load_balancer.defaultj8Egvj.id}",
					"allocation_id": "${alicloud_ens_eip.defaultsGsN4e.id}",
					"instance_type": "SlbInstance",
					"standby":       "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   CHECKSET,
						"allocation_id": CHECKSET,
						"instance_type": "SlbInstance",
						"standby":       "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEnsEipInstanceAttachmentMap6608 = map[string]string{
	"status":        CHECKSET,
	"instance_type": CHECKSET,
}

func AlicloudEnsEipInstanceAttachmentBasicDependence6608(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network" "defaultmIE0nO" {
  network_name  = var.name
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "defaultDiOqwH" {
  cidr_block    = "10.0.0.0/24"
  vswitch_name  = var.name
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.defaultmIE0nO.id
}

resource "alicloud_ens_load_balancer" "defaultj8Egvj" {
  vswitch_id         = alicloud_ens_vswitch.defaultDiOqwH.id
  ens_region_id      = var.ens_region_id
  network_id         = alicloud_ens_vswitch.defaultDiOqwH.network_id
  load_balancer_spec = "elb.s1.small"
  payment_type       = "PayAsYouGo"
}

resource "alicloud_ens_eip" "defaultsGsN4e" {
  bandwidth            = "5"
  eip_name             = var.name
  ens_region_id        = var.ens_region_id
  internet_charge_type = "95BandwidthByMonth"
  payment_type         = "PayAsYouGo"
}


`, name)
}

// Case 实例绑定EIP_20240423 6589
func TestAccAliCloudEnsEipInstanceAttachment_basic6589(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_eip_instance_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsEipInstanceAttachmentMap6589)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsEipInstanceAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%senseipinstanceattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsEipInstanceAttachmentBasicDependence6589)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":   "${alicloud_ens_instance.defaultXKjq1W.id}",
					"allocation_id": "${alicloud_ens_eip.defaultsGsN4e.id}",
					"instance_type": "EnsInstance",
					"standby":       "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   CHECKSET,
						"allocation_id": CHECKSET,
						"instance_type": "EnsInstance",
						"standby":       "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEnsEipInstanceAttachmentMap6589 = map[string]string{
	"status":        CHECKSET,
	"instance_type": CHECKSET,
}

func AlicloudEnsEipInstanceAttachmentBasicDependence6589(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_instance" "defaultXKjq1W" {
  system_disk {
    size     = "20"
    category = "cloud_efficiency"
  }
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.stiny"
  password                   = "12345678abcABC"
  status                     = "Running"
  amount                     = "1"
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = var.name
  auto_use_coupon            = "true"
  instance_charge_strategy   = "PriceHighPriority"
  ens_region_id              = var.ens_region_id
  period_unit                = "Month"
}

resource "alicloud_ens_eip" "defaultsGsN4e" {
  bandwidth            = "5"
  eip_name             = var.name
  ens_region_id        = var.ens_region_id
  internet_charge_type = "95BandwidthByMonth"
  payment_type         = "PayAsYouGo"
}


`, name)
}

// Case LB绑定EIP_20240426 6608  raw
func TestAccAliCloudEnsEipInstanceAttachment_basic6608_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_eip_instance_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsEipInstanceAttachmentMap6608)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsEipInstanceAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%senseipinstanceattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsEipInstanceAttachmentBasicDependence6608)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":   "${alicloud_ens_load_balancer.defaultj8Egvj.id}",
					"allocation_id": "${alicloud_ens_eip.defaultsGsN4e.id}",
					"instance_type": "SlbInstance",
					"standby":       "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   CHECKSET,
						"allocation_id": CHECKSET,
						"instance_type": "SlbInstance",
						"standby":       "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Case 实例绑定EIP_20240423 6589  raw
func TestAccAliCloudEnsEipInstanceAttachment_basic6589_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_eip_instance_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsEipInstanceAttachmentMap6589)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsEipInstanceAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%senseipinstanceattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsEipInstanceAttachmentBasicDependence6589)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":   "${alicloud_ens_instance.defaultXKjq1W.id}",
					"allocation_id": "${alicloud_ens_eip.defaultsGsN4e.id}",
					"instance_type": "EnsInstance",
					"standby":       "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   CHECKSET,
						"allocation_id": CHECKSET,
						"instance_type": "EnsInstance",
						"standby":       "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test Ens EipInstanceAttachment. <<< Resource test cases, automatically generated.
