package alicloud

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
// The test lives in the top-level alicloud package and uses the SDK v2 test
// harness (helper/resource) like the other 1200+ acceptance tests in this
// provider. ProtoV5ProviderFactories still serves the muxed (Framework + SDK v2)
// provider so the Framework-native data source is reachable; importing
// terraform-plugin-testing/helper/resource alongside the SDK v2 one re-registers
// the -sweep flag and panics on init, so the SDK v2 harness is used here.
func TestAccAlicloudThreatDetectionComponentsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
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
