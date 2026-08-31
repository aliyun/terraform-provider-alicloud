package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Quotas TemplateQuota. >>> Resource test cases, automatically generated.
// Case 3099
func TestAccAlicloudQuotasTemplateQuota_basic3099(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_template_quota.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasTemplateQuotaMap3099)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasTemplateQuota")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotastemplatequota%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuotasTemplateQuotaBasicDependence3099)
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
					"quota_action_code": "q_desktop-count",
					"product_code":      "gws",
					"dimensions": []map[string]interface{}{
						{
							"key":   "regionId",
							"value": "cn-hangzhou",
						},
					},
					"desire_value":   "1002",
					"quota_category": "CommonQuota",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code": "q_desktop-count",
						"product_code":      "gws",
						"desire_value":      "1002",
						"dimensions.#":      "1",
						"quota_category":    "CommonQuota",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notice_type": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notice_type": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env_language": "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_language": "zh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_category": "CommonQuota",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_category": "CommonQuota",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desire_value": "1001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desire_value": "1001",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desire_value": "1002",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desire_value": "1002",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_action_code": "q_desktop-count",
					"product_code":      "gws",
					"notice_type":       "3",
					"dimensions": []map[string]interface{}{
						{
							"key":   "regionId",
							"value": "cn-hangzhou",
						},
					},
					"desire_value":   "1001",
					"env_language":   "zh",
					"quota_category": "CommonQuota",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code": "q_desktop-count",
						"product_code":      "gws",
						"notice_type":       "3",
						"dimensions.#":      "1",
						"desire_value":      "1001",
						"env_language":      "zh",
						"quota_category":    "CommonQuota",
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

var AlicloudQuotasTemplateQuotaMap3099 = map[string]string{
	"env_language": CHECKSET,
	"notice_type":  CHECKSET,
}

func AlicloudQuotasTemplateQuotaBasicDependence3099(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3298
func TestAccAlicloudQuotasTemplateQuota_basic3298(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_template_quota.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasTemplateQuotaMap3298)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasTemplateQuota")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotastemplatequota%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuotasTemplateQuotaBasicDependence3298)
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
					"quota_action_code": "vpc_whitelist/ha_vip_whitelist",
					"desire_value":      "1",
					"product_code":      "vpc",
					"quota_category":    "WhiteListLabel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code": "vpc_whitelist/ha_vip_whitelist",
						"desire_value":      "1",
						"product_code":      "vpc",
						"quota_category":    "WhiteListLabel",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective_time": "2023-05-22T16:00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective_time": "2023-05-22T16:00:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_category": "WhiteListLabel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_category": "WhiteListLabel",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notice_type": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notice_type": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expire_time": "2023-05-30T16:00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"expire_time": "2023-05-30T16:00:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env_language": "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_language": "zh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desire_value": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desire_value": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_action_code": "vpc_whitelist/ha_vip_whitelist",
					"effective_time":    "2023-05-22T16:00:00Z",
					"quota_category":    "WhiteListLabel",
					"notice_type":       "3",
					"expire_time":       "2023-05-30T16:00:00Z",
					"desire_value":      "1",
					"env_language":      "zh",
					"product_code":      "vpc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code": "vpc_whitelist/ha_vip_whitelist",
						"effective_time":    "2023-05-22T16:00:00Z",
						"quota_category":    "WhiteListLabel",
						"notice_type":       "3",
						"expire_time":       "2023-05-30T16:00:00Z",
						"desire_value":      "1",
						"env_language":      "zh",
						"product_code":      "vpc",
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

var AlicloudQuotasTemplateQuotaMap3298 = map[string]string{
	"env_language": CHECKSET,
	"notice_type":  CHECKSET,
}

func AlicloudQuotasTemplateQuotaBasicDependence3298(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3099  twin
func TestAccAlicloudQuotasTemplateQuota_basic3099_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_template_quota.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasTemplateQuotaMap3099)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasTemplateQuota")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotastemplatequota%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuotasTemplateQuotaBasicDependence3099)
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
					"quota_action_code": "q_desktop-count",
					"product_code":      "gws",
					"notice_type":       "3",
					"dimensions": []map[string]interface{}{
						{
							"key":   "regionId",
							"value": "cn-hangzhou",
						},
					},
					"desire_value":   "1002",
					"env_language":   "zh",
					"quota_category": "CommonQuota",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code": "q_desktop-count",
						"product_code":      "gws",
						"notice_type":       "3",
						"dimensions.#":      "1",
						"desire_value":      "1002",
						"env_language":      "zh",
						"quota_category":    "CommonQuota",
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

// Case 3298  twin
func TestAccAlicloudQuotasTemplateQuota_basic3298_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_template_quota.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasTemplateQuotaMap3298)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasTemplateQuota")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotastemplatequota%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuotasTemplateQuotaBasicDependence3298)
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
					"quota_action_code": "vpc_whitelist/ha_vip_whitelist",
					"effective_time":    "2023-05-22T16:00:00Z",
					"quota_category":    "WhiteListLabel",
					"notice_type":       "3",
					"expire_time":       "2023-05-30T16:00:00Z",
					"desire_value":      "1",
					"env_language":      "zh",
					"product_code":      "vpc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code": "vpc_whitelist/ha_vip_whitelist",
						"effective_time":    "2023-05-22T16:00:00Z",
						"quota_category":    "WhiteListLabel",
						"notice_type":       "3",
						"expire_time":       "2023-05-30T16:00:00Z",
						"desire_value":      "1",
						"env_language":      "zh",
						"product_code":      "vpc",
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

// Test Quotas TemplateQuota. <<< Resource test cases, automatically generated.
