package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_alb_server_group",
		&resource.Sweeper{
			Name: "alicloud_alb_server_group",
			F:    testSweepAlbServerGroup,
		})
}

func testSweepAlbServerGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListServerGroups"
	request := map[string]interface{}{}

	request["MaxResults"] = PageSizeXLarge

	var response map[string]interface{}
	conn, err := client.NewAlbClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
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

		resp, err := jsonpath.Get("$.ServerGroups", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.ServerGroups", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["ServerGroupName"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["ServerGroupName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Alb Server Group: %s", item["ServerGroupName"].(string))
				continue
			}
			action := "DeleteServerGroup"
			request := map[string]interface{}{
				"ServerGroupId": item["ServerGroupId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Alb Server Group (%s): %s", item["ServerGroupName"].(string), err)
			}
			log.Printf("[INFO] Delete Alb Server Group success: %s ", item["ServerGroupName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudALBServerGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_server_group.default"
	ra := resourceAttrInit(resourceId, AlicloudALBServerGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBServerGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol":          "HTTP",
					"vpc_id":            "${data.alicloud_vpcs.default.vpcs.0.id}",
					"server_group_name": "${var.name}",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
					"sticky_session_config": []map[string]interface{}{
						{
							"sticky_session_enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":                "HTTP",
						"server_group_name":       name,
						"sticky_session_config.#": "1",
						"health_check_config.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": "tf-testAcc-new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": "tf-testAcc-new",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "Wlc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "Wlc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"health_check_connect_port": "46325",
							"health_check_enabled":      "true",
							"health_check_host":         "tf-testAcc.com",
							"health_check_codes":        []string{"http_2xx", "http_3xx", "http_4xx"},
							"health_check_http_version": "HTTP1.1",
							"health_check_interval":     "2",
							"health_check_method":       "HEAD",
							"health_check_path":         "/tf-testAcc",
							"health_check_protocol":     "HTTP",
							"health_check_timeout":      "5",
							"healthy_threshold":         "3",
							"unhealthy_threshold":       "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol": "HTTPS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol": "HTTPS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"sticky_session_config": []map[string]interface{}{
						{
							"cookie_timeout":         "2000",
							"sticky_session_enabled": "true",
							"sticky_session_type":    "Insert",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sticky_session_config": []map[string]interface{}{
						{
							"cookie":                 "tf-testAcc",
							"sticky_session_enabled": "true",
							"sticky_session_type":    "Server",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc7",
						"For":     "Tftestacc7",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc7",
						"tags.For":     "Tftestacc7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": "${var.name}",
					"scheduler":         "Wrr",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
					"sticky_session_config": []map[string]interface{}{
						{
							"sticky_session_enabled": "false",
						},
					},
					"tags": map[string]string{
						"Created": "tfTestAcc99",
						"For":     "Tftestacc99",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name":       name,
						"scheduler":               "Wrr",
						"health_check_config.#":   "1",
						"sticky_session_config.#": "1",
						"tags.%":                  "2",
						"tags.Created":            "tfTestAcc99",
						"tags.For":                "Tftestacc99",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAlicloudALBServerGroup_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_server_group.default"
	ra := resourceAttrInit(resourceId, AlicloudALBServerGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBServerGroupBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol":          "HTTP",
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"server_group_name": "${var.name}",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
					"sticky_session_config": []map[string]interface{}{
						{
							"sticky_session_enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":                "HTTP",
						"server_group_name":       name,
						"sticky_session_config.#": "1",
						"health_check_config.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"description": "tf-testAcc",
							"port":        "80",
							"server_id":   "${alicloud_instance.instance.id}",
							"server_ip":   "${alicloud_instance.instance.private_ip}",
							"server_type": "Ecs",
							"weight":      "10",
						},
						{
							"description": "tf-testAcc",
							"port":        "8080",
							"server_id":   "${alicloud_instance.instance.id}",
							"server_ip":   "${alicloud_instance.instance.private_ip}",
							"server_type": "Ecs",
							"weight":      "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"description": "tf-testAcc-update",
							"port":        "80",
							"server_id":   "${alicloud_instance.instance.id}",
							"server_ip":   "${alicloud_instance.instance.private_ip}",
							"server_type": "Ecs",
							"weight":      "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"description": "tf-testAcc-update",
							"port":        "80",
							"server_id":   "${alicloud_instance.instance.id}",
							"server_ip":   "${alicloud_instance.instance.private_ip}",
							"server_type": "Ecs",
							"weight":      "10",
						},
						{
							"description": "tf-testAcc-update-8056",
							"port":        "8056",
							"server_id":   "${alicloud_instance.instance.id}",
							"server_ip":   "${alicloud_instance.instance.private_ip}",
							"server_type": "Ecs",
							"weight":      "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAlicloudALBServerGroup_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_server_group.default"
	ra := resourceAttrInit(resourceId, AlicloudALBServerGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBServerGroupBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol":          "HTTP",
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"server_group_name": "${var.name}",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"dry_run":           "false",
					"sticky_session_config": []map[string]interface{}{
						{
							"cookie":                 "tf-testAcc",
							"sticky_session_enabled": "true",
							"sticky_session_type":    "Server",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":                "HTTP",
						"server_group_name":       name,
						"sticky_session_config.#": "1",
						"health_check_config.#":   "1",
						"resource_group_id":       CHECKSET,
						"dry_run":                 "false",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAlicloudALBServerGroup_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_server_group.default"
	ra := resourceAttrInit(resourceId, AlicloudALBServerGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBServerGroupBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol":          "HTTP",
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"server_group_name": "${var.name}",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"dry_run":           "false",
					"sticky_session_config": []map[string]interface{}{
						{
							"cookie_timeout":         "2000",
							"sticky_session_enabled": "true",
							"sticky_session_type":    "Insert",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":                "HTTP",
						"server_group_name":       name,
						"sticky_session_config.#": "1",
						"health_check_config.#":   "1",
						"resource_group_id":       CHECKSET,
						"dry_run":                 "false",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudALBServerGroupMap0 = map[string]string{
	"tags.%":            NOSET,
	"dry_run":           NOSET,
	"resource_group_id": CHECKSET,
	"status":            CHECKSET,
	"scheduler":         CHECKSET,
	"vpc_id":            CHECKSET,
}

func AlicloudALBServerGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_resource_manager_resource_groups" "default" {}

`, name)

}

func AlicloudALBServerGroupBasicDependence1(name string) string {
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
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_instance" "instance" {
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  instance_name              = var.name
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = local.vswitch_id
}

`, name)

}

func AlicloudALBServerGroupBasicDependence3(name string) string {
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
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu"
  most_recent = true
  owners      = "system"
}

data "alicloud_vpcs" "default" {
 name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_zones.default.zones[0].id
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_instance" "instance" {
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  instance_name              = var.name
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = local.vswitch_id
}

data "alicloud_resource_manager_resource_groups" "default" {}

`, name)

}
