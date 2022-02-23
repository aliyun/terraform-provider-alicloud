package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/imdario/mergo"
)

func init() {
	resource.AddTestSweepers("alicloud_vpc", &resource.Sweeper{
		Name: "alicloud_vpc",
		F:    testSweepVpcs,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_vswitch",
			"alicloud_nat_gateway",
			"alicloud_security_group",
			"alicloud_ots_instance",
			"alicloud_router_interface",
			"alicloud_route_table",
			"alicloud_cen_instance",
			"alicloud_edas_cluster",
			"alicloud_edas_k8s_cluster",
			"alicloud_network_acl",
			"alicloud_cs_kubernetes",
		},
	})
}

func testSweepVpcs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	vpcIds := make([]string, 0)
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	action := "DescribeVpcs"
	var response map[string]interface{}
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve VPC in service list: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Vpcs.Vpc", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Vpcs.Vpc", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			skip := true
			item := v.(map[string]interface{})
			// Skip the default vpc
			if v, ok := item["IsDefault"].(bool); ok && v {
				continue
			}
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["VpcName"])), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping VPC: %v (%v)", item["VpcName"], item["VpcId"])
				continue
			}
			vpcIds = append(vpcIds, fmt.Sprint(item["VpcId"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, id := range vpcIds {
		log.Printf("[INFO] Deleting VPC: (%s)", id)
		action := "DeleteVpc"
		request := map[string]interface{}{
			"VpcId":    id,
			"RegionId": client.RegionId,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(time.Minute*10, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete VPC (%s): %v", id, err)
			continue
		}
	}
	return nil
}

func TestAccAlicloudVpc_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VpcIpv6SupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_cidrs": []string{"106.11.62.0/24"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_cidrs.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "1",
						"tags.Created": "TF",
					}),
				),
			},
			//{
			//	ResourceName:            resourceId,
			//	ImportState:             true,
			//	ImportStateVerify:       true,
			//	ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6"},
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"cidr_block": "172.16.0.0/16",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"cidr_block": "172.16.0.0/16",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"vpc_name": name,
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"vpc_name": name,
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"description": name,
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"description": name,
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"secondary_cidr_blocks": []string{"10.0.0.0/8"},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"secondary_cidr_blocks.#": "1",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"resource_group_id": CHECKSET,
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"tags": map[string]string{
			//			"Created": "TF",
			//			"For":     "Test",
			//		},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"tags.%":       "2",
			//			"tags.Created": "TF",
			//			"tags.For":     "Test",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"tags": map[string]string{
			//			"Created": "TF-update",
			//			"For":     "Test-update",
			//		},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"tags.%":       "2",
			//			"tags.Created": "TF-update",
			//			"tags.For":     "Test-update",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"cidr_block":            "172.16.0.0/12",
			//		"vpc_name":              name + "update",
			//		"description":           name + "update",
			//		"secondary_cidr_blocks": []string{"10.0.0.0/8"},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"cidr_block":              "172.16.0.0/12",
			//			"vpc_name":                name + "update",
			//			"description":             name + "update",
			//			"secondary_cidr_blocks.#": "1",
			//		}),
			//	),
			//},
		},
	})
}

var AlicloudVpcMap = map[string]string{
	"status":          CHECKSET,
	"router_id":       CHECKSET,
	"router_table_id": CHECKSET,
	"route_table_id":  CHECKSET,
	"ipv6_cidr_block": "",
}

func AlicloudVpcBasicDependence(name string) string {
	return fmt.Sprintf(`

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
`)
}

func TestAccAlicloudVpc_enableIpv6(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VpcIpv6SupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_cidrs":  []string{"106.11.62.0/24"},
					"enable_ipv6": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_cidrs.#": "1",
						"enable_ipv6":  "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_name": name,
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
					"secondary_cidr_blocks": []string{"10.0.0.0/8"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_blocks.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
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
				Config: testAccConfig(map[string]interface{}{
					"vpc_name":              name + "update",
					"description":           name + "update",
					"secondary_cidr_blocks": []string{"10.0.0.0/8"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_name":                name + "update",
						"description":             name + "update",
						"secondary_cidr_blocks.#": "1",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudVpc_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VpcIpv6SupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ipv6":       "true",
					"vpc_name":          name,
					"description":       name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"dry_run":           "false",
					"user_cidrs":        []string{"106.11.62.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ipv6":       "true",
						"vpc_name":          name,
						"description":       name,
						"resource_group_id": CHECKSET,
						"dry_run":           "false",
						"user_cidrs.#":      "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6"},
			},
		},
	})
}

func TestAccAlicloudVpc_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VpcIpv6SupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ipv6":       "true",
					"name":              name,
					"description":       name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"dry_run":           "false",
					"user_cidrs":        []string{"106.11.62.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ipv6":       "true",
						"name":              name,
						"description":       name,
						"resource_group_id": CHECKSET,
						"dry_run":           "false",
						"user_cidrs.#":      "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6"},
			},
		},
	})
}

var AlicloudVpcMap1 = map[string]string{
	"status":          CHECKSET,
	"router_id":       CHECKSET,
	"router_table_id": CHECKSET,
	"route_table_id":  CHECKSET,
	"ipv6_cidr_block": CHECKSET,
}

func AlicloudVpcBasicDependence1(name string) string {
	return fmt.Sprintf(`

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
`)
}

// 1. 资源生成字段,在schema没定义
// 2. tags的生成类型有问题, 目前生成为[]map[string]interface{}类型, 实际应该为map[string]interface{]
// 3. attributes 未定义
// 4. 重试错误码,重试后取后一个
// 5. 目前部分资源有异步等待,生成的时候,需要把对应target 状态在查询返回结果里配置下
func TestAccAlicloudVpc_unit1(t *testing.T) {
	resourceName := "alicloud_vpc"
	p := Provider().(*schema.Provider).ResourcesMap
	dUpdate, _ := schema.InternalMap(p[resourceName].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p[resourceName].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	attributes := map[string]interface{}{
		"cidr_block":            "cidr_block",
		"description":           "description",
		"dry_run":               false,
		"enable_ipv6":           false,
		"ipv6_cidr_block":       "ipv6_cidr_block",
		"resource_group_id":     "resource_group_id",
		"user_cidrs":            []interface{}{"user_cidrs_1", "user_cidrs_2"},
		"vpc_name":              "vpc_name",
		"secondary_cidr_blocks": []interface{}{"secondary_cidr_blocks_1", "secondary_cidr_blocks_2"},
		"tags": map[string]interface{}{
			"tagkey1": "tagkey1",
			"tagkey2": "tagkey2",
		},
	}
	for key, value := range attributes {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = dUpdate.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}

	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// todo: 1. 创建字段没映射
		// todo: 2. ReadMockResponse字段生成配套的字段目前为 必填,选填字段+ computed为true字段
		// DescribeVpcs
		"Vpcs": map[string]interface{}{
			"Vpc": []interface{}{
				map[string]interface{}{
					// 属性
					"CidrBlock":       "cidr_block",
					"Description":     "description",
					"Ipv6CidrBlock":   "ipv6_cidr_block",
					"ResourceGroupId": "resource_group_id",
					"UserCidrs": map[string]interface{}{
						"UserCidr": []interface{}{
							"user_cidrs_1", "user_cidrs_2",
						},
					},
					"VpcName": "vpc_name",
					"SecondaryCidrBlocks": map[string]interface{}{
						"SecondaryCidrBlock": []interface{}{
							"secondary_cidr_blocks_1", "secondary_cidr_blocks_2",
						},
					},

					"CreationTime": time.Now().Unix(),
					"IsDefault":    false,
					"RegionId":     "",
					"VRouterId":    "",
					// 纯出参
					"Status": "Available",

					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"Key":   "tagkey1",
								"Value": "tagkey1",
							},
							map[string]interface{}{
								"Key":   "tagkey2",
								"Value": "tagkey2",
							},
						},
					},
					"VSwitchIds": map[string]interface{}{
						"VSwitchId": []interface{}{},
					},

					"VpcId": "",
				},
			},
		},

		//DescribeRouteTableList
		"Code": "200",
		"RouterTableList": map[string]interface{}{
			"RouterTableListType": []interface{}{
				map[string]interface{}{
					"RouteTableType":  "System",
					"ResourceGroupId": "resource_group_id",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateVpc
		"VRouterId":       "CreateVpc_value",
		"VpcId":           "CreateVpc_value",
		"RouteTableId":    "CreateVpc_value",
		"ResourceGroupId": "resource_group_id",
	}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := map[string]func(errorCode string) (map[string]interface{}, error){
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
		"ServiceError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:    String(errorCode),
				Data:    String(errorCode),
				Message: String(errorCode),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage(resourceName, errorCode))
		},
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			err = mergo.Merge(&ReadMockResponse, operationMockResponse, mergo.WithOverride)
			if err != nil {
				log.Printf("[ERROR] the map merge error.")
			}
		}
		return ReadMockResponse, nil
	}

	// Create
	t.Run("Create", func(t *testing.T) {
		// Client Error Mock
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudVpcCreate(dCreate, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		ReadMockResponseDiff = map[string]interface{}{
			// changed fields
			"VRouterId":       "CreateVpc_value",
			"VpcId":           "CreateVpc_value",
			"RouteTableId":    "CreateVpc_value",
			"ResourceGroupId": "resource_group_id",

			// DescribeVpcs
			"Vpcs": map[string]interface{}{
				"Vpc": []interface{}{
					map[string]interface{}{
						// 属性
						"CidrBlock":     "cidr_block",
						"Description":   "description",
						"Ipv6CidrBlock": "ipv6_cidr_block",
						"UserCidrs": map[string]interface{}{
							"UserCidr": []interface{}{
								"user_cidrs_1", "user_cidrs_2",
							},
						},
						"VpcName": "vpc_name",
						"SecondaryCidrBlocks": map[string]interface{}{
							"SecondaryCidrBlock": []interface{}{
								"secondary_cidr_blocks_1", "secondary_cidr_blocks_2",
							},
						},

						"CreationTime": time.Now().Unix(),
						"IsDefault":    false,
						"RegionId":     "",
						// 纯出参
						"Status": "Available",

						"Tags": map[string]interface{}{
							"Tag": []interface{}{
								map[string]interface{}{
									"Key":   "tagkey1",
									"Value": "tagkey1",
								},
								map[string]interface{}{
									"Key":   "tagkey2",
									"Value": "tagkey2",
								},
							},
						},
						"VSwitchIds": map[string]interface{}{
							"VSwitchId": []interface{}{},
						},

						// ==============
						"VpcId":           "CreateVpc_value",
						"ResourceGroupId": "resource_group_id",
						"VRouterId":       "CreateVpc_value",
						"RouteTableId":    "CreateVpc_value",
					},
				},
			},

			//DescribeRouteTableList
			"Code": "200",
			"RouterTableList": map[string]interface{}{
				"RouterTableListType": []interface{}{
					map[string]interface{}{
						"RouteTableType":  "System",
						"ResourceGroupId": "resource_group_id",
						"RouteTableId":    "CreateVpc_value",
					},
				},
			},
		}

		errCodeList := []string{"TaskConflict", "Throttling", "UnknownError", "NonRetryableError", "nil"}
		for index, errorCode := range errCodeList {
			//for _, errorCode := range []string{"nil"} {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				defer func() {
					if index <= len(errCodeList)-2 {
						errorCode = errCodeList[len(errCodeList)-2]
					}
				}()
				switch errorCode {
				case "nil":
					if *action == "CreateVpc" {
						return CreateMockResponse, nil
					}
					return successResponseMock(ReadMockResponseDiff)
				case "NonRetryableError":
					return failedResponseMock["NoRetryError"](errorCode)
				default:
					return failedResponseMock["RetryError"](errorCode)
				}
			})
			err := resourceAlicloudVpcCreate(dCreate, rawClient)
			patches.Reset()
			if errorCode != "nil" {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				for key, _ := range attributes {
					assert.False(t, dCreate.HasChange(key))
				}
			}
		}
	})

	// Update
	t.Run("Update", func(t *testing.T) {
		// Client Error Mock
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudVpcUpdate(dUpdate, rawClient)
		patches.Reset()
		assert.NotNil(t, err)

		// MoveResourceGroup
		diff := terraform.NewInstanceDiff()
		// todo: region_id 需要去掉
		//diff.SetAttribute("region_id", &terraform.ResourceAttrDiff{Old: d.Get("region_id").(string), New: "MoveResourceGroup_value"})
		diff.SetAttribute("resource_group_id", &terraform.ResourceAttrDiff{Old: dUpdate.Get("resource_group_id").(string), New: "MoveResourceGroup_value"})
		dUpdate, _ = schema.InternalMap(p[resourceName].Schema).Data(dCreate.State(), diff)
		ReadMockResponseDiff, err = successResponseMock(map[string]interface{}{
			// DescribeVpcAttribute Response
			"ResourceGroupId": "MoveResourceGroup_value",
			// DescribeVpcs Response
			"Vpcs": map[string]interface{}{
				"Vpc": []interface{}{
					map[string]interface{}{
						// 属性
						"CidrBlock":     "cidr_block",
						"Description":   "description",
						"Ipv6CidrBlock": "ipv6_cidr_block",
						"UserCidrs": map[string]interface{}{
							"UserCidr": []interface{}{
								"user_cidrs_1", "user_cidrs_2",
							},
						},
						"VpcName": "vpc_name",
						"SecondaryCidrBlocks": map[string]interface{}{
							"SecondaryCidrBlock": []interface{}{
								"secondary_cidr_blocks_1", "secondary_cidr_blocks_2",
							},
						},

						"CreationTime": time.Now().Unix(),
						"IsDefault":    false,
						"RegionId":     "",
						"VRouterId":    "CreateVpc_value",
						// 纯出参
						"Status": "Available",

						"Tags": map[string]interface{}{
							"Tag": []interface{}{
								map[string]interface{}{
									"Key":   "tagkey1",
									"Value": "tagkey1",
								},
								map[string]interface{}{
									"Key":   "tagkey2",
									"Value": "tagkey2",
								},
							},
						},
						"VSwitchIds": map[string]interface{}{
							"VSwitchId": []interface{}{},
						},
						"VpcId":        "CreateVpc_value",
						"RouteTableId": "CreateVpc_value",

						"ResourceGroupId": "MoveResourceGroup_value",
					},
				},
			},

			//DescribeRouteTableList
			"Code": "200",
			"RouterTableList": map[string]interface{}{
				"RouterTableListType": []interface{}{
					map[string]interface{}{
						"RouteTableType":  "System",
						"ResourceGroupId": "MoveResourceGroup_value",
						"RouteTableId":    "CreateVpc_value",
					},
				},
			},
		})
		if err != nil {
			log.Printf("[ERROR] the map merge error.")
		}

		errCodeList := []string{"Throttling", "NonRetryableError", "nil"}
		for index, errorCode := range errCodeList {
			//for _, errorCode := range []string{"nil"} {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				defer func() {
					if index <= len(errCodeList)-2 {
						errorCode = errCodeList[len(errCodeList)-2]
					}
				}()
				switch errorCode {
				case "nil":
					return successResponseMock(ReadMockResponseDiff)
				case "NonRetryableError":
					return failedResponseMock["NoRetryError"](errorCode)
				default:
					return failedResponseMock["RetryError"](errorCode)
				}
			})
			err := resourceAlicloudVpcUpdate(dUpdate, rawClient)
			patches.Reset()
			if errorCode != "nil" {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, dUpdate.Get("resource_group_id"), ReadMockResponse["ResourceGroupId"])
			}
		}

		// ModifyVpcAttribute
		diff = terraform.NewInstanceDiff()
		diff.SetAttribute("cidr_block", &terraform.ResourceAttrDiff{Old: dUpdate.Get("cidr_block").(string), New: "ModifyVpcAttribute_value"})
		diff.SetAttribute("description", &terraform.ResourceAttrDiff{Old: dUpdate.Get("description").(string), New: "ModifyVpcAttribute_value"})
		diff.SetAttribute("enable_ipv6", &terraform.ResourceAttrDiff{Old: fmt.Sprint(dUpdate.Get("enable_ipv6")), New: fmt.Sprint(!dUpdate.Get("enable_ipv6").(bool))})
		diff.SetAttribute("ipv6_cidr_block", &terraform.ResourceAttrDiff{Old: dUpdate.Get("ipv6_cidr_block").(string), New: "ModifyVpcAttribute_value"})
		diff.SetAttribute("vpc_name", &terraform.ResourceAttrDiff{Old: dUpdate.Get("vpc_name").(string), New: "ModifyVpcAttribute_value"})
		diff.SetAttribute("name", &terraform.ResourceAttrDiff{Old: dUpdate.Get("vpc_name").(string), New: "ModifyVpcAttribute_value"})
		dUpdate, _ = schema.InternalMap(p[resourceName].Schema).Data(dUpdate.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			// DescribeVpcAttribute Response
			"CidrBlock":     "ModifyVpcAttribute_value",
			"Description":   "ModifyVpcAttribute_value",
			"Ipv6CidrBlock": "ModifyVpcAttribute_value",
			"VpcName":       "ModifyVpcAttribute_value",

			// DescribeVpcs Response
			"Vpcs": map[string]interface{}{
				"Vpc": []interface{}{
					map[string]interface{}{
						// 属性
						"CidrBlock":     "ModifyVpcAttribute_value",
						"Description":   "ModifyVpcAttribute_value",
						"Ipv6CidrBlock": "ModifyVpcAttribute_value",
						"UserCidrs": map[string]interface{}{
							"UserCidr": []interface{}{
								"user_cidrs_1", "user_cidrs_2",
							},
						},
						"VpcName": "ModifyVpcAttribute_value",
						"SecondaryCidrBlocks": map[string]interface{}{
							"SecondaryCidrBlock": []interface{}{
								"secondary_cidr_blocks_1", "secondary_cidr_blocks_2",
							},
						},

						"CreationTime": time.Now().Unix(),
						"IsDefault":    false,
						"RegionId":     "",
						"VRouterId":    "CreateVpc_value",
						// 纯出参
						"Status": "Available",

						"Tags": map[string]interface{}{
							"Tag": []interface{}{
								map[string]interface{}{
									"Key":   "tagkey1",
									"Value": "tagkey1",
								},
								map[string]interface{}{
									"Key":   "tagkey2",
									"Value": "tagkey2",
								},
							},
						},
						"VSwitchIds": map[string]interface{}{
							"VSwitchId": []interface{}{},
						},
						"VpcId":        "CreateVpc_value",
						"RouteTableId": "CreateVpc_value",

						"ResourceGroupId": "MoveResourceGroup_value",
					},
				},
			},

			//DescribeRouteTableList
			"Code": "200",
			"RouterTableList": map[string]interface{}{
				"RouterTableListType": []interface{}{
					map[string]interface{}{
						"RouteTableType":  "System",
						"ResourceGroupId": "MoveResourceGroup_value",
						"RouteTableId":    "CreateVpc_value",
					},
				},
			},
		}

		errCodeList = []string{"Throttling", "NonRetryableError", "nil"}
		for index, errorCode := range errCodeList {
			//for _, errorCode := range []string{"nil"} {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				defer func() {
					if index <= len(errCodeList)-2 {
						errorCode = errCodeList[len(errCodeList)-2]
					}
				}()
				switch errorCode {
				case "nil":
					return successResponseMock(ReadMockResponseDiff)
				case "NonRetryableError":
					return failedResponseMock["NoRetryError"](errorCode)
				default:
					return failedResponseMock["RetryError"](errorCode)
				}
			})
			err := resourceAlicloudVpcUpdate(dUpdate, rawClient)
			patches.Reset()
			if errorCode != "nil" {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, dUpdate.Get("cidr_block"), ReadMockResponse["CidrBlock"])
				assert.Equal(t, dUpdate.Get("description"), ReadMockResponse["Description"])
				assert.Equal(t, dUpdate.Get("ipv6_cidr_block"), ReadMockResponse["Ipv6CidrBlock"])
				assert.Equal(t, dUpdate.Get("vpc_name"), ReadMockResponse["VpcName"])
				assert.Equal(t, dUpdate.Get("name"), ReadMockResponse["VpcName"])
			}
		}

		// UnassociateVpcCidrBlock
		diff = terraform.NewInstanceDiff()
		// todo .# 未生成
		diff.SetAttribute("secondary_cidr_blocks.#", &terraform.ResourceAttrDiff{Old: "2", New: "1"})
		diff.SetAttribute("secondary_cidr_blocks.0", &terraform.ResourceAttrDiff{Old: dUpdate.Get("secondary_cidr_blocks").([]interface{})[0].(string), New: "UnassociateVpcCidrBlock_value_1"})
		diff.SetAttribute("secondary_cidr_blocks.1", &terraform.ResourceAttrDiff{Old: dUpdate.Get("secondary_cidr_blocks").([]interface{})[1].(string), New: ""})
		/// todo: dUpdate.State()
		dUpdate, _ = schema.InternalMap(p[resourceName].Schema).Data(dUpdate.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			"SecondaryCidrBlocks": map[string]interface{}{
				"SecondaryCidrBlock": []interface{}{
					"secondary_cidr_blocks_1",
				},
			},

			// DescribeVpcs Response
			"Vpcs": map[string]interface{}{
				"Vpc": []interface{}{
					map[string]interface{}{
						// 属性
						"CidrBlock":     "ModifyVpcAttribute_value",
						"Description":   "ModifyVpcAttribute_value",
						"Ipv6CidrBlock": "ModifyVpcAttribute_value",
						"UserCidrs": map[string]interface{}{
							"UserCidr": []interface{}{
								"user_cidrs_1", "user_cidrs_2",
							},
						},
						"VpcName": "ModifyVpcAttribute_value",
						"SecondaryCidrBlocks": map[string]interface{}{
							"SecondaryCidrBlock": []interface{}{
								"secondary_cidr_blocks_1",
							},
						},

						"CreationTime": time.Now().Unix(),
						"IsDefault":    false,
						"RegionId":     "",
						"VRouterId":    "CreateVpc_value",
						// 纯出参
						"Status": "Available",

						"Tags": map[string]interface{}{
							"Tag": []interface{}{
								map[string]interface{}{
									"Key":   "tagkey1",
									"Value": "tagkey1",
								},
								map[string]interface{}{
									"Key":   "tagkey2",
									"Value": "tagkey2",
								},
							},
						},
						"VSwitchIds": map[string]interface{}{
							"VSwitchId": []interface{}{},
						},
						"VpcId":        "CreateVpc_value",
						"RouteTableId": "CreateVpc_value",

						"ResourceGroupId": "MoveResourceGroup_value",
					},
				},
			},

			//DescribeRouteTableList
			"Code": "200",
			"RouterTableList": map[string]interface{}{
				"RouterTableListType": []interface{}{
					map[string]interface{}{
						"RouteTableType":  "System",
						"ResourceGroupId": "MoveResourceGroup_value",
						"RouteTableId":    "CreateVpc_value",
					},
				},
			},
		}

		errCodeList = []string{"Throttling", "NonRetryableError", "nil"}
		for index, errorCode := range errCodeList {
			//for _, errorCode := range []string{"nil"} {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				defer func() {
					if index <= len(errCodeList)-2 {
						errorCode = errCodeList[len(errCodeList)-2]
					}
				}()
				switch errorCode {
				case "nil":
					return successResponseMock(ReadMockResponseDiff)
				case "NonRetryableError":
					return failedResponseMock["NoRetryError"](errorCode)
				default:
					return failedResponseMock["RetryError"](errorCode)
				}
			})
			err := resourceAlicloudVpcUpdate(dUpdate, rawClient)
			patches.Reset()
			if errorCode != "nil" {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, strconv.Itoa(len(dUpdate.Get("secondary_cidr_blocks").([]interface{}))), strconv.Itoa(len(ReadMockResponse["SecondaryCidrBlocks"].(map[string]interface{})["SecondaryCidrBlock"].([]interface{}))))
				assert.Equal(t, dUpdate.Get("secondary_cidr_blocks").([]interface{})[0], ReadMockResponse["SecondaryCidrBlocks"].(map[string]interface{})["SecondaryCidrBlock"].([]interface{})[0])
			}
		}
		// AssociateVpcCidrBlock
		diff = terraform.NewInstanceDiff()
		diff.SetAttribute("secondary_cidr_blocks.#", &terraform.ResourceAttrDiff{Old: "1", New: "2"})
		// todo: key生成错误
		diff.SetAttribute("secondary_cidr_blocks.0", &terraform.ResourceAttrDiff{Old: dUpdate.Get("secondary_cidr_blocks").([]interface{})[0].(string), New: "AssociateVpcCidrBlock_value_1"})
		diff.SetAttribute("secondary_cidr_blocks.1", &terraform.ResourceAttrDiff{Old: "", New: "AssociateVpcCidrBlock_value_2"})
		dUpdate, _ = schema.InternalMap(p[resourceName].Schema).Data(dUpdate.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{

			"SecondaryCidrBlocks": map[string]interface{}{
				"SecondaryCidrBlock": []interface{}{
					"AssociateVpcCidrBlock_value_1",
					"AssociateVpcCidrBlock_value_2",
				},
			},
			// DescribeVpcs Response
			"Vpcs": map[string]interface{}{
				"Vpc": []interface{}{
					map[string]interface{}{
						// 属性
						"CidrBlock":     "ModifyVpcAttribute_value",
						"Description":   "ModifyVpcAttribute_value",
						"Ipv6CidrBlock": "ModifyVpcAttribute_value",
						"UserCidrs": map[string]interface{}{
							"UserCidr": []interface{}{
								"user_cidrs_1", "user_cidrs_2",
							},
						},
						"VpcName": "ModifyVpcAttribute_value",
						"SecondaryCidrBlocks": map[string]interface{}{
							"SecondaryCidrBlock": []interface{}{
								"AssociateVpcCidrBlock_value_1",
								"AssociateVpcCidrBlock_value_2",
							},
						},

						"CreationTime": time.Now().Unix(),
						"IsDefault":    false,
						"RegionId":     "",
						"VRouterId":    "CreateVpc_value",
						// 纯出参
						"Status": "Available",

						"Tags": map[string]interface{}{
							"Tag": []interface{}{
								map[string]interface{}{
									"Key":   "tagkey1",
									"Value": "tagkey1",
								},
								map[string]interface{}{
									"Key":   "tagkey2",
									"Value": "tagkey2",
								},
							},
						},
						"VSwitchIds": map[string]interface{}{
							"VSwitchId": []interface{}{},
						},
						"VpcId":        "CreateVpc_value",
						"RouteTableId": "CreateVpc_value",

						"ResourceGroupId": "MoveResourceGroup_value",
					},
				},
			},

			//DescribeRouteTableList
			"Code": "200",
			"RouterTableList": map[string]interface{}{
				"RouterTableListType": []interface{}{
					map[string]interface{}{
						"RouteTableType":  "System",
						"ResourceGroupId": "MoveResourceGroup_value",
						"RouteTableId":    "CreateVpc_value",
					},
				},
			},
		}

		errCodeList = []string{"Throttling", "NonRetryableError", "nil"}
		for index, errorCode := range errCodeList {
			//for _, errorCode := range []string{"nil"} {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				defer func() {
					if index <= len(errCodeList)-2 {
						errorCode = errCodeList[len(errCodeList)-2]
					}
				}()
				switch errorCode {
				case "nil":
					return successResponseMock(ReadMockResponseDiff)
				case "NonRetryableError":
					return failedResponseMock["NoRetryError"](errorCode)
				default:
					return failedResponseMock["RetryError"](errorCode)
				}
			})
			err := resourceAlicloudVpcUpdate(dUpdate, rawClient)
			patches.Reset()
			if errorCode != "nil" {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, strconv.Itoa(len(dUpdate.Get("secondary_cidr_blocks").([]interface{}))), strconv.Itoa(len(ReadMockResponse["SecondaryCidrBlocks"].(map[string]interface{})["SecondaryCidrBlock"].([]interface{}))))
				assert.Equal(t, dUpdate.Get("secondary_cidr_blocks").([]interface{})[0], ReadMockResponse["SecondaryCidrBlocks"].(map[string]interface{})["SecondaryCidrBlock"].([]interface{})[0])
				assert.Equal(t, dUpdate.Get("secondary_cidr_blocks").([]interface{})[1], ReadMockResponse["SecondaryCidrBlocks"].(map[string]interface{})["SecondaryCidrBlock"].([]interface{})[1])

			}
		}

		// UnTagResources
		diff = terraform.NewInstanceDiff()
		diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "2", New: "1"})
		diff.SetAttribute("tags.tagkey1", &terraform.ResourceAttrDiff{Old: "tagkey1", New: "tagkey1"})
		diff.SetAttribute("tags.tagkey2", &terraform.ResourceAttrDiff{Old: "tagkey2", New: ""})

		dUpdate, _ = schema.InternalMap(p[resourceName].Schema).Data(dUpdate.State(), diff)

		ReadMockResponseDiff = map[string]interface{}{
			"Tags": map[string]interface{}{
				"Tag": []interface{}{
					map[string]interface{}{
						"Key":   "tagkey1",
						"Value": "tagkey1",
					},
				},
			},
			// DescribeVpcAttribute Response
			// DescribeVpcs Response
			"Vpcs": map[string]interface{}{
				"Vpc": []interface{}{
					map[string]interface{}{
						// 属性
						"CidrBlock":     "ModifyVpcAttribute_value",
						"Description":   "ModifyVpcAttribute_value",
						"Ipv6CidrBlock": "ModifyVpcAttribute_value",
						"UserCidrs": map[string]interface{}{
							"UserCidr": []interface{}{
								"user_cidrs_1", "user_cidrs_2",
							},
						},
						"VpcName": "ModifyVpcAttribute_value",
						"SecondaryCidrBlocks": map[string]interface{}{
							"SecondaryCidrBlock": []interface{}{
								"AssociateVpcCidrBlock_value_1",
								"AssociateVpcCidrBlock_value_2",
							},
						},

						"CreationTime": time.Now().Unix(),
						"IsDefault":    false,
						"RegionId":     "",
						"VRouterId":    "CreateVpc_value",
						// 纯出参
						"Status": "Available",

						"Tags": map[string]interface{}{
							"Tag": []interface{}{
								map[string]interface{}{
									"Key":   "tagkey1",
									"Value": "tagkey1",
								},
							},
						},
						"VSwitchIds": map[string]interface{}{
							"VSwitchId": []interface{}{},
						},
						"VpcId":        "CreateVpc_value",
						"RouteTableId": "CreateVpc_value",

						"ResourceGroupId": "MoveResourceGroup_value",
					},
				},
			},

			//DescribeRouteTableList
			"Code": "200",
			"RouterTableList": map[string]interface{}{
				"RouterTableListType": []interface{}{
					map[string]interface{}{
						"RouteTableType":  "System",
						"ResourceGroupId": "MoveResourceGroup_value",
						"RouteTableId":    "CreateVpc_value",
					},
				},
			},
		}

		errCodeList = []string{"Throttling", "NonRetryableError", "nil"}
		for index, errorCode := range errCodeList {
			//for _, errorCode := range []string{"nil"} {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				defer func() {
					if index <= len(errCodeList)-2 {
						errorCode = errCodeList[len(errCodeList)-2]
					}
				}()
				switch errorCode {
				case "nil":
					return successResponseMock(ReadMockResponseDiff)
				case "NonRetryableError":
					return failedResponseMock["NoRetryError"](errorCode)
				default:
					return failedResponseMock["RetryError"](errorCode)
				}
			})
			err := resourceAlicloudVpcUpdate(dUpdate, rawClient)
			patches.Reset()
			if errorCode != "nil" {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, strconv.Itoa(len(dUpdate.Get("tags").(map[string]interface{}))), strconv.Itoa(len(ReadMockResponse["Tags"].(map[string]interface{})["Tag"].([]interface{}))))
			}
		}

		// TagResources
		diff = terraform.NewInstanceDiff()
		diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "2", New: "1"})
		diff.SetAttribute("tags.tagkey1", &terraform.ResourceAttrDiff{Old: "tagkey1", New: "tagkey1"})
		diff.SetAttribute("tags.tagkey2", &terraform.ResourceAttrDiff{Old: "", New: "tagkey2"})

		dUpdate, _ = schema.InternalMap(p["alicloud_vpc"].Schema).Data(dUpdate.State(), diff)
		ReadMockResponse = map[string]interface{}{
			"Tags": map[string]interface{}{
				"Tag": []interface{}{
					map[string]interface{}{
						"Key":   "tagkey1",
						"Value": "tagkey1",
					},
					map[string]interface{}{
						"Key":   "tagkey2",
						"Value": "tagkey2",
					},
				},
			},
			// DescribeVpcAttribute Response
			// DescribeVpcs Response
			"Vpcs": map[string]interface{}{
				"Vpc": []interface{}{
					map[string]interface{}{
						// 属性
						"CidrBlock":     "ModifyVpcAttribute_value",
						"Description":   "ModifyVpcAttribute_value",
						"Ipv6CidrBlock": "ModifyVpcAttribute_value",
						"UserCidrs": map[string]interface{}{
							"UserCidr": []interface{}{
								"user_cidrs_1", "user_cidrs_2",
							},
						},
						"VpcName": "ModifyVpcAttribute_value",
						"SecondaryCidrBlocks": map[string]interface{}{
							"SecondaryCidrBlock": []interface{}{
								"AssociateVpcCidrBlock_value_1",
								"AssociateVpcCidrBlock_value_2",
							},
						},

						"CreationTime": time.Now().Unix(),
						"IsDefault":    false,
						"RegionId":     "",
						"VRouterId":    "CreateVpc_value",
						// 纯出参
						"Status": "Available",

						"Tags": map[string]interface{}{
							"Tag": []interface{}{
								map[string]interface{}{
									"Key":   "tagkey1",
									"Value": "tagkey1",
								},
								map[string]interface{}{
									"Key":   "tagkey2",
									"Value": "tagkey2",
								},
							},
						},
						"VSwitchIds": map[string]interface{}{
							"VSwitchId": []interface{}{},
						},
						"VpcId":        "CreateVpc_value",
						"RouteTableId": "CreateVpc_value",

						"ResourceGroupId": "MoveResourceGroup_value",
					},
				},
			},

			//DescribeRouteTableList
			"Code": "200",
			"RouterTableList": map[string]interface{}{
				"RouterTableListType": []interface{}{
					map[string]interface{}{
						"RouteTableType":  "System",
						"ResourceGroupId": "MoveResourceGroup_value",
						"RouteTableId":    "CreateVpc_value",
					},
				},
			},
		}

		errCodeList = []string{"Throttling", "NonRetryableError", "nil"}
		for index, errorCode := range errCodeList {
			//for _, errorCode := range []string{"nil"} {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				defer func() {
					if index <= len(errCodeList)-2 {
						errorCode = errCodeList[len(errCodeList)-2]
					}
				}()
				switch errorCode {
				case "nil":
					return successResponseMock(ReadMockResponse)
				case "NonRetryableError":
					return failedResponseMock["NoRetryError"](errorCode)
				default:
					return failedResponseMock["RetryError"](errorCode)
				}
			})
			err := resourceAlicloudVpcUpdate(dUpdate, rawClient)
			patches.Reset()
			if errorCode != "nil" {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, strconv.Itoa(len(dUpdate.Get("tags").(map[string]interface{}))), strconv.Itoa(len(ReadMockResponse["Tags"].(map[string]interface{})["Tag"].([]interface{}))))
			}
		}
	})

	// Delete
	t.Run("Delete", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudVpcDelete(dUpdate, rawClient)
		patches.Reset()
		assert.NotNil(t, err)

		// todo:
		errCodeList := []string{"Forbidden.VpcNotFound", "InvalidVpcID.NotFound"}
		for _, errorCode := range errCodeList {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				switch errorCode {
				//case :
				//	return successResponseMock(map[string]interface{}{})
				case "Forbidden.VpcNotFound", "InvalidVpcID.NotFound", "InvalidVpcId.NotFound", "NoRetryError":
					return failedResponseMock["NoRetryError"](errorCode)
				default:
					return failedResponseMock["RetryError"](errorCode)
				}
			})
			err := resourceAlicloudVpcDelete(dUpdate, rawClient)
			patches.Reset()
			assert.Nil(t, err)
		}
	})

	t.Run("Read", func(t *testing.T) {
		ReadMockResponseDiff = map[string]interface{}{
			// DescribeVpcAttribute Response
			// DescribeVpcs Response
			"Vpcs": map[string]interface{}{
				"Vpc": []interface{}{},
			},

			//DescribeRouteTableList
			"Code": "200",
			"RouterTableList": map[string]interface{}{
				"RouterTableListType": []interface{}{
					map[string]interface{}{
						"RouteTableType":  "System",
						"ResourceGroupId": "MoveResourceGroup_value",
						"RouteTableId":    "CreateVpc_value",
					},
				},
			},
		}
		errCodeList := []string{"NoRetryError", "InvalidVpcID.NotFound"}
		for _, errorCode := range errCodeList {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				switch errorCode {
				case "InvalidVpcID.NotFound":
					return failedResponseMock["NotFoundError"](errorCode)
				case "NoRetryError":
					return failedResponseMock["NoRetryError"](errorCode)
				default:
					return failedResponseMock["RetryError"](errorCode)
				}
			})
			err := resourceAlicloudVpcRead(dUpdate, rawClient)
			patches.Reset()
			if err != nil {
				isNotFoundError := false
				for _, notFoundErrorCode := range []string{"Forbidden.VpcNotFound", "InvalidVpcID.NotFound"} {
					if errorCode == notFoundErrorCode {
						isNotFoundError = true
						break
					}
				}
				if isNotFoundError {
					assert.Nil(t, err)
				} else {
					assert.NotNil(t, err)
				}

			} else {
				assert.Nil(t, err)
			}

		}

	})

}
