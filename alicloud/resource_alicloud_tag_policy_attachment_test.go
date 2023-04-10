package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccAlicloudTagPolicyAttachment_basic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alicloud_tag_policy_attachment.default"
	basicMap := map[string]string{
		"policy_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssAlbServerGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTagPolicyAttachment("", rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTagPolicyAttachment(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccTagPolicyAttachment-%d"
	}
	resource "alicloud_tag_policy_attachment" "example" {
	  policy_id         = "p-e707a329d3c047XXX"
	  target_id         = "151266687691****"  
	  target_type       = "USER"
	}
`, common, rand)
}
