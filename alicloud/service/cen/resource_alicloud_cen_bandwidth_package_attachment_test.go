package cen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cen"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCenBandwidthPackageAttachment_basic(t *testing.T) {
	var cenBwp cbn.CenBandwidthPackage

	resourceId := "alicloud_cen_bandwidth_package_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, cenBandwidthPackageAttachmentMap)

	serviceFunc := func() interface{} {
		return cen.NewCenService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &cenBwp, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sCenBandwidthPackageAttachment-%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCenBandwidthPackageAttachmentConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
	"instance_id":          acctest.CHECKSET,
	"bandwidth_package_id": acctest.CHECKSET,
}
