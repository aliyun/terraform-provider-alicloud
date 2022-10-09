package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAlicloudGPDBDBInstancePlan_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	checkoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudGPDBDBInstancePlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGPDBDBInstancePlanBasicDependence0)
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
					"db_instance_plan_name": "${var.name}",
					"plan_desc":             "${var.name}",
					"plan_type":             "PauseResume",
					"plan_schedule_type":    "Regular",
					"plan_start_date":       "${var.plan_start_date}",
					"plan_end_date":         "${var.plan_end_date}",
					"plan_config": []interface{}{
						map[string]interface{}{
							"resume": []interface{}{
								map[string]interface{}{
									"plan_cron_time": "0 0 0 1/1 * ? ",
								},
							},
							"pause": []interface{}{
								map[string]interface{}{
									"plan_cron_time": "0 0 10 1/1 * ? ",
								},
							},
						},
					},
					"db_instance_id": "${data.alicloud_gpdb_instances.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_plan_name":                 name,
						"plan_desc":                             name,
						"plan_type":                             "PauseResume",
						"plan_schedule_type":                    "Regular",
						"plan_start_date":                       CHECKSET,
						"plan_end_date":                         CHECKSET,
						"plan_config.#":                         "1",
						"plan_config.0.resume.#":                "1",
						"plan_config.0.resume.0.plan_cron_time": "0 0 0 1/1 * ? ",
						"plan_config.0.pause.#":                 "1",
						"plan_config.0.pause.0.plan_cron_time":  "0 0 10 1/1 * ? ",
						"db_instance_id":                        CHECKSET,
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

var AlicloudGPDBDBInstancePlanMap0 = map[string]string{
	"plan_type":             CHECKSET,
	"plan_id":               CHECKSET,
	"db_instance_id":        CHECKSET,
	"plan_config.#":         CHECKSET,
	"plan_schedule_type":    CHECKSET,
	"plan_desc":             CHECKSET,
	"db_instance_plan_name": CHECKSET,
	"status":                CHECKSET,
}

func AlicloudGPDBDBInstancePlanBasicDependence0(name string) string {
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")

	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "plan_start_date" {
  default = "%v"
}

variable "plan_end_date" {
  default = "%v"
}

data "alicloud_gpdb_instances" "default" {	
	name_regex = "default-NODELETING"
}
`, name, planStartDate, planEndDate)
}

func TestAccAlicloudGPDBDBInstancePlan_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	checkoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudGPDBDBInstancePlanMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGPDBDBInstancePlanBasicDependence0)
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
					"db_instance_plan_name": "${var.name}",
					"plan_type":             "PauseResume",
					"plan_schedule_type":    "Regular",
					"plan_config": []interface{}{
						map[string]interface{}{
							"resume": []interface{}{
								map[string]interface{}{
									"plan_cron_time": "0 0 0 1/1 * ? ",
								},
							},
							"pause": []interface{}{
								map[string]interface{}{
									"plan_cron_time": "0 0 10 1/1 * ? ",
								},
							},
						},
					},
					"db_instance_id": "${data.alicloud_gpdb_instances.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_plan_name":                 name,
						"plan_type":                             "PauseResume",
						"plan_schedule_type":                    "Regular",
						"plan_config.#":                         "1",
						"plan_config.0.resume.#":                "1",
						"plan_config.0.resume.0.plan_cron_time": "0 0 0 1/1 * ? ",
						"plan_config.0.pause.#":                 "1",
						"plan_config.0.pause.0.plan_cron_time":  "0 0 10 1/1 * ? ",
						"db_instance_id":                        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_plan_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_plan_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_desc": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_desc": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_start_date": "${var.plan_start_date}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_start_date": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_end_date": "${var.plan_end_date}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_end_date": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_config": []interface{}{
						map[string]interface{}{
							"resume": []interface{}{
								map[string]interface{}{
									"plan_cron_time": "0 0 1 1/1 * ? ",
								},
							},
							"pause": []interface{}{
								map[string]interface{}{
									"plan_cron_time": "0 0 11 1/1 * ? ",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_config.#":                         "1",
						"plan_config.0.resume.#":                "1",
						"plan_config.0.resume.0.plan_cron_time": "0 0 1 1/1 * ? ",
						"plan_config.0.pause.#":                 "1",
						"plan_config.0.pause.0.plan_cron_time":  "0 0 11 1/1 * ? ",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "cancel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "cancel",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "active",
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

var AlicloudGPDBDBInstancePlanMap1 = map[string]string{
	"db_instance_plan_name": CHECKSET,
	"plan_schedule_type":    CHECKSET,
	"plan_type":             CHECKSET,
	"status":                CHECKSET,
	"db_instance_id":        CHECKSET,
	"plan_config.#":         CHECKSET,
}

func TestAccAlicloudGPDBDBInstancePlan_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	checkoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudGPDBDBInstancePlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGPDBDBInstancePlanBasicDependence0)
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
					"db_instance_plan_name": "${var.name}",
					"plan_desc":             "${var.name}",
					"plan_type":             "Resize",
					"plan_schedule_type":    "Regular",
					"plan_start_date":       "${var.plan_start_date}",
					"plan_end_date":         "${var.plan_end_date}",
					"plan_config": []interface{}{
						map[string]interface{}{
							"scale_out": []interface{}{
								map[string]interface{}{
									"segment_node_num": "4",
									"plan_cron_time":   "0 0 0 1/1 * ? ",
								},
							},
							"scale_in": []interface{}{
								map[string]interface{}{
									"segment_node_num": "2",
									"plan_cron_time":   "0 0 10 1/1 * ? ",
								},
							},
						},
					},
					"db_instance_id": "${data.alicloud_gpdb_instances.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_plan_name":     name,
						"plan_desc":                 name,
						"plan_type":                 "Resize",
						"plan_schedule_type":        "Regular",
						"plan_start_date":           CHECKSET,
						"plan_end_date":             CHECKSET,
						"plan_config.#":             "1",
						"plan_config.0.scale_out.#": "1",
						"plan_config.0.scale_out.0.plan_cron_time":   "0 0 0 1/1 * ? ",
						"plan_config.0.scale_out.0.segment_node_num": "4",
						"plan_config.0.scale_in.#":                   "1",
						"plan_config.0.scale_in.0.plan_cron_time":    "0 0 10 1/1 * ? ",
						"plan_config.0.scale_in.0.segment_node_num":  "2",
						"db_instance_id":                             CHECKSET,
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

func TestUnitAccAlicloudGpdbDbInstancePlan(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"db_instance_plan_name": "CreateGpdbDbInstancePlanValue",
		"plan_desc":             "CreateGpdbDbInstancePlanValue",
		"plan_type":             "CreateGpdbDbInstancePlanValue",
		"plan_schedule_type":    "CreateGpdbDbInstancePlanValue",
		"plan_start_date":       "CreateGpdbDbInstancePlanValue",
		"plan_end_date":         "CreateGpdbDbInstancePlanValue",
		"plan_config": []interface{}{
			map[string]interface{}{
				"resume": []interface{}{
					map[string]interface{}{
						"plan_cron_time": "CreateGpdbDbInstancePlanValue",
						"execute_time":   "CreateGpdbDbInstancePlanValue",
					},
				},
				"pause": []interface{}{
					map[string]interface{}{
						"plan_cron_time": "CreateGpdbDbInstancePlanValue",
						"execute_time":   "CreateGpdbDbInstancePlanValue",
					},
				},
				"scale_in": []interface{}{
					map[string]interface{}{
						"plan_cron_time":   "CreateGpdbDbInstancePlanValue",
						"execute_time":     "CreateGpdbDbInstancePlanValue",
						"segment_node_num": "2",
					},
				},
				"scale_out": []interface{}{
					map[string]interface{}{
						"plan_cron_time":   "CreateGpdbDbInstancePlanValue",
						"execute_time":     "CreateGpdbDbInstancePlanValue",
						"segment_node_num": "2",
					},
				},
			},
		},
		"db_instance_id": "CreateGpdbDbInstancePlanValue",
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
		"Items": map[string]interface{}{
			"PlanList": []interface{}{
				map[string]interface{}{
					"DBInstanceId":     "CreateGpdbDbInstancePlanValue",
					"PlanId":           "CreateGpdbDbInstancePlanValue",
					"PlanScheduleType": "CreateGpdbDbInstancePlanValue",
					"PlanType":         "CreateGpdbDbInstancePlanValue",
					"PlanDesc":         "CreateGpdbDbInstancePlanValue",
					"PlanName":         "CreateGpdbDbInstancePlanValue",
					"PlanStartDate":    "CreateGpdbDbInstancePlanValue",
					"PlanEndDate":      "CreateGpdbDbInstancePlanValue",
					"PlanConfig":       "{\"resume\":{\"planCronTime\":\"CreateGpdbDbInstancePlanValue\",\"executeTime\":\"CreateGpdbDbInstancePlanValue\"},\"pause\":{\"planCronTime\":\"CreateGpdbDbInstancePlanValue\",\"executeTime\":\"CreateGpdbDbInstancePlanValue\"},\"scaleOut\":{\"planCronTime\":\"CreateGpdbDbInstancePlanValue\",\"executeTime\":\"CreateGpdbDbInstancePlanValue\",\"segmentNodeNum\":\"2\"},\"scaleIn\":{\"planCronTime\":\"CreateGpdbDbInstancePlanValue\",\"executeTime\":\"CreateGpdbDbInstancePlanValue\",\"segmentNodeNum\":\"2\"}}",
					"PlanStatus":       "active",
				},
			},
		},
		"Status": "success",
	}
	CreateMockResponse := map[string]interface{}{
		"PlanId": "CreateGpdbDbInstancePlanValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_gpdb_db_instance_plan", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGpdbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudGpdbDbInstancePlanCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDBInstancePlan" {
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
		err := resourceAlicloudGpdbDbInstancePlanCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGpdbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudGpdbDbInstancePlanUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"db_instance_plan_name": "UpdateGpdbDbInstancePlanValue",
		"plan_desc":             "UpdateGpdbDbInstancePlanValue",
		"plan_start_date":       "UpdateGpdbDbInstancePlanValue",
		"plan_end_date":         "UpdateGpdbDbInstancePlanValue",
		"status":                "cancel",
		"plan_config": []interface{}{
			map[string]interface{}{
				"resume": []interface{}{
					map[string]interface{}{
						"plan_cron_time": "UpdateGpdbDbInstancePlanValue",
						"execute_time":   "UpdateGpdbDbInstancePlanValue",
					},
				},
				"pause": []interface{}{
					map[string]interface{}{
						"plan_cron_time": "UpdateGpdbDbInstancePlanValue",
						"execute_time":   "UpdateGpdbDbInstancePlanValue",
					},
				},
				"scale_in": []interface{}{
					map[string]interface{}{
						"plan_cron_time":   "UpdateGpdbDbInstancePlanValue",
						"execute_time":     "UpdateGpdbDbInstancePlanValue",
						"segment_node_num": "2",
					},
				},
				"scale_out": []interface{}{
					map[string]interface{}{
						"plan_cron_time":   "UpdateGpdbDbInstancePlanValue",
						"execute_time":     "UpdateGpdbDbInstancePlanValue",
						"segment_node_num": "2",
					},
				},
			},
		},
	}
	diff, err := newInstanceDiff("alicloud_gpdb_db_instance_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Items": map[string]interface{}{
			"PlanList": []interface{}{
				map[string]interface{}{
					"PlanDesc":      "UpdateGpdbDbInstancePlanValue",
					"PlanName":      "UpdateGpdbDbInstancePlanValue",
					"PlanStartDate": "UpdateGpdbDbInstancePlanValue",
					"PlanEndDate":   "UpdateGpdbDbInstancePlanValue",
					"PlanConfig":    "{\"resume\":{\"planCronTime\":\"UpdateGpdbDbInstancePlanValue\",\"executeTime\":\"UpdateGpdbDbInstancePlanValue\"},\"pause\":{\"planCronTime\":\"UpdateGpdbDbInstancePlanValue\",\"executeTime\":\"UpdateGpdbDbInstancePlanValue\"},\"scaleOut\":{\"planCronTime\":\"UpdateGpdbDbInstancePlanValue\",\"executeTime\":\"UpdateGpdbDbInstancePlanValue\",\"segmentNodeNum\":\"2\"},\"scaleIn\":{\"planCronTime\":\"UpdateGpdbDbInstancePlanValue\",\"executeTime\":\"UpdateGpdbDbInstancePlanValue\",\"segmentNodeNum\":\"2\"}}",
					"PlanStatus":    "cancel",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateDBInstancePlan" {
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
			if *action == "SetDBInstancePlanStatus" {
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
			return ReadMockResponse, nil
		})
		err := resourceAlicloudGpdbDbInstancePlanUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_gpdb_db_instance_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDBInstancePlans" {
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
		err := resourceAlicloudGpdbDbInstancePlanRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGpdbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudGpdbDbInstancePlanDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_gpdb_db_instance_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDBInstancePlan" {
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
		err := resourceAlicloudGpdbDbInstancePlanDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
