package arms_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/arms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_arms_alert_robot", &resource.Sweeper{
		Name: "alicloud_arms_alert_robot",
		F:    testSweepArmsAlertRobot,
	})
}

func testSweepArmsAlertRobot(region string) error {
	rawClient, err := armsTestClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	action := "DescribeIMRobots"
	request := make(map[string]interface{})
	request["Page"] = 1
	request["Size"] = helper.PageSizeXLarge
	response, err := client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
	if err != nil {
		log.Printf("[ERROR] %s failed error: %v", action, err)
		return nil
	}
	resp, err := jsonpath.Get("$.PageBean.AlertIMRobots", response)
	if err != nil {
		log.Printf("[ERROR] %s error: %v", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := fmt.Sprint(item["RobotName"])
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping arms alert robot: %s ", name)
			continue
		}
		log.Printf("[INFO] delete arms alert robot: %s ", name)
		action = "DeleteIMRobot"
		request = map[string]interface{}{
			"RobotId": fmt.Sprint(item["RobotId"]),
		}
		_, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] %s failed error: %v", action, err)
		}
	}
	return nil
}

func TestAccAliCloudARMSAlertRobot_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_alert_robot.default"
	ra := acctest.ResourceAttrInit(resourceId, ArmsAlertRobotMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsAlertRobot")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccArmsAlertRobot%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, ArmsAlertRobotBasicdependence)
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
					"alert_robot_name": "${var.name}",
					"robot_type":       "wechat",
					"robot_addr":       "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=1c704e23",
					"daily_noc":        true,
					"daily_noc_time":   "09:00,17:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_robot_name": name,
						"robot_type":       "wechat",
						"robot_addr":       "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=1c704e23",
						"daily_noc":        "true",
						"daily_noc_time":   "09:00,17:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_robot_name": "${var.name}_update",
					"robot_type":       "dingding",
					"robot_addr":       "https://oapi.dingtalk.com/robot/send?access_token=1c704e23",
					"daily_noc":        false,
					"daily_noc_time":   "09:30,17:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_robot_name": name + "_update",
						"robot_type":       "dingding",
						"robot_addr":       "https://oapi.dingtalk.com/robot/send?access_token=1c704e23",
						"daily_noc":        "false",
						"daily_noc_time":   "09:30,17:00",
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

var ArmsAlertRobotMap = map[string]string{}

func ArmsAlertRobotBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudARMSAlertRobot(t *testing.T) {
	p := unitTestProviderResources()
	d, _ := schema.InternalMap(p["alicloud_arms_alert_robot"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_arms_alert_robot"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	dMock, _ := schema.InternalMap(p["alicloud_arms_alert_robot"].Schema).Data(nil, nil)
	dMock.MarkNewResource()
	for key, value := range map[string]interface{}{
		"alert_robot_name": "alert_robot_name",
		"robot_type":       "wechat",
		"robot_addr":       "robot_addr",
		"daily_noc":        true,
		"daily_noc_time":   "09:30,17:00",
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
			"AlertIMRobots": []interface{}{
				map[string]interface{}{
					"Type":         "wechat",
					"RobotId":      "MockContactId",
					"RobotName":    "robot_name",
					"RobotAddr":    "robot_addr",
					"CreateTime":   "2024-11-19 10:37:05",
					"DailyNoc":     true,
					"DailyNocTime": "09:30,17:00",
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
			return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_arms_alert_robot", "MockContactRobotId"))
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
			result["RobotId"] = "MockRobotId"
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
		err := arms.ResourceAlertRobotCreate(d, rawClient)
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
		err := arms.ResourceAlertRobotCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		target := reflect.TypeOf(&client.Client{})
		data := func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		}

		patches := gomonkey.ApplyMethod(target, "DoRequest", data)
		err := arms.ResourceAlertRobotCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("MockRobotId")
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

		err := arms.ResourceAlertRobotUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateAlertRobotAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"alert_robot_name"} {
			switch p["alicloud_arms_alert_robot"].Schema[key].Type {
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
		resourceData1, _ := schema.InternalMap(p["alicloud_arms_alert_robot"].Schema).Data(nil, diff)
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
		err := arms.ResourceAlertRobotUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateAlertRobotNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"alert_robot_name"} {
			switch p["alicloud_arms_alert_robot"].Schema[key].Type {
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
		resourceData1, _ := schema.InternalMap(p["alicloud_arms_alert_robot"].Schema).Data(nil, diff)
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
		err := arms.ResourceAlertRobotUpdate(resourceData1, rawClient)
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
		err := arms.ResourceAlertRobotDelete(d, rawClient)
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
		err := arms.ResourceAlertRobotDelete(d, rawClient)
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
		err := arms.ResourceAlertRobotDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeArmsAlertRobotNotFound", func(t *testing.T) {
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
		err := arms.ResourceAlertRobotRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeArmsAlertRobotAbnormal", func(t *testing.T) {
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
		err := arms.ResourceAlertRobotRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}
