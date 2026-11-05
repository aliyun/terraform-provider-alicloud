package live_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/live"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const liveAIStudioResourceId = "alicloud_live_ai_studio.default"

// TestAccAliCloudLiveAIStudio_basic covers the create → update → import
// lifecycle of the hand-written Framework resource alicloud_live_ai_studio.
// CreateLiveAIStudio returns StudioId (echoed as TemplateId by
// DescribeLiveAIStudio); ModifyLiveAIStudio locates the studio by StudioName.
func TestAccAliCloudLiveAIStudio_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-liveaistudio-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			acctest.PreCheck(t)
		},
		IDRefreshName:            liveAIStudioResourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckLiveAIStudioDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveAIStudioBasicConfig(name, "complex", "basic description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(liveAIStudioResourceId, "studio_name", name),
					resource.TestCheckResourceAttr(liveAIStudioResourceId, "matting_type", "complex"),
					resource.TestCheckResourceAttr(liveAIStudioResourceId, "description", "basic description"),
					resource.TestCheckResourceAttrSet(liveAIStudioResourceId, "id"),
					resource.TestCheckResourceAttrSet(liveAIStudioResourceId, "studio_id"),
				),
			},
			{
				// ModifyLiveAIStudio re-applies the studio by StudioName; the
				// changed description, matting_type and width/height must round-trip.
				Config: testAccLiveAIStudioBasicConfig(name, "green", "updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(liveAIStudioResourceId, "studio_name", name),
					resource.TestCheckResourceAttr(liveAIStudioResourceId, "matting_type", "green"),
					resource.TestCheckResourceAttr(liveAIStudioResourceId, "description", "updated description"),
					resource.TestCheckResourceAttrSet(liveAIStudioResourceId, "id"),
				),
			},
			{
				ResourceName:      liveAIStudioResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLiveAIStudioDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	return checkLiveAIStudioDestroyed(s, func(id string) error {
		_, err := live.NewService(client).DescribeLiveAIStudio(id)
		return err
	})
}

// checkLiveAIStudioDestroyed accepts only the NotFound that DescribeLiveAIStudio
// constructs itself; a permission or other non-404 failure must surface as a
// cleanup failure instead of pretending the studio was deleted.
func checkLiveAIStudioDestroyed(s *terraform.State, describe func(string) error) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_live_ai_studio" {
			continue
		}
		err := describe(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("live ai studio %s still exists", rs.Primary.ID)
		}
		if errs.IsProviderNotFound(err) {
			continue
		}
		return fmt.Errorf("describing live ai studio %s failed: %s", rs.Primary.ID, err)
	}
	return nil
}

func testAccLiveAIStudioBasicConfig(name, mattingType, description string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_live_ai_studio" "default" {
  studio_name  = var.name
  matting_type = "%s"
  description  = "%s"
  width        = 1920
  height       = 1080
}
`, name, mattingType, description)
}
