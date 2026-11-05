package aligreen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	aligreen "github.com/aliyun/terraform-provider-alicloud/alicloud/service/aligreen"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Aligreen KeywordLib. >>> Resource test cases, automatically generated.
// Case 文本库 7323
func TestAccAliCloudAligreenKeywordLib_basic7323(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_keyword_lib.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenKeywordLibMap7323)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenKeywordLib")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc_ag%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenKeywordLibBasicDependence7323)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_type":    "TEXT",
					"keyword_lib_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_type":    "TEXT",
						"keyword_lib_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"keyword_lib_name": name + "_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"keyword_lib_name": name + "_u",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"category":         "BLACK",
					"resource_type":    "TEXT",
					"lib_type":         "textKeyword",
					"keyword_lib_name": name + "_u",
					"match_mode":       "fuzzy",
					"language":         "cn",
					"biz_types": []string{
						"test_007"},
					"lang":   "cn",
					"enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":         "BLACK",
						"resource_type":    "TEXT",
						"lib_type":         "textKeyword",
						"keyword_lib_name": name + "_u",
						"match_mode":       "fuzzy",
						"language":         "cn",
						"biz_types.#":      "1",
						"lang":             "cn",
						"enable":           "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

var AlicloudAligreenKeywordLibMap7323 = map[string]string{}

func AlicloudAligreenKeywordLibBasicDependence7323(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_aligreen_biz_type" "defaultMn8sVK" {
  biz_type_name = var.name
  cite_template = true
  industry_info = "社交-注册信息-昵称"
}


`, name)
}

// Case 文本库 7323  twin
func TestAccAliCloudAligreenKeywordLib_basic7323_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_keyword_lib.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenKeywordLibMap7323)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenKeywordLib")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc_ag%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenKeywordLibBasicDependence7323)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"category":         "BLACK",
					"resource_type":    "TEXT",
					"lib_type":         "textKeyword",
					"keyword_lib_name": name,
					"match_mode":       "fuzzy",
					"language":         "cn",
					"biz_types": []string{
						"test_007"},
					"lang":   "cn",
					"enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":         "BLACK",
						"resource_type":    "TEXT",
						"lib_type":         "textKeyword",
						"keyword_lib_name": name,
						"match_mode":       "fuzzy",
						"language":         "cn",
						"biz_types.#":      "1",
						"lang":             "cn",
						"enable":           "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

// Case 文本库 7323  raw
func TestAccAliCloudAligreenKeywordLib_basic7323_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_keyword_lib.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenKeywordLibMap7323)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenKeywordLib")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc_ag%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenKeywordLibBasicDependence7323)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"category":         "BLACK",
					"resource_type":    "TEXT",
					"lib_type":         "textKeyword",
					"keyword_lib_name": name,
					"match_mode":       "fuzzy",
					"language":         "cn",
					"biz_types": []string{
						"test_007"},
					"lang":   "cn",
					"enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":         "BLACK",
						"resource_type":    "TEXT",
						"lib_type":         "textKeyword",
						"keyword_lib_name": name,
						"match_mode":       "fuzzy",
						"language":         "cn",
						"biz_types.#":      "1",
						"lang":             "cn",
						"enable":           "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"keyword_lib_name": name + "_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"keyword_lib_name": name + "_u",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

// Case 文本库 collection-order coverage
// Covers the configurable multi-member TypeList attribute biz_types (ForceNew):
// apply the original member order, then a convergent apply with the members
// reordered. Schema ForceNew triggers a replacement, which the collection-order
// checker allows.
func TestAccAliCloudAligreenKeywordLib_collectionOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_keyword_lib.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenKeywordLibMap7323)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenKeywordLib")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc_ag%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenKeywordLibBasicDependence7323)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_type":    "TEXT",
					"keyword_lib_name": name,
					"biz_types": []string{
						"test_007", "test_008"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_type":    "TEXT",
						"keyword_lib_name": name,
						"biz_types.#":      "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"biz_types": []string{
						"test_008", "test_007"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_types.#": "2",
					}),
				),
			},
		},
	})
}

// Test Aligreen KeywordLib. <<< Resource test cases, automatically generated.
