package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCenBandwidthPackageAttachment_basic(t *testing.T) {
	var cenBwp cbn.CenBandwidthPackage

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cen_bandwidth_package_attachment.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCenBandwidthPackageAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCenBandwidthPackageAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthPackageAttachmentExists("alicloud_cen_bandwidth_package_attachment.foo", &cenBwp),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cenService := CenService{client}

		instance, err := cenService.DescribeCenBandwidthPackageById(rs.Primary.ID)
		if err != nil {
			return err
		}

		*cenBwp = instance
		return nil
	}
}

func testAccCheckCenBandwidthPackageAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	cenService := CenService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_bandwidth_package_attachment" {
			continue
		}

		instance, err := cenService.DescribeCenBandwidthPackageById(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("CEN %s bandwidth package %s still attach", instance.CenIds.CenId[0], instance.CenBandwidthPackageId)
		}
	}

	return nil
}

const testAccCenBandwidthPackageAttachmentConfig = `
resource "alicloud_cen_instance" "cen" {
     name = "tf-testAccCenBandwidthPackageAttachmentConfig"
     description = "tf-testAccCenBandwidthPackageAttachmentDescription"
}

resource "alicloud_cen_bandwidth_package" "bwp" {
    bandwidth = 20
    geographic_region_ids = [
		"China",
		"Asia-Pacific"]
}

resource "alicloud_cen_bandwidth_package_attachment" "foo" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    bandwidth_package_id = "${alicloud_cen_bandwidth_package.bwp.id}"
}
`
