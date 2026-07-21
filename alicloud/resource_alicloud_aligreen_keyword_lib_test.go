package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Aligreen KeywordLib. >>> Resource test cases, automatically generated.
// Case 文本库 7323
func TestAccAliCloudAligreenKeywordLib_basic7323(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_keyword_lib.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenKeywordLibMap7323)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenKeywordLib")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc_ag%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenKeywordLibBasicDependence7323)
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
	ra := resourceAttrInit(resourceId, AlicloudAligreenKeywordLibMap7323)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenKeywordLib")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc_ag%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenKeywordLibBasicDependence7323)
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
	ra := resourceAttrInit(resourceId, AlicloudAligreenKeywordLibMap7323)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenKeywordLib")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc_ag%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenKeywordLibBasicDependence7323)
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

// Test Aligreen KeywordLib. <<< Resource test cases, automatically generated.
