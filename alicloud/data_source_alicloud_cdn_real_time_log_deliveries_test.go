package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCDNRealTimeLogDeliveriesDataSource(t *testing.T) {
	rand := acctest.RandInt()

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCdnRealTimeLogDeliveriesDataSourceName(rand, map[string]string{
			"domain": `"${alicloud_cdn_real_time_log_delivery.default.id}"`,
			"status": `"online"`,
		}),
		fakeConfig: testAccCheckAlicloudCdnRealTimeLogDeliveriesDataSourceName(rand, map[string]string{
			"domain": `"${alicloud_cdn_real_time_log_delivery.default.id}"`,
			"status": `"offline"`,
		}),
	}

	var existAlicloudCdnRealTimeLogDeliveriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"deliveries.#":        "1",
			"deliveries.0.domain": fmt.Sprintf("tf-testaccrealtimelogdeliveries-%d.example.com", rand),
		}
	}
	var fakeAlicloudCdnRealTimeLogDeliveriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"deliveries.#": "0",
		}
	}
	var AlicloudCdnRealTimeLogDeliveriesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cdn_real_time_log_deliveries.default",
		existMapFunc: existAlicloudCdnRealTimeLogDeliveriesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCdnRealTimeLogDeliveriesDataSourceNameMapFunc,
	}
	AlicloudCdnRealTimeLogDeliveriesCheckInfo.dataSourceTestCheck(t, rand, statusConf)
}
func testAccCheckAlicloudCdnRealTimeLogDeliveriesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testaccrealtimelogdeliveries-%d"
}
variable "domain_name" {	
	default = "tf-testaccrealtimelogdeliveries-%d.example.com"
}
resource "alicloud_cdn_domain_new" "default" {
  domain_name = var.domain_name
  cdn_type = "web"
  scope = "overseas"
  sources {
	 content = "www.aliyuntest.com"
	 type = "domain"
	 priority = 20
	 port = 80
	 weight = 10
  }
}
resource "alicloud_log_project" "default" {
  name        = var.name
  description = var.name
}
resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.name
  name                  = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
resource "alicloud_cdn_real_time_log_delivery" "default" {
  domain = alicloud_cdn_domain_new.default.domain_name
  project = alicloud_log_project.default.name
  logstore = alicloud_log_store.default.name
  sls_region = "%s"
}
data "alicloud_cdn_real_time_log_deliveries" "default" {	
	%s
}
`, rand, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
