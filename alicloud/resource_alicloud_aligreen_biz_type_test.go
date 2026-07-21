package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Aligreen BizType. >>> Resource test cases, automatically generated.
// Case 规则管理_副本1721974125938 7330
func TestAccAliCloudAligreenBizType_basic7330(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_biz_type.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenBizTypeMap7330)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenBizType")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_aligreen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenBizTypeBasicDependence7330)
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
					"biz_type_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_type_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "资源测试用例",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "资源测试用例",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "更新一下",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "更新一下",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"biz_type_name":   name + "_u",
					"description":     "资源测试用例",
					"cite_template":   "true",
					"industry_info":   "社交-注册信息-昵称",
					"biz_type_import": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_type_name":   name + "_u",
						"description":     "资源测试用例",
						"cite_template":   "true",
						"industry_info":   "社交-注册信息-昵称",
						"biz_type_import": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"biz_type_import"},
			},
		},
	})
}

var AlicloudAligreenBizTypeMap7330 = map[string]string{}

func AlicloudAligreenBizTypeBasicDependence7330(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 规则管理 7322
func TestAccAliCloudAligreenBizType_basic7322(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_biz_type.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenBizTypeMap7322)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenBizType")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_aligreen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenBizTypeBasicDependence7322)
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
					"biz_type_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_type_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "资源测试用例",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "资源测试用例",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "更新一下",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "更新一下",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"biz_type_name":   name + "_u",
					"description":     "资源测试用例",
					"cite_template":   "true",
					"industry_info":   "社交-注册信息-昵称",
					"biz_type_import": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_type_name":   name + "_u",
						"description":     "资源测试用例",
						"cite_template":   "true",
						"industry_info":   "社交-注册信息-昵称",
						"biz_type_import": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"biz_type_import"},
			},
		},
	})
}

var AlicloudAligreenBizTypeMap7322 = map[string]string{}

func AlicloudAligreenBizTypeBasicDependence7322(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 规则管理_副本1721974125938 7330  twin
func TestAccAliCloudAligreenBizType_basic7330_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_biz_type.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenBizTypeMap7330)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenBizType")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_aligreen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenBizTypeBasicDependence7330)
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
					"biz_type_name":   name,
					"description":     "资源测试用例",
					"cite_template":   "true",
					"industry_info":   "社交-注册信息-昵称",
					"biz_type_import": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_type_name":   name,
						"description":     "资源测试用例",
						"cite_template":   "true",
						"industry_info":   "社交-注册信息-昵称",
						"biz_type_import": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"biz_type_import"},
			},
		},
	})
}

// Case 规则管理 7322  twin
func TestAccAliCloudAligreenBizType_basic7322_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_biz_type.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenBizTypeMap7322)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenBizType")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_aligreen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenBizTypeBasicDependence7322)
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
					"biz_type_name":   name,
					"description":     "资源测试用例",
					"cite_template":   "true",
					"industry_info":   "社交-注册信息-昵称",
					"biz_type_import": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_type_name":   name,
						"description":     "资源测试用例",
						"cite_template":   "true",
						"industry_info":   "社交-注册信息-昵称",
						"biz_type_import": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"biz_type_import"},
			},
		},
	})
}

// Case 规则管理_副本1721974125938 7330  raw
func TestAccAliCloudAligreenBizType_basic7330_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_biz_type.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenBizTypeMap7330)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenBizType")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_aligreen%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenBizTypeBasicDependence7330)
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
					"biz_type_name":   name,
					"description":     "资源测试用例",
					"cite_template":   "true",
					"industry_info":   "社交-注册信息-昵称",
					"biz_type_import": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_type_name":   name,
						"description":     "资源测试用例",
						"cite_template":   "true",
						"industry_info":   "社交-注册信息-昵称",
						"biz_type_import": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "更新一下",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "更新一下",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"biz_type_import"},
			},
		},
	})
}

// Case 规则管理 7322  raw
func TestAccAliCloudAligreenBizType_basic7322_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_biz_type.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenBizTypeMap7322)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenBizType")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_aligreenbiztype%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenBizTypeBasicDependence7322)
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
					"biz_type_name":   name,
					"description":     "资源测试用例",
					"cite_template":   "true",
					"industry_info":   "社交-注册信息-昵称",
					"biz_type_import": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_type_name":   name,
						"description":     "资源测试用例",
						"cite_template":   "true",
						"industry_info":   "社交-注册信息-昵称",
						"biz_type_import": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "更新一下",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "更新一下",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"biz_type_import"},
			},
		},
	})
}

// Test Aligreen BizType. <<< Resource test cases, automatically generated.
