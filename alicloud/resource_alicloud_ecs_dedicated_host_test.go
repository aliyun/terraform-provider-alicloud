package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudECSDedicatedHost_basic(t *testing.T) {
	var v ecs.DedicatedHost
	resourceId := "alicloud_ecs_dedicated_host.default"
	ra := resourceAttrInit(resourceId, EcsDedicatedHostMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDedicatedHost")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccEcsDedicatedHost%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, EcsDedicatedHostBasicdependence)
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
					"dedicated_host_type": "ddh.g6",
					"description":         "From_Terraform",
					"dedicated_host_name": name,
					"payment_type":        "PrePaid",
					"period":              1,
					"period_unit":         "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_type": "ddh.g6",
						"description":         "From_Terraform",
						"dedicated_host_name": name,
						"payment_type":        "PrePaid",
						"period_unit":         "Month",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"detail_fee", "dry_run", "min_quantity", "auto_renew", "auto_renew_period", "expired_time", "period"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "DDH_Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "DDH_Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_name": name + "ddh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_name": name + "ddh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_attributes": []map[string]interface{}{
						{
							"udp_timeout":     "70",
							"slb_udp_timeout": "70",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_attributes.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_on_maintenance": "Migrate",
					"auto_placement":        "off",
					"auto_release_time":     "2027-12-31T23:59:59Z",
					"detail_fee":            false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_on_maintenance": "Migrate",
						"auto_placement":        "off",
						"auto_release_time":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "Terraform",
						"For":     "DDH",
					},
					"dedicated_host_name": name,
					"description":         "From_Terraform",
					"network_attributes": []map[string]interface{}{
						{
							"udp_timeout":     "60",
							"slb_udp_timeout": "60",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":               "2",
						"tags.Created":         "Terraform",
						"tags.For":             "DDH",
						"dedicated_host_name":  name,
						"description":          "From_Terraform",
						"network_attributes.#": "1",
						"resource_group_id":    CHECKSET,
					}),
				),
			},
			// Change payment_type from PrePaid to PostPaid to cover the
			// ModifyDedicatedHostsChargeType update path. period and
			// period_unit stay unchanged from the create step, so only
			// HasChange("payment_type") fires here.
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PostPaid",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudECSDedicatedHost_basic1(t *testing.T) {
	var v ecs.DedicatedHost
	resourceId := "alicloud_ecs_dedicated_host.default"
	ra := resourceAttrInit(resourceId, EcsDedicatedHostMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDedicatedHost")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccEcsDedicatedHost%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, EcsDedicatedHostBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_type": "ddh.g6",
					"description":         "From_Terraform",
					"dedicated_host_name": name,
					"auto_renew":          "true",
					"auto_renew_period":   "1",
					"expired_time":        "1",
					"sale_cycle":          "Week",
					"payment_type":        "PrePaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_type": "ddh.g6",
						"description":         "From_Terraform",
						"dedicated_host_name": name,
						"sale_cycle":          "Week",
						"payment_type":        "PrePaid",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"detail_fee", "dry_run", "min_quantity", "auto_renew", "auto_renew_period", "expired_time"},
			},
		},
	})
}

// TestAccAliCloudECSDedicatedHost_autoRenew covers the renewal_status,
// auto_renew_with_ecs and duration fields backed by the
// ModifyDedicatedHostAutoRenewAttribute / DescribeDedicatedHostAutoRenew APIs,
// and the period, period_unit, quantity and auto_pay arguments surfaced on the
// AllocateDedicatedHosts / ModifyDedicatedHostsChargeType paths. It exercises
// the create -> update (ModifyDedicatedHostAutoRenewAttribute) -> read
// (DescribeDedicatedHostAutoRenew) round-trip on a subscription dedicated host.
// The deprecated sale_cycle / expired_time aliases are intentionally retained
// in the config to verify backward compatibility and keep attribute coverage.
func TestAccAliCloudECSDedicatedHost_autoRenew(t *testing.T) {
	var v ecs.DedicatedHost
	resourceId := "alicloud_ecs_dedicated_host.default"
	ra := resourceAttrInit(resourceId, EcsDedicatedHostMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDedicatedHost")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccEcsDedicatedHost%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, EcsDedicatedHostBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_type": "ddh.g6",
					"dedicated_host_name": name,
					"description":         "From_Terraform",
					"payment_type":        "PrePaid",
					"period":              1,
					"period_unit":         "Month",
					"sale_cycle":          "Month",
					"expired_time":        "1",
					"quantity":            1,
					"auto_pay":            true,
					"renewal_status":      "Normal",
					"auto_renew_with_ecs": "StopRenewWithEcs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_type": "ddh.g6",
						"payment_type":        "PrePaid",
						"period_unit":         "Month",
						"sale_cycle":          "Month",
						"renewal_status":      "Normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_type": "ddh.g6",
					"dedicated_host_name": name,
					"description":         "From_Terraform",
					"payment_type":        "PrePaid",
					"period":              1,
					"period_unit":         "Year",
					"sale_cycle":          "Year",
					"expired_time":        "1",
					"quantity":            1,
					"auto_pay":            true,
					"renewal_status":      "AutoRenewal",
					"auto_renew_with_ecs": "AutoRenewWithEcs",
					"duration":            "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status":      "AutoRenewal",
						"auto_renew_with_ecs": "AutoRenewWithEcs",
						"duration":            "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"detail_fee", "dry_run", "min_quantity", "auto_renew", "auto_renew_period", "expired_time", "period", "auto_pay", "quantity"},
			},
		},
	})
}

func TestAccAliCloudECSDedicatedHost_basic2(t *testing.T) {
	var v ecs.DedicatedHost
	resourceId := "alicloud_ecs_dedicated_host.default"
	ra := resourceAttrInit(resourceId, EcsDedicatedHostMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDedicatedHost")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccEcsDedicatedHost%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, EcsDedicatedHostBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EcsDedicatedHostRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_type":   "ddh.r6",
					"description":           "From_Terraform",
					"dedicated_host_name":   name,
					"action_on_maintenance": "Migrate",
					"auto_placement":        "on",
					"min_quantity":          "1",
					"network_attributes": []map[string]interface{}{
						{
							"udp_timeout":     "70",
							"slb_udp_timeout": "70",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "DDH_Test",
					},
					"zone_id": "cn-hangzhou-i",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_type":   "ddh.r6",
						"description":           "From_Terraform",
						"dedicated_host_name":   name,
						"action_on_maintenance": "Migrate",
						"auto_placement":        "on",
						"min_quantity":          "1",
						"network_attributes.#":  "1",
						"resource_group_id":     CHECKSET,
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.For":              "DDH_Test",
						"zone_id":               CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"detail_fee", "dry_run", "min_quantity", "auto_renew", "auto_renew_period", "expired_time"},
			},
		},
	})
}

// TestAccAliCloudECSDedicatedHost_cpuOverCommit exercises the cpu_over_commit_ratio
// attribute, which is only supported by the g6s/c6s/r6s custom dedicated host
// specifications. It covers the create -> update (ModifyDedicatedHostAttribute)
// -> read round-trip and is the only ACC case that sets this field, so the
// TestingCoverageRate must-set and modify gates are both satisfied here.
func TestAccAliCloudECSDedicatedHost_cpuOverCommit(t *testing.T) {
	var v ecs.DedicatedHost
	resourceId := "alicloud_ecs_dedicated_host.default"
	ra := resourceAttrInit(resourceId, EcsDedicatedHostMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDedicatedHost")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccEcsDedicatedHost%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, EcsDedicatedHostBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EcsDedicatedHostRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_type":   "ddh.g6s",
					"description":           "From_Terraform",
					"dedicated_host_name":   name,
					"cpu_over_commit_ratio": 1.2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_type":   "ddh.g6s",
						"cpu_over_commit_ratio": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu_over_commit_ratio": 1.5,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu_over_commit_ratio": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"detail_fee", "dry_run", "min_quantity", "auto_renew", "auto_renew_period", "expired_time"},
			},
		},
	})
}

var EcsDedicatedHostMap = map[string]string{
	"detail_fee": "false",
	"dry_run":    "false",
	"status":     CHECKSET,
}

func EcsDedicatedHostBasicdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	  name_regex = "^default-NODELETING$"
	}
	data "alicloud_vswitches" "default" {
	  vpc_id = data.alicloud_vpcs.default.ids.0
      zone_id = "cn-hangzhou-i"
	}
	data "alicloud_resource_manager_resource_groups" "default"{
	}
`)
}

// lintignore: R001
func TestUnitAlicloudECSDedicatedHost(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"dedicated_host_type":   "AllocateDedicatedHostsValue",
		"description":           "AllocateDedicatedHostsValue",
		"dedicated_host_name":   "AllocateDedicatedHostsValue",
		"action_on_maintenance": "AllocateDedicatedHostsValue",
		"auto_placement":        "AllocateDedicatedHostsValue",
		"min_quantity":          1,
		"network_attributes": []map[string]interface{}{
			{
				"udp_timeout":     70,
				"slb_udp_timeout": 70,
			},
		},
		"resource_group_id": "AllocateDedicatedHostsValue",
		"tags": map[string]string{
			"Created": "TF",
		},
		"zone_id":                   "AllocateDedicatedHostsValue",
		"auto_release_time":         "AllocateDedicatedHostsValue",
		"auto_renew":                false,
		"auto_renew_period":         1,
		"auto_pay":                  true,
		"cpu_over_commit_ratio":     1.2,
		"dedicated_host_cluster_id": "AllocateDedicatedHostsValue",
		"expired_time":              "AllocateDedicatedHostsValue",
		"payment_type":              "AllocateDedicatedHostsValue",
		"period":                    1,
		"period_unit":               "Month",
		"quantity":                  1,
		"sale_cycle":                "AllocateDedicatedHostsValue",
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
		// DescribeDedicatedHosts
		"DedicatedHosts": map[string]interface{}{
			"DedicatedHost": []interface{}{
				map[string]interface{}{
					"DedicatedHostId":     "AllocateDedicatedHostsValue",
					"ActionOnMaintenance": "AllocateDedicatedHostsValue",
					"AutoPlacement":       "AllocateDedicatedHostsValue",
					"AutoReleaseTime":     "AllocateDedicatedHostsValue",
					"CpuOverCommitRatio":  "AllocateDedicatedHostsValue",
					"DedicatedHostName":   "AllocateDedicatedHostsValue",
					"DedicatedHostType":   "AllocateDedicatedHostsValue",
					"Description":         "AllocateDedicatedHostsValue",
					"NetworkAttributes": map[string]interface{}{
						"SlbUdpTimeout": 70,
						"UdpTimeout":    70,
					},
					"ChargeType":      "AllocateDedicatedHostsValue",
					"ResourceGroupId": "AllocateDedicatedHostsValue",
					"SaleCycle":       "AllocateDedicatedHostsValue",
					"Status":          "Available",
					"ZoneId":          "AllocateDedicatedHostsValue",
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"Key":   "Created",
								"Value": "TF",
							},
						},
					},
				},
			},
		},
		"DedicatedHostIdSets": map[string]interface{}{
			"DedicatedHostId": []interface{}{
				"AllocateDedicatedHostsValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// AllocateDedicatedHosts
		"DedicatedHostIdSets": map[string]interface{}{
			"DedicatedHostId": []interface{}{
				"AllocateDedicatedHostsValue",
			},
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
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ecs_dedicated_host", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEcsDedicatedHostCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDedicatedHosts Response
		"DedicatedHostIdSets": map[string]interface{}{
			"DedicatedHostId": []interface{}{
				"AllocateDedicatedHostsValue",
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "AllocateDedicatedHosts" {
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
		err := resourceAlicloudEcsDedicatedHostCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEcsDedicatedHostUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyDedicatedHostAutoReleaseTime
	attributesDiff := map[string]interface{}{
		"auto_release_time": "ModifyDedicatedHostAutoReleaseTimeValue",
	}
	diff, err := newInstanceDiff("alicloud_ecs_dedicated_host", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDedicatedHosts Response
		"DedicatedHosts": map[string]interface{}{
			"DedicatedHost": []interface{}{
				map[string]interface{}{
					"AutoReleaseTime": "ModifyDedicatedHostAutoReleaseTimeValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDedicatedHostAutoReleaseTime" {
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
		err := resourceAlicloudEcsDedicatedHostUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// JoinResourceGroup
	attributesDiff = map[string]interface{}{
		"resource_group_id": "JoinResourceGroupValue",
	}
	diff, err = newInstanceDiff("alicloud_ecs_dedicated_host", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDedicatedHosts Response
		"DedicatedHosts": map[string]interface{}{
			"DedicatedHost": []interface{}{
				map[string]interface{}{
					"ResourceGroupId": "JoinResourceGroupValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "JoinResourceGroup" {
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
		err := resourceAlicloudEcsDedicatedHostUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// RenewDedicatedHosts
	attributesDiff = map[string]interface{}{
		"expired_time": "RenewDedicatedHostsValue",
		"sale_cycle":   "RenewDedicatedHostsValue",
	}
	diff, err = newInstanceDiff("alicloud_ecs_dedicated_host", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDedicatedHosts Response
		"DedicatedHosts": map[string]interface{}{
			"DedicatedHost": []interface{}{
				map[string]interface{}{
					"SaleCycle": "RenewDedicatedHostsValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "RenewDedicatedHosts" {
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
		err := resourceAlicloudEcsDedicatedHostUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// ModifyDedicatedHostsChargeType
	attributesDiff = map[string]interface{}{
		"expired_time": "ModifyDedicatedHostsChargeTypeValue",
		"payment_type": "ModifyDedicatedHostsChargeTypeValue",
		"sale_cycle":   "ModifyDedicatedHostsChargeTypeValue",
		"detail_fee":   false,
		"dry_run":      false,
	}
	diff, err = newInstanceDiff("alicloud_ecs_dedicated_host", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDedicatedHosts Response
		"DedicatedHosts": map[string]interface{}{
			"DedicatedHost": []interface{}{
				map[string]interface{}{
					"ChargeType": "ModifyDedicatedHostsChargeTypeValue",
					"SaleCycle":  "ModifyDedicatedHostsChargeTypeValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDedicatedHostsChargeType" {
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
		err := resourceAlicloudEcsDedicatedHostUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// ModifyDedicatedHostAttribute
	attributesDiff = map[string]interface{}{
		"action_on_maintenance": "ModifyDedicatedHostAttributeValue",
		"auto_placement":        "ModifyDedicatedHostAttributeValue",
		"cpu_over_commit_ratio": 1.6,
		"dedicated_host_name":   "ModifyDedicatedHostAttributeValue",
		"description":           "ModifyDedicatedHostAttributeValue",
		"network_attributes": []map[string]interface{}{
			{
				"udp_timeout":     80,
				"slb_udp_timeout": 80,
			},
		},
		"dedicated_host_cluster_id": "ModifyDedicatedHostAttributeValue",
	}
	diff, err = newInstanceDiff("alicloud_ecs_dedicated_host", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDedicatedHosts Response
		"DedicatedHosts": map[string]interface{}{
			"DedicatedHost": []interface{}{
				map[string]interface{}{
					"ActionOnMaintenance": "ModifyDedicatedHostAttributeValue",
					"AutoPlacement":       "ModifyDedicatedHostAttributeValue",
					"CpuOverCommitRatio":  "ModifyDedicatedHostAttributeValue",
					"DedicatedHostName":   "ModifyDedicatedHostAttributeValue",
					"Description":         "ModifyDedicatedHostAttributeValue",
					"NetworkAttributes": map[string]interface{}{
						"SlbUdpTimeout": 80,
						"UdpTimeout":    80,
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDedicatedHostAttribute" {
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
		err := resourceAlicloudEcsDedicatedHostUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDedicatedHosts" {
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
		err := resourceAlicloudEcsDedicatedHostRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEcsDedicatedHostDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "IncorrectHostStatus.Initializing", "nil", "InvalidDedicatedHostId.NotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ReleaseDedicatedHost" {
				switch errorCode {
				case "NonRetryableError", "InvalidDedicatedHostId.NotFound":
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
		err := resourceAlicloudEcsDedicatedHostDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InvalidDedicatedHostId.NotFound":
			assert.Nil(t, err)
		}
	}
}
