package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudStorageGatewayStocksDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	allConf := dataSourceTestAccConfig{

		existConfig: testAccCheckAlicloudCloudStorageGatewayStocksSourceConfig(rand, map[string]string{
			"gateway_class": `"Advanced"`,
		}),
		fakeConfig: "",
	}

	var existCloudStorageGatewayStocksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"stocks.#":                             CHECKSET,
			"stocks.0.zone_id":                     CHECKSET,
			"stocks.0.available_gateway_classes.#": CHECKSET,
		}
	}

	var fakeCloudStorageGatewayStocksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"stocks.#": "0",
		}
	}

	var CloudStorageGatewayStocksRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_storage_gateway_stocks.default",
		existMapFunc: existCloudStorageGatewayStocksMapFunc,
		fakeMapFunc:  fakeCloudStorageGatewayStocksMapFunc,
	}

	CloudStorageGatewayStocksRecordsCheckInfo.dataSourceTestCheck(t, rand, allConf)

}

func testAccCheckAlicloudCloudStorageGatewayStocksSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_cloud_storage_gateway_stocks" "default"{
%s
}
`, strings.Join(pairs, "\n   "))
	return config
}
