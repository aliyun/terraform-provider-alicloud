package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"testing"
)

func TestAccAliCloudECIOpenAPiImageCacheBasic(t *testing.T) {
	var v eci.DescribeImageCachesImageCache0 //.DescribeImageCachesResponse
	resourceId := "alicloud_eciopenapi_image_cache.default"
	ra := resourceAttrInit(resourceId, testAccCheckKeyValueInMapsForECI)

	serviceFunc := func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("testECIImage%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceECIOpenAPIImageCachesBasicConfigDependence)
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
					"image_cache_name":  "name",
					"images":            []string{"ubuntu_18_04_x64_20G_alibase_20200220.vhd","centos_6_10_x64_20G_alibase_20200214.vhd"},
					"vswitch_id":        "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"retention_days":    "7",
					"image_cache_size":  "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_days":   "7",
						"image_cache_size": "20",
					}),
				),
			},
		},
	})
}

var testAccCheckKeyValueInMapsForECI = map[string]string{
	"image_cache_name":  CHECKSET,
	"images.#":          CHECKSET,
	"vswitch_id":        CHECKSET,
	"security_group_id": CHECKSET,
}

func resourceECIOpenAPIImageCachesBasicConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_vpcs" "default" {
  is_default = true
}
data "alicloud_vswitches" "default" {
  ids = [data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0]
}
resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}
`, name)
}
