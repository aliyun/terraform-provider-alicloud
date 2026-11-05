package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ESA OriginPool. >>> Resource test cases, automatically generated.
// Case 0
func TestAccAliCloudESAOriginPool_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_origin_pool.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESAOriginPoolMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaOriginPool")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("bcd%d.com", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESAOriginPoolBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":          "${data.alicloud_esa_sites.default.sites.0.id}",
					"origin_pool_name": name,
					"origins": []map[string]interface{}{
						{
							"type":    "S3",
							"address": "test.s3.com",
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
							"name":   name,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":          acctest.CHECKSET,
						"origin_pool_name": name,
						"origins.#":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "true",
					}),
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
							"weight":            "60",
							"name":              "origin11",
							"ip_version_policy": "ipv4_first",
						},
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
							"weight":            "70",
							"name":              "origin22",
							"ip_version_policy": "ipv6_first",
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
							"weight":            "30",
							"name":              "origin3",
							"ip_version_policy": "follow",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"origins.#": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

func TestAccAliCloudESAOriginPool_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_origin_pool.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESAOriginPoolMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaOriginPool")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("bcd%d.com", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESAOriginPoolBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":          "${data.alicloud_esa_sites.default.sites.0.id}",
					"origin_pool_name": name,
					"enabled":          "true",
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
							"weight":            "60",
							"name":              "origin11",
							"ip_version_policy": "ipv4_first",
						},
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
							"weight":            "70",
							"name":              "origin22",
							"ip_version_policy": "ipv6_first",
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
							"weight":            "30",
							"name":              "origin3",
							"ip_version_policy": "follow",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":          acctest.CHECKSET,
						"origin_pool_name": name,
						"enabled":          "true",
						"origins.#":        "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

var AliCloudESAOriginPoolMap0 = map[string]string{
	"enabled":        acctest.CHECKSET,
	"origin_pool_id": acctest.CHECKSET,
}

func AliCloudESAOriginPoolBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}
`, name)
}

// Test ESA OriginPool. <<< Resource test cases, automatically generated.
