// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var elasticsearchLogstashMap = map[string]string{
	"description":                 CHECKSET,
	"network_config.#":            "1",
	"network_config.0.type":       "vpc",
	"network_config.0.vpc_id":     CHECKSET,
	"network_config.0.vswitch_id": CHECKSET,
	"network_config.0.vs_area":    CHECKSET,
	"node_amount":                 CHECKSET,
	"node_spec.#":                 "1",
	"node_spec.0.disk_type":       "cloud_efficiency",
	"node_spec.0.spec":            "elasticsearch.sn1ne.large",
	"node_spec.0.disk":            "20",
	"payment_type":                CHECKSET,
	"resource_group_id":           CHECKSET,
	"status":                      CHECKSET,
	"tags.%":                      REMOVEKEY,
	"version":                     "7.4_with_X-Pack",
}

func TestAccAliCloudElasticsearchLogstash_basic(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccESLogstash%d", rand)

	resourceId := "alicloud_elasticsearch_logstash.default"
	ra := resourceAttrInit(resourceId, elasticsearchLogstashMap)
	serviceFunc := func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, nil, serviceFunc, "DescribeElasticsearchLogstash")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceElasticsearchLogstashConfigDependence)

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
					"version":           "7.4_with_X-Pack",
					"payment_type":      "Subscription",
					"node_amount":       "1",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"node_spec": []map[string]interface{}{
						{"disk_type": "cloud_efficiency", "spec": "elasticsearch.sn1ne.large", "disk": "20"},
					},
					"network_config": []map[string]interface{}{
						{"type": "vpc", "vpc_id": "${alicloud_vpc.default.id}", "vswitch_id": "${alicloud_vswitch.default.id}", "vs_area": "${data.alicloud_elasticsearch_zones.default.zones.0.id}"},
					},
					"payment_info": []map[string]interface{}{
						{"duration": "1", "pricing_cycle": "Month", "auto_renew": false, "auto_renew_duration": "1"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                 "tf-acc-create-test-1",
						"version":                     "7.4_with_X-Pack",
						"payment_type":                "Subscription",
						"node_amount":                 "1",
						"resource_group_id":           CHECKSET,
						"status":                      CHECKSET,
						"create_time":                 CHECKSET,
						"updated_at":                  CHECKSET,
						"network_config.#":            "1",
						"network_config.0.type":       "vpc",
						"network_config.0.vpc_id":     CHECKSET,
						"network_config.0.vswitch_id": CHECKSET,
						"network_config.0.vs_area":    CHECKSET,
						"node_spec.#":                 "1",
						"node_spec.0.disk_type":       "cloud_efficiency",
						"node_spec.0.spec":            "elasticsearch.sn1ne.large",
						"node_spec.0.disk":            "20",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-acc-update-test-2",
					"node_amount": "2",
					"node_spec": []map[string]interface{}{
						{"disk_type": "cloud_efficiency", "spec": "elasticsearch.sn1ne.large", "disk": "20"},
					},
					"payment_info": []map[string]interface{}{
						{"duration": "1", "pricing_cycle": "Month", "auto_renew": false, "auto_renew_duration": "1"},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "tf-acc-update-test-2",
						"node_amount":  "2",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

func resourceElasticsearchLogstashConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_elasticsearch_zones" "default" {}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_elasticsearch_zones.default.zones.0.id
}
`, name)
}
