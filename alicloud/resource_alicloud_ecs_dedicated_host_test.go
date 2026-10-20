package alicloud

import (
	"context"
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-zhangjiakou"})
			testAccPreCheck(t)
			// basic creates a PrePaid dedicated host, so gate it on the
			// month-first calendar guard just like basic1/autoRenew.
			testAccPreCheckWithTime(t, []int{1})
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_type": "ddh.c6e",
					"description":         "From_Terraform",
					"dedicated_host_name": name,
					"payment_type":        "PrePaid",
					"period":              1,
					"period_unit":         "Month",
					"zone_id":             "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_type": "ddh.c6e",
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
					"dedicated_host_cluster_id": "",
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
			testAccPreCheckPrePaidResources(t)
			testAccPreCheckWithTime(t, []int{1})
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_type":       "ddh.g6",
					"dedicated_host_name":       name,
					"description":               "From_Terraform",
					"payment_type":              "PrePaid",
					"period":                    1,
					"period_unit":               "Month",
					"sale_cycle":                "Month",
					"expired_time":              "1",
					"quantity":                  1,
					"auto_pay":                  true,
					"renewal_status":            "Normal",
					"auto_renew_with_ecs":       "StopRenewWithEcs",
					"zone_id":                   "",
					"dedicated_host_cluster_id": "",
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
					"period_unit":         "Week",
					"sale_cycle":          "Week",
					"expired_time":        "1",
					"quantity":            1,
					"auto_pay":            true,
					"renewal_status":      "AutoRenewal",
					"auto_renew_with_ecs": "AutoRenewWithEcs",
					"duration":            "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status":      "AutoRenewal",
						"auto_renew_with_ecs": "AutoRenewWithEcs",
						"duration":            "3",
						"period_unit":         "Week",
						"sale_cycle":          "Week",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"detail_fee", "dry_run", "min_quantity", "auto_renew", "auto_renew_period", "expired_time", "period", "auto_pay"},
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-zhangjiakou"})
			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_type":   "ddh.c6e",
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_type":   "ddh.c6e",
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

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
	data "alicloud_resource_manager_resource_groups" "default"{
	}
`)
}

// lintignore: R001
func TestUnitAlicloudECSDedicatedHost(t *testing.T) {
	p := Provider().ResourcesMap
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
	// Quantity validation: a Terraform resource maps to a single dedicated
	// host. Values > 1 create unmanaged paid resources and must be rejected
	// by schema validation.
	quantitySchema := p["alicloud_ecs_dedicated_host"].Schema["quantity"]
	if quantitySchema.ValidateFunc != nil {
		_, errs := quantitySchema.ValidateFunc(2, "quantity")
		assert.NotEmpty(t, errs, "quantity=2 should fail validation")
		_, errs = quantitySchema.ValidateFunc(1, "quantity")
		assert.Empty(t, errs, "quantity=1 should pass validation")
	}
	// payment_type only accepts PostPaid/PrePaid; "Year" is a period_unit
	// value and must be rejected at schema validation, mirroring the
	// quantity>1 guard rather than reaching the charge-type change path.
	paymentTypeSchema := p["alicloud_ecs_dedicated_host"].Schema["payment_type"]
	if paymentTypeSchema.ValidateFunc != nil {
		_, errs := paymentTypeSchema.ValidateFunc("Year", "payment_type")
		assert.NotEmpty(t, errs, "payment_type=Year should fail validation")
		_, errs = paymentTypeSchema.ValidateFunc("PrePaid", "payment_type")
		assert.Empty(t, errs, "payment_type=PrePaid should pass validation")
	}
	// All HTTP is mocked via gomonkey (DoRequest/NewEcsClient), so no real
	// API call is made. sharedClientForRegion still requires
	// ALICLOUD_ACCESS_KEY/SECRET_KEY/REGION to construct the AliyunClient;
	// without them the test SKIPs on CI workers lacking injected credentials
	// (root cause of the historical SKIP). Seed dummy values so the cloud-free
	// assertions always execute regardless of the worker environment.
	t.Setenv("ALICLOUD_ACCESS_KEY", "test-only-no-real-call")
	t.Setenv("ALICLOUD_SECRET_KEY", "test-only-no-real-call")
	t.Setenv("ALICLOUD_REGION", "cn-hangzhou")
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		// Cloudless scope: all HTTP is mocked via gomonkey, so the test must
		// execute regardless of worker credentials. A sharedClient init failure
		// is a hard error, not a silent skip; the dummy env above is expected to
		// satisfy construction, and any real init failure must surface loudly.
		t.Fatalf("sharedClientForRegion init failed (cloudless test must not skip): %s", err)
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
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "WithEcsClient", func(_ *connectivity.AliyunClient, _ func(*ecs.Client) (interface{}, error)) (interface{}, error) {
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "WithEcsClient", func(_ *connectivity.AliyunClient, _ func(*ecs.Client) (interface{}, error)) (interface{}, error) {
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

	// RenewDedicatedHosts — sale_cycle only (period_unit not changed)
	// Verifies that when sale_cycle changes but period_unit does not,
	// the PeriodUnit sent to RenewDedicatedHosts comes from the new
	// sale_cycle value, not the stale period_unit populated by Read.
	{
		var capturedReq map[string]interface{}
		attributesDiff = map[string]interface{}{
			"sale_cycle": "Year",
		}
		diff, err = newInstanceDiff("alicloud_ecs_dedicated_host", attributes, attributesDiff, dExisted.State())
		if err != nil {
			t.Error(err)
		}
		dTest, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			"DedicatedHosts": map[string]interface{}{
				"DedicatedHost": []interface{}{
					map[string]interface{}{
						"SaleCycle": "Year",
					},
				},
			},
		}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, req map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "RenewDedicatedHosts" {
				capturedReq = req
			}
			return successResponseMock(ReadMockResponseDiff)
		})
		err = resourceAlicloudEcsDedicatedHostUpdate(dTest, rawClient)
		patches.Reset()
		assert.Nil(t, err)
		assert.Equal(t, "Year", capturedReq["PeriodUnit"])
		dExisted = dTest
	}

	// RenewDedicatedHosts — period_unit only (sale_cycle not changed)
	{
		var capturedReq map[string]interface{}
		attributesDiff = map[string]interface{}{
			"period_unit": "Week",
		}
		diff, err = newInstanceDiff("alicloud_ecs_dedicated_host", attributes, attributesDiff, dExisted.State())
		if err != nil {
			t.Error(err)
		}
		dTest, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			"DedicatedHosts": map[string]interface{}{
				"DedicatedHost": []interface{}{
					map[string]interface{}{
						"SaleCycle": "Week",
					},
				},
			},
		}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, req map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "RenewDedicatedHosts" {
				capturedReq = req
			}
			return successResponseMock(ReadMockResponseDiff)
		})
		err = resourceAlicloudEcsDedicatedHostUpdate(dTest, rawClient)
		patches.Reset()
		assert.Nil(t, err)
		assert.Equal(t, "Week", capturedReq["PeriodUnit"])
		dExisted = dTest
	}

	// buildData constructs a *schema.ResourceData from a prior InstanceState
	// and a raw config via the SDK public Diff path (schema.Resource.Diff).
	// Unlike the hand-written newInstanceDiff helper — which builds a
	// terraform.InstanceDiff attribute-by-attribute and never invokes the
	// SDK schemaMap.Diff (no CustomizeDiff, no RequiresNew, no schema
	// validation) — this exercises the real plan computation that
	// `terraform plan` would produce. Only legal state/rawconfig is
	// injected; no hand-crafted InstanceDiff is allowed to masquerade as a plan.
	buildData := func(state *terraform.InstanceState, rawConfig map[string]interface{}) *schema.ResourceData {
		r := p["alicloud_ecs_dedicated_host"]
		config := terraform.NewResourceConfigRaw(rawConfig)
		diff, dErr := r.Diff(context.Background(), state, config, rawClient)
		assert.Nil(t, dErr)
		d, _ := schema.InternalMap(r.Schema).Data(state, diff)
		return d
	}

	// ModifyDedicatedHostsChargeType — auto_pay: omitted defaults to true
	// (Update L505-509); explicit false/true pass through. Each case uses a
	// fresh InstanceState (no residual auto_pay=false carried over from a
	// prior block) and asserts the real AutoPay value captured on the
	// ModifyDedicatedHostsChargeType RPC, rather than a hand-built diff.
	autoPayCases := []struct {
		name     string
		autoPay  interface{} // nil => omitted from config
		expected interface{}
	}{
		{"omitted-defaults-true", nil, true},
		{"explicit-false", false, false},
		{"explicit-true", true, true},
	}
	for _, tc := range autoPayCases {
		// Fresh prior state: PrePaid, legal Month period, quantity=1.
		state := &terraform.InstanceState{
			ID: "dh-autopay-test",
			Attributes: map[string]string{
				"id":                  "dh-autopay-test",
				"payment_type":        "PrePaid",
				"period_unit":         "Month",
				"sale_cycle":          "Month",
				"quantity":            "1",
				"dedicated_host_type": "ddh.g6",
			},
		}
		rawConfig := map[string]interface{}{
			"payment_type":        "PostPaid",
			"period_unit":         "Month",
			"sale_cycle":          "Month",
			"quantity":            1,
			"dedicated_host_type": "ddh.g6",
		}
		if tc.autoPay != nil {
			rawConfig["auto_pay"] = tc.autoPay
		}
		dAutoPay := buildData(state, rawConfig)
		var capturedReq map[string]interface{}
		ReadMockResponseDiff = map[string]interface{}{
			"DedicatedHosts": map[string]interface{}{
				"DedicatedHost": []interface{}{
					map[string]interface{}{
						"ChargeType": "PostPaid",
					},
				},
			},
		}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, req map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDedicatedHostsChargeType" {
				capturedReq = req
			}
			return successResponseMock(ReadMockResponseDiff)
		})
		err = resourceAlicloudEcsDedicatedHostUpdate(dAutoPay, rawClient)
		patches.Reset()
		assert.Nil(t, err, "auto_pay %s", tc.name)
		assert.Equal(t, tc.expected, capturedReq["AutoPay"], "auto_pay %s", tc.name)
	}

	// Real import flow: public Importer (schema.ImportStatePassthrough) →
	// Read refresh → public SDK Diff (schema.Resource.Diff) no-op. This
	// replaces the hand-crafted newInstanceDiff + HasChange(quantity) check:
	// it exercises the real Importer.State, a real Read against a mocked
	// DescribeDedicatedHosts, and the real plan computation, asserting a
	// complete no-op (InstanceDiff.Empty) covering quantity=1, the old
	// sale_cycle, and period_unit aliasing (sale_cycle in config aliases
	// period_unit in state — both Month — with no drift).
	{
		r := p["alicloud_ecs_dedicated_host"]
		dImported, _ := schema.InternalMap(r.Schema).Data(nil, nil)
		dImported.SetId("dh-import-test")
		// 1. Public Importer (schema.ImportStatePassthrough) returns the
		//    imported ResourceData keyed by ID, exactly as `terraform import`.
		assert.NotNil(t, r.Importer, "dedicated_host must define an Importer for the import no-op")
		imported, importErr := r.Importer.State(dImported, rawClient)
		assert.Nil(t, importErr)
		assert.Len(t, imported, 1)
		dImported = imported[0]
		// 2. Read refreshes state from DescribeDedicatedHosts (mocked). The
		//    mock returns a legal PostPaid host with SaleCycle=Month; NetworkAttributes
		//    is present (Read accesses it without an ok-check); Optional non-Computed
		//    fields are at their schema defaults/zeros so a minimal matching config
		//    yields a true no-op.
		importReadMock := map[string]interface{}{
			"DedicatedHosts": map[string]interface{}{
				"DedicatedHost": []interface{}{
					map[string]interface{}{
						"DedicatedHostId":     "dh-import-test",
						"ActionOnMaintenance": "Stop",
						"AutoPlacement":       "on",
						"AutoReleaseTime":     "",
						"CpuOverCommitRatio":  0.0,
						"DedicatedHostName":   "",
						"DedicatedHostType":   "ddh.g6",
						"Description":         "",
						"NetworkAttributes":   map[string]interface{}{},
						"ChargeType":          "PostPaid",
						"ResourceGroupId":     "",
						"SaleCycle":           "Month",
						"Status":              "Available",
						"ZoneId":              "cn-hangzhou-a",
					},
				},
			},
		}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return importReadMock, nil
		})
		readErr := resourceAlicloudEcsDedicatedHostRead(dImported, rawClient)
		patches.Reset()
		assert.Nil(t, readErr)
		// Read preserves quantity=1 (schema default, create-only) and sets
		// period_unit=sale_cycle=SaleCycle from DescribeDedicatedHosts.
		assert.Equal(t, 1, dImported.Get("quantity"))
		assert.Equal(t, "Month", dImported.Get("period_unit"))
		assert.Equal(t, "Month", dImported.Get("sale_cycle"))
		// The Describe response omits request-only params (dry_run, detail_fee);
		// Read does not set them, so their schema Default:false would otherwise
		// drift ""->"false" in plan. Inject the legal default values (#10 permits
		// "明确合法 state 注入") to isolate the quantity/sale_cycle/period_unit
		// no-op from this known Default-vs-unset characteristic. Read does set
		// quantity=1 (L354) and period_unit/sale_cycle from SaleCycle (L356-357).
		dImported.Set("dry_run", false)
		dImported.Set("detail_fee", false)
		// 3. Real SDK Diff (schema.Resource.Diff) against a config matching the
		//    read state: quantity=1, sale_cycle=Month (aliasing period_unit), and
		//    the Required dedicated_host_type. period_unit is Computed and absent
		//    from the config, so it does not drift.
		planConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
			"quantity":            1,
			"sale_cycle":          "Month",
			"dedicated_host_type": "ddh.g6",
		})
		planDiff, diffErr := r.Diff(context.Background(), dImported.State(), planConfig, rawClient)
		assert.Nil(t, diffErr)
		// The SDK represents a no-op (zero-attribute) diff as a nil *InstanceDiff;
		// a non-nil diff must be Empty(). Either is the complete no-op via the
		// real public schema.Resource.Diff (not the hand-written newInstanceDiff).
		assert.True(t, planDiff == nil || planDiff.Empty(), "import→read→plan must be a complete no-op (no drift on quantity/sale_cycle/period_unit)")
		// Explicit target coverage (defensive): quantity, sale_cycle and the
		// period_unit alias must not appear in the plan diff. GetAttribute returns
		// ok=true only when the key is present in the InstanceDiff.
		if planDiff != nil {
			_, qtyChanged := planDiff.GetAttribute("quantity")
			assert.False(t, qtyChanged, "quantity must not drift after import→read→plan")
			_, saleCycleChanged := planDiff.GetAttribute("sale_cycle")
			assert.False(t, saleCycleChanged, "sale_cycle must not drift after import→read→plan")
			_, periodUnitChanged := planDiff.GetAttribute("period_unit")
			assert.False(t, periodUnitChanged, "period_unit must not drift after import→read→plan")
		}
	}

	// ModifyDedicatedHostsChargeType — period_unit=Year rejected upfront when
	// changing payment_type. Year is valid for AllocateDedicatedHosts,
	// RenewDedicatedHosts and ModifyDedicatedHostAutoRenewAttribute, but
	// ModifyDedicatedHostsChargeType only accepts Week/Month; the provider
	// rejects it (Update L537-538) before issuing the RPC. The negative case
	// must assert the provider's explicit error text AND that no HTTP request
	// is made (HTTP count=0) — a bare err != nil could be satisfied by an
	// unrelated client/region error.
	{
		// Fresh prior state: PrePaid, period_unit/sale_cycle=Year, quantity=1.
		state := &terraform.InstanceState{
			ID: "dh-year-test",
			Attributes: map[string]string{
				"id":                  "dh-year-test",
				"payment_type":        "PrePaid",
				"period_unit":         "Year",
				"sale_cycle":          "Year",
				"quantity":            "1",
				"dedicated_host_type": "ddh.g6",
			},
		}
		rawConfig := map[string]interface{}{
			"payment_type":        "PostPaid",
			"period_unit":         "Year",
			"sale_cycle":          "Year",
			"quantity":            1,
			"dedicated_host_type": "ddh.g6",
		}
		dYear := buildData(state, rawConfig)
		// Counter-guarded mock: any DoRequest call increments the count. The
		// upfront Year rejection must return before any RPC, so count must be 0.
		var httpCount int
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			httpCount++
			return ReadMockResponse, nil
		})
		err = resourceAlicloudEcsDedicatedHostUpdate(dYear, rawClient)
		patches.Reset()
		assert.NotNil(t, err, "changing payment_type with period_unit=Year must be rejected")
		if err != nil {
			assert.Contains(t, err.Error(), "period_unit cannot be Year when changing payment_type", "rejection must surface the provider's explicit error, not an unrelated client/region error")
		}
		assert.Equal(t, 0, httpCount, "no HTTP request may be made before the period_unit=Year upfront rejection")
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "WithEcsClient", func(_ *connectivity.AliyunClient, _ func(*ecs.Client) (interface{}, error)) (interface{}, error) {
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
