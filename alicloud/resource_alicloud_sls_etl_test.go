package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Sls Etl. >>> Resource test cases, automatically generated.
// Case test1 10468
func TestAccAliCloudSlsEtl_basic10468(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_etl.default"
	ra := resourceAttrInit(resourceId, AliCloudSlsEtlMap10468)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsEtl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSlsEtlBasicDependence10468)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project": "${alicloud_log_project.defaulthhAPo6.id}",
					"configuration": []map[string]interface{}{
						{
							"script":   "* | extend a=1",
							"lang":     "SPL",
							"role_arn": "test-role-arn",
							"sink": []map[string]interface{}{
								{
									"name":     "11111",
									"endpoint": "cn-hangzhou-intranet.log.aliyuncs.com",
									"project":  "gy-hangzhou-huolang-1",
									"logstore": "gy-rm2",
									"datasets": []string{
										"__UNNAMED__"},
									"role_arn": "test-role-arn",
								},
							},
							"logstore":  "${alicloud_log_store.defaultzWKLkp.name}",
							"from_time": "1706771697",
							"to_time":   "1738394097",
						},
					},
					"job_name":     "etl-1740472705-185721",
					"display_name": "etl-1740472705-185721",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":         CHECKSET,
						"configuration.#": "1",
						"job_name":        "etl-1740472705-185721",
						"display_name":    "etl-1740472705-185721",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"script":   "* | extend a=2",
							"lang":     "SPL",
							"logstore": "${alicloud_log_store.defaultzWKLkp.name}",
							"role_arn": "test-role-arn-update",
							"parameters": map[string]interface{}{
								"AK": "11111",
								"SK": "22222",
							},
							"sink": []map[string]interface{}{
								{
									"name":     "11111",
									"endpoint": "https://cn-qingdao.log.aliyuncs.com",
									"project":  "hclcn-qingdao",
									"logstore": "test",
									"datasets": []string{
										"__set__"},
									"role_arn": "test-role-arn-update",
								},
							},
							"from_time": "1706771697",
							"to_time":   "1738394097",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"script":   "* | extend c=2",
							"logstore": "${alicloud_log_store.defaultzWKLkp.name}",
							"role_arn": "test-role-arn",
							"sink": []map[string]interface{}{
								{
									"name":     "11111",
									"endpoint": "https://cn-nanjing.log.aliyuncs.com",
									"project":  "hcl-cn-nanjing",
									"logstore": "test",
									"datasets": []string{
										"__UNNAMED__"},
									"role_arn": "test-role-arn",
								},
								{
									"name":     "22222",
									"endpoint": "https://cn-shenzhen.log.aliyuncs.com",
									"project":  "hclcn-shenzhen",
									"logstore": "test",
									"datasets": []string{
										"__UU__"},
									"role_arn": "test-role-arn",
								},
								{
									"name":     "33333",
									"endpoint": "https://cn-wulanchabu.log.aliyuncs.com",
									"project":  "test",
									"logstore": "test",
									"datasets": []string{
										"qqq"},
									"role_arn": "test-role-arn",
								},
								{
									"name":     "44444",
									"endpoint": "https://cn-huhehaote.log.aliyuncs.com",
									"project":  "test1",
									"logstore": "test1",
									"datasets": []string{
										"aaa"},
									"role_arn": "test-role-arn",
								},
								{
									"name":     "55555",
									"endpoint": "https://cn-zhangjiakou.log.aliyuncs.com",
									"project":  "test2",
									"logstore": "test2",
									"datasets": []string{
										"ggg"},
									"role_arn": "test-role-arn",
								},
							},
							"lang":      "SPL",
							"from_time": "1706771697",
							"to_time":   "1738394097",
							"parameters": map[string]interface{}{
								"AK": "333",
								"SK": "444",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration.#": "1",
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

var AliCloudSlsEtlMap10468 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AliCloudSlsEtlBasicDependence10468(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "defaulthhAPo6" {
  description = "terraform-etl-test-312"
  name        = "terraform-etl-test-625"
}

resource "alicloud_log_store" "defaultzWKLkp" {
  hot_ttl          = "8"
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.defaulthhAPo6.id
  name             = "test"
}


`, name)
}

func TestAccAliCloudSlsEtl_basic10468_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_etl.default"
	ra := resourceAttrInit(resourceId, AliCloudSlsEtlMap10468)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsEtl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSlsEtlBasicDependence10468)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project": "${alicloud_log_project.defaulthhAPo6.id}",
					"configuration": []map[string]interface{}{
						{
							"script":   "* | extend a=2",
							"lang":     "SPL",
							"logstore": "${alicloud_log_store.defaultzWKLkp.name}",
							"role_arn": "test-role-arn-update",
							"parameters": map[string]interface{}{
								"AK": "11111",
								"SK": "22222",
							},
							"sink": []map[string]interface{}{
								{
									"name":     "11111",
									"endpoint": "https://cn-qingdao.log.aliyuncs.com",
									"project":  "hclcn-qingdao",
									"logstore": "test",
									"datasets": []string{
										"__set__"},
									"role_arn": "test-role-arn-update",
								},
							},
							"from_time": "1706771697",
							"to_time":   "1738394097",
						},
					},
					"job_name":     "etl-1740472705-185721",
					"display_name": "etl-1740472705-185721",
					"description":  name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":         CHECKSET,
						"configuration.#": "1",
						"job_name":        "etl-1740472705-185721",
						"display_name":    "etl-1740472705-185721",
						"description":     name,
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

// Test Sls Etl. <<< Resource test cases, automatically generated.
