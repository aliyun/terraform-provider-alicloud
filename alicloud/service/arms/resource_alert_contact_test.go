package arms_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/arms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_arms_alert_contact", &resource.Sweeper{
		Name: "alicloud_arms_alert_contact",
		F:    testSweepArmsAlertContact,
	})
}

func testSweepArmsAlertContact(region string) error {
	rawClient, err := armsTestClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting Alicloud client.")
	}
	return sweepArmsAlertContact(rawClient.(*connectivity.AliyunClient))
}

// sweepArmsAlertContact cleans both test-created and documentation-example
// contacts: the documented example names its contact with the
// "terraform-example" prefix, so that prefix must be recognized as well.
func sweepArmsAlertContact(client *connectivity.AliyunClient) error {
	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
		"terraform-example",
	}

	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = helper.PageSizeLarge
	request["CurrentPage"] = 1
	var response map[string]interface{}
	var err error
	for {
		action := "SearchAlertContact"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(2*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %s", action, err)
			return nil
		}
		resp, err := jsonpath.Get("$.PageBean.Contacts", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.PageBean.Contacts", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["ContactName"])
			skip := true
			if !armsSweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping arms alert contact: %s ", name)
					continue
				}
			}
			log.Printf("[INFO] delete arms alert contact: %s ", name)
			action = "DeleteAlertContact"
			request := map[string]interface{}{
				"ContactId": fmt.Sprint(item["ContactId"]),
				"RegionId":  client.RegionId,
			}

			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete alarm contact (%s): %s", name, err)
			}
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	return nil
}

// TestUnitArmsAlertContactSweeper pins that the alert-contact sweeper, which
// also cleans the contact created by the notification-policy documentation
// example, recognizes the documented "terraform-example" prefix without
// widening cleanup to unrelated contacts.
func TestUnitArmsAlertContactSweeper(t *testing.T) {
	var mu sync.Mutex
	deleted := map[string]bool{}
	// ARMS test client shared with the notification-policy unit tests.
	client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "SearchAlertContact":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"PageBean": map[string]interface{}{"Contacts": []map[string]interface{}{
					{"ContactId": float64(11), "ContactName": "tf-testAcc-contact"},
					{"ContactId": float64(12), "ContactName": "terraform-example"},
					{"ContactId": float64(13), "ContactName": "prod-contact"},
				}},
			})
		case "DeleteAlertContact":
			mu.Lock()
			deleted[r.Form.Get("ContactId")] = true
			mu.Unlock()
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"IsSuccess": true})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})
	if err := sweepArmsAlertContact(client); err != nil {
		t.Fatalf("sweeper returned an error: %s", err)
	}
	if !deleted["11"] || !deleted["12"] {
		t.Fatalf("test and documented-example prefixes must be swept: %v", deleted)
	}
	if deleted["13"] {
		t.Fatalf("unrelated contact must be retained: %v", deleted)
	}
}

func TestAccAliCloudARMSAlertContact_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_alert_contact.default"
	ra := acctest.ResourceAttrInit(resourceId, ArmsAlertContactMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsAlertContact")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccArmsAlertContact%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, ArmsAlertContactBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_contact_name": "${var.name}",
					"email":              "hello.uuuu@aaa.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_contact_name": name,
						"email":              "hello.uuuu@aaa.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_contact_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_contact_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email": "helloupdate.uuuu@aaa.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email": "helloupdate.uuuu@aaa.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"phone_num": "12345678900",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"phone_num": "12345678900",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ding_robot_webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=91f2f6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ding_robot_webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=91f2f6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_noc": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_noc": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_contact_name":     "${var.name}",
					"ding_robot_webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=91f2f7",
					"email":                  "hello.uuuu@aaa.com",
					"phone_num":              "12345678901",
					"system_noc":             "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_contact_name":     name,
						"ding_robot_webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=91f2f7",
						"email":                  "hello.uuuu@aaa.com",
						"phone_num":              "12345678901",
						"system_noc":             "false",
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

var ArmsAlertContactMap = map[string]string{
	"ding_robot_webhook_url": "",
	"email":                  "",
	"phone_num":              "",
	"system_noc":             "false",
}

func ArmsAlertContactBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudARMSAlertContact(t *testing.T) {
	p := unitTestProviderResources()
	d, _ := schema.InternalMap(p["alicloud_arms_alert_contact"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_arms_alert_contact"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"alert_contact_name":     "alert_contact_name",
		"ding_robot_webhook_url": "ding_robot_webhook_url",
		"email":                  "email",
		"phone_num":              "phone_num",
		"system_noc":             true,
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := armsTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"PageBean": map[string]interface{}{
			"Contacts": []interface{}{
				map[string]interface{}{
					"ContactId":   "MockContactId",
					"ContactName": "alert_contact_name",
					"DingRobot":   "ding_robot_webhook_url",
					"Email":       "email",
					"Phone":       "phone_num",
					"SystemNoc":   true,
				},
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       helper.String(errorCode),
				Data:       helper.String(errorCode),
				Message:    helper.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_arms_alert_contact", "MockContactId"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       helper.String(errorCode),
				Data:       helper.String(errorCode),
				Message:    helper.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["ContactId"] = "MockContactId"
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewArmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := arms.ResourceAlertContactCreate(d, rawClient)
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
		err := arms.ResourceAlertContactCreate(d, rawClient)
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
		err := arms.ResourceAlertContactCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("MockContactId")
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewArmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := arms.ResourceAlertContactUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdatePrometheusAlertRuleAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"alert_contact_name", "ding_robot_webhook_url", "email", "phone_num", "system_noc"} {
			switch p["alicloud_arms_alert_contact"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(2)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_arms_alert_contact"].Schema).Data(nil, diff)
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
		err := arms.ResourceAlertContactUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdatePrometheusAlertRuleNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"alert_contact_name", "ding_robot_webhook_url", "email", "phone_num", "system_noc"} {
			switch p["alicloud_arms_alert_contact"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_arms_alert_contact"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
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
			return responseMock["UpdateNormal"]("")
		})
		err := arms.ResourceAlertContactUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewArmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := arms.ResourceAlertContactDelete(d, rawClient)
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
		err := arms.ResourceAlertContactDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
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
			return responseMock["DeleteNormal"]("")
		})
		err := arms.ResourceAlertContactDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeArmsAlertContactNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := arms.ResourceAlertContactRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeArmsAlertContactAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := arms.ResourceAlertContactRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})

}
