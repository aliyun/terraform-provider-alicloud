package gwlb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/gwlb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Gwlb Listener. >>> Resource test cases, automatically generated.
// Case Test Listener 8508
func TestAccAliCloudGwlbListener_basic8508(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gwlb_listener.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGwlbListenerMap8508)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gwlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeListener")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgwlblistener%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGwlbListenerBasicDependence8508)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"listener_description": "test-tf-lsn",
					"server_group_id":      "${alicloud_gwlb_server_group.defaultoAkLbr.id}",
					"load_balancer_id":     "${alicloud_gwlb_load_balancer.defaultQ5setL.id}",
					"dry_run":              "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_description": "test-tf-lsn",
						"server_group_id":      acctest.CHECKSET,
						"load_balancer_id":     acctest.CHECKSET,
						"dry_run":              "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"listener_description": "test-tf-update",
					"server_group_id":      "${alicloud_gwlb_server_group.defaultN4DOzm.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_description": "test-tf-update",
						"server_group_id":      acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"listener_description": "test-tf-lsn",
					"server_group_id":      "${alicloud_gwlb_server_group.defaultoAkLbr.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_description": "test-tf-lsn",
						"server_group_id":      acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudGwlbListenerMap8508 = map[string]string{
	"status":    acctest.CHECKSET,
	"region_id": acctest.CHECKSET,
}

func AlicloudGwlbListenerBasicDependence8508(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

variable "zone_id2" {
  default = "cn-wulanchabu-c"
}

variable "zone_id1" {
  default = "cn-wulanchabu-b"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaulti9Axhl" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default9NaKmL" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id1
  cidr_block   = "10.0.0.0/24"
  vswitch_name = format("%%s1", var.name)
}

resource "alicloud_vswitch" "defaultH4pKT4" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id2
  cidr_block   = "10.0.1.0/24"
  vswitch_name = format("%%s2", var.name)
}

resource "alicloud_gwlb_load_balancer" "defaultQ5setL" {
  vpc_id             = alicloud_vpc.defaulti9Axhl.id
  load_balancer_name = format("%%s3", var.name)
  zone_mappings {
    vswitch_id = alicloud_vswitch.default9NaKmL.id
    zone_id    = var.zone_id1
  }
  address_ip_version = "Ipv4"
}

resource "alicloud_gwlb_server_group" "defaultoAkLbr" {
  scheduler = "5TCH"
  health_check_config {
    health_check_protocol        = "TCP"
    health_check_connect_port    = "80"
    health_check_connect_timeout = "5"
    health_check_enabled         = true
    health_check_interval        = "10"
    healthy_threshold            = "2"
    unhealthy_threshold          = "2"
  }
  protocol          = "GENEVE"
  server_group_type = "Ip"
  connection_drain_config {
    connection_drain_enabled = true
    connection_drain_timeout = "1"
  }
  vpc_id = alicloud_vpc.defaulti9Axhl.id
  servers {
    server_id   = "10.0.0.1"
    server_ip   = "10.0.0.1"
    server_type = "Ip"
  }
  servers {
    server_id   = "10.0.0.2"
    server_ip   = "10.0.0.2"
    server_type = "Ip"
  }
  servers {
    server_id   = "10.0.0.3"
    server_ip   = "10.0.0.3"
    server_type = "Ip"
  }
  server_group_name = format("%%s4", var.name)
}

resource "alicloud_gwlb_server_group" "defaultN4DOzm" {
  scheduler = "5TCH"
  health_check_config {
    health_check_protocol        = "TCP"
    health_check_connect_port    = "80"
    health_check_connect_timeout = "5"
    health_check_enabled         = true
    health_check_interval        = "10"
    healthy_threshold            = "2"
    unhealthy_threshold          = "2"
  }
  protocol          = "GENEVE"
  server_group_type = "Ip"
  connection_drain_config {
    connection_drain_enabled = true
    connection_drain_timeout = "1"
  }
  vpc_id = alicloud_vpc.defaulti9Axhl.id
  servers {
    server_id   = "10.0.0.1"
    server_ip   = "10.0.0.1"
    server_type = "Ip"
  }
  servers {
    server_id   = "10.0.0.2"
    server_ip   = "10.0.0.2"
    server_type = "Ip"
  }
  servers {
    server_id   = "10.0.0.3"
    server_ip   = "10.0.0.3"
    server_type = "Ip"
  }
  server_group_name = format("%%s5", var.name)
}


`, name)
}

// Test Gwlb Listener. <<< Resource test cases, automatically generated.

// TestAccAliCloudGwlbListener_tcpIdleTimeout8508 covers the tcp_idle_timeout
// attribute: server default on omit, custom value on create, and update.
func TestAccAliCloudGwlbListener_tcpIdleTimeout8508(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gwlb_listener.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGwlbListenerMap8508)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gwlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeListener")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgwlblistener%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGwlbListenerBasicDependence8508)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// omit tcp_idle_timeout → server default 350
				Config: testAccConfig(map[string]interface{}{
					"server_group_id":  "${alicloud_gwlb_server_group.defaultoAkLbr.id}",
					"load_balancer_id": "${alicloud_gwlb_load_balancer.defaultQ5setL.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tcp_idle_timeout": "350",
					}),
				),
			},
			{
				// set custom value 120
				Config: testAccConfig(map[string]interface{}{
					"server_group_id":  "${alicloud_gwlb_server_group.defaultoAkLbr.id}",
					"load_balancer_id": "${alicloud_gwlb_load_balancer.defaultQ5setL.id}",
					"tcp_idle_timeout": "120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tcp_idle_timeout": "120",
					}),
				),
			},
			{
				// update to 200
				Config: testAccConfig(map[string]interface{}{
					"server_group_id":  "${alicloud_gwlb_server_group.defaultoAkLbr.id}",
					"load_balancer_id": "${alicloud_gwlb_load_balancer.defaultQ5setL.id}",
					"tcp_idle_timeout": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tcp_idle_timeout": "200",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}
