package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_cms_dynamic_tag_group",
		&resource.Sweeper{
			Name: "alicloud_cms_dynamic_tag_group",
			F:    testSweepCmsDynamicTagGroup,
		})
}

func testSweepCmsDynamicTagGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeDynamicTagRuleList"
	request := map[string]interface{}{}

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
		resp, err := jsonpath.Get("$.TagGroupList.TagGroup", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["TagKey"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping DynamicTagGroup: %s", item["TagKey"].(string))
					continue
				}
			}
			action := "DeleteDynamicTagGroup"
			request := map[string]interface{}{
				"DynamicTagRuleId": item["DynamicTagRuleId"],
			}
			_, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete DynamicTagGroup (%s): %s", item["TagKey"].(string), err)
			}
			log.Printf("[INFO] Delete DynamicTagGroup success: %s ", item["TagKey"].(string))
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudCmsDynamicTagGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_dynamic_tag_group.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsDynamicTagGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsDynamicTagGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCmsDynamicTagGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsDynamicTagGroupBasicDependence0)
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
					"tag_key":            name,
					"contact_group_list": "${alicloud_cms_alarm_contact_group.default.*.id}",
					"match_express": []map[string]interface{}{
						{
							"tag_value":                name + "all",
							"tag_value_match_function": "all",
						},
						{
							"tag_value":                name + "startWith",
							"tag_value_match_function": "startWith",
						},
						{
							"tag_value":                name + "endWith",
							"tag_value_match_function": "endWith",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_key":              name,
						"contact_group_list.#": "6",
						"match_express.#":      "3",
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

func TestAccAliCloudCmsDynamicTagGroup_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_dynamic_tag_group.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsDynamicTagGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsDynamicTagGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCmsDynamicTagGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsDynamicTagGroupBasicDependence0)
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
					"tag_key":                       name,
					"match_express_filter_relation": "or",
					"contact_group_list":            "${alicloud_cms_alarm_contact_group.default.*.id}",
					"template_id_list":              "${alicloud_cms_metric_rule_template.default.*.id}",
					"match_express": []map[string]interface{}{
						{
							"tag_value":                name + "contains",
							"tag_value_match_function": "contains",
						},
						{
							"tag_value":                name + "notContains",
							"tag_value_match_function": "notContains",
						},
						{
							"tag_value":                name + "equals",
							"tag_value_match_function": "equals",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_key":                       name,
						"match_express_filter_relation": "or",
						"contact_group_list.#":          "6",
						"template_id_list.#":            "6",
						"match_express.#":               "3",
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

var AliCloudCmsDynamicTagGroupMap0 = map[string]string{
	"match_express_filter_relation": CHECKSET,
	"status":                        CHECKSET,
}

func AliCloudCmsDynamicTagGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_cms_alarm_contact_group" "default" {
  		count                    = 6
  		alarm_contact_group_name = "${var.name}-${count.index}"
	}

	resource "alicloud_cms_metric_rule_template" "default" {
  		count                     = 6
  		metric_rule_template_name = "${var.name}-${count.index}"
	}
`, name)
}

func TestUnitAliCloudCmsDynamicTagGroup(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cms_dynamic_tag_group"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cms_dynamic_tag_group"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"contact_group_list": []string{"CreateDynamicTagGroupValue", "CreateDynamicTagGroupValue"},
		"tag_key":            "CreateDynamicTagGroupValue",
		"match_express": []map[string]interface{}{
			{
				"tag_value":                "CreateDynamicTagGroupValue",
				"tag_value_match_function": "CreateDynamicTagGroupValue",
			},
		},
		"template_id_list": []string{"CreateDynamicTagGroupValue", "CreateDynamicTagGroupValue"},
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
		// DescribeDynamicTagRuleList
		"TagGroupList": map[string]interface{}{
			"TagGroup": []interface{}{
				map[string]interface{}{
					"DynamicTagRuleId":           "CreateDynamicTagGroupValue",
					"Status":                     "CreateDynamicTagGroupValue",
					"MatchExpressFilterRelation": "CreateDynamicTagGroupValue",
					"TemplateIdList":             []string{"CreateDynamicTagGroupValue", "CreateDynamicTagGroupValue"},
					"MatchExpress": map[string]interface{}{
						"0": []interface{}{
							map[string]interface{}{
								"TagValue":              "CreateDynamicTagGroupValue",
								"TagValueMatchFunction": "CreateDynamicTagGroupValue",
							},
						},
					},
					"TagKey": "CreateDynamicTagGroupValue",
				},
			},
		},
		"Id": "CreateDynamicTagGroupValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateDynamicTagGroup
		"Id": "CreateDynamicTagGroupValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cms_dynamic_tag_group", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudCmsDynamicTagGroupCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDynamicTagRuleList Response
		"Id": "CreateDynamicTagGroupValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDynamicTagGroup" {
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
			return ReadMockResponse, nil
		})
		err := resourceAliCloudCmsDynamicTagGroupCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cms_dynamic_tag_group"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDynamicTagRuleList" {
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
			return ReadMockResponse, nil
		})
		err := resourceAliCloudCmsDynamicTagGroupRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudCmsDynamicTagGroupDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDynamicTagGroup" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"Success": true,
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudCmsDynamicTagGroupDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
