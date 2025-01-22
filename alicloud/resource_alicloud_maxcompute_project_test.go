package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudMaxcomputeProject_basic1968(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_maxcompute_project.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxcomputeProjectMap1968)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxcomputeProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.MaxComputeProjectSupportRegions)
	name := fmt.Sprintf("tf_testaccmp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxcomputeProjectBasicDependence1968)
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
					"default_quota": "默认后付费Quota",
					"project_name":  "${var.name}",
					"comment":       "${var.name}",
					"product_type":  "PayAsYouGo",
					"is_logical":    "false",
					"ip_white_list": []map[string]interface{}{
						{
							"ip_list":     "1.1.1.1,2.2.2.2",
							"vpc_ip_list": "10.10.10.10,11.11.11.11",
						},
					},
					"properties": []map[string]interface{}{
						{
							"allow_full_scan":  "false",
							"enable_decimal2":  "true",
							"retention_days":   "1",
							"sql_metering_max": "0",
							"timezone":         "Asia/Shanghai",
							"type_system":      "2",
							"encryption": []map[string]interface{}{
								{
									"enable":    "true",
									"algorithm": "AESCTR",
									"key":       "f58d854d-7bc0-4a6e-9205-160e10ffedec",
								},
							},
							"table_lifecycle": []map[string]interface{}{
								{
									"type":  "optional",
									"value": "37231",
								},
							},
						},
					},
					"security_properties": []map[string]interface{}{
						{
							"enable_download_privilege":            "false",
							"label_security":                       "true",
							"object_creator_has_access_permission": "true",
							"object_creator_has_grant_permission":  "true",
							"project_protection": []map[string]interface{}{
								{
									"protected": "false",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_quota":         CHECKSET,
						"project_name":          CHECKSET,
						"comment":               CHECKSET,
						"product_type":          CHECKSET,
						"ip_white_list.#":       "1",
						"properties.#":          "1",
						"security_properties.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment":      "${var.name}_u",
					"project_name": "${var.name}_u",
					"is_logical":   "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment":      CHECKSET,
						"project_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"is_logical": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "is_logical"},
			},
		},
	})
}

func TestAccAlicloudMaxcomputeProject_basic1968_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_maxcompute_project.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxcomputeProjectMap1968)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxcomputeProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.MaxComputeProjectSupportRegions)
	name := fmt.Sprintf("tf_testaccmp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxcomputeProjectBasicDependence1968)
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
					"default_quota": "默认后付费Quota",
					"project_name":  "${var.name}",
					"comment":       "${var.name}",
					"product_type":  "PayAsYouGo",
					"is_logical":    "false",
					"ip_white_list": []map[string]interface{}{
						{
							"ip_list":     "1.1.1.1,2.2.2.2",
							"vpc_ip_list": "10.10.10.10,11.11.11.11",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_quota":   CHECKSET,
						"project_name":    CHECKSET,
						"comment":         CHECKSET,
						"product_type":    CHECKSET,
						"ip_white_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"is_logical": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_quota": "os_terrform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_quota": "os_terrform",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "is_logical"},
			},
		},
	})
}

var AlicloudMaxcomputeProjectMap1968 = map[string]string{}

func AlicloudMaxcomputeProjectBasicDependence1968(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

var AlicloudMaxComputeProjectMap7168 = map[string]string{
	"status":       CHECKSET,
	"owner":        CHECKSET,
	"project_name": CHECKSET,
	"create_time":  CHECKSET,
	"type":         CHECKSET,
}

func AlicloudMaxComputeProjectBasicDependence7168(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case Terraform发布#57547636 7168  raw
func TestAccAliCloudMaxComputeProject_basic7168_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_maxcompute_project.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeProjectMap7168)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccmp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeProjectBasicDependence7168)
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
					"default_quota": "默认后付费Quota",
					"project_name":  name,
					"comment":       "terraform测试项目",
					"properties": []map[string]interface{}{
						{
							"type_system":      "2",
							"sql_metering_max": "10240",
							"encryption": []map[string]interface{}{
								{
									"enable":    "true",
									"algorithm": "AESCTR",
									"key":       "f58d854d-7bc0-4a6e-9205-160e10ffedec",
								},
							},
						},
					},
					"status": "AVAILABLE",
					"ip_white_list": []map[string]interface{}{
						{
							"ip_list":     "10.0.0.0/8",
							"vpc_ip_list": "10.0.0.0/8",
						},
					},
					"security_properties": []map[string]interface{}{
						{
							"using_acl":                            "false",
							"using_policy":                         "false",
							"object_creator_has_access_permission": "false",
							"object_creator_has_grant_permission":  "false",
							"label_security":                       "false",
							"enable_download_privilege":            "false",
							"project_protection": []map[string]interface{}{
								{
									"protected":        "true",
									"exception_policy": "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"odps:*\\\"],\\\"Resource\\\":[\\\"acs:odps:*:projects/ludong/tables/*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":[\\\"ALIYUN$ludong@aliyun.com\\\"]}]}",
								},
							},
						},
					},
					"tags": map[string]string{
						"Created": "TF-CI",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_quota": "默认后付费Quota",
						"project_name":  name,
						"comment":       "terraform测试项目",
						"status":        "AVAILABLE",
						"tags.%":        "2",
						"tags.Created":  "TF-CI",
						"tags.For":      "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "READONLY",
					"ip_white_list": []map[string]interface{}{
						{
							"ip_list":     "0.0.0.0/0",
							"vpc_ip_list": "0.0.0.0/0",
						},
					},
					"security_properties": []map[string]interface{}{
						{
							"using_acl":                            "true",
							"using_policy":                         "true",
							"object_creator_has_access_permission": "true",
							"object_creator_has_grant_permission":  "true",
							"label_security":                       "true",
							"enable_download_privilege":            "true",
							"project_protection": []map[string]interface{}{
								{
									"protected":        "true",
									"exception_policy": "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"odps:*\\\"],\\\"Resource\\\":[\\\"acs:odps:*:projects/ludong/tables/*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":[\\\"ALIYUN$ludong@aliyun.com\\\",\\\"ALIYUN$liuhao@aliyun.com\\\"]}]}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "READONLY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "FROZEN",
					"security_properties": []map[string]interface{}{
						{
							"project_protection": []map[string]interface{}{
								{
									"protected": "false",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "FROZEN",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "AVAILABLE",
					"properties": []map[string]interface{}{
						{
							"enable_decimal2":  "false",
							"type_system":      "1",
							"sql_metering_max": "1024",
							"retention_days":   "2",
							"allow_full_scan":  "false",
							"timezone":         "Asia/Shanghai",
							"table_lifecycle": []map[string]interface{}{
								{
									"type":  "optional",
									"value": "37231",
								},
							},
							"encryption": []map[string]interface{}{
								{
									"enable":    "true",
									"algorithm": "AESCTR",
									"key":       "f58d854d-7bc0-4a6e-9205-160e10ffedec",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "AVAILABLE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type"},
			},
		},
	})
}
