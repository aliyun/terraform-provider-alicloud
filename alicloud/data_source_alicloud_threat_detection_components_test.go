package alicloud_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccAlicloudThreatDetectionComponentsDataSource covers the data source
// contract end to end through the muxed provider: a broad regex returns the
// full component set, an exact component_name filter returns one match, and a
// filter that matches nothing returns an empty list rather than erroring on
// Component.External.NotExist.
//
// The test is an external test (package alicloud_test) because it drives the
// data source through acctest.ProtoV5ProviderFactories, the muxed (Framework +
// SDK v2) server; importing acctest from package alicloud would form an
// import cycle (acctest is a non-test file that imports package alicloud for
// Provider()).
//
// The SDK v2 helper/resource harness is used rather than
// terraform-plugin-testing/helper/resource because package alicloud already
// links the SDK v2 helper/resource through non-test files
// (resource_alicloud_polardb_polar_fs*.go); both packages register a -sweep
// flag in init, so a test binary linking both panics with
// "flag redefined: sweep" before any test runs. Staying on the SDK v2 harness
// keeps a single helper/resource in the binary. The SDK v2 helper/acctest is
// imported only for RandIntRange. acctest.PreCheck supplies the exported
// precheck (credentials + region default).
func TestAccAlicloudThreatDetectionComponentsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccThreatDetectionComponentsConfig(map[string]string{
					"name_regex": `".*"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_components.default", "id"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_components.default", "components.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_components.default", "components.0.component_name"),
				),
			},
			{
				Config: testAccThreatDetectionComponentsConfig(map[string]string{
					"component_name": `"Http"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_components.default", "components.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_components.default", "components.0.component_name", "Http"),
				),
			},
			{
				Config: testAccThreatDetectionComponentsConfig(map[string]string{
					"component_name": fmt.Sprintf(`"tf-acc-fake-component-not-exist-%d"`, rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_components.default", "components.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_components.default", "ids.#", "0"),
				),
			},
		},
	})
}

func testAccThreatDetectionComponentsConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := `
data "alicloud_threat_detection_components" "default" {
`
	for _, p := range pairs {
		config += "  " + p + "\n"
	}
	config += `}`
	return config
}
