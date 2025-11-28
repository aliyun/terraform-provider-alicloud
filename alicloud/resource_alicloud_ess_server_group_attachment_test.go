package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudEssServerGroupAttachment_basic_alb(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_server_group_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroupServerGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingGroupServerGroup)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssServerGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":       []string{"alicloud_ess_scaling_configuration.default"},
					"force_attach":     true,
					"weight":           "100",
					"port":             "80",
					"type":             "ALB",
					"server_group_id":  "${alicloud_alb_server_group.default.id}",
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":               CHECKSET,
						"weight":           "100",
						"port":             "80",
						"type":             "ALB",
						"server_group_id":  CHECKSET,
						"scaling_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_attach"},
			},
		},
	})
}

func TestAccAliCloudEssServerGroupAttachment_nonForceAttach_alb(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_server_group_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)

	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroupServerGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingGroupServerGroupNotForceAttach)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssServerGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":       []string{"alicloud_ess_scaling_configuration.default"},
					"force_attach":     false,
					"weight":           "11",
					"port":             "22",
					"type":             "ALB",
					"server_group_id":  "${alicloud_alb_server_group.default.id}",
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":               CHECKSET,
						"weight":           "11",
						"port":             "22",
						"type":             "ALB",
						"server_group_id":  CHECKSET,
						"scaling_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_attach"},
			},
		},
	})
}

func TestAccAliCloudEssServerGroupAttachment_basic_nlb(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_server_group_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)

	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroupServerGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingGroupServerGroup)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssServerGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":       []string{"alicloud_ess_scaling_configuration.default"},
					"force_attach":     true,
					"weight":           "100",
					"port":             "80",
					"type":             "NLB",
					"server_group_id":  "${alicloud_nlb_server_group.default.id}",
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":               CHECKSET,
						"weight":           "100",
						"port":             "80",
						"type":             "NLB",
						"server_group_id":  CHECKSET,
						"scaling_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_attach"},
			},
		},
	})
}

func TestAccAliCloudEssServerGroupAttachment_nonForceAttach_nlb(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_server_group_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)

	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroupServerGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingGroupServerGroupNotForceAttach)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssServerGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":       []string{"alicloud_ess_scaling_configuration.default"},
					"force_attach":     false,
					"weight":           "11",
					"port":             "22",
					"type":             "NLB",
					"server_group_id":  "${alicloud_nlb_server_group.default.id}",
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":               CHECKSET,
						"server_group_id":  CHECKSET,
						"scaling_group_id": CHECKSET,
						"weight":           "11",
						"port":             "22",
						"type":             "NLB",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_attach"},
			},
		},
	})
}

func TestAccAliCloudEssServerGroupAttachment_nonForceAttach_mutil(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_ess_server_group_attachment.default.1"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)

	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingGroupServerGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingGroupServerGroupNotForceAttachMutil)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssServerGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":            "2",
					"depends_on":       []string{"alicloud_ess_scaling_configuration.default"},
					"force_attach":     false,
					"weight":           "11",
					"port":             "22",
					"type":             "ALB",
					"server_group_id":  "${element(alicloud_alb_server_group.default.*.id,count.index)}",
					"scaling_group_id": "${alicloud_ess_scaling_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight":           "11",
						"port":             "22",
						"type":             "ALB",
						"server_group_id":  CHECKSET,
						"scaling_group_id": CHECKSET,
					}),
				),
			},
		},
	})
}

func testAccCheckEssServerGroupsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	essService := EssService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_group" {
			continue
		}

		scalingGroup, err := essService.DescribeEssScalingGroup(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if len(scalingGroup.ServerGroups.ServerGroup) > 0 {
			return WrapError(fmt.Errorf("There are still attached server groups."))
		}
	}
	return nil
}
func testAccEssScalingGroupServerGroupNotForceAttach(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "0"
	  max_size = "2"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}
	data "alicloud_images" "default2" {
		name_regex  = "^win"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_instance_types" "c6" {
      instance_type_family = "ecs.n4"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	data "alicloud_resource_manager_resource_groups" "default" {}
	resource "alicloud_nlb_server_group" "default" {
	  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
	  server_group_name = var.name
	  server_group_type = "Instance"
	  vpc_id            = alicloud_vpc.default.id
	  scheduler         = "Wrr"
	  protocol          = "TCP"
	  health_check {
		health_check_url             = "/test/index.html"
		health_check_domain          = "tf-testAcc.com"
		health_check_enabled         = true
		health_check_type            = "TCP"
		health_check_connect_port    = 0
		healthy_threshold            = 2
		unhealthy_threshold          = 2
		health_check_connect_timeout = 5
		health_check_interval        = 10
		http_check_method            = "GET"
		health_check_http_code       = ["http_2xx", "http_3xx", "http_4xx"]
	  }
	  connection_drain           = true
	  connection_drain_timeout   = 60
	  preserve_client_ip_enabled = true
	  tags = {
		Created = "TF"
	  }
	  address_ip_version = "Ipv4"
	}
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		image_id = data.alicloud_images.default2.images[0].id
		instance_type = data.alicloud_instance_types.c6.instance_types[0].id
		security_group_id = alicloud_security_group.default.id
		force_delete = true
		active = true
		enable = true
	}

	resource "alicloud_alb_server_group" "default" {
		server_group_name = "${var.name}"
		vpc_id = "${alicloud_vpc.default.id}"
		health_check_config {
		  health_check_enabled = "false"
		}
		sticky_session_config {
		  sticky_session_enabled = true
		  cookie                 = "tf-testAcc"
		  sticky_session_type    = "Server"
	  }
	}
	`, EcsInstanceCommonTestCase, name)
}

func testAccEssScalingGroupServerGroupNotForceAttachMutil(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "0"
	  max_size = "2"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}
	data "alicloud_images" "default2" {
		name_regex  = "^win"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_instance_types" "c6" {
      instance_type_family = "ecs.n4"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	data "alicloud_resource_manager_resource_groups" "default" {}
	resource "alicloud_nlb_server_group" "default" {
	  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
	  server_group_name = var.name
	  server_group_type = "Instance"
	  vpc_id            = alicloud_vpc.default.id
	  scheduler         = "Wrr"
	  protocol          = "TCP"
	  health_check {
		health_check_url             = "/test/index.html"
		health_check_domain          = "tf-testAcc.com"
		health_check_enabled         = true
		health_check_type            = "TCP"
		health_check_connect_port    = 0
		healthy_threshold            = 2
		unhealthy_threshold          = 2
		health_check_connect_timeout = 5
		health_check_interval        = 10
		http_check_method            = "GET"
		health_check_http_code       = ["http_2xx", "http_3xx", "http_4xx"]
	  }
	  connection_drain           = true
	  connection_drain_timeout   = 60
	  preserve_client_ip_enabled = true
	  tags = {
		Created = "TF"
	  }
	  address_ip_version = "Ipv4"
	}
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		image_id = data.alicloud_images.default2.images[0].id
		instance_type = data.alicloud_instance_types.c6.instance_types[0].id
		security_group_id = alicloud_security_group.default.id
		force_delete = true
		active = true
		enable = true
	}

	resource "alicloud_alb_server_group" "default" {
		count = 2
        server_group_name = "${var.name}"
		vpc_id = "${alicloud_vpc.default.id}"
		health_check_config {
		  health_check_enabled = "false"
		}
		sticky_session_config {
		  sticky_session_enabled = true
		  cookie                 = "tf-testAcc"
		  sticky_session_type    = "Server"
	  }
	}
	`, EcsInstanceCommonTestCase, name)
}

func testAccEssScalingGroupServerGroup(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	
	resource "alicloud_ess_scaling_group" "default" {
	  min_size = "0"
	  max_size = "2"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}

	data "alicloud_images" "default2" {
		name_regex  = "^win"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_instance_types" "c6" {
      instance_type_family = "ecs.n4"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	data "alicloud_resource_manager_resource_groups" "default" {}
	resource "alicloud_nlb_server_group" "default" {
	  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
	  server_group_name = var.name
	  server_group_type = "Instance"
	  vpc_id            = alicloud_vpc.default.id
	  scheduler         = "Wrr"
	  protocol          = "TCP"
	  health_check {
		health_check_url             = "/test/index.html"
		health_check_domain          = "tf-testAcc.com"
		health_check_enabled         = true
		health_check_type            = "TCP"
		health_check_connect_port    = 0
		healthy_threshold            = 2
		unhealthy_threshold          = 2
		health_check_connect_timeout = 5
		health_check_interval        = 10
		http_check_method            = "GET"
		health_check_http_code       = ["http_2xx", "http_3xx", "http_4xx"]
	  }
	  connection_drain           = true
	  connection_drain_timeout   = 60
	  preserve_client_ip_enabled = true
	  tags = {
		Created = "TF"
	  }
	  address_ip_version = "Ipv4"
	}
	resource "alicloud_ess_scaling_configuration" "default" {
		scaling_group_id = alicloud_ess_scaling_group.default.id
		image_id = data.alicloud_images.default2.images[0].id
		instance_type = data.alicloud_instance_types.c6.instance_types[0].id
		security_group_id = alicloud_security_group.default.id
		force_delete = true
		active = true
		enable = true
	}

	resource "alicloud_alb_server_group" "default" {
		server_group_name = "${var.name}"
		vpc_id = "${alicloud_vpc.default.id}"
		health_check_config {
		  health_check_enabled = "false"
		}
		sticky_session_config {
		  sticky_session_enabled = true
		  cookie                 = "tf-testAcc"
		  sticky_session_type    = "Server"
	  }
	}
	`, EcsInstanceCommonTestCase, name)
}
