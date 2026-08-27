package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudSlsMetricStoresDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSlsMetricStoresDataSourceConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_sls_metric_stores.default", "metric_stores.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_sls_metric_stores.default", "metric_stores.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_sls_metric_stores.default", "metric_stores.0.create_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_sls_metric_stores.default", "metric_stores.0.last_modify_time"),
					resource.TestCheckResourceAttr("data.alicloud_sls_metric_stores.default", "metric_stores.0.metric_store_name", name),
					resource.TestCheckResourceAttr("data.alicloud_sls_metric_stores.default", "metric_stores.0.project_name", "sls-sdk-testp-metricstore"),
					resource.TestCheckResourceAttr("data.alicloud_sls_metric_stores.default", "metric_stores.0.ttl", "7"),
					resource.TestCheckResourceAttr("data.alicloud_sls_metric_stores.default", "metric_stores.0.shard_count", "2"),
					resource.TestCheckResourceAttr("data.alicloud_sls_metric_stores.default", "metric_stores.0.mode", "standard"),
					resource.TestCheckResourceAttr("data.alicloud_sls_metric_stores.default", "metric_stores.0.metric_type", "prometheus"),
				),
			},
			{
				Config: testAccSlsMetricStoresDataSourceConfigNameRegex(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_sls_metric_stores.regex", "metric_stores.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_sls_metric_stores.regex", "metric_stores.0.metric_store_name", name),
				),
			},
		},
	})
}

func testAccSlsMetricStoresDataSourceConfigBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "defaultj3iK32" {
  description = "test"
  name        = "sls-sdk-testp-metricstore"
}

resource "alicloud_sls_metric_store" "default" {
  project_name      = "sls-sdk-testp-metricstore"
  metering_mode     = "ChargeByFunction"
  mode              = "standard"
  metric_type       = "prometheus"
  metric_store_name = var.name
  ttl               = "7"
  shard_count       = "2"
}

data "alicloud_sls_metric_stores" "default" {
  project_name      = alicloud_sls_metric_store.default.project_name
  metric_store_name = alicloud_sls_metric_store.default.metric_store_name
  offset            = 0
  size              = 100
}
`, name)
}

func testAccSlsMetricStoresDataSourceConfigNameRegex(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "defaultj3iK32" {
  description = "test"
  name        = "sls-sdk-testp-metricstore"
}

resource "alicloud_sls_metric_store" "default" {
  project_name      = "sls-sdk-testp-metricstore"
  metering_mode     = "ChargeByFunction"
  mode              = "standard"
  metric_type       = "prometheus"
  metric_store_name = var.name
  ttl               = "7"
  shard_count       = "2"
}

data "alicloud_sls_metric_stores" "regex" {
  project_name = alicloud_sls_metric_store.default.project_name
  name_regex   = "^${var.name}$"
  offset       = 0
  size         = 100
}
`, name)
}
