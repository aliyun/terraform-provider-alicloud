// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test OpenApiExplorer ApiMcpServer. >>> Resource test cases, automatically generated.
// Case resource_ApiMcpServer_test 11977
func TestAccAliCloudOpenApiExplorerApiMcpServer_basic11977(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_open_api_explorer_api_mcp_server.default"
	ra := resourceAttrInit(resourceId, AlicloudOpenApiExplorerApiMcpServerMap11977)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OpenApiExplorerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOpenApiExplorerApiMcpServer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOpenApiExplorerApiMcpServerBasicDependence11977)
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
					"system_tools": []string{
						"FetchRamActionDetails"},
					"description": "创建",
					"prompts": []map[string]interface{}{
						{
							"description": "获取用户定制信息描述",
							"content":     "prompt正文，{{name}}",
							"arguments": []map[string]interface{}{
								{
									"description": "名称信息",
									"required":    "true",
									"name":        "name",
								},
							},
							"name": "获取用户定制信息",
						},
						{
							"description": "获取用户定制信息描述",
							"content":     "prompt正文，{{name}}, {{age}}, {{description}}",
							"arguments": []map[string]interface{}{
								{
									"description": "名称信息",
									"required":    "true",
									"name":        "name",
								},
								{
									"description": "年龄信息",
									"required":    "true",
									"name":        "age",
								},
								{
									"description": "描述信息",
									"required":    "true",
									"name":        "description",
								},
							},
							"name": "获取用户定制信息1",
						},
					},
					"oauth_client_id": "123456",
					"apis": []map[string]interface{}{
						{
							"api_version": "2014-05-26",
							"product":     "Ecs",
							"selectors": []string{
								"DescribeAvailableResource", "DescribeRegions", "DescribeZones"},
						},
						{
							"api_version": "2017-03-21",
							"product":     "vod",
							"selectors": []string{
								"CreateUploadVideo"},
						},
						{
							"api_version": "2014-05-15",
							"product":     "Slb",
							"selectors": []string{
								"DescribeAvailableResource"},
						},
					},
					"instructions": "介绍整个MCP Server的作用",
					"additional_api_descriptions": []map[string]interface{}{
						{
							"api_version":          "2014-05-26",
							"enable_output_schema": "true",
							"api_name":             "DescribeAvailableResource",
							"const_parameters": []map[string]interface{}{
								{
									"value": "cn-hangzhou",
									"key":   "x_mcp_region_id",
								},
								{
									"value": "B1",
									"key":   "a1",
								},
								{
									"value": "b2",
									"key":   "a2",
								},
							},
							"api_override_json":   "{\\\"summary\\\": \\\"本接口支持根据不同请求条件查询实例列表，并关联查询实例的详细信息。\\\"}",
							"product":             "Ecs",
							"execute_cli_command": "false",
						},
						{
							"api_version":          "2014-05-26",
							"enable_output_schema": "true",
							"api_name":             "DescribeRegions",
							"product":              "Ecs",
							"execute_cli_command":  "true",
						},
						{
							"api_version":          "2014-05-26",
							"enable_output_schema": "true",
							"api_name":             "DescribeZones",
							"product":              "Ecs",
							"execute_cli_command":  "true",
						},
					},
					"vpc_whitelists": []string{
						"vpc-testa", "vpc-testb", "vpc-testc"},
					"name":                     name,
					"language":                 "ZH_CN",
					"enable_assume_role":       "true",
					"assume_role_extra_policy": "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Effect\\\":\\\"Allow\\\",\\\"Action\\\":[\\\"ecs:Describe*\\\",\\\"vpc:Describe*\\\",\\\"vpc:List*\\\"],\\\"Resource\\\":\\\"*\\\"}]}",
					"terraform_tools": []map[string]interface{}{
						{
							"description":    "terraform as tool测试",
							"async":          "true",
							"destroy_policy": "NEVER",
							"code":           "variable \\\"name\\\" {\\n  default = \\\"terraform-example\\\"\\n}\\n\\nprovider \\\"alicloud\\\" {\\n  region = \\\"cn-beijing\\\"\\n}\\n\\nresource \\\"alicloud_vpc\\\" \\\"default\\\" {\\n  ipv6_isp    = \\\"BGP\\\"\\n  description = \\\"test\\\"\\n  cidr_block  = \\\"10.0.0.0/8\\\"\\n  vpc_name    = var.name\\n  enable_ipv6 = true\\n}",
							"name":           "tftest",
						},
						{
							"description":    "terraform as tool测试",
							"async":          "true",
							"destroy_policy": "NEVER",
							"code":           "variable \\\"name\\\" {\\n  default = \\\"terraform-example\\\"\\n}\\n\\nprovider \\\"alicloud\\\" {\\n  region = \\\"cn-beijing\\\"\\n}\\n\\nresource \\\"alicloud_vpc\\\" \\\"default\\\" {\\n  ipv6_isp    = \\\"BGP\\\"\\n  description = \\\"test\\\"\\n  cidr_block  = \\\"10.0.0.0/8\\\"\\n  vpc_name    = var.name\\n  enable_ipv6 = true\\n}",
							"name":           "tftest2",
						},
						{
							"description":    "terraform as tool测试",
							"async":          "true",
							"destroy_policy": "NEVER",
							"code":           "variable \\\"name\\\" {\\n  default = \\\"terraform-example\\\"\\n}\\n\\nprovider \\\"alicloud\\\" {\\n  region = \\\"cn-beijing\\\"\\n}\\n\\nresource \\\"alicloud_vpc\\\" \\\"default\\\" {\\n  ipv6_isp    = \\\"BGP\\\"\\n  description = \\\"test\\\"\\n  cidr_block  = \\\"10.0.0.0/8\\\"\\n  vpc_name    = var.name\\n  enable_ipv6 = true\\n}",
							"name":           "tftest3",
						},
					},
					"assume_role_name":            "default-role",
					"public_access":               "on",
					"enable_custom_vpc_whitelist": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_tools.#":                "1",
						"description":                   "创建",
						"prompts.#":                     "2",
						"oauth_client_id":               CHECKSET,
						"apis.#":                        "3",
						"instructions":                  "介绍整个MCP Server的作用",
						"additional_api_descriptions.#": "3",
						"vpc_whitelists.#":              "3",
						"name":                          name,
						"language":                      "ZH_CN",
						"enable_assume_role":            "true",
						"assume_role_extra_policy":      CHECKSET,
						"terraform_tools.#":             "3",
						"assume_role_name":              "default-role",
						"public_access":                 "on",
						"enable_custom_vpc_whitelist":   "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prompts": []map[string]interface{}{
						{
							"description": "获取用户定制信息描述",
							"content":     "prompt正文，{{name}}",
							"arguments": []map[string]interface{}{
								{
									"description": "名称信息",
									"required":    "true",
									"name":        "name",
								},
							},
							"name": "获取用户定制信息",
						},
						{
							"description": "获取用户定制信息描述",
							"content":     "prompt正文，{{name}}, {{age}}",
							"arguments": []map[string]interface{}{
								{
									"description": "名称信息",
									"required":    "true",
									"name":        "name",
								},
								{
									"description": "年龄信息",
									"required":    "true",
									"name":        "age",
								},
							},
							"name": "获取用户定制信息1",
						},
					},
					"apis": []map[string]interface{}{
						{
							"api_version": "2014-05-26",
							"product":     "Ecs",
							"selectors": []string{
								"DescribeAvailableResource", "DescribeRegions", "DescribeZones"},
						},
						{
							"api_version": "2017-03-21",
							"product":     "vod",
							"selectors": []string{
								"CreateUploadVideo"},
						},
					},
					"additional_api_descriptions": []map[string]interface{}{
						{
							"api_version":          "2014-05-26",
							"enable_output_schema": "true",
							"api_name":             "DescribeAvailableResource",
							"const_parameters": []map[string]interface{}{
								{
									"value": "cn-hangzhou",
									"key":   "x_mcp_region_id",
								},
								{
									"value": "B1",
									"key":   "a1",
								},
							},
							"api_override_json":   "{\\\"summary\\\": \\\"本接口支持根据不同请求条件查询实例列表，并关联查询实例的详细信息。\\\"}",
							"product":             "Ecs",
							"execute_cli_command": "false",
						},
						{
							"api_version":          "2014-05-26",
							"enable_output_schema": "true",
							"api_name":             "DescribeRegions",
							"product":              "Ecs",
							"execute_cli_command":  "true",
						},
						{
							"api_version":          "2014-05-26",
							"enable_output_schema": "true",
							"api_name":             "DescribeZones",
							"product":              "Ecs",
							"execute_cli_command":  "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prompts.#":                     "2",
						"apis.#":                        "2",
						"additional_api_descriptions.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"additional_api_descriptions": []map[string]interface{}{
						{
							"api_version":          "2014-05-26",
							"enable_output_schema": "true",
							"api_name":             "DescribeAvailableResource",
							"api_override_json":    "{\\\"summary\\\": \\\"本接口支持根据不同请求条件查询实例列表，并关联查询实例的详细信息。\\\"}",
							"product":              "Ecs",
							"execute_cli_command":  "false",
						},
						{
							"api_version":          "2014-05-26",
							"enable_output_schema": "true",
							"api_name":             "DescribeRegions",
							"product":              "Ecs",
							"execute_cli_command":  "true",
						},
						{
							"api_version":          "2014-05-26",
							"enable_output_schema": "true",
							"api_name":             "DescribeZones",
							"product":              "Ecs",
							"execute_cli_command":  "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"additional_api_descriptions.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_tools": []string{},
					"description":  "updated-description",
					"prompts": []map[string]interface{}{
						{
							"description": "prompt-descript",
							"content":     "prompt-content{{a}}",
							"arguments": []map[string]interface{}{
								{
									"description": "参数a",
									"required":    "false",
									"name":        "a",
								},
							},
							"name": "prompt-update",
						},
					},
					"oauth_client_id": "update-1234",
					"apis": []map[string]interface{}{
						{
							"api_version": "2014-08-15",
							"product":     "Rds",
							"selectors": []string{
								"DescribeInstances"},
						},
						{
							"api_version": "2024-11-30",
							"product":     "OpenAPIExplorer",
							"selectors": []string{
								"GenerateCLICommand", "ListApiDefinitions"},
						},
					},
					"instructions": "新的整个描述信息",
					"additional_api_descriptions": []map[string]interface{}{
						{
							"api_version":          "2014-08-15",
							"enable_output_schema": "true",
							"api_name":             "DescribeInstances",
							"const_parameters": []map[string]interface{}{
								{
									"value": "cn-shanghai",
									"key":   "x_mcp_region_id",
								},
								{
									"value": "B",
									"key":   "a",
								},
								{
									"value": "c1_value",
									"key":   "c1",
								},
							},
							"api_override_json":   "{\\\"summary\\\": \\\"修改的描述信息。\\\"}",
							"product":             "Rds",
							"execute_cli_command": "true",
						},
					},
					"vpc_whitelists": []string{
						"vpc-update-test"},
					"language":                 "EN_US",
					"assume_role_extra_policy": "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Effect\\\":\\\"Allow\\\",\\\"Action\\\":[\\\"ecs:Describ*\\\",\\\"vpc:Describe*\\\",\\\"vpc:List*\\\"],\\\"Resource\\\":\\\"*\\\"}]}",
					"terraform_tools": []map[string]interface{}{
						{
							"description":    "更新时候的描述",
							"async":          "true",
							"destroy_policy": "ON_FAILURE",
							"code":           "variable \\\"name\\\" {\\n  default = \\\"terraform-example-update\\\"\\n}\\n\\nprovider \\\"alicloud\\\" {\\n  region = \\\"cn-beijing\\\"\\n}\\n\\nresource \\\"alicloud_vpc\\\" \\\"default\\\" {\\n  ipv6_isp    = \\\"BGP\\\"\\n  description = \\\"test\\\"\\n  cidr_block  = \\\"10.0.0.0/8\\\"\\n  vpc_name    = var.name\\n  enable_ipv6 = true\\n}",
							"name":           "tftoolupdate",
						},
					},
					"assume_role_name":            "update-role-name",
					"public_access":               "follow",
					"enable_custom_vpc_whitelist": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_tools.#":                "0",
						"description":                   "updated-description",
						"prompts.#":                     "1",
						"oauth_client_id":               "update-1234",
						"apis.#":                        "2",
						"instructions":                  "新的整个描述信息",
						"additional_api_descriptions.#": "1",
						"vpc_whitelists.#":              "1",
						"language":                      "EN_US",
						"assume_role_extra_policy":      CHECKSET,
						"terraform_tools.#":             "1",
						"assume_role_name":              "update-role-name",
						"public_access":                 "follow",
						"enable_custom_vpc_whitelist":   "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"additional_api_descriptions": REMOVEKEY,
					"vpc_whitelists":              []string{},
					"terraform_tools":             REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"additional_api_descriptions.#": "0",
						"vpc_whitelists.#":              "0",
						"terraform_tools.#":             "0",
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

var AlicloudOpenApiExplorerApiMcpServerMap11977 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudOpenApiExplorerApiMcpServerBasicDependence11977(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case test 10799
func TestAccAliCloudOpenApiExplorerApiMcpServer_basic10799(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_open_api_explorer_api_mcp_server.default"
	ra := resourceAttrInit(resourceId, AlicloudOpenApiExplorerApiMcpServerMap10799)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OpenApiExplorerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOpenApiExplorerApiMcpServer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOpenApiExplorerApiMcpServerBasicDependence10799)
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
					"description": "zpp server",
					"language":    "EN_US",
					"system_tools": []string{
						"FetchRamActionDetails"},
					"apis": []map[string]interface{}{
						{
							"api_version": "2014-05-26",
							"product":     "Ecs",
							"selectors": []string{
								"Describe*"},
						},
					},
					"additional_api_descriptions": []map[string]interface{}{
						{
							"api_version":       "2022-02-22",
							"api_name":          "Test",
							"api_override_json": "{}",
							"product":           "Ess",
						},
					},
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                   "zpp server",
						"language":                      "EN_US",
						"system_tools.#":                "1",
						"apis.#":                        "1",
						"additional_api_descriptions.#": "1",
						"name":                          name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-for-update",
					"prompts": []map[string]interface{}{
						{
							"description": "Test",
							"content":     "Test",
							"arguments": []map[string]interface{}{
								{
									"description": "Test",
									"required":    "true",
									"name":        "Test",
								},
							},
							"name": "Test",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-for-update",
						"prompts.#":   "1",
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

var AlicloudOpenApiExplorerApiMcpServerMap10799 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudOpenApiExplorerApiMcpServerBasicDependence10799(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test OpenApiExplorer ApiMcpServer. <<< Resource test cases, automatically generated.
