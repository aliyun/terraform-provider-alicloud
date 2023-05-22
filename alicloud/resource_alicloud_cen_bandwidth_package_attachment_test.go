package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCenBandwidthPackageAttachment_basic(t *testing.T) {
	var cenBwp cbn.CenBandwidthPackage

	resourceId := "alicloud_cen_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, cenBandwidthPackageAttachmentMap)

	serviceFunc := func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &cenBwp, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCenBandwidthPackageAttachment-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenBandwidthPackageAttachmentConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":          "${alicloud_cen_instance.default.id}",
					"bandwidth_package_id": "${alicloud_cen_bandwidth_package.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(cenBandwidthPackageAttachmentMap),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Skip this testcase because of the account cannot purchase non-internal products.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"instance_id":          "${alicloud_cen_instance.default.id}",
			//		"bandwidth_package_id": "${alicloud_cen_bandwidth_package.default.id}",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(cenBandwidthPackageAttachmentMap),
			//		testAccCheckCenBandwidthPackageRegionId(&cenBwp, "China", "Asia-Pacific"),
			//	),
			//},
		},
	})
}

func resourceCenBandwidthPackageAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_cen_instance" "default" {
  		name        = "%s"
  		description = "tf-testAccCenBandwidthPackageAttachmentDescription"
	}

	resource "alicloud_cen_bandwidth_package" "default" {
  		cen_bandwidth_package_name = "%s"
  		bandwidth                  = 5
  		geographic_region_a_id     = "China"
  		geographic_region_b_id     = "China"
	}
`, name, name)
}

var cenBandwidthPackageAttachmentMap = map[string]string{
	"instance_id":          CHECKSET,
	"bandwidth_package_id": CHECKSET,
}
