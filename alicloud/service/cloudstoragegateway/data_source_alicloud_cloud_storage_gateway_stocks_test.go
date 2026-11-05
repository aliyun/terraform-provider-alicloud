package cloudstoragegateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCloudStorageGatewayStocksDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	allConf := acctest.DataSourceTestAccConfig{

		ExistConfig: testAccCheckAlicloudCloudStorageGatewayStocksSourceConfig(rand, map[string]string{
			"gateway_class": `"Advanced"`,
		}),
		FakeConfig: "",
	}

	var existCloudStorageGatewayStocksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"stocks.#":                             acctest.CHECKSET,
			"stocks.0.zone_id":                     acctest.CHECKSET,
			"stocks.0.available_gateway_classes.#": acctest.CHECKSET,
		}
	}

	var fakeCloudStorageGatewayStocksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"stocks.#": "0",
		}
	}

	var CloudStorageGatewayStocksRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cloud_storage_gateway_stocks.default",
		ExistMapFunc: existCloudStorageGatewayStocksMapFunc,
		FakeMapFunc:  fakeCloudStorageGatewayStocksMapFunc,
	}

	CloudStorageGatewayStocksRecordsCheckInfo.DataSourceTestCheck(t, rand, allConf)

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
