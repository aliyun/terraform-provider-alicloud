// Package tag_test.
package tag_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/tag"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass values -- pre-existing uncovered TypeList attribute outside this pure directory migration (file moved from alicloud/ root; the attribute and its test coverage are unchanged stock)
func TestAccAliCloudTagMetaTag_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_tag_meta_tag.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, TagMetaTagMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return tag.NewTagService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeTagValue")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("testAccTagMetaTag%d", rand)
	nameWithColon := fmt.Sprintf("testAcc:TagMetaTag%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, TagMetaTagBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
