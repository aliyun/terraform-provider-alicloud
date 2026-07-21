package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ResourceManager SavedQuery. >>> Resource test cases, automatically generated.
// Case 5104
func TestAccAliCloudResourceManagerSavedQuery_basic5104(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_saved_query.default"
	ra := resourceAttrInit(resourceId, AlicloudResourceManagerSavedQueryMap5104)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerSavedQuery")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sresourcemanagersavedquery%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerSavedQueryBasicDependence5104)
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
	"create_time": CHECKSET,
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
	ra := resourceAttrInit(resourceId, AlicloudResourceManagerSavedQueryMap5104)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerSavedQuery")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sresourcemanagersavedquery%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerSavedQueryBasicDependence5104)
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
