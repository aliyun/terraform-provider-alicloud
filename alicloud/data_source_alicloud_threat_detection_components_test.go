package alicloud_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAlicloudThreatDetectionComponentsDataSource covers the data source
// contract end to end through the muxed provider: a broad regex returns the
// full component set, an exact component_name filter returns one match, and a
// filter that matches nothing returns an empty list rather than erroring on
// Component.External.NotExist.
//
// The test is an external test (package alicloud_test) because it drives the
// data source through acctest.ProtoV5ProviderFactories, the muxed (Framework +
// SDK v2) server; package alicloud's own acceptance tests stay internal and use
// the in-package factories, but importing acctest from package alicloud would
// form an import cycle (acctest is a non-test file that imports package alicloud
// for Provider()). The terraform-plugin-testing harness is used, matching the
// other Framework/ephemeral acceptance tests (ims, kms); only one
// helper/resource may be linked into a test binary because both the SDK v2 and
// plugin-testing variants register a -sweep flag in init, and a binary linking
// both panics with "flag redefined: sweep" before any test runs. The SDK v2
// helper/acctest is imported only for RandIntRange and does not register the
// flag, so it coexists with the plugin-testing harness.
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
