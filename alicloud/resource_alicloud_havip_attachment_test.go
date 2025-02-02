package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_havip_attachment",
		&resource.Sweeper{
			Name: "alicloud_havip_attachment",
			F:    testSweepHavipAttachment,
		})
}

func testSweepHavipAttachment(region string) error {
	rawClient, err := sharedClientForRegionWithBackendRegions(region, true, connectivity.ActiontrailSupportRegions)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeHaVips"
	request := map[string]interface{}{}

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.HaVips.HaVip", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.HaVips.HaVip", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping HaVip Attachment: %s", item["Name"].(string))
				continue
			}
			action := "UnassociateHaVip"
			request := map[string]interface{}{
				"HaVipId":    item["HaVipId"],
				"InstanceId": item["InstanceId"],
				"RegionId":   client.RegionId,
			}
			_, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] Failed to delete HaVip Attachment (%s): %s", item["Name"].(string), err)
			}
			log.Printf("[INFO] Delete HaVip Attachment success: %s ", item["Name"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudVPCHavipAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_havip_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudHavipAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcHaVipAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shavipattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHavipAttachmentBasicDependence0)
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
					"havip_id":    "${alicloud_havip.foo.id}",
					"instance_id": "${alicloud_instance.foo.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"havip_id":      CHECKSET,
						"instance_id":   CHECKSET,
						"instance_type": "EcsInstance",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudVPCHavipAttachment_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_havip_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudHavipAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcHaVipAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shavipattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHavipAttachmentBasicDependence1)
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
					"havip_id":      "${alicloud_havip.foo.id}",
					"instance_id":   "${alicloud_instance.foo.id}",
					"instance_type": "EcsInstance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"havip_id":      CHECKSET,
						"instance_id":   CHECKSET,
						"instance_type": "EcsInstance",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudVPCHavipAttachment_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_havip_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudHavipAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcHaVipAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shavipattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHavipAttachmentBasicDependence2)
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
					"havip_id":      "${alicloud_havip.foo.id}",
					"instance_id":   "${alicloud_ecs_network_interface_attachment.default[0].network_interface_id}",
					"instance_type": "NetworkInterface",
					"force":         "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"havip_id":      CHECKSET,
						"instance_id":   CHECKSET,
						"instance_type": "NetworkInterface",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"force": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudVPCHavipAttachment_basic_multiple_instance_bug_fix(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_havip_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudHavipAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcHaVipAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shavipattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHavipAttachmentBasicDependence_bug_fix)
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
					"havip_id":      "${alicloud_havip.foo.id}",
					"instance_id":   "${alicloud_ecs_network_interface_attachment.default[0].network_interface_id}",
					"instance_type": "NetworkInterface",
					"force":         "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"havip_id":      CHECKSET,
						"instance_id":   CHECKSET,
						"instance_type": "NetworkInterface",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"force": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

var AlicloudHavipAttachmentMap0 = map[string]string{}

func AlicloudHavipAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
   availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count = 1
  memory_size = 2
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu"
  most_recent = true
  owners = "system"
}

resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_havip" "foo" {
  vswitch_id = "${alicloud_vswitch.foo.id}"
  description = "${var.name}"
}

resource "alicloud_security_group" "tf_test_foo" {
  name = "${var.name}"
  description = "foo"
  vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "foo" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  image_id = "${data.alicloud_images.default.images.0.id}"
  # series III
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  system_disk_category = "cloud_efficiency"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
  instance_name = "${var.name}"
  user_data = "echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"
}

`, name)
}

func AlicloudHavipAttachmentBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
   availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count = 1
  memory_size = 2
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu"
  most_recent = true
  owners = "system"
}

resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_havip" "foo" {
  vswitch_id = "${alicloud_vswitch.foo.id}"
  description = "${var.name}"
}

resource "alicloud_security_group" "tf_test_foo" {
  name = "${var.name}"
  description = "foo"
  vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "foo" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  image_id = "${data.alicloud_images.default.images.0.id}"
  # series III
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  system_disk_category = "cloud_efficiency"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
  instance_name = "${var.name}"
  user_data = "echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"
}

`, name)
}

func AlicloudHavipAttachmentBasicDependence2(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g7"
}

data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
  available_instance_type     = data.alicloud_instance_types.default.instance_types.0.id
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

locals {
  vswitch_id = data.alicloud_vswitches.default.ids[0]
}

resource "alicloud_security_group" "default" {
  name        = var.name
  description = var.name
  vpc_id      = data.alicloud_vpcs.default.ids.0
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

resource "alicloud_havip" "foo" {
  vswitch_id = local.vswitch_id
  description = "${var.name}"
}

resource "alicloud_instance" "default" {
  count                = 2
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_name        = var.name
  host_name            = var.name
  image_id             = data.alicloud_images.default.images.0.id
  instance_type        = data.alicloud_instance_types.default.instance_types.0.id
  security_groups      = [alicloud_security_group.default.id]
  vswitch_id           = local.vswitch_id
  system_disk_category = "cloud_essd"
}

resource "alicloud_ecs_network_interface" "default" {
  count                = 2
  network_interface_name = var.name
  vswitch_id             = local.vswitch_id
  security_group_ids     = [alicloud_security_group.default.id]
}

resource "alicloud_ecs_network_interface_attachment" "default" {
  count                = 2
  instance_id          = element(alicloud_instance.default.*.id, count.index)
  network_interface_id = element(alicloud_ecs_network_interface.default.*.id, count.index)
}

`, name)
}

func AlicloudHavipAttachmentBasicDependence_bug_fix(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g7"
}

data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
  available_instance_type     = data.alicloud_instance_types.default.instance_types.0.id
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

locals {
  vswitch_id = data.alicloud_vswitches.default.ids[0]
}

resource "alicloud_security_group" "default" {
  name        = var.name
  description = var.name
  vpc_id      = data.alicloud_vpcs.default.ids.0
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

resource "alicloud_havip" "foo" {
  vswitch_id = local.vswitch_id
  description = "${var.name}"
}

resource "alicloud_instance" "default" {
  count                = 2
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_name        = var.name
  host_name            = var.name
  image_id             = data.alicloud_images.default.images.0.id
  instance_type        = data.alicloud_instance_types.default.instance_types.0.id
  security_groups      = [alicloud_security_group.default.id]
  vswitch_id           = local.vswitch_id
  system_disk_category = "cloud_essd"
}

resource "alicloud_ecs_network_interface" "default" {
  count                = 2
  network_interface_name = var.name
  vswitch_id             = local.vswitch_id
  security_group_ids     = [alicloud_security_group.default.id]
}

resource "alicloud_ecs_network_interface_attachment" "default" {
  count                = 2
  instance_id          = element(alicloud_instance.default.*.id, count.index)
  network_interface_id = element(alicloud_ecs_network_interface.default.*.id, count.index)
}

resource "alicloud_havip_attachment" "default0" {
  havip_id    = alicloud_havip.foo.id
  instance_id = alicloud_ecs_network_interface_attachment.default[1].network_interface_id
  instance_type = "NetworkInterface"
}

`, name)
}
