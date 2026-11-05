package cloudauth_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCloudauthFaceConfig_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testaccCAfaceconfig%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// Cloudauth is Hangzhou-only. The flat-suite gate
			// testAccPreCheckWithRegions(t, true, CloudAuthSupportRegions)
			// converged to cn-hangzhou for every region input, so pin it here.
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccFaceConfigConfig(name+"-biz", name+"-type"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cloudauth_face_config.default", "biz_name", name+"-biz"),
					resource.TestCheckResourceAttr("alicloud_cloudauth_face_config.default", "biz_type", name+"-type"),
				),
			},
			{
				Config: testAccFaceConfigConfig(name+"update", name+"-type"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cloudauth_face_config.default", "biz_name", name+"update"),
				),
			},
			{
				ResourceName:      "alicloud_cloudauth_face_config.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFaceConfigConfig(bizName, bizType string) string {
	return fmt.Sprintf(`
resource "alicloud_cloudauth_face_config" "default" {
  biz_name = "%s"
  biz_type = "%s"
}
`, bizName, bizType)
}
