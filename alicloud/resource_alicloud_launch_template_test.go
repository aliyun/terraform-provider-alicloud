package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_launch_template", &resource.Sweeper{
		Name: "alicloud_launch_template",
		F:    testAlicloudLaunchTemplate,
	})
}

func testAlicloudLaunchTemplate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	request := ecs.CreateDescribeLaunchTemplatesRequest()
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeLaunchTemplates(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, region, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response := raw.(*ecs.DescribeLaunchTemplatesResponse)

	var ids []string
	for _, tpl := range response.LaunchTemplateSets.LaunchTemplateSet {
		if strings.HasPrefix(tpl.LaunchTemplateName, "tf-testAcc") {
			ids = append(ids, tpl.LaunchTemplateId)
		}
	}

	for i := range ids {
		templateRequest := ecs.CreateDeleteLaunchTemplateRequest()
		templateRequest.LaunchTemplateId = ids[i]
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteLaunchTemplate(templateRequest)
		})
		if err != nil || !NotFoundError(err) {
			log.Printf("delete template failed in sweepers, %v", err)
		}
		addDebug(templateRequest.GetActionName(), raw)
	}

	return nil
}

func TestAccAlicloudLaunchTemplateBasic(t *testing.T) {
	var v ecs.LaunchTemplateSet

	resourceId := "alicloud_launch_template.default"
	ra := resourceAttrInit(resourceId, testAccLaunchTemplateCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testaccLaunchTemplateBasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLaunchTemplateConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                          name,
					"description":                   name,
					"image_id":                      data.alicloud_images.default.images.0.id,
					"host_name":                     name,
					"instance_charge_type":          "PrePaid",
					"instance_name":                 name,
					"instance_type":                 data.alicloud_instance_types.default.instance_types.0.id,
					"internet_charge_type":          "PayByBandwidth",
					"internet_max_bandwidth_in":     "5",
					"internet_max_bandwidth_out":    "0",
					"io_optimized":                  "none",
					"key_pair_name":                 name,
					"ram_role_name":                 name,
					"network_type":                  "vpc",
					"security_enhancement_strategy": "Active",
					"spot_price_limit":              "5",
					"spot_strategy":                 "SpotWithPriceLimit",
					"security_group_id":             alicloud_security_group.default.id,
					"system_disk_category":          "cloud_ssd",
					"system_disk_description":       name,
					"system_disk_name":              name,
					"system_disk_size":              "40",
					"resource_group_id":             "rg-zkdfjahg9zxncv0",
					"userdata":                      "xxxxxxxxxxxxxx",
					"vswitch_id":                    "sw-ljkngaksdjfj0nnasdf",
					"vpc_id":                        "vpc-asdfnbg0as8dfk1nb2",
					"zone_id":                       "beijing-a",

					"tags": map[string]string{
						"tag1": "hello",
						"tag2": "world",
					},
					"network_interfaces": []map[string]string{
						{
							"name":              "eth0",
							"description":       "hello1",
							"primary_ip":        "10.0.0.2",
							"security_group_id": "xxxx",
							"vswitch_id":        "xxxxxxx",
						},
					},
					"data_disks": []map[string]string{
						{
							"name":        "disk1",
							"description": "test1",
						},
						{
							"name":        "disk2",
							"description": "test2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                    name,
						"description":             name,
						"host_name":               name,
						"instance_name":           name,
						"key_pair_name":           name,
						"ram_role_name":           name,
						"system_disk_description": name,
						"system_disk_name":        name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"io_optimized": "none",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"io_optimized": "none",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PostPaid",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_change",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"internet_charge_type": "PayByTraffic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_charge_type": "PayByTraffic",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id": "vpc-asdfnbg0as8dfk1nb2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id": "vpc-asdfnbg0as8dfk1nb2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disks": []map[string]string{
						{
							"name":        "disk1",
							"description": "test1",
						},
						{
							"name":        "disk2",
							"description": "test2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disks.#":             "2",
						"data_disks.0.name":        "disk1",
						"data_disks.0.description": "test1",
						"data_disks.1.name":        "disk2",
						"data_disks.1.description": "test2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"key_pair_name": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name": name + "_change",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": "${alicloud_security_group.default.id}_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_name": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_name": name + "_change",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.alicloud_instance_types.default.instance_types.0.id}_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"spot_strategy": "NoSpot",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_strategy": "NoSpot",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"userdata": "xxxxxxxxxxxxx",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"userdata": "xxxxxxxxxxxxx",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id": "sw-ljkngaksdjfj0nnasdf_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id": "sw-ljkngaksdjfj0nnasdf_change",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"image_id": "${data.alicloud_images.default.images.0.id}_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"network_type": "classic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_type": "classic",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"spot_price_limit": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_price_limit": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"ram_role_name": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ram_role_name": name + "_change",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"host_name": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name": name + "_change",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "rg-zkdfjahg9zxncv0_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": "rg-zkdfjahg9zxncv0_change",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_change",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id": "beijing-b",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id": "beijing-b",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"security_enhancement_strategy": "Deactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_enhancement_strategy": "Deactive",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"internet_max_bandwidth_out": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_out": "5",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_description": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_description": name + "_change",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "51",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "51",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "52",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "52",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "53",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "53",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "54",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "54",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"internet_max_bandwidth_in": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_in": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"tag1": "tag1",
						"tag2": "tag2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":    "2",
						"tags.tag1": "tag1",
						"tags.tag2": "tag2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_interfaces": []map[string]string{
						{
							"name":              "eth0",
							"description":       "hello1",
							"primary_ip":        "10.0.0.2",
							"security_group_id": "xxxx",
							"vswitch_id":        "xxxxxxx_change",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interfaces.#":                   "1",
						"network_interfaces.0.name":              "eth0",
						"network_interfaces.0.description":       "hello1",
						"network_interfaces.0.primary_ip":        "10.0.0.2",
						"network_interfaces.0.security_group_id": "xxxx",
						"network_interfaces.0.vswitch_id":        "xxxxxxx_change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                          name,
					"description":                   name,
					"image_id":                      data.alicloud_images.default.images.0.id,
					"host_name":                     name,
					"instance_charge_type":          "PrePaid",
					"instance_name":                 name,
					"instance_type":                 data.alicloud_instance_types.default.instance_types.0.id,
					"internet_charge_type":          "PayByBandwidth",
					"internet_max_bandwidth_in":     "5",
					"internet_max_bandwidth_out":    "0",
					"io_optimized":                  "none",
					"key_pair_name":                 name,
					"ram_role_name":                 name,
					"network_type":                  "vpc",
					"security_enhancement_strategy": "Active",
					"spot_price_limit":              "5",
					"spot_strategy":                 "SpotWithPriceLimit",
					"security_group_id":             alicloud_security_group.default.id,
					"system_disk_category":          "cloud_ssd",
					"system_disk_description":       name,
					"system_disk_name":              name,
					"system_disk_size":              "40",
					"resource_group_id":             "rg-zkdfjahg9zxncv0",
					"userdata":                      "xxxxxxxxxxxxxx",
					"vswitch_id":                    "sw-ljkngaksdjfj0nnasdf",
					"vpc_id":                        "vpc-asdfnbg0as8dfk1nb2",
					"zone_id":                       "beijing-a",

					"tags": map[string]string{
						"tag1": "hello",
						"tag2": "world",
					},
					"network_interfaces": []map[string]string{
						{
							"name":              "eth0",
							"description":       "hello1",
							"primary_ip":        "10.0.0.2",
							"security_group_id": "xxxx",
							"vswitch_id":        "xxxxxxx",
						},
					},
					"data_disks": []map[string]string{
						{
							"name":        "disk1",
							"description": "test1",
						},
						{
							"name":        "disk2",
							"description": "test2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                   name,
						"description":                            name,
						"host_name":                              name,
						"instance_name":                          name,
						"key_pair_name":                          name,
						"ram_role_name":                          name,
						"system_disk_description":                name,
						"system_disk_name":                       name,
						"image_id":                               CHECKSET,
						"instance_charge_type":                   "PrePaid",
						"instance_type":                          CHECKSET,
						"internet_charge_type":                   "PayByBandwidth",
						"internet_max_bandwidth_in":              "5",
						"internet_max_bandwidth_out":             "0",
						"io_optimized":                           "none",
						"network_type":                           "vpc",
						"security_enhancement_strategy":          "Active",
						"spot_price_limit":                       "5",
						"spot_strategy":                          "SpotWithPriceLimit",
						"security_group_id":                      CHECKSET,
						"system_disk_category":                   "cloud_ssd",
						"system_disk_size":                       "40",
						"resource_group_id":                      CHECKSET,
						"userdata":                               "xxxxxxxxxxxxxx",
						"vswitch_id":                             CHECKSET,
						"vpc_id":                                 CHECKSET,
						"zone_id":                                CHECKSET,
						"tags.%":                                 "2",
						"tags.tag1":                              "hello",
						"tags.tag2":                              "world",
						"network_interfaces.#":                   "1",
						"network_interfaces.0.name":              "eth0",
						"network_interfaces.0.description":       "hello1",
						"network_interfaces.0.primary_ip":        "10.0.0.2",
						"network_interfaces.0.security_group_id": "xxxx",
						"network_interfaces.0.vswitch_id":        "xxxxxxx",
						"data_disks.#":                           "2",
						"data_disks.0.name":                      "disk1",
						"data_disks.0.description":               "test1",
						"data_disks.1.name":                      "disk2",
						"data_disks.1.description":               "test2",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudLaunchTemplateMulti(t *testing.T) {
	var v ecs.LaunchTemplateSet

	resourceId := "alicloud_launch_template.default.4"
	ra := resourceAttrInit(resourceId, testAccLaunchTemplateCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testaccLaunchTemplateMulti%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLaunchTemplateConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":                         "5",
					"name":                          name + "${count.index}",
					"description":                   name,
					"image_id":                      data.alicloud_images.default.images.0.id,
					"host_name":                     name,
					"instance_charge_type":          "PrePaid",
					"instance_name":                 name,
					"instance_type":                 data.alicloud_instance_types.default.instance_types.0.id,
					"internet_charge_type":          "PayByBandwidth",
					"internet_max_bandwidth_in":     "5",
					"internet_max_bandwidth_out":    "0",
					"io_optimized":                  "none",
					"key_pair_name":                 name,
					"ram_role_name":                 name,
					"network_type":                  "vpc",
					"security_enhancement_strategy": "Active",
					"spot_price_limit":              "5",
					"spot_strategy":                 "SpotWithPriceLimit",
					"security_group_id":             alicloud_security_group.default.id,
					"system_disk_category":          "cloud_ssd",
					"system_disk_description":       name,
					"system_disk_name":              name,
					"system_disk_size":              "40",
					"resource_group_id":             "rg-zkdfjahg9zxncv0",
					"userdata":                      "xxxxxxxxxxxxxx",
					"vswitch_id":                    "sw-ljkngaksdjfj0nnasdf",
					"vpc_id":                        "vpc-asdfnbg0as8dfk1nb2",
					"zone_id":                       "beijing-a",

					"tags": map[string]string{
						"tag1": "hello",
						"tag2": "world",
					},
					"network_interfaces": []map[string]string{
						{
							"name":              "eth0",
							"description":       "hello1",
							"primary_ip":        "10.0.0.2",
							"security_group_id": "xxxx",
							"vswitch_id":        "xxxxxxx",
						},
					},
					"data_disks": []map[string]string{
						{
							"name":        "disk1",
							"description": "test1",
						},
						{
							"name":        "disk2",
							"description": "test2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                    name + "4",
						"description":             name,
						"host_name":               name,
						"instance_name":           name,
						"key_pair_name":           name,
						"ram_role_name":           name,
						"system_disk_description": name,
						"system_disk_name":        name,
					}),
				),
			},
		},
	})
}
func resourceLaunchTemplateConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*_64"
  most_recent = true
  owners      = "system"
}
resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alicloud_zones.default.zones.0.id
  name              = var.name
}
resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = alicloud_security_group.default.id
  	cidr_ip = "172.16.0.0/24"
}


`, name)
}

var testAccLaunchTemplateCheckMap = map[string]string{
	"image_id":                               CHECKSET,
	"instance_charge_type":                   "PrePaid",
	"instance_type":                          CHECKSET,
	"internet_charge_type":                   "PayByBandwidth",
	"internet_max_bandwidth_in":              "5",
	"internet_max_bandwidth_out":             "0",
	"io_optimized":                           "none",
	"network_type":                           "vpc",
	"security_enhancement_strategy":          "Active",
	"spot_price_limit":                       "5",
	"spot_strategy":                          "SpotWithPriceLimit",
	"security_group_id":                      CHECKSET,
	"system_disk_category":                   "cloud_ssd",
	"system_disk_size":                       "40",
	"resource_group_id":                      CHECKSET,
	"userdata":                               "xxxxxxxxxxxxxx",
	"vswitch_id":                             CHECKSET,
	"vpc_id":                                 CHECKSET,
	"zone_id":                                CHECKSET,
	"tags.%":                                 "2",
	"tags.tag1":                              "hello",
	"tags.tag2":                              "world",
	"network_interfaces.#":                   "1",
	"network_interfaces.0.name":              "eth0",
	"network_interfaces.0.description":       "hello1",
	"network_interfaces.0.primary_ip":        "10.0.0.2",
	"network_interfaces.0.security_group_id": "xxxx",
	"network_interfaces.0.vswitch_id":        "xxxxxxx",
	"data_disks.#":                           "2",
	"data_disks.0.name":                      "disk1",
	"data_disks.0.description":               "test1",
	"data_disks.1.name":                      "disk2",
	"data_disks.1.description":               "test2",
}
