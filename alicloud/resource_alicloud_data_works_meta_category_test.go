package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccAliCloudDataWorksMetaCategory_basic0 covers the full Create → Read → Update →
// (clear comment) → re-import lifecycle for the alicloud_data_works_meta_category resource.
//
// The re-import step (ImportStateVerify with an empty ImportStateVerifyIgnore list)
// guards the json.Number regression fixed alongside this test: RpcPost responses are
// decoded with json.UseNumber, so CreateTime arrives as json.Number; if Read does not
// convert it back into state, create_time is silently dropped and the import step fails
// with a state diff. The category_id and create_time CHECKSET assertions in the create
// and update steps also require both fields to be reliably back-filled by Read.
func TestAccAliCloudDataWorksMetaCategory_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_meta_category.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksMetaCategoryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksPublicService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksMetaCategory")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-metacategory-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksMetaCategoryBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Create a root-level category with name + comment, exercising the
				// CreateMetaCategory ParentId parameter (mapped from parent_category_id).
				Config: testAccConfig(map[string]interface{}{
					"name":               name,
					"comment":            "tf-meta-category-comment",
					"parent_category_id": 0,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name,
						"comment":            "tf-meta-category-comment",
						"parent_category_id": "0",
						"category_id":        CHECKSET,
						"create_time":        CHECKSET,
					}),
				),
			},
			{
				// Update name and comment; parent_category_id is ForceNew so it stays.
				Config: testAccConfig(map[string]interface{}{
					"name":    fmt.Sprintf("%s-updated", name),
					"comment": "tf-meta-category-comment-updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("%s-updated", name),
						"comment":     "tf-meta-category-comment-updated",
						"category_id": CHECKSET,
						"create_time": CHECKSET,
					}),
				),
			},
			{
				// Clear the optional comment to exercise the empty-string update path.
				Config: testAccConfig(map[string]interface{}{
					"name":    fmt.Sprintf("%s-updated", name),
					"comment": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":    fmt.Sprintf("%s-updated", name),
						"comment": "",
					}),
				),
			},
			{
				// Re-import with no ImportStateVerifyIgnore: every computed attribute
				// (including create_time) must round-trip through Read, which catches the
				// json.Number silent-drop regression in CreateTime back-fill.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudDataWorksMetaCategoryMap0 = map[string]string{
	"category_id":        CHECKSET,
	"create_time":        CHECKSET,
	"parent_category_id": CHECKSET,
}

func AlicloudDataWorksMetaCategoryBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}
