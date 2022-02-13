package alicloud

import (
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_pvtz_endpoint",
		&resource.Sweeper{
			Name: "alicloud_pvtz_endpoint",
			F:    testSweepPrivateZoneEndpoint,
		})
}

func testSweepPrivateZoneEndpoint(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}
	action := "DescribeResolverEndpoints"
	request := map[string]interface{}{}

	request["Lang"] = "en"
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	conn, err := client.NewPvtzClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
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

		resp, err := jsonpath.Get("$.Endpoints", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Endpoints", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["Name"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping PrivateZone Endpoint: %s", item["Name"].(string))
				continue
			}
			action := "DeleteResolverEndpoint"
			request := map[string]interface{}{
				"EndpointId": item["Id"],
			}
			request["ClientToken"] = buildClientToken(action)
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				log.Printf("[ERROR] Failed to delete PrivateZone Endpoint (%s): %s", item["Id"].(string), err)
			}
			log.Printf("[INFO] Delete PrivateZone Endpoint success: %s ", item["Id"].(string))
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudPrivateZoneEndpoint_basic0(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.PvtzSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_pvtz_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateZoneEndpointMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePvtzEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateZoneEndpointBasicDependence0)
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
					"vpc_id":            "${alicloud_vpc.default.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"vpc_region_id":     defaultRegionToTest,
					"endpoint_name":     name,
					"ip_configs": []map[string]interface{}{
						{
							"zone_id":    "${alicloud_vswitch.default[0].zone_id}",
							"cidr_block": "${alicloud_vswitch.default[0].cidr_block}",
							"vswitch_id": "${alicloud_vswitch.default[0].id}",
						},
						{
							"zone_id":    "${alicloud_vswitch.default[1].zone_id}",
							"cidr_block": "${alicloud_vswitch.default[1].cidr_block}",
							"vswitch_id": "${alicloud_vswitch.default[1].id}",
						},
						{
							"zone_id":    "${alicloud_vswitch.default[2].zone_id}",
							"cidr_block": "${alicloud_vswitch.default[2].cidr_block}",
							"vswitch_id": "${alicloud_vswitch.default[2].id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":            CHECKSET,
						"security_group_id": CHECKSET,
						"vpc_region_id":     CHECKSET,
						"endpoint_name":     name,
						"ip_configs.#":      "3",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"endpoint_name": name + "update",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"endpoint_name": name + "update",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_configs": []map[string]interface{}{
						{
							"zone_id":    "${alicloud_vswitch.default[2].zone_id}",
							"cidr_block": "${alicloud_vswitch.default[2].cidr_block}",
							"vswitch_id": "${alicloud_vswitch.default[2].id}",
						},
						{
							"zone_id":    "${alicloud_vswitch.default[3].zone_id}",
							"cidr_block": "${alicloud_vswitch.default[3].cidr_block}",
							"vswitch_id": "${alicloud_vswitch.default[3].id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_configs.#": "2",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"endpoint_name": name,
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"endpoint_name": name,
			//		}),
			//	),
			//},
			//{
			//	ResourceName:            resourceId,
			//	ImportState:             true,
			//	ImportStateVerify:       true,
			//	ImportStateVerifyIgnore: []string{},
			//},
		},
	})
}

var AlicloudPrivateZoneEndpointMap0 = map[string]string{}

func AlicloudPrivateZoneEndpointBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_pvtz_resolver_zones" "default" {
  status = "NORMAL"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_security_group" "default" {
  vpc_id      = alicloud_vpc.default.id
  name        = var.name
}

resource "alicloud_vswitch" "default" {
  count      = 4
  vpc_id     = alicloud_vpc.default.id
  cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, count.index)
  zone_id    = data.alicloud_pvtz_resolver_zones.default.zones[count.index].zone_id
}
`, name)
}

func TestAccAlicloudPrivateZoneEndpoint_unit(t *testing.T) {
	resourceName := "alicloud_pvtz_endpoint"
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p[resourceName].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p[resourceName].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	attributes := map[string]interface{}{
		"endpoint_name":     "endpoint_name",
		"security_group_id": "security_group_id",
		"vpc_id":            "vpc_id",
		"vpc_region_id":     "vpc_region_id",
		"ip_configs": []map[string]interface{}{
			{
				"zone_id":    "cn-hangzhou-a",
				"cidr_block": "10.0.0.0/16",
				"vswitch_id": "vswitch_id_0",
				"ip":         "172.16.0.0",
			},
		},
	}
	for key, value := range attributes {
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
		"Status":          "SUCCESS",
		"SecurityGroupId": "security_group_id",
		"Name":            "endpoint_name",
		"VpcId":           "vpc_id",
		"VpcRegionId":     "vpc_region_id",
		"IpConfigs": []interface{}{
			map[string]interface{}{
				"VSwitchId": "vswitch_id_0",
				"AzId":      "cn-hangzhou-a",
				"CidrBlock": "10.0.0.0/16",
				"Ip":        "172.16.0.0",
			},
		},
	}
	ReadMockUpdateResponse := map[string]interface{}{
		"Status":          "SUCCESS",
		"SecurityGroupId": "security_group_id",
		"Name":            "endpoint_name_update",
		"VpcId":           "vpc_id",
		"VpcRegionId":     "vpc_region_id",
		"IpConfigs": []interface{}{
			map[string]interface{}{
				"VSwitchId": "vswitch_id_0_update",
				"AzId":      "cn-hangzhou-a_update",
				"CidrBlock": "10.0.0.0/16_update",
				"Ip":        "172.16.0.0_update",
			},
		},
	}
	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["EndpointId"] = "MockId"
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockUpdateResponse
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

	t.Run("Create", func(t *testing.T) {
		// Client Error
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPvtzClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err = resourceAlicloudPvtzEndpointCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)

		retryFlag, noRetryFlag, abnormal := true, true, true
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		for {
			err = resourceAlicloudPvtzEndpointCreate(dCreate, rawClient)
			if abnormal {
				abnormal = false
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				for key, _ := range attributes {
					assert.False(t, dCreate.HasChange(key))
				}
				break
			}
		}
		patches.Reset()
	})

	t.Run("Update", func(t *testing.T) {
		rand1 := acctest.RandIntRange(1000, 9999)
		// Client Error
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPvtzClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err = resourceAlicloudPvtzEndpointUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)

		// DoRequest
		retryFlag, noRetryFlag, abnormal := false, false, false
		targetFunc := "UpdateResolverEndpoint"

		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if targetFunc == *action && retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if targetFunc == *action && noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})

		diff := terraform.NewInstanceDiff()

		for _, key := range []string{"endpoint_name", "ip_configs"} {
			switch p[resourceName].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(2)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			case schema.TypeSet:
				diff.SetAttribute(fmt.Sprintf("%s.#", key), &terraform.ResourceAttrDiff{Old: "1", New: "1"})
				for _, ipConfig := range d.Get(key).(*schema.Set).List() {
					ipConfigArg := ipConfig.(map[string]interface{})
					for field, _ := range p[resourceName].Schema[key].Elem.(*schema.Resource).Schema {
						diff.SetAttribute(fmt.Sprintf("%s.%d.%s", key, rand1, field), &terraform.ResourceAttrDiff{Old: ipConfigArg[field].(string), New: ipConfigArg[field].(string) + "_update"})
					}
				}
			}
		}



		dCreate, err := schema.InternalMap(p[resourceName].Schema).Data(dCreate.State(), diff)
		//dCreate.

		dCreate.SetId("MockId")

		for {
			err = resourceAlicloudPvtzEndpointUpdate(dCreate, rawClient)
			if abnormal {
				abnormal = false
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				for key, _ := range attributes {
					assert.False(t, dCreate.HasChange(key))
				}
				break
			}
		}
		patches.Reset()

	})

}
