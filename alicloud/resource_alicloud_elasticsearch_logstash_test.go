package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Elasticsearch Logstash. >>> Resource test cases, automatically generated.
// Case 3472
func TestAccAlicloudElasticsearchLogstash_basic3472(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_logstash.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchLogstashMap3472)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchLogstash")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%selasticsearchlogstash%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchLogstashBasicDependence3472)
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
					"version": "7.4_with_X-Pack",
					"node_spec": []map[string]interface{}{
						{
							"disk_type": "cloud_efficiency",
							"spec":      "elasticsearch.sn1ne.large",
							"disk":      "20",
						},
					},
					"network_config": []map[string]interface{}{
						{
							"type":       "vpc",
							"vpc_id":     "${alicloud_vpc.defaultZFtcRh.id}",
							"vswitch_id": "${alicloud_vswitch.defaultMiMSn6.id}",
							"vs_area":    "cn-hangzhou-i",
						},
					},
					"node_amount": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version":     "7.4_with_X-Pack",
						"node_amount": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-acc-create-test-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-acc-create-test-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${alicloud_resource_manager_resource_group.defaultWUcLGe.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{
							"disk_type": "cloud_efficiency",
							"spec":      "elasticsearch.sn1ne.large",
							"disk":      "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_amount": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_amount": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn2ne.large",
							"disk":      "30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-test-updatedescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-test-updatedescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-acc-create-test-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-acc-create-test-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{
							"disk_type": "cloud_efficiency",
							"spec":      "elasticsearch.sn1ne.large",
							"disk":      "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn2ne.large",
							"disk":      "30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_spec": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "tf-acc-create-test-1",
					"resource_group_id": "${alicloud_resource_manager_resource_group.defaultWUcLGe.id}",
					"version":           "7.4_with_X-Pack",
					"node_spec": []map[string]interface{}{
						{
							"disk_type": "cloud_efficiency",
							"spec":      "elasticsearch.sn1ne.large",
							"disk":      "20",
						},
					},
					"network_config": []map[string]interface{}{
						{
							"type":       "vpc",
							"vpc_id":     "${alicloud_vpc.defaultZFtcRh.id}",
							"vswitch_id": "${alicloud_vswitch.defaultMiMSn6.id}",
							"vs_area":    "cn-hangzhou-i",
						},
					},
					"payment_type": "PayAsYouGo",
					"node_amount":  "1",
					"payment_info": []map[string]interface{}{
						{},
					},
					"tags": []map[string]interface{}{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "tf-acc-create-test-1",
						"resource_group_id": CHECKSET,
						"version":           "7.4_with_X-Pack",
						"payment_type":      "PayAsYouGo",
						"node_amount":       "1",
						"tags.#":            "0",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudElasticsearchLogstashMap3472 = map[string]string{
	"status":       CHECKSET,
	"create_time":  CHECKSET,
	"payment_type": CHECKSET,
}

func AlicloudElasticsearchLogstashBasicDependence3472(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultZFtcRh" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "defaultMiMSn6" {
  vpc_id       = alicloud_vpc.defaultZFtcRh.id
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "${var.name}1"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_resource_manager_resource_group" "defaultWUcLGe" {
  display_name        = "rdktest05"
  resource_group_name = "${var.name}2"
}


`, name)
}

// Case 3472  twin
func TestAccAlicloudElasticsearchLogstash_basic3472_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_logstash.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchLogstashMap3472)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchLogstash")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%selasticsearchlogstash%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchLogstashBasicDependence3472)
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
					"description":       "tf-acc-create-test-1",
					"resource_group_id": "${alicloud_resource_manager_resource_group.defaultWUcLGe.id}",
					"version":           "7.4_with_X-Pack",
					"node_spec": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn2ne.large",
							"disk":      "30",
						},
					},
					"network_config": []map[string]interface{}{
						{
							"type":       "vpc",
							"vpc_id":     "${alicloud_vpc.defaultZFtcRh.id}",
							"vswitch_id": "${alicloud_vswitch.defaultMiMSn6.id}",
							"vs_area":    "cn-hangzhou-i",
						},
					},
					"payment_type": "PayAsYouGo",
					"node_amount":  "1",
					"payment_info": []map[string]interface{}{
						{},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "tf-acc-create-test-1",
						"resource_group_id": CHECKSET,
						"version":           "7.4_with_X-Pack",
						"payment_type":      "PayAsYouGo",
						"node_amount":       "1",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
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

// Test Elasticsearch Logstash. <<< Resource test cases, automatically generated.
