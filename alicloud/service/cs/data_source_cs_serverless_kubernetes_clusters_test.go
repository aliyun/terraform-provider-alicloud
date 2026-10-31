package cs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCSServerlessKubernetesClustersDataSource(t *testing.T) {

	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cs_serverless_kubernetes_clusters.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testaccserverlessk8s-%d", rand),
		dataSourceCSServerlessKubernetesClustersConfigDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_serverless_kubernetes.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_serverless_kubernetes.default.id}-fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_serverless_kubernetes.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_serverless_kubernetes.default.name}-fake",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_serverless_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_serverless_kubernetes.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_serverless_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_serverless_kubernetes.default.name}-fake",
		}),
	}

	var existCSServerlessKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"ids.0":                          acctest.CHECKSET,
			"names.#":                        "1",
			"names.0":                        acctest.REGEXMATCH + fmt.Sprintf("tf-testaccserverlessk8s-%d", rand),
			"clusters.#":                     "1",
			"clusters.0.id":                  acctest.CHECKSET,
			"clusters.0.name":                acctest.REGEXMATCH + fmt.Sprintf("tf-testaccserverlessk8s-%d", rand),
			"clusters.0.security_group_id":   acctest.CHECKSET,
			"clusters.0.vpc_id":              acctest.CHECKSET,
			"clusters.0.vswitch_id":          acctest.CHECKSET,
			"clusters.0.deletion_protection": acctest.CHECKSET,
			"clusters.0.tags.%":              acctest.CHECKSET,
			"clusters.0.endpoint_public_access_enabled": acctest.CHECKSET,
			"clusters.0.connections.#":                  acctest.CHECKSET,
		}
	}

	var fakeCSServerlessKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"clusters.#": "0",
		}
	}

	var csServerlessKubernetesClustersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existCSServerlessKubernetesClustersMapFunc,
		FakeMapFunc:  fakeCSServerlessKubernetesClustersMapFunc,
	}
	preCheck := func() {
	}
	csServerlessKubernetesClustersCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceCSServerlessKubernetesClustersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = cidrsubnet(alicloud_vpc.vpc.cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_cs_serverless_kubernetes" "default" {
  name_prefix                    = "${var.name}"
  vpc_id                         = alicloud_vpc.vpc.id
  vswitch_ids                    = [alicloud_vswitch.vswitch.id]
  new_nat_gateway                = false
  cluster_spec                   = "ack.pro.small"
  endpoint_public_access_enabled = true
  private_zone                   = false
  deletion_protection            = false
  service_cidr                   = "10.0.1.0/24"
  load_balancer_spec             = "slb.s2.small"
  tags = {
    "k-aa" : "v-aa"
    "k-bb" : "v-aa",
  }
}
`, name)
}
