package alicloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudECSInvocation_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_invocation.default"
	ra := resourceAttrInit(resourceId, AliCloudECSInvocationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsInvocation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsinvocation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudECSInvocationBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": []string{"${alicloud_instance.default.id}"},
					"command_id":  "${alicloud_ecs_command.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id.#": "1",
						"command_id":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"windows_password_name"},
			},
		},
	})
}

func TestAccAliCloudECSInvocation_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_invocation.default"
	ra := resourceAttrInit(resourceId, AliCloudECSInvocationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsInvocation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsinvocation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudECSInvocationBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": []string{"${alicloud_instance.default.id}"},
					"command_id":  "${alicloud_ecs_command.default.id}",
					"repeat_mode": "Period",
					"timed":       "true",
					"frequency":   "0 15 10 ? * *",
					"parameters": map[string]string{
						"name": "exampleName",
					},
					"username":              "root",
					"windows_password_name": "axtSecretPassword",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id.#":   "1",
						"command_id":      CHECKSET,
						"repeat_mode":     "Period",
						"timed":           "true",
						"frequency":       "0 15 10 ? * *",
						"parameters.%":    "1",
						"parameters.name": "exampleName",
						"username":        "root",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"windows_password_name"},
			},
		},
	})
}

func TestAccAliCloudECSInvocation_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_invocation.default"
	ra := resourceAttrInit(resourceId, AliCloudECSInvocationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsInvocation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsinvocation%d", defaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEcsInvocationBasic2Config(name, "tf-testacc"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command_id":           CHECKSET,
						"resource_tag.#":       "1",
						"resource_tag.0.key":   "tf-testacc",
						"resource_tag.0.value": "ecs-invocation",
					}),
				),
			},
			{
				// Changing the resource_tag key forces replacement (the nested
				// key/value fields are ForceNew), exercising the resource_tag
				// TypeSet replacement path.
				Config: testAccAliCloudEcsInvocationBasic2Config(name, "tf-testacc2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command_id":           CHECKSET,
						"resource_tag.#":       "1",
						"resource_tag.0.key":   "tf-testacc2",
						"resource_tag.0.value": "ecs-invocation",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"windows_password_name", "resource_tag"},
			},
		},
	})
}

var AliCloudECSInvocationMap0 = map[string]string{
	"instance_id.#": CHECKSET,
}

func AliCloudECSInvocationBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_essd"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.u1"
  		image_id             = data.alicloud_images.default.images.0.id
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_essd"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = var.name
	}

	resource "alicloud_ecs_command" "default" {
  		name            = var.name
  		command_content = "bHMK"
  		description     = "For Terraform Test"
  		type            = "RunShellScript"
  		working_dir     = "/root"
	}
`, name)
}

func AliCloudECSInvocationBasicDependence1(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_essd"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.u1"
  		image_id             = data.alicloud_images.default.images.0.id
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_essd"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = var.name
	}

	resource "alicloud_ecs_command" "default" {
  		name             = var.name
  		command_content  = "ZWNobyBoZWxsbyx7e25hbWV9fQ=="
  		description      = "For Terraform Test"
  		type             = "RunShellScript"
  		working_dir      = "/root"
  		enable_parameter = true
	}
`, name)
}

func AliCloudECSInvocationBasicDependence2(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_essd"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.u1"
  		image_id             = data.alicloud_images.default.images.0.id
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_essd"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = var.name
  		tags = {
    			"tf-testacc"  = "ecs-invocation"
    			"tf-testacc2" = "ecs-invocation"
  		}
	}

	resource "alicloud_ecs_command" "default" {
  		name            = var.name
  		command_content = "bHMK"
  		description     = "For Terraform Test"
  		type            = "RunShellScript"
  		working_dir     = "/root"
	}
`, name)
}

// testAccAliCloudEcsInvocationBasic2Config builds the full HCL for basic2.
// The invocation uses resource_tag (no instance_id reference), so there is no
// implicit dependency edge on alicloud_instance.default; depends_on forces the
// tagged instance to reach Running before InvokeCommand queries by ResourceTag,
// avoiding InvalidResourceTag.InstanceNotFound.
func testAccAliCloudEcsInvocationBasic2Config(name string, tagKey string) string {
	return fmt.Sprintf(`%s

	resource "alicloud_ecs_invocation" "default" {
  command_id = alicloud_ecs_command.default.id
  resource_tag {
    key   = "%s"
    value = "ecs-invocation"
  }
  depends_on = [alicloud_instance.default]
}
`, AliCloudECSInvocationBasicDependence2(name), tagKey)
}

// AliCloudECSInvocationReorderDependence creates two distinct Running ECS
// instances so the TypeList instance_id attribute can exercise a convergent
// reorder: apply [first, second], then apply the same config with instance_id
// reordered to [second, first] and require the post-apply plan to be empty.
// instance_id is ForceNew, so the reorder triggers a replacement; the resource
// Read reorders the DescribeInvocations response to the prior-state order, so
// the post-replacement refresh converges instead of perpetually drifting.
func AliCloudECSInvocationReorderDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_essd"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.u1"
  		image_id             = data.alicloud_images.default.images.0.id
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "first" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_essd"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = "${var.name}-first"
	}

	resource "alicloud_instance" "second" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_essd"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = "${var.name}-second"
	}

	resource "alicloud_ecs_command" "default" {
  		name            = var.name
  		command_content = "bHMK"
  		description     = "For Terraform Test"
  		type            = "RunShellScript"
  		working_dir     = "/root"
	}
`, name)
}

func TestAccAliCloudECSInvocation_instanceIDReorder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_invocation.default"
	ra := resourceAttrInit(resourceId, AliCloudECSInvocationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsInvocation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsinvoreorder%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudECSInvocationReorderDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": []string{
						"${alicloud_instance.first.id}",
						"${alicloud_instance.second.id}",
					},
					"command_id": "${alicloud_ecs_command.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id.#": "2",
						"command_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": []string{
						"${alicloud_instance.second.id}",
						"${alicloud_instance.first.id}",
					},
					"command_id": "${alicloud_ecs_command.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id.#": "2",
						"command_id":    CHECKSET,
					}),
				),
			},
		},
	})
}

// AliCloudECSInvocationFieldDependence creates one Running ECS instance placed
// in a resource group plus an OSS bucket, so the create-only / response
// properties resource_group_id, working_dir and oss_output_delivery can be
// exercised with real values rather than placeholder state.
func AliCloudECSInvocationFieldDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_essd"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.u1"
  		image_id             = data.alicloud_images.default.images.0.id
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_essd"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = var.name
  		resource_group_id          = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	}

	resource "alicloud_ecs_command" "default" {
  		name            = var.name
  		command_content = "bHMK"
  		description     = "For Terraform Test"
  		type            = "RunShellScript"
  		working_dir     = "/root"
	}

	resource "alicloud_oss_bucket" "default" {
  		bucket = "${var.name}-output"
  		acl    = "private"
  		force_destroy = true
	}
`, name)
}

func TestAccAliCloudECSInvocation_fieldCoverage(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_invocation.default"
	ra := resourceAttrInit(resourceId, AliCloudECSInvocationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsInvocation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsinvofield%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudECSInvocationFieldDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":         []string{"${alicloud_instance.default.id}"},
					"command_id":          "${alicloud_ecs_command.default.id}",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"working_dir":         "/root",
					"oss_output_delivery": "oss://${alicloud_oss_bucket.default.bucket}/output",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id.#":       "1",
						"command_id":          CHECKSET,
						"resource_group_id":   CHECKSET,
						"working_dir":         "/root",
						"oss_output_delivery": CHECKSET,
					}),
				),
			},
		},
	})
}

// TestAccAliCloudECSInvocation_instanceTagSwitch exercises the bidirectional
// switch between the mutually-exclusive instance_id and resource_tag
// selectors. Each switch flips the Create path between InstanceId and
// ResourceTag; because both are ForceNew the switch is a replacement. The
// tagged instance is shared across steps, so only the invocation selector
// changes.
func TestAccAliCloudECSInvocation_instanceTagSwitch(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_invocation.default"
	ra := resourceAttrInit(resourceId, AliCloudECSInvocationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsInvocation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsinvoswitch%d", defaultRegionToTest, rand)
	instanceIDConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudECSInvocationBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: instanceIDConfig(map[string]interface{}{
					"instance_id": []string{"${alicloud_instance.default.id}"},
					"command_id":  "${alicloud_ecs_command.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id.#": "1",
						"command_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccAliCloudEcsInvocationBasic2Config(name, "tf-testacc"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command_id":           CHECKSET,
						"resource_tag.#":       "1",
						"resource_tag.0.key":   "tf-testacc",
						"resource_tag.0.value": "ecs-invocation",
					}),
				),
			},
			{
				Config: instanceIDConfig(map[string]interface{}{
					"instance_id": []string{"${alicloud_instance.default.id}"},
					"command_id":  "${alicloud_ecs_command.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id.#":        "1",
						"command_id":           CHECKSET,
						"resource_tag.#":       REMOVEKEY,
						"resource_tag.0.key":   REMOVEKEY,
						"resource_tag.0.value": REMOVEKEY,
					}),
				),
			},
		},
	})
}

// lintignore: R001
func TestUnitAliCloudEcsInvocation(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ecs_invocation"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"instance_id":           []string{"CreateECSInvocationValue"},
		"command_id":            "CreateECSInvocationValue",
		"repeat_mode":           "Once",
		"timed":                 true,
		"frequency":             "CreateECSInvocationValue",
		"username":              "CreateECSInvocationValue",
		"windows_password_name": "CreateECSInvocationValue",
	}
	if err := dInit.Set("instance_id", attributes["instance_id"]); err != nil {
		t.Fatalf("set instance_id: %v", err)
	}
	if err := dInit.Set("command_id", attributes["command_id"]); err != nil {
		t.Fatalf("set command_id: %v", err)
	}
	if err := dInit.Set("repeat_mode", attributes["repeat_mode"]); err != nil {
		t.Fatalf("set repeat_mode: %v", err)
	}
	if err := dInit.Set("timed", attributes["timed"]); err != nil {
		t.Fatalf("set timed: %v", err)
	}
	if err := dInit.Set("frequency", attributes["frequency"]); err != nil {
		t.Fatalf("set frequency: %v", err)
	}
	if err := dInit.Set("username", attributes["username"]); err != nil {
		t.Fatalf("set username: %v", err)
	}
	if err := dInit.Set("windows_password_name", attributes["windows_password_name"]); err != nil {
		t.Fatalf("set windows_password_name: %v", err)
	}

	// Canned DescribeInvocations response object. InvokeId matches the
	// InvokeId returned by the InvokeCommand handler so EcsService.DescribeEcsInvocation's
	// id check passes; Parameters is a JSON string because Read unmarshals
	// object["Parameters"].(string). InvocationStatus is a WaitForState target
	// so Create's state refresh resolves on the first poll.
	readObject := map[string]interface{}{
		"InvokeId":         "CreateECSInvocationValue",
		"CommandId":        "CreateECSInvocationValue",
		"RepeatMode":       "Once",
		"Frequency":        "CreateECSInvocationValue",
		"Parameters":       "{}",
		"Timed":            true,
		"Username":         "CreateECSInvocationValue",
		"InvocationStatus": "Success",
		"WorkingDir":       "/root",
		"InvokeInstances": map[string]interface{}{
			"InvokeInstance": []interface{}{
				map[string]interface{}{"InstanceId": "CreateECSInvocationValue"},
			},
		},
	}

	// ecsInvocationRpcHandler dispatches the three ECS RPC actions through the
	// real RpcPost + EcsService path against a credential-free httptest server,
	// replacing the former gomonkey patches on NewEcsClient/DoRequest that
	// rotted (panic "retrieve method by name failed") once the connectivity
	// client dropped the NewEcsClient accessor.
	ecsInvocationRpcHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "InvokeCommand":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"InvokeId": "CreateECSInvocationValue"})
		case "DescribeInvocations":
			_ = json.NewEncoder(w).Encode(invocationDescribeResponse(readObject))
		case "StopInvocation":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
	meta := ecsInvocationTestClient(t, ecsInvocationRpcHandler)

	// Create: InvokeCommand -> SetId -> WaitForState(DescribeInvocations) ->
	// Read(DescribeInvocations), all hitting the local server end-to-end.
	if err := resourceAlicloudEcsInvocationCreate(dInit, meta); err != nil {
		t.Fatalf("create: %v", err)
	}
	assert.Equal(t, "CreateECSInvocationValue", dInit.Id())

	// Read backfills the computed fields from the canned response.
	if err := resourceAlicloudEcsInvocationRead(dInit, meta); err != nil {
		t.Fatalf("read: %v", err)
	}
	assert.Equal(t, "CreateECSInvocationValue", dInit.Get("command_id"))
	assert.Equal(t, "Success", dInit.Get("status"))

	// Delete: StopInvocation succeeds; the resource is dropped from state.
	if err := resourceAlicloudEcsInvocationDelete(dInit, meta); err != nil {
		t.Fatalf("delete: %v", err)
	}

	// NonRetryableError: InvokeCommand returns a 400 with a non-throttling
	// code; NeedRetry is false, so Create surfaces the error instead of
	// retrying or reaching WaitForState.
	errMeta := ecsInvocationTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Code":      "InvalidParameter",
			"Message":   "CreateECSInvocationValue",
			"RequestId": "CreateECSInvocationValue",
		})
	})
	dErr, _ := schema.InternalMap(p["alicloud_ecs_invocation"].Schema).Data(nil, nil)
	dErr.MarkNewResource()
	_ = dErr.Set("instance_id", attributes["instance_id"])
	_ = dErr.Set("command_id", attributes["command_id"])
	_ = dErr.Set("repeat_mode", attributes["repeat_mode"])
	_ = dErr.Set("timed", attributes["timed"])
	_ = dErr.Set("frequency", attributes["frequency"])
	_ = dErr.Set("username", attributes["username"])
	_ = dErr.Set("windows_password_name", attributes["windows_password_name"])
	if err := resourceAlicloudEcsInvocationCreate(dErr, errMeta); err == nil {
		t.Fatalf("expected Create to fail on a non-retryable InvokeCommand error")
	}
}

// TestUnitEcsInvocationReadErrorPreservesID guards the Read error contract that
// the original (gomonkey) TestUnitAliCloudECSInvocation covered and the httptest
// rewrite dropped: a non-retryable DescribeInvocations error (HTTP 400 with a
// non-throttling code, NeedRetry=false) must surface from Read WITHOUT clearing
// the resource id. Only an authoritative NotFoundError removes the id; a generic
// service error must preserve it so the caller can retry or report. The handler
// verifies the real RpcPost round-trip carried Action=DescribeInvocations and the
// queried InvokeId equal to the resource id.
func TestUnitEcsInvocationReadErrorPreservesID(t *testing.T) {
	p := Provider().ResourcesMap
	dExisted, _ := schema.InternalMap(p["alicloud_ecs_invocation"].Schema).Data(nil, nil)
	dExisted.SetId("invoke-read-err")
	if err := dExisted.Set("command_id", "cmd-test"); err != nil {
		t.Fatalf("set command_id: %v", err)
	}
	if err := dExisted.Set("instance_id", []interface{}{"i-test"}); err != nil {
		t.Fatalf("set instance_id: %v", err)
	}

	var seenAction, seenInvokeId string
	handler := func(w http.ResponseWriter, r *http.Request) {
		seenAction = r.Form.Get("Action")
		seenInvokeId = r.Form.Get("InvokeId")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Code":      "InvalidParameter",
			"Message":   "boom",
			"RequestId": "r-read-err",
		})
	}
	meta := ecsInvocationTestClient(t, handler)

	err := resourceAlicloudEcsInvocationRead(dExisted, meta)
	assert.NotNil(t, err, "Read must surface a non-retryable DescribeInvocations error")
	// Only NotFoundError clears the id; a generic service error must preserve it.
	assert.Equal(t, "invoke-read-err", dExisted.Id())
	// Verify the real request targeted DescribeInvocations with the resource id.
	assert.Equal(t, "DescribeInvocations", seenAction)
	assert.Equal(t, "invoke-read-err", seenInvokeId)
}

// TestUnitEcsInvocationReadNotFoundRemoves guards the authoritative NotFound
// removal path dropped by the httptest rewrite. DescribeInvocations returns an
// empty Invocation list (the canonical "resource gone" response);
// EcsService.DescribeEcsInvocation treats it as NotFoundError, so Read clears
// the id and returns nil, letting Terraform drop the resource from state. Verifies
// Action=DescribeInvocations and the queried InvokeId equals the resource id.
func TestUnitEcsInvocationReadNotFoundRemoves(t *testing.T) {
	p := Provider().ResourcesMap
	dExisted, _ := schema.InternalMap(p["alicloud_ecs_invocation"].Schema).Data(nil, nil)
	dExisted.SetId("invoke-gone")
	if err := dExisted.Set("command_id", "cmd-test"); err != nil {
		t.Fatalf("set command_id: %v", err)
	}
	if err := dExisted.Set("instance_id", []interface{}{"i-test"}); err != nil {
		t.Fatalf("set instance_id: %v", err)
	}

	var seenAction, seenInvokeId string
	handler := func(w http.ResponseWriter, r *http.Request) {
		seenAction = r.Form.Get("Action")
		seenInvokeId = r.Form.Get("InvokeId")
		w.Header().Set("Content-Type", "application/json")
		// Empty Invocation list -> service treats the response as NotFound.
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Invocations": map[string]interface{}{
				"Invocation": []interface{}{},
			},
		})
	}
	meta := ecsInvocationTestClient(t, handler)

	err := resourceAlicloudEcsInvocationRead(dExisted, meta)
	assert.Nil(t, err, "Read must swallow an authoritative NotFound and return nil")
	assert.Equal(t, "", dExisted.Id(), "NotFound must clear the resource id so Terraform removes it from state")
	assert.Equal(t, "DescribeInvocations", seenAction)
	assert.Equal(t, "invoke-gone", seenInvokeId)
}

// TestUnitEcsInvocationDeleteErrorNotCleaned guards the Delete hard-error path
// dropped by the httptest rewrite. A non-retryable StopInvocation error (HTTP 400,
// non-throttling code) must surface from Delete: the success path logs a WARN and
// returns nil, but a hard service error must return non-nil so the resource is NOT
// misreported as destroyed. Verifies Action=StopInvocation and InvokeId equal to
// the resource id.
func TestUnitEcsInvocationDeleteErrorNotCleaned(t *testing.T) {
	p := Provider().ResourcesMap
	dExisted, _ := schema.InternalMap(p["alicloud_ecs_invocation"].Schema).Data(nil, nil)
	dExisted.SetId("invoke-del-err")
	if err := dExisted.Set("command_id", "cmd-test"); err != nil {
		t.Fatalf("set command_id: %v", err)
	}
	if err := dExisted.Set("instance_id", []interface{}{"i-test"}); err != nil {
		t.Fatalf("set instance_id: %v", err)
	}

	var seenAction, seenInvokeId string
	handler := func(w http.ResponseWriter, r *http.Request) {
		seenAction = r.Form.Get("Action")
		seenInvokeId = r.Form.Get("InvokeId")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Code":      "InvalidParameter",
			"Message":   "cannot stop",
			"RequestId": "r-del-err",
		})
	}
	meta := ecsInvocationTestClient(t, handler)

	err := resourceAlicloudEcsInvocationDelete(dExisted, meta)
	assert.NotNil(t, err, "Delete must surface a non-retryable StopInvocation error, not treat it as cleaned")
	assert.Equal(t, "invoke-del-err", dExisted.Id(), "on Delete error the id is preserved so the resource is retried, not dropped")
	assert.Equal(t, "StopInvocation", seenAction)
	assert.Equal(t, "invoke-del-err", seenInvokeId)
}

// TestUnitEcsInvocationCreateThrottlingRetry guards the bounded retry contract
// dropped by the httptest rewrite: InvokeCommand returns a transient throttling
// error (HTTP 400, Code "Throttling", NeedRetry=true) on the first call and
// succeeds on the second, so Create retries once and reaches SetId + WaitForState
// + Read. It verifies TWO real InvokeCommand round-trips (transient error then
// success) carrying the configured CommandId, and that the resource id is set
// from the InvokeCommand response after the retry. The single incrementalWait
// backoff is the production retry semantics retained from the original gomonkey
// test; no long-wait framework is introduced.
func TestUnitEcsInvocationCreateThrottlingRetry(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ecs_invocation"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"instance_id":           []string{"i-throttle"},
		"command_id":            "cmd-throttle",
		"repeat_mode":           "Once",
		"timed":                 true,
		"frequency":             "CreateECSInvocationValue",
		"username":              "CreateECSInvocationValue",
		"windows_password_name": "CreateECSInvocationValue",
	}
	if err := dInit.Set("instance_id", attributes["instance_id"]); err != nil {
		t.Fatalf("set instance_id: %v", err)
	}
	if err := dInit.Set("command_id", attributes["command_id"]); err != nil {
		t.Fatalf("set command_id: %v", err)
	}
	if err := dInit.Set("repeat_mode", attributes["repeat_mode"]); err != nil {
		t.Fatalf("set repeat_mode: %v", err)
	}
	if err := dInit.Set("timed", attributes["timed"]); err != nil {
		t.Fatalf("set timed: %v", err)
	}
	if err := dInit.Set("frequency", attributes["frequency"]); err != nil {
		t.Fatalf("set frequency: %v", err)
	}
	if err := dInit.Set("username", attributes["username"]); err != nil {
		t.Fatalf("set username: %v", err)
	}
	if err := dInit.Set("windows_password_name", attributes["windows_password_name"]); err != nil {
		t.Fatalf("set windows_password_name: %v", err)
	}

	readObject := map[string]interface{}{
		"InvokeId":         "invoke-throttle",
		"CommandId":        "cmd-throttle",
		"RepeatMode":       "Once",
		"Frequency":        "CreateECSInvocationValue",
		"Parameters":       "{}",
		"Timed":            true,
		"Username":         "CreateECSInvocationValue",
		"InvocationStatus": "Success",
		"WorkingDir":       "/root",
		"InvokeInstances": map[string]interface{}{
			"InvokeInstance": []interface{}{
				map[string]interface{}{"InstanceId": "i-throttle"},
			},
		},
	}

	var invokeCount int
	var lastCommandId string
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "InvokeCommand":
			invokeCount++
			lastCommandId = r.Form.Get("CommandId")
			if invokeCount == 1 {
				// transient throttling -> NeedRetry true -> one bounded retry
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"Code":      "Throttling",
					"Message":   "too many requests",
					"RequestId": "r-throttle",
				})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"InvokeId": "invoke-throttle"})
		case "DescribeInvocations":
			_ = json.NewEncoder(w).Encode(invocationDescribeResponse(readObject))
		case "StopInvocation":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
	meta := ecsInvocationTestClient(t, handler)

	if err := resourceAlicloudEcsInvocationCreate(dInit, meta); err != nil {
		t.Fatalf("Create must succeed after a bounded throttling retry: %v", err)
	}
	assert.Equal(t, 2, invokeCount, "InvokeCommand must be called twice: throttling then success")
	assert.Equal(t, "cmd-throttle", lastCommandId, "the InvokeCommand request must carry the configured CommandId")
	assert.Equal(t, "invoke-throttle", dInit.Id(), "the resource id must be set from the InvokeCommand response after retry")
}

// TestUnitEcsInvocationCreateOnlyNoDrift is a regression guard for the
// resource_group_id / working_dir / oss_output_delivery schema design. It uses
// a valid config (instance_id selected) so it exercises the real SDK diff path
// rather than bypassing AtLeastOneOf validation.
//
// Background: DescribeInvocations returns WorkingDir and OssOutputDelivery as
// response properties, but ResourceGroupId is only a query filter and is NOT
// part of the Invocation response object (ContentEncoding is a query control
// parameter, not a create parameter nor a response property, and is not modeled
// on this resource). An earlier version modeled resource_group_id and
// content_encoding as Optional+Computed+ForceNew and unconditionally called
// d.Set(..., object[...]) in Read. Because the response omits those keys,
// d.Set received nil and, for a Computed ForceNew field, that produced
// RequiresNew on every refresh — `terraform plan` after `terraform refresh`
// always wanted to replace the resource.
//
// The fix: resource_group_id is no longer Computed (create-only, preserved
// from config) and Read only sets it when the response actually carries the
// key; working_dir/oss_output_delivery are set defensively too. This test
// asserts a refresh against an unchanged config does not plan replacement,
// while changing the create-only resource_group_id still does.
func TestUnitEcsInvocationCreateOnlyNoDrift(t *testing.T) {
	r := resourceAlicloudEcsInvocation()

	config := map[string]interface{}{
		"command_id":          "cmd-test",
		"instance_id":         []interface{}{"i-test"},
		"resource_group_id":   "rg-custom",
		"working_dir":         "/root",
		"oss_output_delivery": "acs:oss://bucket/prefix",
	}
	// Build the post-Read state. The fixed Read preserves resource_group_id
	// (create-only, absent from the DescribeInvocations response) instead of
	// clearing it to nil, and fills the Computed fields from the response. The
	// Computed fields below mirror a realistic refresh so that omitting them
	// from config does not by itself trigger replacement.
	data := schema.TestResourceDataRaw(t, r.Schema, config)
	data.SetId("invoke-test")
	if err := data.Set("instance_id", []interface{}{"i-test"}); err != nil {
		t.Fatalf("set instance_id: %v", err)
	}
	if err := data.Set("repeat_mode", "Once"); err != nil {
		t.Fatalf("set repeat_mode: %v", err)
	}
	if err := data.Set("timed", false); err != nil {
		t.Fatalf("set timed: %v", err)
	}
	if err := data.Set("username", "root"); err != nil {
		t.Fatalf("set username: %v", err)
	}
	if err := data.Set("status", "Success"); err != nil {
		t.Fatalf("set status: %v", err)
	}
	state := data.State()

	// Refresh against the same config must be a complete no-op: not merely
	// lacking RequiresNew, but an empty diff so create-only fields preserved
	// across Read do not surface as spurious changes.
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if diff != nil && len(diff.Attributes) > 0 {
		t.Fatalf("expected empty diff on refresh, got: %#v", diff.Attributes)
	}

	// Changing the create-only resource_group_id must still plan replacement,
	// proving the field remains ForceNew and was not silently dropped.
	changed := map[string]interface{}{
		"command_id":          "cmd-test",
		"instance_id":         []interface{}{"i-test"},
		"resource_group_id":   "rg-other",
		"working_dir":         "/root",
		"oss_output_delivery": "acs:oss://bucket/prefix",
	}
	diff, err = r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(changed), nil)
	if err != nil {
		t.Fatalf("changed diff: %v", err)
	}
	if diff == nil || !diff.RequiresNew() {
		t.Fatalf("expected replacement when resource_group_id changes, got diff=%#v", diff)
	}
}

// TestUnitEcsInvocationReadPreservesCreateOnly exercises the real Read path
// with a DescribeInvocations response that omits ResourceGroupId (a query
// filter, not a response property) but carries WorkingDir/OssOutputDelivery.
// ecsInvocationTestClient stands up a credential-free AliyunClient whose ECS
// transport is pointed at an in-process httptest server, mirroring the
// ecs-snapshot unit-test construction. The DescribeInvocations handler is
// supplied by the caller, so the real resourceAlicloudEcsInvocationRead and
// EcsService.DescribeEcsInvocation run end-to-end (RpcPost -> local server)
// without cloud credentials, gomonkey binary patching, or an empty AliyunClient
// that SIGSEGVs on an un-patched code path.
func ecsInvocationTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("ecs", endpoint)
	return client
}

// invocationDescribeResponse wraps a canned Invocation object in the
// DescribeInvocations RPC envelope ({"Invocations":{"Invocation":[object]}})
// expected by EcsService.DescribeEcsInvocation's jsonpath lookup
// ("$.Invocations.Invocation"). InvokeId must match the resource id or the
// service treats the response as NotFound.
func invocationDescribeResponse(object map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"Invocations": map[string]interface{}{
			"Invocation": []interface{}{object},
		},
	}
}

// ecsInvocationDescribeHandler returns an httptest handler that serves the
// DescribeInvocations action from a caller-held response variable. The
// variable is captured by reference so the caller can swap the canned
// Invocation between reads (e.g. to return InvokeInstances in a different
// order); each DescribeInvocations request serializes with the test goroutine
// through the synchronous RpcPost round-trip, so the read/write pair is
// happens-before ordered and race-free.
func ecsInvocationDescribeHandler(t *testing.T, response *map[string]interface{}) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("Action") != "DescribeInvocations" {
			t.Errorf("unexpected ECS action in Read: %s", r.Form.Get("Action"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.NewEncoder(w).Encode(invocationDescribeResponse(*response)); err != nil {
			t.Error(err)
		}
	}
}

// It runs without cloud credentials (a credential-free httptest server serves
// DescribeInvocations) and verifies the fixed Read preserves the configured
// resource_group_id instead of clearing it (which previously caused RequiresNew
// drift), sets the returned output fields, and leaves a subsequent no-change
// plan stable.
func TestUnitEcsInvocationReadPreservesCreateOnly(t *testing.T) {
	p := Provider().ResourcesMap
	dExisted, _ := schema.InternalMap(p["alicloud_ecs_invocation"].Schema).Data(nil, nil)
	dExisted.SetId("invoke-test")
	if err := dExisted.Set("command_id", "cmd-test"); err != nil {
		t.Fatalf("set command_id: %v", err)
	}
	if err := dExisted.Set("instance_id", []interface{}{"i-test"}); err != nil {
		t.Fatalf("set instance_id: %v", err)
	}
	if err := dExisted.Set("resource_group_id", "rg-custom"); err != nil {
		t.Fatalf("set resource_group_id: %v", err)
	}
	if err := dExisted.Set("working_dir", "/root"); err != nil {
		t.Fatalf("set working_dir: %v", err)
	}
	if err := dExisted.Set("oss_output_delivery", "acs:oss://bucket/prefix"); err != nil {
		t.Fatalf("set oss_output_delivery: %v", err)
	}

	// DescribeInvocations response: WorkingDir/OssOutputDelivery present,
	// ResourceGroupId intentionally absent (query filter, not a response field).
	// InvokeInstances carries the same InstanceId as the config so the Computed
	// instance_id backfills to the configured value instead of clearing it (which
	// would otherwise force replacement on the very refresh this test guards).
	readObject := map[string]interface{}{
		"InvokeId":          "invoke-test",
		"CommandId":         "cmd-test",
		"RepeatMode":        "Once",
		"Frequency":         "",
		"Parameters":        "{}",
		"Timed":             false,
		"Username":          "root",
		"InvocationStatus":  "Success",
		"WorkingDir":        "/root",
		"OssOutputDelivery": "acs:oss://bucket/prefix",
		"InvokeInstances": map[string]interface{}{
			"InvokeInstance": []interface{}{
				map[string]interface{}{"InstanceId": "i-test"},
			},
		},
	}
	meta := ecsInvocationTestClient(t, ecsInvocationDescribeHandler(t, &readObject))

	// Read runs end-to-end through the real EcsService.DescribeEcsInvocation and
	// RpcPost against the local server; the create-only resource_group_id is
	// absent from the response, so the fixed Read must preserve the configured
	// value instead of clearing it to nil.
	if err := resourceAlicloudEcsInvocationRead(dExisted, meta); err != nil {
		t.Fatalf("read: %v", err)
	}
	// resource_group_id is create-only and absent from the response; the fixed
	// Read must preserve the configured value, not clear it to nil.
	assert.Equal(t, "rg-custom", dExisted.Get("resource_group_id"))
	assert.Equal(t, "/root", dExisted.Get("working_dir"))
	assert.Equal(t, "acs:oss://bucket/prefix", dExisted.Get("oss_output_delivery"))

	// A subsequent no-change plan must not replace the resource.
	r := resourceAlicloudEcsInvocation()
	config := map[string]interface{}{
		"command_id":          "cmd-test",
		"instance_id":         []interface{}{"i-test"},
		"resource_group_id":   "rg-custom",
		"working_dir":         "/root",
		"oss_output_delivery": "acs:oss://bucket/prefix",
	}
	diff, err := r.Diff(context.Background(), dExisted.State(), terraform.NewResourceConfigRaw(config), nil)
	assert.NoError(t, err)
	if diff != nil && len(diff.Attributes) > 0 {
		t.Fatalf("expected empty diff after preserving create-only fields, got: %#v", diff.Attributes)
	}
}

// TestUnitEcsInvocationImportCompat exercises the real import path and the
// instance_id <-> resource_tag selection semantics, all without cloud
// credentials (DescribeInvocations is served by a credential-free httptest
// server, so no SKIP):
//
//   - Import: the resource's real Importer (ImportStatePassthrough wired via
//     State) is invoked through r.Importer.State from an empty state carrying
//     only the ID, then Read backfills the recoverable fields (instance_id,
//     working_dir, oss_output_delivery) from the response. A subsequent plan
//     against a matching legacy instance_id config must be a complete no-op
//     (empty diff, not merely !RequiresNew), proving import does not force
//     replacement and that no GetOk(instance_id) guard suppressed backfill.
//     resource_group_id is a create-only query filter absent from the response;
//     it stays empty and must not drift into replacement.
//   - Two consecutive reads with the InvokeInstances response order swapped
//     between reads produce no drift: Read reorders the response to follow the
//     prior state so an unchanged membership stays stable across refreshes.
//   - Switching a legacy instance_id config to resource_tag (and back) must
//     force replacement: they are mutually exclusive ForceNew selectors.
//     resource_tag is write-only and not backfilled by Read, so consecutive
//     reads against the same resource_tag config do not drift.
func TestUnitEcsInvocationImportCompat(t *testing.T) {
	r := resourceAlicloudEcsInvocation()

	readObject := func(instanceIDs ...string) map[string]interface{} {
		obj := map[string]interface{}{
			"InvokeId":          "invoke-test",
			"CommandId":         "cmd-test",
			"RepeatMode":        "Once",
			"Frequency":         "",
			"Parameters":        "{}",
			"Timed":             false,
			"Username":          "root",
			"InvocationStatus":  "Success",
			"WorkingDir":        "/root",
			"OssOutputDelivery": "acs:oss://bucket/prefix",
			// ResourceGroupId intentionally absent: create-only query filter, not
			// part of the Invocation response object.
		}
		if len(instanceIDs) > 0 {
			instances := make([]interface{}, 0, len(instanceIDs))
			for _, id := range instanceIDs {
				instances = append(instances, map[string]interface{}{"InstanceId": id})
			}
			obj["InvokeInstances"] = map[string]interface{}{
				"InvokeInstance": instances,
			}
		}
		return obj
	}

	// Mutable response so consecutive reads can swap the InvokeInstances order.
	response := readObject("i-a", "i-b")
	// A credential-free httptest server serves DescribeInvocations from the
	// mutable response; the real Read/Importer run end-to-end against it.
	meta := ecsInvocationTestClient(t, ecsInvocationDescribeHandler(t, &response))

	// --- import: empty state + only ID, invoke the real Importer, then Read ---
	dSeed := r.Data(nil)
	dSeed.SetId("invoke-test")
	imported, err := r.Importer.State(dSeed, meta)
	if err != nil || len(imported) != 1 {
		t.Fatalf("import via real Importer (r.Importer.State): %#v, %v", imported, err)
	}
	dImported := imported[0]
	if err := resourceAlicloudEcsInvocationRead(dImported, meta); err != nil {
		t.Fatalf("import read: %v", err)
	}
	// instance_id was backfilled from the response (no GetOk guard that would
	// leave it empty and break legacy imports); output fields carry the
	// non-empty response values; resource_group_id stays empty because the
	// response omits the create-only query filter and must not drift.
	assert.Equal(t, []interface{}{"i-a", "i-b"}, dImported.Get("instance_id"))
	assert.Equal(t, "/root", dImported.Get("working_dir"))
	assert.Equal(t, "acs:oss://bucket/prefix", dImported.Get("oss_output_delivery"))
	assert.Empty(t, dImported.Get("resource_group_id"))

	// A matching legacy instance_id config must be a complete no-op plan.
	instanceCfg := map[string]interface{}{
		"command_id":          "cmd-test",
		"instance_id":         []interface{}{"i-a", "i-b"},
		"oss_output_delivery": "acs:oss://bucket/prefix",
	}
	assertEmptyDiff := func(t *testing.T, state *terraform.InstanceState, cfg map[string]interface{}, msg string) {
		t.Helper()
		diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(cfg), nil)
		assert.NoError(t, err)
		if diff != nil && len(diff.Attributes) > 0 {
			t.Fatalf("%s: expected empty diff, got: %#v", msg, diff.Attributes)
		}
	}
	assertEmptyDiff(t, dImported.State(), instanceCfg, "legacy instance_id import")

	// --- two consecutive reads with the response order swapped: no drift ---
	// First Read used response [i-a, i-b] with an empty reference, so state is
	// [i-a, i-b]. Re-read with the response returned in reverse order [i-b, i-a];
	// the fixed Read reorders it to follow the prior state, so the membership
	// and ordering stay stable and the plan remains a no-op.
	response = readObject("i-b", "i-a")
	if err := resourceAlicloudEcsInvocationRead(dImported, meta); err != nil {
		t.Fatalf("reordered read: %v", err)
	}
	assert.Equal(t, []interface{}{"i-a", "i-b"}, dImported.Get("instance_id"))
	assertEmptyDiff(t, dImported.State(), instanceCfg, "reordered consecutive read")
	// A third read restoring the original order must also stay a no-op.
	response = readObject("i-a", "i-b")
	if err := resourceAlicloudEcsInvocationRead(dImported, meta); err != nil {
		t.Fatalf("third read: %v", err)
	}
	assert.Equal(t, []interface{}{"i-a", "i-b"}, dImported.Get("instance_id"))
	assertEmptyDiff(t, dImported.State(), instanceCfg, "third consecutive read")

	// --- instance_id -> resource_tag switch must force replacement ---
	tagCfg := map[string]interface{}{
		"command_id":          "cmd-test",
		"resource_tag":        []interface{}{map[string]interface{}{"key": "k", "value": "v"}},
		"oss_output_delivery": "acs:oss://bucket/prefix",
	}
	diff, err := r.Diff(context.Background(), dImported.State(), terraform.NewResourceConfigRaw(tagCfg), nil)
	assert.NoError(t, err)
	if diff == nil || !diff.RequiresNew() {
		t.Fatalf("instance_id->resource_tag switch: expected replacement, got diff=%#v", diff)
	}

	// --- resource_tag explicit mode: consecutive reads do not drift ---
	dTag := r.Data(nil)
	dTag.SetId("invoke-test")
	if err := dTag.Set("command_id", "cmd-test"); err != nil {
		t.Fatalf("set command_id: %v", err)
	}
	if err := dTag.Set("resource_tag", []interface{}{map[string]interface{}{"key": "k", "value": "v"}}); err != nil {
		t.Fatalf("set resource_tag: %v", err)
	}
	response = readObject("i-tagged")
	if err := resourceAlicloudEcsInvocationRead(dTag, meta); err != nil {
		t.Fatalf("tag read: %v", err)
	}
	// resource_tag is write-only and preserved (Read does not backfill it);
	// instance_id is backfilled from the response without drifting resource_tag.
	// resource_tag is a TypeSet, so dTag.Get returns a *schema.Set; compare with
	// set semantics (ElementsMatch on Set.List) instead of a raw slice, which
	// would always mismatch the *schema.Set type.
	tagSet, ok := dTag.Get("resource_tag").(*schema.Set)
	if !ok {
		t.Fatalf("resource_tag should be a *schema.Set, got %T", dTag.Get("resource_tag"))
	}
	assert.ElementsMatch(t, []interface{}{map[string]interface{}{"key": "k", "value": "v"}}, tagSet.List())
	assert.Equal(t, []interface{}{"i-tagged"}, dTag.Get("instance_id"))
	assertEmptyDiff(t, dTag.State(), tagCfg, "resource_tag consecutive read")
	// Switching the resource_tag selector back to instance_id must replace.
	diff, err = r.Diff(context.Background(), dTag.State(), terraform.NewResourceConfigRaw(instanceCfg), nil)
	assert.NoError(t, err)
	if diff == nil || !diff.RequiresNew() {
		t.Fatalf("resource_tag->instance_id switch: expected replacement, got diff=%#v", diff)
	}
}

// TestUnitEcsInvocationInstanceIDOrderDrift guards the instance_id ordering
// absorption in Read without cloud credentials (DescribeInvocations is served
// by a credential-free httptest server, so no SKIP):
//
//   - A refresh whose response returns the same membership in a different
//     order is a no-op: Read reorders the response to follow the prior state.
//   - An explicit change to the configured instance_id List order still forces
//     replacement, because the reorder reference is the prior state, not the
//     new config (TypeList is order-sensitive).
//   - A real membership change (add/remove) forces replacement; the reorder
//     only absorbs ordering drift, not content drift.
//   - The first read of an ID-only import (no prior ordering) backfills from
//     the response order without a GetOk guard that would suppress backfill.
func TestUnitEcsInvocationInstanceIDOrderDrift(t *testing.T) {
	r := resourceAlicloudEcsInvocation()

	readObject := func(instanceIDs ...string) map[string]interface{} {
		obj := map[string]interface{}{
			"InvokeId":          "invoke-test",
			"CommandId":         "cmd-test",
			"RepeatMode":        "Once",
			"Frequency":         "",
			"Parameters":        "{}",
			"Timed":             false,
			"Username":          "root",
			"InvocationStatus":  "Success",
			"WorkingDir":        "/root",
			"OssOutputDelivery": "acs:oss://bucket/prefix",
		}
		if len(instanceIDs) > 0 {
			instances := make([]interface{}, 0, len(instanceIDs))
			for _, id := range instanceIDs {
				instances = append(instances, map[string]interface{}{"InstanceId": id})
			}
			obj["InvokeInstances"] = map[string]interface{}{
				"InvokeInstance": instances,
			}
		}
		return obj
	}
	cfg := func(ids ...string) map[string]interface{} {
		list := make([]interface{}, 0, len(ids))
		for _, id := range ids {
			list = append(list, id)
		}
		return map[string]interface{}{
			"command_id":          "cmd-test",
			"instance_id":         list,
			"oss_output_delivery": "acs:oss://bucket/prefix",
		}
	}

	response := readObject("i-one", "i-two")
	// A credential-free httptest server serves DescribeInvocations from the
	// mutable response; the real Read runs end-to-end against it.
	meta := ecsInvocationTestClient(t, ecsInvocationDescribeHandler(t, &response))

	assertEmptyDiff := func(state *terraform.InstanceState, c map[string]interface{}, msg string) {
		t.Helper()
		diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(c), nil)
		assert.NoError(t, err)
		if diff != nil && len(diff.Attributes) > 0 {
			t.Fatalf("%s: expected empty diff, got: %#v", msg, diff.Attributes)
		}
	}
	assertRequiresNew := func(state *terraform.InstanceState, c map[string]interface{}, msg string) {
		t.Helper()
		diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(c), nil)
		assert.NoError(t, err)
		if diff == nil || !diff.RequiresNew() {
			t.Fatalf("%s: expected replacement, got diff=%#v", msg, diff)
		}
	}

	// Establish prior state: explicit config + first Read. The reference is
	// empty on the first read, so instance_id backfills in response order.
	d := r.Data(nil)
	d.SetId("invoke-test")
	if err := d.Set("command_id", "cmd-test"); err != nil {
		t.Fatalf("set command_id: %v", err)
	}
	if err := d.Set("instance_id", []interface{}{"i-one", "i-two"}); err != nil {
		t.Fatalf("set instance_id: %v", err)
	}
	if err := d.Set("oss_output_delivery", "acs:oss://bucket/prefix"); err != nil {
		t.Fatalf("set oss_output_delivery: %v", err)
	}
	if err := resourceAlicloudEcsInvocationRead(d, meta); err != nil {
		t.Fatalf("first read: %v", err)
	}
	assert.Equal(t, []interface{}{"i-one", "i-two"}, d.Get("instance_id"))

	// (a) same membership, response order swapped -> no-op (absorbed).
	response = readObject("i-two", "i-one")
	if err := resourceAlicloudEcsInvocationRead(d, meta); err != nil {
		t.Fatalf("swapped read: %v", err)
	}
	assert.Equal(t, []interface{}{"i-one", "i-two"}, d.Get("instance_id"))
	assertEmptyDiff(d.State(), cfg("i-one", "i-two"), "same membership swapped order")

	// (b) explicit change to the configured List order -> ForceNew: the reorder
	// reference is the prior state [i-one, i-two], not the new config
	// [i-two, i-one], so the TypeList order mismatch is preserved.
	assertRequiresNew(d.State(), cfg("i-two", "i-one"), "explicit config order change")

	// (c) real membership change -> ForceNew (i-two removed, i-three added).
	response = readObject("i-one", "i-three")
	if err := resourceAlicloudEcsInvocationRead(d, meta); err != nil {
		t.Fatalf("membership-change read: %v", err)
	}
	assert.Equal(t, []interface{}{"i-one", "i-three"}, d.Get("instance_id"))
	assertRequiresNew(d.State(), cfg("i-one", "i-two"), "real membership change vs original config")

	// (d) ID-only import: empty prior ordering -> backfill from response order,
	// no GetOk guard suppressing backfill. The matching config is a no-op.
	dImport := r.Data(nil)
	dImport.SetId("invoke-test")
	response = readObject("i-x", "i-y")
	if err := resourceAlicloudEcsInvocationRead(dImport, meta); err != nil {
		t.Fatalf("import read: %v", err)
	}
	assert.Equal(t, []interface{}{"i-x", "i-y"}, dImport.Get("instance_id"))
	assertEmptyDiff(dImport.State(), cfg("i-x", "i-y"), "ID-only import backfill")
}
