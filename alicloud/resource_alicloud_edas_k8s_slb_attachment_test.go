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
	roa "github.com/alibabacloud-go/tea-roa/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

const ClusterTypeK8s = 5

func init() {
	resource.AddTestSweepers(
		"alicloud_edas_k8s_slb_attachment",
		&resource.Sweeper{
			Name: "alicloud_edas_k8s_slb_attachment",
			F:    testSweepEDASK8sSlbAttachment,
		})
}

func testSweepEDASK8sSlbAttachment(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tftestAcc",
	}

	appListReq := edas.CreateListApplicationRequest()
	appListReq.RegionId = region
	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListApplication(appListReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to list edas app: %s", err)
		return nil
	}
	listAppResponse, _ := raw.(edas.ListApplicationResponse)
	if listAppResponse.Code != 200 {
		log.Printf("[ERROR] Failed to list edas app: %v", listAppResponse)
		return nil
	}

	var appIdList []string
	for _, v := range listAppResponse.ApplicationList.Application {
		if ClusterTypeK8s != v.ClusterType {
			continue
		}
		skip := true
		appId := v.AppId
		app, err := edasService.DescribeEdasK8sApplication(appId)
		if err != nil {
			log.Printf("[ERROR] Failed to get edas k8s app, id: %s, err: %v", appId, err)
			continue
		}
		for _, pre := range prefixes {
			if strings.HasPrefix(app.Name, pre) {
				skip = false
				break
			}
		}

		if skip {
			log.Printf("[INFO] Skipping edas k8s app: %s", app.Name)
		} else {
			appIdList = append(appIdList, appId)
		}
	}

	for _, appId := range appIdList {
		log.Printf("[INFO] Deleting edas k8s app: %s", appId)
		request := edas.CreateDeleteK8sApplicationRequest()
		request.RegionId = client.RegionId
		request.AppId = appId

		wait := incrementalWait(1*time.Second, 2*time.Second)
		err := resource.Retry(time.Minute*2, func() *resource.RetryError {
			raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
				return edasClient.DeleteK8sApplication(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RoaRequest, request)
			response := raw.(*edas.DeleteK8sApplicationResponse)
			if response.Code != 200 {
				return resource.NonRetryableError(Error("[ERROR] Delete k8s application failed for %s", response.Message))
			}
			return edasService.WaitForChangeOrderFinishedNonRetryable(appId, response.ChangeOrderId, 3*time.Minute)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, appId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func TestUnitAccAlicloudEDASK8sSlbAttachment(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_edas_k8s_slb_attachment"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_edas_k8s_slb_attachment"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"app_id": "mock_app_id",
		"slb_configs": []map[string]interface{}{
			{
				"type":          "internet",
				"name":          "create_slb_attachment",
				"scheduler":     "rr",
				"specification": "slb.s1.small",
				"port_mappings": []map[string]interface{}{
					{
						"loadbalancer_protocol": "TCP",
						"service_port": []map[string]interface{}{
							{
								"port":        80,
								"protocol":    "TCP",
								"target_port": 8080,
							},
						},
					},
				},
			},
		},
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
		"body": map[string]interface{}{
			"Application": map[string]interface{}{
				"AppId":   "mock_app_id",
				"SlbInfo": "[{\"addressType\":\"internet\",\"externalTrafficPolicy\":\"Local\",\"ip\":\"mock_ip\",\"name\":\"create_slb_attachment\",\"portMappings\":[{\"loadBalancerProtocol\":\"TCP\",\"servicePort\":{\"port\":80,\"protocol\":\"TCP\",\"targetPort\":8080}}],\"scheduler\":\"rr\",\"serviceType\":\"LoadBalancer\",\"specification\":\"slb.s1.small\"}]",
			},
			"Code": 200,
		},
	}

	CreateMockResponse := map[string]interface{}{
		"body": map[string]interface{}{
			"ChangeOrderId": "mockChangeOrderId",
			"Code":          200,
		},
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	hasChangeOrderMockResponse := map[string]interface{}{
		"body": map[string]interface{}{
			"ChangeOrderList": map[string]interface{}{
				"ChangeOrder": []interface{}{
					map[string]interface{}{
						"Status": 2,
					},
				},
			},
			"Code": 200,
		},
	}
	changeOrderMockResponse := map[string]interface{}{
		"body": map[string]interface{}{
			"changeOrderInfo": map[string]interface{}{
				"Status": 2,
			},
			"Code": 200,
		},
	}

	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_edas_k8s_slb_attachment", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEdasClient", func(_ *connectivity.AliyunClient) (*roa.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEdasK8sSlbAttachmentCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, action *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "/pop/v5/k8s/acs/k8s_slb_binding" {
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
			} else if *action == "/pop/v5/changeorder/change_order_list" {
				return hasChangeOrderMockResponse, nil
			} else if *action == "/pop/v5/changeorder/change_order_info" {
				return changeOrderMockResponse, nil
			} else if *action == "/pop/v5/app/app_info" {
				return ReadMockResponse, nil
			}
			return nil, nil
		})
		err := resourceAlicloudEdasK8sSlbAttachmentCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_edas_k8s_slb_attachment"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEdasClient", func(_ *connectivity.AliyunClient) (*roa.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEdasK8sSlbAttachmentUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"slb_configs": []map[string]interface{}{
			{
				"type":          "internet",
				"name":          "create_slb_attachment",
				"scheduler":     "wrr",
				"specification": "slb.s1.small",
				"port_mappings": []map[string]interface{}{
					{
						"loadbalancer_protocol": "TCP",
						"service_port": []map[string]interface{}{
							{
								"port":        81,
								"protocol":    "TCP",
								"target_port": 8081,
							},
						},
					},
				},
			},
		},
	}
	diff, err := newInstanceDiff("alicloud_edas_k8s_slb_attachment", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_edas_k8s_slb_attachment"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"body": map[string]interface{}{
			"Application": map[string]interface{}{
				"AppId":   "mock_app_id",
				"SlbInfo": "[{\"addressType\":\"internet\",\"externalTrafficPolicy\":\"Local\",\"ip\":\"mock_ip\",\"name\":\"create_slb_attachment\",\"portMappings\":[{\"loadBalancerProtocol\":\"TCP\",\"servicePort\":{\"port\":81,\"protocol\":\"TCP\",\"targetPort\":8081}}],\"scheduler\":\"wrr\",\"serviceType\":\"LoadBalancer\",\"specification\":\"slb.s1.small\"}]",
			},
			"Code": 200,
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, action *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "/pop/v5/k8s/acs/k8s_slb_binding" {
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
			} else if *action == "/pop/v5/changeorder/change_order_list" {
				return hasChangeOrderMockResponse, nil
			} else if *action == "/pop/v5/changeorder/change_order_info" {
				return changeOrderMockResponse, nil
			} else if *action == "/pop/v5/app/app_info" {
				return ReadMockResponse, nil
			}
			return nil, nil
		})
		err := resourceAlicloudEdasK8sSlbAttachmentUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_edas_k8s_slb_attachment"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_edas_k8s_slb_attachment", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_edas_k8s_slb_attachment"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, action *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "/pop/v5/app/app_info" {
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
			return nil, nil
		})
		err := resourceAlicloudEdasK8sSlbAttachmentRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}
	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEdasClient", func(_ *connectivity.AliyunClient) (*roa.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEdasK8sSlbAttachmentDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_edas_k8s_slb_attachment", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_edas_k8s_slb_attachment"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "failCode"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, action *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "/pop/v5/k8s/acs/k8s_slb_binding" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				case "failCode":
					return map[string]interface{}{
						"body": map[string]interface{}{
							"Code": 400,
						},
					}, nil
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			} else if *action == "/pop/v5/changeorder/change_order_list" {
				return hasChangeOrderMockResponse, nil
			} else if *action == "/pop/v5/changeorder/change_order_info" {
				return changeOrderMockResponse, nil
			}
			return nil, nil
		})
		err := resourceAlicloudEdasK8sSlbAttachmentDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "failCode":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
		}
	}
}

func TestAccAlicloudEDASK8sSlbAttachment_basic(t *testing.T) {
	var v []map[string]interface{}

	resourceId := "alicloud_edas_k8s_slb_attachment.default"

	ra := resourceAttrInit(resourceId, edasK8sSlbAttachmentBasicMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edask8s-slb-attach-%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEDASK8sSlbAttachmentConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEdasSlbAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id": "${alicloud_edas_k8s_application.default.id}",
					"slb_configs": []map[string]interface{}{
						{
							"type":      "internet",
							"scheduler": "rr",
							"port_mappings": []map[string]interface{}{
								{
									"loadbalancer_protocol": "TCP",
									"service_port": []map[string]interface{}{
										{
											"port":        "80",
											"protocol":    "TCP",
											"target_port": "18081",
										},
									},
								},
								{
									"loadbalancer_protocol": "TCP",
									"service_port": []map[string]interface{}{
										{
											"port":        "8080",
											"protocol":    "TCP",
											"target_port": "18081",
										},
									},
								},
							},
						},
						{
							"type":      "intranet",
							"scheduler": "rr",
							"port_mappings": []map[string]interface{}{
								{
									"loadbalancer_protocol": "TCP",
									"service_port": []map[string]interface{}{
										{
											"port":        "80",
											"protocol":    "TCP",
											"target_port": "18081",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_id":        CHECKSET,
						"slb_configs.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"slb_configs": []map[string]interface{}{
						{
							"type":      "internet",
							"scheduler": "rr",
							"port_mappings": []map[string]interface{}{
								{
									"loadbalancer_protocol": "TCP",
									"service_port": []map[string]interface{}{
										{
											"port":        "81",
											"protocol":    "TCP",
											"target_port": "18081",
										},
									},
								},
								{
									"loadbalancer_protocol": "TCP",
									"service_port": []map[string]interface{}{
										{
											"port":        "8081",
											"protocol":    "TCP",
											"target_port": "18081",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"slb_configs.#": "1",
					}),
				),
			},
		},
	})
}

var edasK8sSlbAttachmentBasicMap = map[string]string{
	"app_id": CHECKSET,
}

func testAccCheckEdasSlbAttachmentDestroy(*terraform.State) error {
	return nil
}

func resourceEDASK8sSlbAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}

		data "alicloud_zones" default {
		  available_resource_creation = "VSwitch"
		}
		
		data "alicloud_instance_types" "default" {
		  availability_zone = data.alicloud_zones.default.zones.0.id
		  cpu_core_count = 4
		  memory_size = 8
		  kubernetes_node_role = "Worker"
		}
		
		resource "alicloud_vpc" "default" {
		  name = var.name
		  cidr_block = "10.1.0.0/21"
		}
		
		resource "alicloud_vswitch" "default" {
		  name = var.name
		  vpc_id = alicloud_vpc.default.id
		  cidr_block = "10.1.1.0/24"
		  availability_zone = data.alicloud_zones.default.zones.0.id
		}
		
		resource "alicloud_cs_managed_kubernetes" "default" {
		  name_prefix          = var.name
		  cluster_spec         = "ack.pro.small"
		  worker_vswitch_ids   = [alicloud_vswitch.default.id]
		  new_nat_gateway      = false
		  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
		  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
		  slb_internet_enabled = true
		}

		resource "alicloud_key_pair" "default" {
		  key_pair_name = var.name
		}

		resource "alicloud_cs_kubernetes_node_pool" "default" {
		  name                 = "desired_size"
		  cluster_id           = alicloud_cs_managed_kubernetes.default.id
		  vswitch_ids          = [alicloud_vswitch.default.id]
		  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
		  system_disk_category = "cloud_efficiency"
		  system_disk_size     = 40
		  key_name             = alicloud_key_pair.default.key_name
		  desired_size         = 2
		}
		
		resource "alicloud_edas_k8s_cluster" "default" {
		  cs_cluster_id = alicloud_cs_managed_kubernetes.default.id
		}

		resource "alicloud_edas_k8s_application" "default" {
          application_name = var.name
          cluster_id = alicloud_edas_k8s_cluster.default.id
          package_type = "FatJar"
          package_url = "http://edas-bj.oss-cn-beijing.aliyuncs.com/prod/demo/SPRING_CLOUD_PROVIDER.jar"
          jdk = "Open JDK 8"
          replicas = "1"
        }
		`, name)
}
