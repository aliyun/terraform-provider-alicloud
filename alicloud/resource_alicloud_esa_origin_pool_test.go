package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA OriginPool. >>> Resource test cases, automatically generated.
// Case originpool_test
func TestAccAliCloudESAOriginPooloriginpool_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_origin_pool.default"
	ra := resourceAttrInit(resourceId, AliCloudESAOriginPooloriginpool_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaOriginPool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("bcd%d.com", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAOriginPooloriginpool_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"origins": []map[string]interface{}{

						{
							"type":    "OSS",
							"address": "zhouwei123.oss-cn-beijing.aliyuncs.com",
							"header":  "{\\\"Host\\\":[\\\"zhouwei123.oss-cn-beijing.aliyuncs.com\\\"]}",
							"enabled": "true",
							"auth_conf": []map[string]interface{}{
								{
									"secret_key": "bd8tjba5lXxxxxiRXFIBvoCIfJIL2WJ",
									"auth_type":  "private_cross_account",
									"access_key": "LTAI5tGLgmPe1wFwpX8645BF",
								},
							},
							"weight": "50",
							"name":   "origin1",
						},

						{
							"type":    "S3",
							"address": "test.s3.com",
							"header":  "{\\\"Host\\\":[\\\"example1.com\\\"]}",
							"enabled": "true",
							"auth_conf": []map[string]interface{}{
								{
									"secret_key": "bd8tjba5lXxxxxiRXFIBvoCIfJIL2WJ",
									"version":    "v2",
									"region":     "us-east-1",
									"auth_type":  "private",
									"access_key": "LTAI5tGLgmPe1wFwpX8645BF",
								},
							},
							"weight": "50",
							"name":   "origin2",
						},

						{
							"type":    "S3",
							"address": "test1111.s3.com",
							"header":  "{\\\"Host\\\":[\\\"example1111.com\\\"]}",
							"enabled": "true",
							"auth_conf": []map[string]interface{}{
								{
									"secret_key": "bd8tjba5lXxxxxiRXFIBvoCIfJIL2WJ",
									"version":    "v2",
									"region":     "us-east-1",
									"auth_type":  "private",
									"access_key": "LTAI5tGLgmPe1wFwpX8645BF",
								},
							},
							"weight": "30",
							"name":   "origin3",
						},
					},
					"site_id":          "${alicloud_esa_site.resource_Site_OriginPool_test.id}",
					"origin_pool_name": "testoriginpool",
					"enabled":          "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"origins": []map[string]interface{}{

						{
							"type":    "OSS",
							"address": "zhouwei456.oss-cn-beijing.aliyuncs.com",
							"header":  "{\\\"Host\\\":[\\\"zhouwei456.oss-cn-beijing.aliyuncs.com\\\"]}",
							"enabled": "false",
							"auth_conf": []map[string]interface{}{
								{
									"secret_key": "bd8tjba5lXxxxxiRXFIBvoCIfJIL2WJ",
									"auth_type":  "private_cross_account",
									"access_key": "LTAI5tGLgmPe1wFwpX8645BF",
								},
							},
							"weight": "60",
							"name":   "origin11",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"origins": []map[string]interface{}{

						{
							"type":    "S3",
							"address": "test11.s3.com",
							"header":  "{\\\"Host\\\":[\\\"example11.com\\\"]}",
							"enabled": "false",
							"auth_conf": []map[string]interface{}{
								{
									"secret_key": "bd8tjba5lXxxxxiRXFIBvoCIfJIL2WJ",
									"version":    "v4",
									"region":     "us-east-11",
									"auth_type":  "private",
									"access_key": "LTAI5tGLgmPe1wFwpX8645BF",
								},
							},
							"weight": "70",
							"name":   "origin22",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"origins"},
			},
		},
	})
}

var AliCloudESAOriginPooloriginpool_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAOriginPooloriginpool_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_OriginPool_test" {
  site_name   = var.name
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Test ESA OriginPool. <<< Resource test cases, automatically generated.
