package alb_test

// collection-order: bypass health_check_codes -- pre-existing TypeList attribute; pure directory migration relocates the file into alicloud/service/alb without changing schema, CRUD or list-order handling — stock attribute not added or changed by this PR

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudAlbHealthCheckTemplate_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_health_check_template.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAlbHealthCheckTemplateMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewAlbService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlbHealthCheckTemplate")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbhealthchecktemplate%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAlbHealthCheckTemplateBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_template_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_template_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_template_name": name + "测试",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_template_name": name + "测试",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_connect_port": "8080",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_connect_port": "8080",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_host": "tf-testAcc.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_host": "tf-testAcc.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_http_version": "HTTP1.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_http_version": "HTTP1.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_method": "GET",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_method": "GET",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_path": "/tf-testAcc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_path": "/tf-testAcc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_timeout": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_timeout": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"healthy_threshold": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"unhealthy_threshold": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unhealthy_threshold": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_codes": []string{"http_2xx", "http_3xx", "http_4xx"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_codes.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_host":         acctest.REMOVEKEY,
					"health_check_http_version": acctest.REMOVEKEY,
					"health_check_method":       acctest.REMOVEKEY,
					"health_check_path":         acctest.REMOVEKEY,
					"health_check_protocol":     "TCP",
					"health_check_codes":        acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_host":         acctest.REMOVEKEY,
						"health_check_http_version": acctest.REMOVEKEY,
						"health_check_method":       acctest.REMOVEKEY,
						"health_check_path":         acctest.REMOVEKEY,
						"health_check_protocol":     "TCP",
						"health_check_codes.#":      "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_host":         "tf-testAcc.com",
					"health_check_http_version": "HTTP1.1",
					"health_check_method":       "GET",
					"health_check_path":         "/tf-testAcc",
					"health_check_protocol":     "HTTP",
					"health_check_codes":        []string{"http_2xx", "http_3xx", "http_4xx"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_host":         "tf-testAcc.com",
						"health_check_http_version": "HTTP1.1",
						"health_check_method":       "GET",
						"health_check_path":         "/tf-testAcc",
						"health_check_protocol":     "HTTP",
						"health_check_codes.#":      "3",
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

func TestAccAliCloudAlbHealthCheckTemplate_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_health_check_template.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAlbHealthCheckTemplateMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewAlbService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlbHealthCheckTemplate")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbhealthchecktemplate%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAlbHealthCheckTemplateBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_template_name": name,
					"health_check_connect_port":  "8080",
					"health_check_host":          "tf-testAcc.com",
					"health_check_http_version":  "HTTP1.0",
					"health_check_interval":      "20",
					"health_check_method":        "GET",
					"health_check_path":          "/tf-testAcc",
					"health_check_protocol":      "HTTP",
					"health_check_timeout":       "60",
					"healthy_threshold":          "6",
					"unhealthy_threshold":        "6",
					"health_check_codes":         []string{"http_2xx", "http_3xx", "http_4xx"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_template_name": name,
						"health_check_connect_port":  "8080",
						"health_check_host":          "tf-testAcc.com",
						"health_check_http_version":  "HTTP1.0",
						"health_check_interval":      "20",
						"health_check_method":        "GET",
						"health_check_path":          "/tf-testAcc",
						"health_check_protocol":      "HTTP",
						"health_check_timeout":       "60",
						"healthy_threshold":          "6",
						"unhealthy_threshold":        "6",
						"health_check_codes.#":       "3",
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

func TestAccAliCloudAlbHealthCheckTemplate_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_health_check_template.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAlbHealthCheckTemplateMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewAlbService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlbHealthCheckTemplate")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbhealthchecktemplate%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAlbHealthCheckTemplateBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_template_name": name,
					"health_check_protocol":      "TCP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_template_name": name,
						"health_check_protocol":      "TCP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_template_name": name + "测试",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_template_name": name + "测试",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_connect_port": "8080",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_connect_port": "8080",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_timeout": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_timeout": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"healthy_threshold": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"unhealthy_threshold": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unhealthy_threshold": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_host":         "tf-testAcc.com",
					"health_check_http_version": "HTTP1.1",
					"health_check_method":       "GET",
					"health_check_path":         "/tf-testAcc",
					"health_check_protocol":     "HTTP",
					"health_check_codes":        []string{"http_2xx", "http_3xx", "http_4xx"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_host":         "tf-testAcc.com",
						"health_check_http_version": "HTTP1.1",
						"health_check_method":       "GET",
						"health_check_path":         "/tf-testAcc",
						"health_check_protocol":     "HTTP",
						"health_check_codes.#":      "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_host":         acctest.REMOVEKEY,
					"health_check_http_version": acctest.REMOVEKEY,
					"health_check_method":       acctest.REMOVEKEY,
					"health_check_path":         acctest.REMOVEKEY,
					"health_check_protocol":     "TCP",
					"health_check_codes":        acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_host":         acctest.REMOVEKEY,
						"health_check_http_version": acctest.REMOVEKEY,
						"health_check_method":       acctest.REMOVEKEY,
						"health_check_path":         acctest.REMOVEKEY,
						"health_check_protocol":     "TCP",
						"health_check_codes.#":      "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
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

func TestAccAliCloudAlbHealthCheckTemplate_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_health_check_template.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAlbHealthCheckTemplateMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewAlbService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlbHealthCheckTemplate")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbhealthchecktemplate%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAlbHealthCheckTemplateBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_template_name": name,
					"health_check_connect_port":  "8080",
					"health_check_interval":      "20",
					"health_check_protocol":      "TCP",
					"health_check_timeout":       "60",
					"healthy_threshold":          "6",
					"unhealthy_threshold":        "6",
					"resource_group_id":          "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_template_name": name,
						"health_check_connect_port":  "8080",
						"health_check_interval":      "20",
						"health_check_protocol":      "TCP",
						"health_check_timeout":       "60",
						"healthy_threshold":          "6",
						"unhealthy_threshold":        "6",
						"resource_group_id":          acctest.CHECKSET,
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

var AliCloudAlbHealthCheckTemplateMap0 = map[string]string{
	"health_check_connect_port": acctest.CHECKSET,
	"health_check_host":         acctest.CHECKSET,
	"health_check_http_version": acctest.CHECKSET,
	"health_check_interval":     acctest.CHECKSET,
	"health_check_method":       acctest.CHECKSET,
	"health_check_path":         acctest.CHECKSET,
	"health_check_protocol":     acctest.CHECKSET,
	"health_check_timeout":      acctest.CHECKSET,
	"healthy_threshold":         acctest.CHECKSET,
	"unhealthy_threshold":       acctest.CHECKSET,
	"health_check_codes.#":      acctest.CHECKSET,
}

var AliCloudAlbHealthCheckTemplateMap1 = map[string]string{
	"health_check_connect_port": acctest.CHECKSET,
	"health_check_interval":     acctest.CHECKSET,
	"health_check_protocol":     acctest.CHECKSET,
	"health_check_timeout":      acctest.CHECKSET,
	"healthy_threshold":         acctest.CHECKSET,
	"unhealthy_threshold":       acctest.CHECKSET,
}

func AliCloudAlbHealthCheckTemplateBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

`, name)
}

// lintignore: R001
func TestUnitAliCloudAlbHealthCheckTemplate(t *testing.T) {
	p := map[string]*schema.Resource{
		"alicloud_alb_health_check_template": alb.ResourceHealthCheckTemplate(),
	}
	resourceSchema := p["alicloud_alb_health_check_template"]
	d, _ := schema.InternalMap(p["alicloud_alb_health_check_template"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_alb_health_check_template"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"health_check_template_name": "health_check_template_name",
		"health_check_protocol":      "HTTP",
		"dry_run":                    false,
		"health_check_codes":         []string{"http_3xx", "http_4xx"},
		"health_check_connect_port":  8080,
		"health_check_host":          "www.test.com",
		"health_check_http_version":  "HTTP1.0",
		"health_check_interval":      2,
		"health_check_method":        "GET",
		"health_check_path":          "/test",
		"health_check_timeout":       50,
		"healthy_threshold":          5,
		"unhealthy_threshold":        5,
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"HealthCheckTemplateId":   "MockHealthCheckTemplateId",
		"HealthCheckCodes":        "health_check_codes",
		"HealthCheckConnectPort":  8080,
		"HealthCheckInterval":     2,
		"HealthCheckProtocol":     "HTTP",
		"HealthCheckMethod":       "GET",
		"HealthCheckPath":         "/test",
		"HealthCheckHost":         "www.test.com",
		"HealthCheckHttpVersion":  "HTTP1.0",
		"HealthCheckTemplateName": "health_check_template_name",
		"HealthCheckTimeout":      50,
		"HealthyThreshold":        5,
		"UnhealthyThreshold":      5,
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
			return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_alb_health_check_template", "MockHealthCheckTemplateId"))
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
			result["HealthCheckTemplateId"] = "MockHealthCheckTemplateId"
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceSchema.Create(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("IdempotenceProcessing")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceSchema.Create(d, rawClient)
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
		err := resourceSchema.Create(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("MockHealthCheckTemplateId")
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceSchema.Update(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateHealthCheckTemplateAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"health_check_codes", "health_check_connect_port", "health_check_host", "health_check_http_version", "health_check_interval", "health_check_method", "health_check_path", "health_check_protocol", "health_check_template_name", "health_check_timeout", "healthy_threshold", "unhealthy_threshold", "dry_run"} {
			switch p["alicloud_alb_health_check_template"].Schema[key].Type {
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
			case schema.TypeList:
				diff.Attributes["health_check_codes.#"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["health_check_codes.0"] = &terraform.ResourceAttrDiff{Old: "", New: "http_3xx"}
				diff.Attributes["health_check_codes.1"] = &terraform.ResourceAttrDiff{Old: "", New: "http_4xx"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_alb_health_check_template"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("IdempotenceProcessing")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceSchema.Update(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateHealthCheckTemplateAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"health_check_codes", "health_check_connect_port", "health_check_host", "health_check_http_version", "health_check_interval", "health_check_method", "health_check_path", "health_check_protocol", "health_check_template_name", "health_check_timeout", "healthy_threshold", "unhealthy_threshold", "dry_run"} {
			switch p["alicloud_alb_health_check_template"].Schema[key].Type {
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
		resourceData1, _ := schema.InternalMap(p["alicloud_alb_health_check_template"].Schema).Data(nil, diff)
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
		err := resourceSchema.Update(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceSchema.Delete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("IdempotenceProcessing")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceSchema.Delete(d, rawClient)
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
		err := resourceSchema.Delete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
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
		err := resourceSchema.Delete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeAlbHealthCheckTemplateNotFound", func(t *testing.T) {
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
		err := resourceSchema.Read(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeAlbHealthCheckTemplateAbnormal", func(t *testing.T) {
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
		err := resourceSchema.Read(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

// Test Alb HealthCheckTemplate. >>> Resource test cases, automatically generated.
// Case HCT_test241220_1 9647
func TestAccAliCloudAlbHealthCheckTemplate_basic9647(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_health_check_template.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAlbHealthCheckTemplateMap9647)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewAlbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlbHealthCheckTemplate")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAlbHealthCheckTemplateBasicDependence9647)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval":      "2",
					"unhealthy_threshold":        "3",
					"health_check_template_name": name,
					"health_check_host":          "$SERVER_IP",
					"health_check_path":          "/testtf",
					"health_check_http_version":  "HTTP1.1",
					"health_check_timeout":       "5",
					"health_check_connect_port":  "0",
					"health_check_codes": []string{
						"http_2xx"},
					"health_check_method":   "HEAD",
					"healthy_threshold":     "3",
					"health_check_protocol": "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval":      "2",
						"unhealthy_threshold":        "3",
						"health_check_template_name": name,
						"health_check_host":          "$SERVER_IP",
						"health_check_path":          "/testtf",
						"health_check_http_version":  "HTTP1.1",
						"health_check_timeout":       "5",
						"health_check_connect_port":  "0",
						"health_check_codes.#":       "1",
						"health_check_method":        "HEAD",
						"healthy_threshold":          "3",
						"health_check_protocol":      "HTTP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval":      "5",
					"unhealthy_threshold":        "5",
					"health_check_template_name": name + "_update",
					"health_check_host":          "www.test.com",
					"health_check_path":          "/path",
					"health_check_http_version":  "HTTP1.0",
					"health_check_timeout":       "10",
					"health_check_connect_port":  "80",
					"health_check_codes": []string{
						"http_3xx"},
					"health_check_method": "GET",
					"healthy_threshold":   "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval":      "5",
						"unhealthy_threshold":        "5",
						"health_check_template_name": name + "_update",
						"health_check_host":          "www.test.com",
						"health_check_path":          "/path",
						"health_check_http_version":  "HTTP1.0",
						"health_check_timeout":       "10",
						"health_check_connect_port":  "80",
						"health_check_codes.#":       "1",
						"health_check_method":        "GET",
						"healthy_threshold":          "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval":     "2",
					"unhealthy_threshold":       "3",
					"health_check_host":         "$SERVER_IP",
					"health_check_http_version": "HTTP1.1",
					"health_check_timeout":      "5",
					"health_check_connect_port": "0",
					"health_check_method":       "HEAD",
					"healthy_threshold":         "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval":     "2",
						"unhealthy_threshold":       "3",
						"health_check_host":         "$SERVER_IP",
						"health_check_http_version": "HTTP1.1",
						"health_check_timeout":      "5",
						"health_check_connect_port": "0",
						"health_check_method":       "HEAD",
						"healthy_threshold":         "3",
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudAlbHealthCheckTemplateMap9647 = map[string]string{}

func AlicloudAlbHealthCheckTemplateBasicDependence9647(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case HCT_test241220_3 9691
func TestAccAliCloudAlbHealthCheckTemplate_basic9691(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_health_check_template.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAlbHealthCheckTemplateMap9691)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewAlbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlbHealthCheckTemplate")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAlbHealthCheckTemplateBasicDependence9691)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval":      "2",
					"unhealthy_threshold":        "3",
					"health_check_template_name": name,
					"health_check_host":          "$SERVER_IP",
					"health_check_path":          "/testtf",
					"health_check_http_version":  "HTTP1.1",
					"health_check_timeout":       "5",
					"health_check_connect_port":  "0",
					"health_check_codes": []string{
						"http_2xx", "http_3xx", "http_4xx", "http_5xx"},
					"health_check_method":   "HEAD",
					"healthy_threshold":     "3",
					"health_check_protocol": "HTTP",
					"dry_run":               "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval":      "2",
						"unhealthy_threshold":        "3",
						"health_check_template_name": name,
						"health_check_host":          "$SERVER_IP",
						"health_check_path":          "/testtf",
						"health_check_http_version":  "HTTP1.1",
						"health_check_timeout":       "5",
						"health_check_connect_port":  "0",
						"health_check_codes.#":       "4",
						"health_check_method":        "HEAD",
						"healthy_threshold":          "3",
						"health_check_protocol":      "HTTP",
						"dry_run":                    "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval":      "5",
					"unhealthy_threshold":        "5",
					"health_check_template_name": name + "_update",
					"health_check_host":          "www.test.com",
					"health_check_path":          "/path",
					"health_check_http_version":  "HTTP1.0",
					"health_check_timeout":       "10",
					"health_check_connect_port":  "80",
					"health_check_codes": []string{
						"http_2xx"},
					"health_check_method":   "GET",
					"healthy_threshold":     "5",
					"health_check_protocol": "HTTPS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval":      "5",
						"unhealthy_threshold":        "5",
						"health_check_template_name": name + "_update",
						"health_check_host":          "www.test.com",
						"health_check_path":          "/path",
						"health_check_http_version":  "HTTP1.0",
						"health_check_timeout":       "10",
						"health_check_connect_port":  "80",
						"health_check_codes.#":       "1",
						"health_check_method":        "GET",
						"healthy_threshold":          "5",
						"health_check_protocol":      "HTTPS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval":     "2",
					"unhealthy_threshold":       "3",
					"health_check_host":         "$SERVER_IP",
					"health_check_http_version": "HTTP1.1",
					"health_check_timeout":      "5",
					"health_check_connect_port": "0",
					"health_check_method":       "HEAD",
					"healthy_threshold":         "3",
					"health_check_protocol":     "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval":     "2",
						"unhealthy_threshold":       "3",
						"health_check_host":         "$SERVER_IP",
						"health_check_http_version": "HTTP1.1",
						"health_check_timeout":      "5",
						"health_check_connect_port": "0",
						"health_check_method":       "HEAD",
						"healthy_threshold":         "3",
						"health_check_protocol":     "HTTP",
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudAlbHealthCheckTemplateMap9691 = map[string]string{}

func AlicloudAlbHealthCheckTemplateBasicDependence9691(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Alb HealthCheckTemplate. <<< Resource test cases, automatically generated.
