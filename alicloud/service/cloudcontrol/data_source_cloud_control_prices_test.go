package cloudcontrol_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// CloudControl price data sources are Hangzhou-only; the flat suite's
// testAccPreCheckWithRegions converged to cn-hangzhou, so pin it here.
func TestAccAliCloudCloudControlPriceDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCloudControlPriceSourceConfig(rand, map[string]string{
					"product":       `"SLB"`,
					"resource_code": `"LoadBalancer"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cloud_control_prices.default", "prices.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_prices.default", "prices.0.original_price"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_prices.default", "prices.0.discount_price"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_prices.default", "prices.0.currency"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_prices.default", "prices.0.module_details.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_prices.default", "prices.0.promotion_details.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_prices.default", "prices.0.trade_price"),
				),
			},
		},
	})
}

func TestAccAliCloudCloudControlPriceDataSourceJson(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCloudControlPriceSourceJsonConfig(rand, map[string]string{
					"product":       `"ECS"`,
					"resource_code": `"Instance"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cloud_control_prices.default", "prices.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_prices.default", "prices.0.original_price"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_prices.default", "prices.0.discount_price"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_prices.default", "prices.0.currency"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_prices.default", "prices.0.trade_price"),
				),
			},
		},
	})
}

func testAccCheckAlicloudCloudControlPriceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCloudControlPrice%d"
}


data "alicloud_cloud_control_prices" "default" {
    desire_attributes = {
      AddressType = "internet"
      PaymentType = "PayAsYouGo"
    }
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

func testAccCheckAlicloudCloudControlPriceSourceJsonConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCloudControlPrice%d"
}


data "alicloud_cloud_control_prices" "default" {
    desire_attributes_json = jsonencode({
      InstanceType = "ecs.g7.large"
      SystemDisk = {
        Category = "cloud_essd"
      }
    })
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
