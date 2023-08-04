package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudTagMetaTag_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_tag_meta_tag.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	ra := resourceAttrInit(resourceId, TagMetaTagMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &TagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeTagValue")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("testAccTagMetaTag%d", rand)
	nameWithColon := fmt.Sprintf("testAcc:TagMetaTag%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, TagMetaTagBasicdependence)
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
					"key":    name,
					"values": []string{"test"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key":      name,
						"values.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key":    nameWithColon,
					"values": []string{"test"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key":      nameWithColon,
						"values.#": "1",
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

var TagMetaTagMap = map[string]string{}

func TagMetaTagBasicdependence(name string) string {
	return ""
}
