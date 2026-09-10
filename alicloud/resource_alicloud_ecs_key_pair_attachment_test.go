package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudECSKeyPairAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_key_pair_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsKeyPairAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsKeyPairAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEcsKeyPairAttachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsKeyPairAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"key_pair_name": "${alicloud_ecs_key_pair.default.id}",
					"instance_ids":  "${alicloud_instance.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name":  CHECKSET,
						"instance_ids.#": "5",
					}),
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

func TestAccAliCloudECSKeyPairAttachment_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_key_pair_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsKeyPairAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsKeyPairAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEcsKeyPairAttachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsKeyPairAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"key_pair_name": "${alicloud_ecs_key_pair.default.id}",
					"instance_ids":  "${alicloud_instance.default.*.id}",
					"force":         "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name":  CHECKSET,
						"instance_ids.#": "5",
					}),
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

func TestAccAliCloudECSKeyPairAttachment_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_key_pair_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsKeyPairAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsKeyPairAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEcsKeyPairAttachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsKeyPairAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"key_name":     "${alicloud_ecs_key_pair.default.id}",
					"instance_ids": "${alicloud_instance.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_name":       CHECKSET,
						"instance_ids.#": "5",
					}),
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

func TestAccAliCloudECSKeyPairAttachment_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_key_pair_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsKeyPairAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsKeyPairAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudEcsKeyPairAttachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsKeyPairAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"key_name":     "${alicloud_ecs_key_pair.default.id}",
					"instance_ids": "${alicloud_instance.default.*.id}",
					"force":        "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_name":       CHECKSET,
						"instance_ids.#": "5",
					}),
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

var AliCloudEcsKeyPairAttachmentMap0 = map[string]string{}

func AliCloudEcsKeyPairAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_images" "default" {
  		name_regex = "^ubuntu_18.*64"
  		owners     = "system"
	}

	data "alicloud_instance_types" "default" {
  		image_id          = data.alicloud_images.default.images.0.id
	}
	
	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}
	
	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
	}
	
	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "default" {
  		count                      = 5
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_efficiency"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = var.name
	}

	resource "alicloud_ecs_key_pair" "default" {
  		key_pair_name = var.name
	}
`, name)
}
