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

// CloudControl product data sources are Hangzhou-only; the flat suite's
// testAccPreCheckWithRegions converged to cn-hangzhou, so pin it here.
func TestAccAlicloudCloudControlProductDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// ids filter, matching
			{
				Config: testAccCheckAlicloudCloudControlProductSourceConfig(rand, map[string]string{
					"ids": `["Live"]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cloud_control_products.default", "products.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_products.default", "products.0.product_name"),
					resource.TestCheckResourceAttrSet("data.alicloud_cloud_control_products.default", "products.0.product_code"),
				),
			},
			// ids filter, no match
			{
				Config: testAccCheckAlicloudCloudControlProductSourceConfig(rand, map[string]string{
					"ids": `["Live_fake"]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cloud_control_products.default", "products.#", "0"),
				),
			},
		},
	})
}

func testAccCheckAlicloudCloudControlProductSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCloudControlProduct%d"
}


data "alicloud_cloud_control_products" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
