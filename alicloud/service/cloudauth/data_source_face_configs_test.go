package cloudauth_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCloudAuthFaceConfigDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testaccCAfaceconfig%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// Cloudauth is Hangzhou-only (CloudAuthSupportRegions); the flat
			// suite gate testAccPreCheckWithRegions converged to cn-hangzhou
			// for every region input, so pin it here.
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// ids filter, matching
			{
				Config: testAccFaceConfigsConfig(name, fmt.Sprintf(`ids = ["${alicloud_cloudauth_face_config.default.id}"]`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "configs.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "configs.0.biz_type", name),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "configs.0.biz_name", name),
				),
			},
			// ids filter, no match
			{
				Config: testAccFaceConfigsConfig(name, fmt.Sprintf(`ids = ["${alicloud_cloudauth_face_config.default.id}_fake"]`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "names.#", "0"),
				),
			},
			// name_regex filter, matching
			{
				Config: testAccFaceConfigsConfig(name, fmt.Sprintf(`name_regex = "${alicloud_cloudauth_face_config.default.biz_name}"`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "configs.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "configs.0.biz_type", name),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "configs.0.biz_name", name),
				),
			},
			// name_regex filter, no match
			{
				Config: testAccFaceConfigsConfig(name, fmt.Sprintf(`name_regex = "${alicloud_cloudauth_face_config.default.biz_name}_fake"`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "names.#", "0"),
				),
			},
			// both filters, matching
			{
				Config: testAccFaceConfigsConfig(name, fmt.Sprintf(`ids = ["${alicloud_cloudauth_face_config.default.id}"]
  name_regex = "${alicloud_cloudauth_face_config.default.biz_name}"`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "configs.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "configs.0.biz_type", name),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "configs.0.biz_name", name),
				),
			},
			// both filters, no match
			{
				Config: testAccFaceConfigsConfig(name, fmt.Sprintf(`ids = ["${alicloud_cloudauth_face_config.default.id}_fake"]
  name_regex = "${alicloud_cloudauth_face_config.default.biz_name}_fake"`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_cloudauth_face_configs.default", "names.#", "0"),
				),
			},
		},
	})
}

func testAccFaceConfigsConfig(name string, filter string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_cloudauth_face_config" "default" {
	biz_type = var.name
	biz_name = var.name
}

data "alicloud_cloudauth_face_configs" "default" {
	%s
}
`, name, filter)
}
