package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVpcsDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()
	initVswitchConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${alicloud_vswitch.default.id}"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccVpcsdatasource%d"`, rand),
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccVpcsdatasource%d_fake"`, rand),
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_vpc.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_vpc.default.id}_fake" ]`,
		}),
	}
	cidrBlockConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"cidr_block": `"172.16.0.0/12"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"cidr_block": `"172.16.0.0/0"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"status":     `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"status":     `"Pending"`,
		}),
	}
	idDefaultConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"is_default": `"false"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"is_default": `"true"`,
		}),
	}
	vswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"vswitch_id": `"${alicloud_vswitch.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"vswitch_id": `"${alicloud_vswitch.default.id}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"tags": `{
							Created = "TF"
							For 	= "acceptance test"
					  }`,
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"tags": `{
							Created = "TF-fake"
							For 	= "acceptance test-fake"
					  }`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${var.name}"`,
			"resource_group_id": fmt.Sprintf(`"%s"`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID")),
		}),
	}
	pagingConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"vpc_name":    `"${var.name}"`,
			"page_number": `1`,
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"vpc_name":    `"${var.name}"`,
			"page_number": `2`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${var.name}"`,
			"ids":               `[ "${alicloud_vpc.default.id}" ]`,
			"cidr_block":        `"172.16.0.0/12"`,
			"status":            `"Available"`,
			"is_default":        `"false"`,
			"vswitch_id":        `"${alicloud_vswitch.default.id}"`,
			"resource_group_id": fmt.Sprintf(`"%s"`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID")),
			"vpc_name":          `"${var.name}"`,
			"page_number":       `1`,
		}),
		fakeConfig: testAccCheckAlicloudVpcsDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${var.name}"`,
			"ids":               `[ "${alicloud_vpc.default.id}" ]`,
			"cidr_block":        `"172.16.0.0/16"`,
			"status":            `"Available"`,
			"is_default":        `"false"`,
			"vswitch_id":        `"${alicloud_vswitch.default.id}_fake"`,
			"resource_group_id": fmt.Sprintf(`"%s"`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID")),
			"vpc_name":          `"${var.name}"`,
			"page_number":       `2`,
		}),
	}

	vpcsCheckInfo.dataSourceTestCheck(t, rand, initVswitchConf, nameRegexConf, idsConf, cidrBlockConf, statusConf, idDefaultConf, vswitchIdConf, tagsConf, resourceGroupIdConf, pagingConf, allConf)
}

func testAccCheckAlicloudVpcsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccVpcsdatasource%d"
}

resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}"
  cidr_block = "172.16.0.0/12"
  tags 		= {
		Created = "TF"
		For 	= "acceptance test"
  }
  resource_group_id = "%s"
}

data "alicloud_zones" "default" {

}

resource "alicloud_vswitch" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
	vpc_id = "${alicloud_vpc.default.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

data "alicloud_vpcs" "default" {
	enable_details = true
  %s
}
`, rand, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"), strings.Join(pairs, "\n  "))
	return config
}

var existVpcsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                 "1",
		"names.#":               "1",
		"vpcs.#":                "1",
		"total_count":           CHECKSET,
		"vpcs.0.id":             CHECKSET,
		"vpcs.0.region_id":      CHECKSET,
		"vpcs.0.status":         "Available",
		"vpcs.0.vpc_name":       fmt.Sprintf("tf-testAccVpcsdatasource%d", rand),
		"vpcs.0.vswitch_ids.#":  "1",
		"vpcs.0.cidr_block":     "172.16.0.0/12",
		"vpcs.0.vrouter_id":     CHECKSET,
		"vpcs.0.router_id":      CHECKSET,
		"vpcs.0.route_table_id": CHECKSET,
		"vpcs.0.description":    "",
		"vpcs.0.is_default":     "false",
		"vpcs.0.creation_time":  CHECKSET,
	}
}

var fakeVpcsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"vpcs.#":  "0",
	}
}

var vpcsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_vpcs.default",
	existMapFunc: existVpcsMapFunc,
	fakeMapFunc:  fakeVpcsMapFunc,
}

func TestAccAlicloudVpcsDataSourceBasic_unit(t *testing.T) {
	p := Provider().(*schema.Provider).DataSourcesMap
	d, _ := schema.InternalMap(p["alicloud_vpcs"].Schema).Data(nil, nil)
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rand := acctest.RandIntRange(1000000, 9999999)
	attributes := map[string]interface{}{
		"dhcp_options_set_id": "dhcp_options_set_id",
		"vswitch_id":          "vswitch_id",
		"dry_run":             false,
		"is_default":          false,
		"resource_group_id":   "resource_group_id",
		"name_regex":          "vpc_name",
		"vpc_name":            "vpc_name",
		"vpc_owner_id":        rand,
		"ids":                 []string{"MockId"},
		"output_file":         "file.json",
		"tags": map[string]interface{}{
			"TagKey1": "TagValue1",
			"TagKey2": "TagValue2",
		},
		"cidr_block": "cidr_block",
		"status":     "Available",
	}
	for key, value := range attributes {
		err := d.Set(key, value)
		assert.Nil(t, err)
	}
	ReadMockResponse := map[string]interface{}{
		"TotalCount": "5",
		"Vpcs": map[string]interface{}{
			"Vpc": []interface{}{
				map[string]interface{}{
					"VpcId": "MockId",
					"UserCidrs": map[string]interface{}{
						"UserCidr": []interface{}{"UserCidr"},
					},
					"CidrBlock":     "cidr_block",
					"Description":   "description",
					"Ipv6CidrBlock": "ipv6_cidr_block",
					"VRouterId":     "v_router_id",
					"SecondaryCidrBlocks": map[string]interface{}{
						"SecondaryCidrBlock": []interface{}{"secondary_cidr_blocks"},
					},
					"Status": "Available",
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"Key":   "TagKey1",
								"Value": "TagValue1",
							},
							map[string]interface{}{
								"Key":   "TagKey2",
								"Value": "TagValue2",
							},
						},
					},
					"VSwitchIds": map[string]interface{}{
						"VSwitchId": []interface{}{"vswitch_id"},
					},
					"UserCidr": "user_cidrs",
					"VpcName":  "vpc_name",
				},
				map[string]interface{}{
					"VpcId": "MockId1",
					"UserCidrs": map[string]interface{}{
						"UserCidr": []interface{}{"UserCidr1"},
					},
					"CidrBlock":     "cidr_block1",
					"Description":   "description1",
					"Ipv6CidrBlock": "ipv6_cidr_block1",
					"VRouterId":     "v_router_id1",
					"SecondaryCidrBlocks": map[string]interface{}{
						"SecondaryCidrBlock": []interface{}{"secondary_cidr_blocks1"},
					},
					"Status": "Available",
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"Key":   "TagKey1",
								"Value": "TagValue1",
							},
							map[string]interface{}{
								"Key":   "TagKey2",
								"Value": "TagValue2",
							},
						},
					},
					"VSwitchIds": map[string]interface{}{
						"VSwitchId": []interface{}{"vswitch_id"},
					},
					"UserCidr": "user_cidrs",
					"VpcName":  "vpc_name",
				},
				map[string]interface{}{
					"VpcId": "MockId2",
					"UserCidrs": map[string]interface{}{
						"UserCidr": []interface{}{"UserCidr1"},
					},
					"CidrBlock":     "cidr_block",
					"Description":   "description1",
					"Ipv6CidrBlock": "ipv6_cidr_block1",
					"VRouterId":     "v_router_id1",
					"SecondaryCidrBlocks": map[string]interface{}{
						"SecondaryCidrBlock": []interface{}{"secondary_cidr_blocks1"},
					},
					"Status": "Available",
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"Key":   "TagKey1",
								"Value": "TagValue1",
							},
							map[string]interface{}{
								"Key":   "TagKey2",
								"Value": "TagValue2",
							},
						},
					},
					"VSwitchIds": map[string]interface{}{
						"VSwitchId": []interface{}{"vswitch_id1"},
					},
					"UserCidr": "user_cidrs",
					"VpcName":  "vpc_name1",
				},
				map[string]interface{}{
					"VpcId": "MockId3",
					"UserCidrs": map[string]interface{}{
						"UserCidr": []interface{}{"UserCidr1"},
					},
					"CidrBlock":     "cidr_block",
					"Description":   "description",
					"Ipv6CidrBlock": "ipv6_cidr_block",
					"VRouterId":     "v_router_id",
					"SecondaryCidrBlocks": map[string]interface{}{
						"SecondaryCidrBlock": []interface{}{"secondary_cidr_blocks"},
					},
					"Status": "Available",
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"Key":   "TagKey1",
								"Value": "TagValue1",
							},
							map[string]interface{}{
								"Key":   "TagKey2",
								"Value": "TagValue2",
							},
						},
					},
					"VSwitchIds": map[string]interface{}{
						"VSwitchId": []interface{}{"vswitch_id"},
					},
					"UserCidr": "user_cidrs",
					"VpcName":  "vpc_name",
				},
				map[string]interface{}{
					"VpcId": "MockId4",
					"UserCidrs": map[string]interface{}{
						"UserCidr": []interface{}{"UserCidr"},
					},
					"CidrBlock":     "cidr_block",
					"Description":   "description",
					"Ipv6CidrBlock": "ipv6_cidr_block",
					"VRouterId":     "v_router_id",
					"SecondaryCidrBlocks": map[string]interface{}{
						"SecondaryCidrBlock": []interface{}{"secondary_cidr_blocks"},
					},
					"Status": "Pending",
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"Key":   "TagKey1",
								"Value": "TagValue1",
							},
							map[string]interface{}{
								"Key":   "TagKey2",
								"Value": "TagValue2",
							},
						},
					},
					"VSwitchIds": map[string]interface{}{
						"VSwitchId": []interface{}{"vswitch_id"},
					},
					"UserCidr": "user_cidrs",
					"VpcName":  "vpc_name",
				},
			},
		},

		"Code": "200",
		"RouterTableList": map[string]interface{}{
			"RouterTableListType": []interface{}{
				map[string]interface{}{
					"RouteTableType": "System",
					"RouteTableId":   "route_table_id",
				},
			},
		},
	}
	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:    String(errorCode),
				Data:    String(errorCode),
				Message: String(errorCode),
			}
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:    String(errorCode),
				Data:    String(errorCode),
				Message: String(errorCode),
			}
		},
	}

	// //enable_details = false
	t.Run("ReadSummary", func(t *testing.T) {
		// Cilent Error
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := dataSourceAlicloudVpcsRead(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)

		noRetryFlag,abnormal := true,true
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		for {
			err = dataSourceAlicloudVpcsRead(d, rawClient)
			if abnormal{
				abnormal = false
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				break
			}
		}

		//ReadNormalSummaryPage
		attributes["page_number"] = 1
		attributes["page_size"] = PageSizeLarge
		d.Set("page_number", 1)
		d.Set("page_number", PageSizeLarge)
		noRetryFlag = false
		err = dataSourceAlicloudVpcsRead(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// enable_details = true
	t.Run("ReadNormalDetail", func(t *testing.T) {
		attributes["enable_details"] = true
		d.Set("enable_details", true)
		noRetryFlag,abnormal := true,true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})

		for {
			err = dataSourceAlicloudVpcsRead(d, rawClient)
			if abnormal{
				abnormal = false
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				break
			}
		}
		patches.Reset()
	})

	t.Run("ReadResponseParseAbnormal", func(t *testing.T) {
		noRetryFlag := false
		patchesParse := gomonkey.ApplyFunc(jsonpath.Get, func(path string, value interface{}) (interface{}, error) {
			return nil, WrapError(fmt.Errorf("parsing error"))
		})
		patchesDoRequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := dataSourceAlicloudVpcsRead(d, rawClient)
		patchesDoRequest.Reset()
		patchesParse.Reset()
		assert.NotNil(t, err)
	})
	t.Run("ReadDescribeRouteTableListAbnormal", func(t *testing.T) {
		attributes["enable_details"] = true
		d.Set("enable_details", true)
		noRetryFlag := false
		parsesDescribeRouteTableList := gomonkey.ApplyMethod(reflect.TypeOf(&VpcService{}), "DescribeRouteTableList", func(_ *VpcService, id string) (object map[string]interface{}, err error) {
			return nil, WrapError(fmt.Errorf("vpcService.DescribeRouteTableList"))
		})
		patchesDoRequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := dataSourceAlicloudVpcsRead(d, rawClient)
		patchesDoRequest.Reset()
		parsesDescribeRouteTableList.Reset()
		assert.NotNil(t, err)
	})
	t.Run("ReadSetFieldAbnormal", func(t *testing.T) {
		targetField := ""
		noRetryFlag := false
		patchesDoRequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&schema.ResourceData{}), "Set", func(_ *schema.ResourceData, key string, value interface{}) error {
			if key == targetField  {
				return WrapError(fmt.Errorf("terraform setting error"))
			}
			return nil
		})
		for _, field := range []string{"ids", "names","vpcs","total_count"}{
			targetField = field
			err = dataSourceAlicloudVpcsRead(d, rawClient)
			assert.NotNil(t, err)
		}
		patchesDoRequest.Reset()
		patches.Reset()
	})
}
