package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_vpc_prefix_list",
		&resource.Sweeper{
			Name: "alicloud_vpc_prefix_list",
			F:    testSweepVpcPrefixList,
		})
}

func testSweepVpcPrefixList(region string) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListPrefixLists"
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge
	var response map[string]interface{}

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.PrefixLists", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.PrefixLists", action, err)
			return nil
		}
		if result, ok := resp.([]interface{}); ok {
			for _, v := range result {
				item := v.(map[string]interface{})

				skip := true
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["PrefixListName"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Vpc Prefix List: %s", item["PrefixListName"].(string))
					continue
				}
				action := "DeleteVpcPrefixList"
				request := map[string]interface{}{
					"PrefixListId": item["PrefixListId"],
					"RegionId":     client.RegionId,
				}
				_, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
				if err != nil {
					log.Printf("[ERROR] Failed to delete Vpc Prefix List (%s): %s", item["PrefixListName"].(string), err)
				}
				log.Printf("[INFO] Delete Vpc Prefix List success: %s ", item["PrefixListName"].(string))
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudVPCPrefixList_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_prefix_list.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCPrefixListMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcPrefixList")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcprefixlist%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCPrefixListBasicDependence0)
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
					"ip_version":  "IPV4",
					"max_entries": "50",
					"entrys": []map[string]interface{}{
						{
							"cidr":        "192.168.0.0/16",
							"description": "${var.name}",
						},
					},
					"prefix_list_name":        "${var.name}",
					"prefix_list_description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_version":              "IPV4",
						"max_entries":             "50",
						"entrys.#":                "1",
						"prefix_list_name":        name,
						"prefix_list_description": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudVPCPrefixListMap0 = map[string]string{
	"ip_version":  CHECKSET,
	"entrys.#":    CHECKSET,
	"max_entries": CHECKSET,
	"status":      CHECKSET,
}

func AlicloudVPCPrefixListBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func TestAccAlicloudVPCPrefixList_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_prefix_list.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCPrefixListMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcPrefixList")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcprefixlist%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCPrefixListBasicDependence1)
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
					"entrys": []map[string]interface{}{
						{
							"cidr":        "192.168.0.0/16",
							"description": "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"entrys.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prefix_list_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prefix_list_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_entries": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_entries": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entrys": []map[string]interface{}{
						{
							"cidr":        "192.168.1.0/24",
							"description": "${var.name}_update1",
						},
						{
							"cidr":        "192.168.2.0/24",
							"description": "${var.name}_update2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"entrys.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entrys": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"entrys.#": "0",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudVPCPrefixList_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_prefix_list.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCPrefixListMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcPrefixList")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcprefixlist%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCPrefixListBasicDependence1)
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
					"prefix_list_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prefix_list_name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudVPCPrefixListMap1 = map[string]string{
	"max_entries": CHECKSET,
	"entrys.#":    CHECKSET,
	"ip_version":  CHECKSET,
	"status":      CHECKSET,
}

func AlicloudVPCPrefixListBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func TestUnitAccAlicloudVpcPrefixList(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_vpc_prefix_list"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_vpc_prefix_list"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"ip_version":  "CreateVpcPrefixListValue",
		"max_entries": 50,
		"entrys": []map[string]interface{}{
			{
				"cidr":        "CreateVpcPrefixListValue",
				"description": "CreateVpcPrefixListValue",
			},
		},
		"prefix_list_name":        "CreateVpcPrefixListValue",
		"prefix_list_description": "CreateVpcPrefixListValue",
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
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
		"PrefixLists": []interface{}{
			map[string]interface{}{
				"Status":                "Created",
				"MaxEntries":            50,
				"PrefixListDescription": "CreateVpcPrefixListValue",
				"PrefixListId":          "VpcPrefixListId",
				"IpVersion":             "CreateVpcPrefixListValue",
				"PrefixListName":        "CreateVpcPrefixListValue",
			},
		},
	}
	ReadEntriesMockResponse := map[string]interface{}{
		"PrefixListEntry": []interface{}{
			map[string]interface{}{
				"Description":  "CreateVpcPrefixListValue",
				"Cidr":         "CreateVpcPrefixListValue",
				"PrefixListId": "VpcPrefixListId",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"PrefixListId": "VpcPrefixListId",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_vpc_prefix_list", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	successEntriesResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadEntriesMockResponse, operationMockResponse)
		}
		return ReadEntriesMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudVpcPrefixListCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateVpcPrefixList" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			if *action == "GetVpcPrefixListEntries" {
				return ReadEntriesMockResponse, nil
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudVpcPrefixListCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_vpc_prefix_list"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudVpcPrefixListUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"max_entries": 40,
		"entrys": []map[string]interface{}{
			{
				"cidr":        "UpdateVpcPrefixListValue",
				"description": "UpdateVpcPrefixListValue",
			},
		},
		"prefix_list_name":        "UpdateVpcPrefixListValue",
		"prefix_list_description": "UpdateVpcPrefixListValue",
	}
	diff, err := newInstanceDiff("alicloud_vpc_prefix_list", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_vpc_prefix_list"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"PrefixLists": []interface{}{
			map[string]interface{}{
				"Status":                "Created",
				"MaxEntries":            40,
				"PrefixListDescription": "UpdateVpcPrefixListValue",
				"PrefixListName":        "UpdateVpcPrefixListValue",
			},
		},
	}
	ReadEntriesMockResponseDiff := map[string]interface{}{
		"PrefixListEntry": []interface{}{
			map[string]interface{}{
				"Description": "UpdateVpcPrefixListValue",
				"Cidr":        "UpdateVpcPrefixListValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyVpcPrefixList" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			if *action == "GetVpcPrefixListEntries" {
				successEntriesResponseMock(ReadEntriesMockResponseDiff)
				return ReadEntriesMockResponse, nil
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudVpcPrefixListUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_vpc_prefix_list"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	diff, err = newInstanceDiff("alicloud_vpc_prefix_list", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_vpc_prefix_list"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ListPrefixLists" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			if *action == "GetVpcPrefixListEntries" {
				return ReadEntriesMockResponse, nil
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudVpcPrefixListRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudVpcPrefixListDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_vpc_prefix_list", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_vpc_prefix_list"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteVpcPrefixList" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudVpcPrefixListDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// Test Vpc PrefixList. >>> Resource test cases, automatically generated.
// Case 3120
func TestAccAlicloudVpcPrefixList_basic3120(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_prefix_list.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcPrefixListMap3120)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcPrefixList")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sVpcPrefixList%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcPrefixListBasicDependence3120)
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
					"prefix_list_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prefix_list_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_entries": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_entries": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${alicloud_resource_manager_resource_group.defaultRg.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prefix_list_description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prefix_list_description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entrys": []map[string]interface{}{
						{
							"cidr":        "192.168.0.0/16",
							"description": "test",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"entrys.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_entries": "40",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_entries": "40",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${alicloud_resource_manager_resource_group.changeRg.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prefix_list_description": "testupdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prefix_list_description": "testupdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prefix_list_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prefix_list_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entrys": []map[string]interface{}{
						{
							"cidr":        "10.0.0.0/8",
							"description": "testupdate",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"entrys.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entrys": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"entrys.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_entries":             "50",
					"resource_group_id":       "${alicloud_resource_manager_resource_group.defaultRg.id}",
					"prefix_list_description": "test",
					"ip_version":              "IPV4",
					"prefix_list_name":        name + "_update",
					"entrys": []map[string]interface{}{
						{
							"cidr":        "192.168.0.0/16",
							"description": "test",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_entries":             "50",
						"resource_group_id":       CHECKSET,
						"prefix_list_description": "test",
						"ip_version":              "IPV4",
						"prefix_list_name":        name + "_update",
						"entrys.#":                "1",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudVpcPrefixListMap3120 = map[string]string{}

func AlicloudVpcPrefixListBasicDependence3120(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_resource_manager_resource_group" "defaultRg" {
  display_name        = "tf-testacc-chenyi"
  resource_group_name = var.name
}

resource "alicloud_resource_manager_resource_group" "changeRg" {
  display_name        = "tf-testacc-chenyi-change"
  resource_group_name = "${var.name}1"
}


`, name)
}

// Case 3120  twin
func TestAccAlicloudVpcPrefixList_basic3120_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_prefix_list.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcPrefixListMap3120)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcPrefixList")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sVpcPrefixList%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcPrefixListBasicDependence3120)
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
					"max_entries":             "40",
					"resource_group_id":       "${alicloud_resource_manager_resource_group.changeRg.id}",
					"prefix_list_description": "testupdate",
					"ip_version":              "IPV4",
					"prefix_list_name":        name,
					"entrys": []map[string]interface{}{
						{
							"cidr":        "10.0.0.0/8",
							"description": "testupdate",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_entries":             "40",
						"resource_group_id":       CHECKSET,
						"prefix_list_description": "testupdate",
						"ip_version":              "IPV4",
						"prefix_list_name":        name,
						"entrys.#":                "1",
						"tags.%":                  "2",
						"tags.Created":            "TF",
						"tags.For":                "Test",
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

// Test Vpc PrefixList. <<< Resource test cases, automatically generated.
