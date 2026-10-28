package nlb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nlb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudNlbServerGroupServerAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group_server_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupServerAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroupServerAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%snlbservergroupserverattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupServerAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"server_type":     "Ecs",
					"server_id":       "${alicloud_instance.default.id}",
					"server_group_id": "${alicloud_nlb_server_group.default.id}",
					"port":            "80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_type":     "Ecs",
						"server_id":       acctest.CHECKSET,
						"server_group_id": acctest.CHECKSET,
						"port":            "80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"weight": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"weight": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"weight": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight": "0",
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

func TestAccAliCloudNlbServerGroupServerAttachment_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group_server_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupServerAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroupServerAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%snlbservergroupserverattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupServerAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"server_type":     "Ecs",
					"server_id":       "${alicloud_instance.default.id}",
					"description":     name,
					"port":            "80",
					"server_group_id": "${alicloud_nlb_server_group.default.id}",
					"weight":          "0",
					"server_ip":       "${alicloud_instance.default.private_ip}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_type":     "Ecs",
						"server_id":       acctest.CHECKSET,
						"description":     name,
						"port":            "80",
						"server_group_id": acctest.CHECKSET,
						"weight":          "0",
						"server_ip":       acctest.CHECKSET,
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

var AliCloudNlbServerGroupServerAttachmentMap0 = map[string]string{
	"server_ip": acctest.CHECKSET,
	"status":    acctest.CHECKSET,
	"weight":    acctest.CHECKSET,
	"zone_id":   acctest.CHECKSET,
}

func AliCloudNlbServerGroupServerAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones[0].id
}


resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  instance_name              = var.name
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = data.alicloud_vswitches.default.ids[0]
}

resource "alicloud_nlb_server_group" "default" {
  server_group_name = var.name
  server_group_type = "Instance"
  vpc_id            = data.alicloud_vpcs.default.ids.0
  scheduler         = "Wrr"
  protocol          = "TCP"
  health_check {
    health_check_enabled         = false
  }
  connection_drain           = true
  connection_drain_timeout   = 60
  preserve_client_ip_enabled = true
  address_ip_version = "Ipv4"
}
`, name)
}

func TestAccAliCloudNlbServerGroupServerAttachment_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group_server_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupServerAttachmentMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroupServerAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%snlbservergroupserverattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupServerAttachmentBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"server_type":     "Ip",
					"server_id":       "10.0.0.0",
					"port":            "80",
					"server_group_id": "${alicloud_nlb_server_group.default.id}",
					"server_ip":       "10.0.0.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_type":     "Ip",
						"server_id":       "10.0.0.0",
						"port":            "80",
						"server_group_id": acctest.CHECKSET,
						"server_ip":       "10.0.0.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"weight": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"weight": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"weight": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight": "0",
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

func TestAccAliCloudNlbServerGroupServerAttachment_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group_server_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupServerAttachmentMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroupServerAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%snlbservergroupserverattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupServerAttachmentBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"server_type":     "Ip",
					"server_id":       "10.0.0.0",
					"port":            "80",
					"server_group_id": "${alicloud_nlb_server_group.default.id}",
					"server_ip":       "10.0.0.0",
					"description":     name,
					"weight":          "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_type":     "Ip",
						"server_id":       "10.0.0.0",
						"port":            "80",
						"server_group_id": acctest.CHECKSET,
						"server_ip":       "10.0.0.0",
						"description":     name,
						"weight":          "0",
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

var AliCloudNlbServerGroupServerAttachmentMap1 = map[string]string{
	"status": acctest.CHECKSET,
	"weight": acctest.CHECKSET,
}

func AliCloudNlbServerGroupServerAttachmentBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_nlb_server_group" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  server_group_name = var.name
  server_group_type = "Ip"
  vpc_id            = data.alicloud_vpcs.default.ids.0
  scheduler         = "Wrr"
  protocol          = "TCP"
  health_check {
    health_check_enabled = false
  }
  address_ip_version = "Ipv4"
}
`, name)
}

// lintignore: R001
