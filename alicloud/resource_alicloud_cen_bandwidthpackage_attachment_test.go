package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudCen_BandwidthPackage_Attachment_basic(t *testing.T) {
	var cenBwp cbn.CenBandwidthPackage

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cen_bandwidthpackage_attachment.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCenBandwidthPackageAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenBandwidthPackageAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthPackageAttachmentExists("alicloud_cen_bandwidthpackage_attachment.foo", &cenBwp),
				),
			},
		},
	})
}

func testAccCheckCenBandwidthPackageAttachmentExists(n string, cenBwp *cbn.CenBandwidthPackage) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CenBandwidthPackage ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)

		cenBwpId, cenId, err := getCenIdAndAnotherId(rs.Primary.ID)
		if err != nil {
			return err
		}

		instance, err := client.DescribeCenBandwidthPackageById(cenBwpId, cenId)
		if err != nil {
			return err
		}

		if instance.Status != string(InUse) {
			return fmt.Errorf("CEN id %s CEN bwp id %s status error", cenBwpId, cenId)
		}

		*cenBwp = instance
		return nil
	}
}

func testAccCheckCenBandwidthPackageAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_bandwidthpackage_attachment" {
			continue
		}

		cenBwpId, cenId, err := getCenIdAndAnotherId(rs.Primary.ID)
		if err != nil {
			return err
		}

		instance, err := client.DescribeCenBandwidthPackageById(cenBwpId, cenId)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if instance.CenBandwidthPackageId != "" {
			return fmt.Errorf("CEN Bandwidth Package %s still bind", instance.CenBandwidthPackageId)
		}
	}

	return nil
}

const testAccCenBandwidthPackageAttachmentConfig = `
resource "alicloud_cen" "cen" {
     name = "terraform-01"
     description = "terraform01"
}

resource "alicloud_cen_bandwidthpackage" "bwp" {
    bandwidth = 20
    geographic_region_id = [
		"China",
		"Asia-Pacific"]
}

resource "alicloud_cen_bandwidthpackage_attachment" "foo" {
    cen_id = "${alicloud_cen.cen.id}"
    cen_bandwidthpackage_id = "${alicloud_cen_bandwidthpackage.bwp.id}"
}
`
