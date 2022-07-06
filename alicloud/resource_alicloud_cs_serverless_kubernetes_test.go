package alicloud

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testRegionForCSSeverless = "cn-beijing"
)

func TestAccAlicloudCSServerlessKubernetes_basic(t *testing.T) {
	var timeZoneMap = map[string]string{
		"eu-central-1": "Europe/London",
		"cn-hangzhou":  "Asia/Shanghai",
		"cn-shanghai":  "Asia/Shanghai",
		"cn-beijing":   "Asia/Shanghai",
	}

	var regionId string
	if v := os.Getenv("ALICLOUD_REGION"); v == testRegionForCSSeverless {
		regionId = v
	} else {
		log.Printf("[INFO] Test: Using %s as test region", testRegionForCSSeverless)
		regionId = testRegionForCSSeverless
	}

	var timeZone string
	if v, ok := timeZoneMap[regionId]; ok {
		timeZone = v
	}

	var v *cs.ServerlessClusterResponse
	resourceId := "alicloud_cs_serverless_kubernetes.default"
	ra := resourceAttrInit(resourceId, csServerlessKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccserverlesskubernetes-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSServerlessKubernetesConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ServerlessKubernetesSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                           name,
					"vpc_id":                         "${alicloud_vpc.default.id}",
					"vswitch_ids":                    []string{"${alicloud_vswitch.default.id}"},
					"new_nat_gateway":                "true",
					"deletion_protection":            "false",
					"enable_rrsa":                    "true",
					"endpoint_public_access_enabled": "true",
					"load_balancer_spec":             "slb.s2.small",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"tags": map[string]string{
						"Platform": "TF",
					},
					"service_cidr":            "172.21.0.0/20",
					"service_discovery_types": []string{"PrivateZone"},
					"logging_type":            "SLS",
					"time_zone":               timeZone,
					"cluster_spec":            "ack.pro.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"deletion_protection":            "false",
						"new_nat_gateway":                "true",
						"endpoint_public_access_enabled": "true",
						"resource_group_id":              CHECKSET,
						"vswitch_ids.#":                  "1",
						"cluster_spec":                   "ack.pro.small",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"load_balancer_spec", "endpoint_public_access_enabled", "force_update",
					"new_nat_gateway", "private_zone", "zone_id", "vswitch_ids", "service_cidr", "service_discovery_types", "logging_type", "time_zone"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Platform": "TF",
						"Env":      "Pre",
					},
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":              "2",
						"tags.Platform":       "TF",
						"tags.Env":            "Pre",
						"deletion_protection": "false",
					}),
				),
			},
		},
	})
}

func resourceCSServerlessKubernetesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_eci_zones" "default" {}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "default" {
	vpc_name   = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
	vswitch_name      = "${var.name}"
	vpc_id            = "${alicloud_vpc.default.id}"
	cidr_block        = "10.1.1.0/24"
	availability_zone = "${data.alicloud_eci_zones.default.zones.0.zone_ids.0}"
}
`, name)
}

var csServerlessKubernetesBasicMap = map[string]string{
	"new_nat_gateway":                "true",
	"deletion_protection":            "false",
	"endpoint_public_access_enabled": "true",
	"force_update":                   "false",
}
