package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaBandwidthPackageAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudGaBandwidthPackageAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", AlicloudGaBandwidthPackageAttachmentBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":       "${data.alicloud_ga_accelerators.default.ids.0}",
					"bandwidth_package_id": "${data.alicloud_ga_bandwidth_packages.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":       CHECKSET,
						"bandwidth_package_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accelerator_id"},
			},
		},
	})
}

var AlicloudGaBandwidthPackageAttachmentMap = map[string]string{
	"accelerators.#": CHECKSET,
	"status":         "binded",
}

func AlicloudGaBandwidthPackageAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
data "alicloud_ga_accelerators" "default"{
}
data "alicloud_ga_bandwidth_packages" "default"{
}`)
}
