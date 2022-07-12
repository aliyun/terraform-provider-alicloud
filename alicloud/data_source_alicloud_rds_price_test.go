package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"strings"
	"testing"
)

func TestAccAlicloudRdsPriceDataSource(t *testing.T) {
	rand := acctest.RandInt()

	testConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsPriceDataSourceName(rand, map[string]string{
			"engine_version":           `"13.0"`,
			"db_instance_class":        `"${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}"`,
			"db_instance_storage":      `"20"`,
			"quantity":                 `"10"`,
			"engine":                   `"PostgreSQL"`,
			"commodity_code":           `"bards"`,
			"pay_type":                 `"Prepaid"`,
			"used_time":                `"1"`,
			"time_type":                `"Year"`,
			"instance_used_type":       `"0"`,
			"order_type":               `"BUY"`,
			"db_instance_storage_type": `"cloud_essd"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsPriceDataSourceName(rand, map[string]string{
			"engine_version":           `"13.0"`,
			"db_instance_class":        `"${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}"`,
			"db_instance_storage":      `"20"`,
			"quantity":                 `"1"`,
			"engine":                   `"PostgreSQL"`,
			"commodity_code":           `"bards"`,
			"pay_type":                 `"Prepaid"`,
			"used_time":                `"1"`,
			"time_type":                `"Year"`,
			"instance_used_type":       `"0"`,
			"order_type":               `"BUY"`,
			"db_instance_storage_type": `"cloud_essd"`,
		}),
	}
	var existAlicloudRdsPriceDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"price.#": "1",
		}
	}
	var fakeAlicloudRdsPriceDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"price.#": "1",
		}
	}
	var alicloudRdsLogsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_rds_price.default",
		existMapFunc: existAlicloudRdsPriceDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudRdsPriceDataSourceNameMapFunc,
	}
	alicloudRdsLogsCheckInfo.dataSourceTestCheck(t, rand, testConf)
}

func testAccCheckAlicloudRdsPriceDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-rds-price"
}

data "alicloud_db_zones" "default" {
  engine                   = "PostgreSQL"
  engine_version           = "13.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "PostgreSQL"
  engine_version           = "13.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_rds_price" "default" {
  %s
}`, strings.Join(pairs, "\n"))
}
