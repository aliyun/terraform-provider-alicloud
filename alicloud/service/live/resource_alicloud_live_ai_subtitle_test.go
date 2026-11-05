package live_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/live"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudLiveAISubtitle_basic0(t *testing.T) {
	resourceId := "alicloud_live_ai_subtitle.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	name := fmt.Sprintf("tf-testacc%sliveaisub%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckLiveAISubtitleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudLiveAISubtitleBasic0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "subtitle_name", name),
					resource.TestCheckResourceAttr(resourceId, "src_language", "zh-CN"),
					resource.TestCheckResourceAttr(resourceId, "dst_language", "en-US"),
					resource.TestCheckResourceAttr(resourceId, "font_name", "KaiTi"),
					resource.TestCheckResourceAttr(resourceId, "word_per_line", "20"),
					resource.TestCheckResourceAttr(resourceId, "font_size_normalized", "0.037"),
					resource.TestCheckResourceAttr(resourceId, "max_lines", "2"),
					resource.TestCheckResourceAttr(resourceId, "description", "tf-testacc live ai subtitle"),
					// position_normalized is a list; verify both coordinates refresh from Describe.
					resource.TestCheckResourceAttr(resourceId, "position_normalized.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "position_normalized.0", "0.32"),
					resource.TestCheckResourceAttr(resourceId, "position_normalized.1", "0.27"),
				),
			},
			{
				Config: testAccAliCloudLiveAISubtitleBasic0Update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "dst_language", "ru-RU"),
					resource.TestCheckResourceAttr(resourceId, "max_lines", "3"),
					resource.TestCheckResourceAttr(resourceId, "description", "tf-testacc updated"),
					resource.TestCheckResourceAttr(resourceId, "word_per_line", "25"),
					// Unchanged fields stay as configured after the update refresh.
					resource.TestCheckResourceAttr(resourceId, "src_language", "zh-CN"),
					resource.TestCheckResourceAttr(resourceId, "font_name", "KaiTi"),
					resource.TestCheckResourceAttr(resourceId, "position_normalized.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "position_normalized.0", "0.32"),
					resource.TestCheckResourceAttr(resourceId, "position_normalized.1", "0.27"),
					resource.TestCheckResourceAttr(resourceId, "font_size_normalized", "0.037"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"copy_from"},
			},
		},
	})
}

func testAccCheckLiveAISubtitleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_live_ai_subtitle" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := live.NewService(client).DescribeLiveAISubtitle(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing live ai subtitle %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("live ai subtitle %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAliCloudLiveAISubtitleBasic0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_live_ai_subtitle" "default" {
  subtitle_name        = var.name
  src_language         = "zh-CN"
  dst_language         = "en-US"
  font_name            = "KaiTi"
  position_normalized   = [0.32, 0.27]
  word_per_line        = 20
  font_size_normalized = 0.037
  max_lines            = 2
  description          = "tf-testacc live ai subtitle"
}
`, name)
}

func testAccAliCloudLiveAISubtitleBasic0Update(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_live_ai_subtitle" "default" {
  subtitle_name        = var.name
  src_language         = "zh-CN"
  dst_language         = "ru-RU"
  font_name            = "KaiTi"
  position_normalized   = [0.32, 0.27]
  word_per_line        = 25
  font_size_normalized = 0.037
  max_lines            = 3
  description          = "tf-testacc updated"
}
`, name)
}
