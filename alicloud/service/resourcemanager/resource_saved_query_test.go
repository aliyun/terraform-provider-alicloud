package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/resourcemanager"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ResourceManager SavedQuery. >>> Resource test cases, automatically generated.
// Case 5104
func TestAccAliCloudResourceManagerSavedQuery_basic5104(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_saved_query.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudResourceManagerSavedQueryMap5104)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerSavedQuery")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sresourcemanagersavedquery%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudResourceManagerSavedQueryBasicDependence5104)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"expression":       "select * from resources limit 1;",
					"saved_query_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"expression":       "select * from resources limit 1;",
						"saved_query_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expression": "select * from resources limit 1;",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"expression": "select * from resources limit 1;",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"saved_query_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"saved_query_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test_desc_2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test_desc_2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expression": "select",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"expression": "select",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"saved_query_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"saved_query_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      "test_desc",
					"expression":       "select * from resources limit 1;",
					"saved_query_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      "test_desc",
						"expression":       "select * from resources limit 1;",
						"saved_query_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudResourceManagerSavedQueryMap5104 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudResourceManagerSavedQueryBasicDependence5104(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 5104  twin
func TestAccAliCloudResourceManagerSavedQuery_basic5104_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_saved_query.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudResourceManagerSavedQueryMap5104)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerSavedQuery")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sresourcemanagersavedquery%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudResourceManagerSavedQueryBasicDependence5104)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      "test_desc_2",
					"expression":       "select",
					"saved_query_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      "test_desc_2",
						"expression":       "select",
						"saved_query_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test ResourceManager SavedQuery. <<< Resource test cases, automatically generated.
