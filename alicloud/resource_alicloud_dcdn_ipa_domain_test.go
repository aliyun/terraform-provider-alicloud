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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_dcdn_ipa_domain",
		&resource.Sweeper{
			Name: "alicloud_dcdn_ipa_domain",
			F:    testSweepDcdnIpaDomain,
		})
}
func testSweepDcdnIpaDomain(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.DCDNSupportRegions) {
		log.Printf("[INFO] Skipping Dcdn Ipa Domain unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDcdnIpaUserDomains"
	request := map[string]interface{}{}

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	conn, err := aliyunClient.NewDcdnClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
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

		resp, err := jsonpath.Get("$.Domains.PageData", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Domains.PageData", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DomainName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Dcdn Ipa Domain: %s", item["DomainName"].(string))
				continue
			}
			action := "DeleteDcdnIpaDomain"
			request := map[string]interface{}{
				"DomainName": item["DomainName"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Dcdn Ipa Domain (%s): %s", item["DomainName"].(string), err)
			}
			log.Printf("[INFO] Delete Dcdn Ipa Domain success: %s ", item["DomainName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}
func TestAccAlicloudDCDNIpaDomain_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_ipa_domain.default"
	checkoutSupportedRegions(t, true, connectivity.DCDNSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudDCDNIpaDomainMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnIpaDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacccn-%d.xiaozhu.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDCDNIpaDomainBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEnterpriseAccountEnabled(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": "${var.domain_name}",
					"sources": []map[string]interface{}{
						{
							"content":  "www.xiaozhu.com",
							"port":     "8898",
							"priority": "20",
							"type":     "domain",
							"weight":   "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"domain_name":       name,
						"sources.#":         "1",
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
					"status": "offline",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "offline",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "online",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "online",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sources": []map[string]interface{}{
						{
							"content":  "3.3.3.3",
							"port":     "7261",
							"priority": "20",
							"type":     "ipaddr",
							"weight":   "7",
						},
						{
							"content":  "5.3.3.3",
							"port":     "7221",
							"priority": "20",
							"type":     "ipaddr",
							"weight":   "9",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.#": "2",
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
func TestAccAlicloudDCDNIpaDomain_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_ipa_domain.default"
	checkoutSupportedRegions(t, true, connectivity.DCDNSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudDCDNIpaDomainMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnIpaDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacccn-%d.xiaozhu.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDCDNIpaDomainBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEnterpriseAccountEnabled(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": "${var.domain_name}",
					"sources": []map[string]interface{}{
						{
							"content":  "www.xiaozhu.com",
							"port":     "8898",
							"priority": "20",
							"type":     "domain",
							"weight":   "10",
						},
					},
					"scope":             "global",
					"status":            "online",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"domain_name":       name,
						"sources.#":         "1",
						"status":            "online",
						"scope":             "global",
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

var AlicloudDCDNIpaDomainMap0 = map[string]string{
	"domain_name":       CHECKSET,
	"resource_group_id": CHECKSET,
	"status":            CHECKSET,
	"sources.#":         CHECKSET,
}

func AlicloudDCDNIpaDomainBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "domain_name" {	
	default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
	name_regex = "default"
}
`, name)
}

func TestUnitAlicloudDCDNIpaDomain(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"domain_name": "domain_name",
		"sources": []map[string]interface{}{
			{
				"content":  "content",
				"port":     8888,
				"priority": "priority",
				"type":     "type",
				"weight":   10,
			},
		},
		"scope":             "scope",
		"resource_group_id": "resource_group_id",
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"DomainDetail": map[string]interface{}{
			"Cname":           "domain_name",
			"ResourceGroupId": "resource_group_id",
			"DomainStatus":    "online",
			"DomainName":      "domain_name",
			"Sources": map[string]interface{}{
				"Source": []interface{}{
					map[string]interface{}{
						"Enabled":  "online",
						"Port":     8888,
						"Type":     "type",
						"Content":  "content",
						"Priority": "priority",
						"Weight":   "10",
					},
				},
			},
			"GmtModified": "gmtModified",
			"GmtCreated":  "gmtCreated",
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_dcdn_ipa_domain", "xxx_id"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateOfflineStatusNormal": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"DomainDetail": map[string]interface{}{
					"Cname":           "domain_name",
					"ResourceGroupId": "resource_group_id",
					"DomainStatus":    "offline",
					"DomainName":      "domain_name",
					"Sources": map[string]interface{}{
						"Source": []interface{}{
							map[string]interface{}{
								"Enabled":  "offline",
								"Port":     8888,
								"Type":     "type",
								"Content":  "content",
								"Priority": "priority",
								"Weight":   "10",
							},
						},
					},
					"GmtModified": "gmtModified",
					"GmtCreated":  "gmtCreated",
				},
			}
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDcdnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudDcdnIpaDomainCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("CreateRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNoRetryErrorError", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["CreateNormal"]("")
		})
		patchDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DcdnService{}), "DescribeDcdnIpaDomain", func(*DcdnService, string) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := resourceAlicloudDcdnIpaDomainCreate(d, rawClient)
		patches.Reset()
		patchDescribe.Reset()
		assert.NotNil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("domain_name")

	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDcdnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAlicloudDcdnIpaDomainUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"resource_group_id"} {
			switch p["alicloud_dcdn_ipa_domain"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyAttributeNoRetryErrorAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"resource_group_id"} {
			switch p["alicloud_dcdn_ipa_domain"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DcdnService{}), "DescribeDcdnIpaDomain", func(*DcdnService, string) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData, rawClient)
		patches.Reset()
		patcheDescribe.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"resource_group_id", "sources"} {
			switch p["alicloud_dcdn_ipa_domain"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)})
			case schema.TypeSet:
				diff.SetAttribute("sources.0.content", &terraform.ResourceAttrDiff{Old: "content", New: "content_update"})
				diff.SetAttribute("sources.0.port", &terraform.ResourceAttrDiff{Old: "8888", New: "8888"})
				diff.SetAttribute("sources.0.priority", &terraform.ResourceAttrDiff{Old: "priority", New: "priority"})
				diff.SetAttribute("sources.0.type", &terraform.ResourceAttrDiff{Old: "type", New: "type"})
				diff.SetAttribute("sources.0.weight", &terraform.ResourceAttrDiff{Old: "10", New: "9"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("UpdateModifyStatusOfflineAttributeNotFoundErrorAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.SetAttribute("status", &terraform.ResourceAttrDiff{Old: "online", New: "offline"})
		resourceData1, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["NotFoundError"]("ResourceNotfound")
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyStatusOfflineAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.SetAttribute("status", &terraform.ResourceAttrDiff{Old: "online", New: "offline"})
		resourceData1, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		patchDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DcdnService{}), "DescribeDcdnIpaDomain", func(*DcdnService, string) (map[string]interface{}, error) {
			object := map[string]interface{}{
				"Cname":           "domain_name",
				"ResourceGroupId": "resource_group_id",
				"DomainStatus":    "online",
				"DomainName":      "domain_name",
				"Sources": map[string]interface{}{
					"Source": []interface{}{
						map[string]interface{}{
							"Enabled":  "online",
							"Port":     8888,
							"Type":     "type",
							"Content":  "content",
							"Priority": "priority",
							"Weight":   "10",
						},
					},
				},
				"GmtModified": "gmtModified",
				"GmtCreated":  "gmtCreated",
			}
			return object, nil
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData1, rawClient)
		patches.Reset()
		patchDescribe.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyStatusOnlineAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.SetAttribute("status", &terraform.ResourceAttrDiff{Old: "offline", New: "online"})

		resourceData1, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		patchDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DcdnService{}), "DescribeDcdnIpaDomain", func(*DcdnService, string) (map[string]interface{}, error) {
			object := map[string]interface{}{
				"Cname":           "domain_name",
				"ResourceGroupId": "resource_group_id",
				"DomainStatus":    "offline",
				"DomainName":      "domain_name",
				"Sources": map[string]interface{}{
					"Source": []interface{}{
						map[string]interface{}{
							"Enabled":  "offline",
							"Port":     8888,
							"Type":     "type",
							"Content":  "content",
							"Priority": "priority",
							"Weight":   "10",
						},
					},
				},
				"GmtModified": "gmtModified",
				"GmtCreated":  "gmtCreated",
			}
			return object, nil
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData1, rawClient)
		patches.Reset()
		patchDescribe.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyStatusOfflineAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.SetAttribute("status", &terraform.ResourceAttrDiff{Old: "online", New: "offline"})
		resourceData1, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		patchDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DcdnService{}), "DcdnIpaDomainStateRefreshFunc", func(*DcdnService, string, []string) resource.StateRefreshFunc {
			return func() (interface{}, string, error) {
				object := map[string]interface{}{
					"Cname":           "domain_name",
					"ResourceGroupId": "resource_group_id",
					"DomainStatus":    "offline",
					"DomainName":      "domain_name",
					"Sources": map[string]interface{}{
						"Source": []interface{}{
							map[string]interface{}{
								"Enabled":  "offline",
								"Port":     8888,
								"Type":     "type",
								"Content":  "content",
								"Priority": "priority",
								"Weight":   "10",
							},
						},
					},
					"GmtModified": "gmtModified",
					"GmtCreated":  "gmtCreated",
				}
				return object, "offline", nil
			}
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData1, rawClient)
		patches.Reset()
		patchDescribe.Reset()
		assert.Nil(t, err)
	})
	t.Run("UpdateModifyStatusOnlineAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.SetAttribute("status", &terraform.ResourceAttrDiff{Old: "offline", New: "online"})
		resourceData1, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateOfflineStatusNormal"]("")
		})
		patchDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DcdnService{}), "DcdnIpaDomainStateRefreshFunc", func(*DcdnService, string, []string) resource.StateRefreshFunc {
			return func() (interface{}, string, error) {
				object := map[string]interface{}{
					"Cname":           "domain_name",
					"ResourceGroupId": "resource_group_id",
					"DomainStatus":    "online",
					"DomainName":      "domain_name",
					"Sources": map[string]interface{}{
						"Source": []interface{}{
							map[string]interface{}{
								"Enabled":  "online",
								"Port":     8888,
								"Type":     "type",
								"Content":  "content",
								"Priority": "priority",
								"Weight":   "10",
							},
						},
					},
					"GmtModified": "gmtModified",
					"GmtCreated":  "gmtCreated",
				}
				return object, "online", nil
			}
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData1, rawClient)
		patches.Reset()
		patchDescribe.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDcdnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudDcdnIpaDomainDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		patchDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DcdnService{}), "DescribeDcdnIpaDomain", func(*DcdnService, string) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := resourceAlicloudDcdnIpaDomainDelete(d, rawClient)
		patches.Reset()
		patchDescribe.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteNonRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeNotFound", func(t *testing.T) {
		patchRequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainRead(d, rawClient)
		patchRequest.Reset()
		assert.Nil(t, err)
	})
	t.Run("ReadDescribeAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})
}
