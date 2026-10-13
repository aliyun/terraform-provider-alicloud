package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DataWorks DataAssetTag. >>> Resource test cases, automatically generated.
func TestAccAliCloudDataWorksDataAssetTag_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_data_asset_tag.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksDataAssetTagMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksDataAssetTag")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksDataAssetTagBasicDependence)
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
					"key":         name,
					"value_type":  "String",
					"description": "tf_desc",
					"values":      []string{"value1", "value2"},
					"managers":    []string{"manager1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key":         name,
						"value_type":  "String",
						"description": "tf_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf_desc_update",
					"values":      []string{"value1", "value2", "value3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf_desc_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"managers": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"category", "create_time", "modify_time"},
			},
		},
	})
}

var AlicloudDataWorksDataAssetTagMap = map[string]string{}

func AlicloudDataWorksDataAssetTagBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// Test DataWorks DataAssetTag. <<< Resource test cases, automatically generated.
