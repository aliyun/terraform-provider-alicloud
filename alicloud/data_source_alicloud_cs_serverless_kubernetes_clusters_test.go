package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCSServerlessKubernetesClustersDataSource(t *testing.T) {
	prevRegion := os.Getenv("ALICLOUD_REGION")
	os.Setenv("ALICLOUD_REGION", "ap-southeast-1")
	defer os.Setenv("ALICLOUD_REGION", prevRegion)

	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cs_serverless_kubernetes_clusters.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testaccserverlessk8s-%d", rand),
		dataSourceCSServerlessKubernetesClustersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_serverless_kubernetes.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_serverless_kubernetes.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_serverless_kubernetes.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_serverless_kubernetes.default.name}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_serverless_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_serverless_kubernetes.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_serverless_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_serverless_kubernetes.default.name}-fake",
		}),
	}

	var existCSServerlessKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"ids.0":                          CHECKSET,
			"names.#":                        "1",
			"names.0":                        REGEXMATCH + fmt.Sprintf("tf-testaccserverlessk8s-%d", rand),
			"clusters.#":                     "1",
			"clusters.0.id":                  CHECKSET,
			"clusters.0.name":                REGEXMATCH + fmt.Sprintf("tf-testaccserverlessk8s-%d", rand),
			"clusters.0.security_group_id":   CHECKSET,
			"clusters.0.vpc_id":              CHECKSET,
			"clusters.0.vswitch_id":          CHECKSET,
			"clusters.0.deletion_protection": CHECKSET,
			"clusters.0.tags.%":              CHECKSET,
			"clusters.0.endpoint_public_access_enabled": CHECKSET,
			"clusters.0.connections.%":                  CHECKSET,
		}
	}

	var fakeCSServerlessKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"clusters.#": "0",
		}
	}

	var csServerlessKubernetesClustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCSServerlessKubernetesClustersMapFunc,
		fakeMapFunc:  fakeCSServerlessKubernetesClustersMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.ServerlessKubernetesSupportedRegions)
	}
	csServerlessKubernetesClustersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceCSServerlessKubernetesClustersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_eci_zones" "default" {
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
	vpc_id   = data.alicloud_vpcs.default.ids.0
	zone_id  = data.alicloud_eci_zones.default.zones.0.zone_ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_eci_zones.default.zones.0.zone_ids.0
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_cs_serverless_kubernetes" "default" {
  name_prefix = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
  vswitch_ids = [local.vswitch_id]
  new_nat_gateway = true
  cluster_spec = "ack.pro.small" 
  endpoint_public_access_enabled = true
  private_zone = false
  deletion_protection = false
  service_cidr = "10.0.1.0/24"
  load_balancer_spec = "slb.s2.small"
  tags = {
		"k-aa":"v-aa"
		"k-bb":"v-aa",
  }
}
`, name)
}
