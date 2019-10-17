package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCenInstanceAttachment_importBasic(t *testing.T) {
	resourceName := "alicloud_cen_instance_attachment.default"
	rand := acctest.RandIntRange(1000000, 9999999)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenInstanceAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCenInstanceAttachmentBasic(rand, defaultRegionToTest),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCenInstanceAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	cenService := CenService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_instance_attachment" {
			continue
		}

		instance, err := cenService.DescribeCenInstanceAttachment(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("CEN %s child instance %s still attach", instance.CenId, instance.ChildInstanceId)
	}

	return nil
}
