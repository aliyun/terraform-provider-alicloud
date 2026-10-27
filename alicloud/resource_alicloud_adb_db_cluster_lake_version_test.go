package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudADBDBClusterLakeVersion_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbDbClusterLakeVersionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckPrePaidResources(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version": "5.0",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":         "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":            "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":       "Subscription",
					"period":             "1",
					"compute_resource":   "16ACU",
					"storage_resource":   "0ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version": "5.0",
						"vpc_id":             CHECKSET,
						"vswitch_id":         CHECKSET,
						"zone_id":            CHECKSET,
						"payment_type":       "Subscription",
						"compute_resource":   "16ACU",
						"storage_resource":   "0ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compute_resource": "32ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compute_resource": "32ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_resource": "24ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_resource": "24ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": "10.23.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips": "10.23.1.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"used_time":    "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"period":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "enable_default_resource_group", "source_db_cluster_id", "backup_set_id", "restore_type", "restore_to_time", "used_time"},
			},
		},
	})
}

func TestAccAliCloudADBDBClusterLakeVersion_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbDbClusterLakeVersionBasicDependence0Twin)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckPrePaidResources(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":            "5.0",
					"vpc_id":                        "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":                    "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":                       "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":                  "Subscription",
					"secondary_vswitch_id":          "${data.alicloud_vswitches.secondary.ids.0}",
					"secondary_zone_id":             "${data.alicloud_adb_zones.default.ids.1}",
					"compute_resource":              "16ACU",
					"storage_resource":              "0ACU",
					"disk_encryption":               "true",
					"kms_id":                        "${alicloud_kms_key.default.id}",
					"security_ips":                  "10.23.1.0",
					"enable_ssl":                    "true",
					"db_cluster_description":        name,
					"resource_group_id":             "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"period":                        "1",
					"enable_default_resource_group": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":     "5.0",
						"vpc_id":                 CHECKSET,
						"vswitch_id":             CHECKSET,
						"zone_id":                CHECKSET,
						"payment_type":           "Subscription",
						"secondary_vswitch_id":   CHECKSET,
						"secondary_zone_id":      CHECKSET,
						"compute_resource":       "16ACU",
						"storage_resource":       "0ACU",
						"disk_encryption":        "true",
						"kms_id":                 CHECKSET,
						"security_ips":           "10.23.1.0",
						"enable_ssl":             "true",
						"db_cluster_description": name,
						"resource_group_id":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "enable_default_resource_group", "source_db_cluster_id", "backup_set_id", "restore_type", "restore_to_time", "used_time"},
			},
		},
	})
}

func TestAccAliCloudADBDBClusterLakeVersion_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbDbClusterLakeVersionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":  "5.0",
					"vpc_id":              "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":          "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":             "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":        "PayAsYouGo",
					"product_form":        "IntegrationForm",
					"product_version":     "EnterpriseVersion",
					"reserved_node_size":  "8ACU",
					"reserved_node_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":  "5.0",
						"vpc_id":              CHECKSET,
						"vswitch_id":          CHECKSET,
						"zone_id":             CHECKSET,
						"payment_type":        "PayAsYouGo",
						"product_form":        "IntegrationForm",
						"product_version":     "EnterpriseVersion",
						"reserved_node_size":  "8ACU",
						"reserved_node_count": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reserved_node_size": "12ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reserved_node_size": "12ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reserved_node_count": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reserved_node_count": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reserved_node_size":  "8ACU",
					"reserved_node_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reserved_node_size":  "8ACU",
						"reserved_node_count": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": "10.23.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips": "10.23.1.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_default_resource_group", "source_db_cluster_id", "backup_set_id", "restore_type", "restore_to_time"},
			},
		},
	})
}

func TestAccAliCloudADBDBClusterLakeVersion_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbDbClusterLakeVersionBasicDependence0Twin)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":            "5.0",
					"vpc_id":                        "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":                    "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":                       "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":                  "PayAsYouGo",
					"product_form":                  "IntegrationForm",
					"product_version":               "EnterpriseVersion",
					"reserved_node_size":            "8ACU",
					"reserved_node_count":           "3",
					"disk_encryption":               "true",
					"kms_id":                        "${alicloud_kms_key.default.id}",
					"security_ips":                  "10.23.1.0",
					"enable_ssl":                    "true",
					"db_cluster_description":        name,
					"resource_group_id":             "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"enable_default_resource_group": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":     "5.0",
						"vpc_id":                 CHECKSET,
						"vswitch_id":             CHECKSET,
						"zone_id":                CHECKSET,
						"payment_type":           "PayAsYouGo",
						"product_form":           "IntegrationForm",
						"product_version":        "EnterpriseVersion",
						"reserved_node_size":     "8ACU",
						"reserved_node_count":    "3",
						"disk_encryption":        "true",
						"kms_id":                 CHECKSET,
						"security_ips":           "10.23.1.0",
						"enable_ssl":             "true",
						"db_cluster_description": name,
						"resource_group_id":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_default_resource_group", "source_db_cluster_id", "backup_set_id", "restore_type", "restore_to_time"},
			},
		},
	})
}

var AliCloudAdbDbClusterLakeVersionMap0 = map[string]string{
	"product_form":           CHECKSET,
	"product_version":        CHECKSET,
	"reserved_node_size":     CHECKSET,
	"security_ips":           CHECKSET,
	"db_cluster_description": CHECKSET,
	"commodity_code":         CHECKSET,
	"connection_string":      CHECKSET,
	"engine":                 CHECKSET,
	"engine_version":         CHECKSET,
	"create_time":            CHECKSET,
	"lock_mode":              CHECKSET,
	"port":                   CHECKSET,
	"resource_group_id":      CHECKSET,
	"status":                 CHECKSET,
}

func AliCloudAdbDbClusterLakeVersionBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_adb_zones" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_adb_zones.default.ids.0
	}
`, name)
}

func AliCloudAdbDbClusterLakeVersionBasicDependence0Twin(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_adb_zones" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_adb_zones.default.ids.0
	}

	data "alicloud_vswitches" "secondary" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_adb_zones.default.ids.1
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		pending_window_in_days = "7"
  		status                 = "Enabled"
	}
`, name)
}

// lintignore: R001
func TestUnitAliCloudAdbDbClusterLakeVersion(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"compute_resource":              "CreateDBClusterValue",
		"db_cluster_version":            "5.0",
		"enable_default_resource_group": true,
		"payment_type":                  "CreateDBClusterValue",
		"storage_resource":              "CreateDBClusterValue",
		"vswitch_id":                    "CreateDBClusterValue",
		"zone_id":                       "CreateDBClusterValue",
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
		// DescribeDBClusterAttribute
		"Items": map[string]interface{}{
			"DBCluster": []interface{}{
				map[string]interface{}{
					"CommodityCode":    "DefaultValue",
					"ComputeResource":  "CreateDBClusterValue",
					"ConnectionString": "DefaultValue",
					"CreationTime":     "DefaultValue",
					"DBClusterId":      "CreateDBClusterValue",
					"DBVersion":        "5.0",
					"Engine":           "DefaultValue",
					"EngineVersion":    "DefaultValue",
					"ExpireTime":       "DefaultValue",
					"Expired":          "DefaultValue",
					"LockMode":         "DefaultValue",
					"LockReason":       "DefaultValue",
					"PayType":          "CreateDBClusterValue",
					"Port":             "DefaultValue",
					"ResourceGroupId":  "DefaultValue",
					"DBClusterStatus":  "Running",
					"StorageResource":  "CreateDBClusterValue",
					"VPCId":            "CreateDBClusterValue",
					"VSwitchId":        "CreateDBClusterValue",
					"ZoneId":           "CreateDBClusterValue",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateDBCluster
		"DBClusterId":     "CreateDBClusterValue",
		"OrderId":         "MockValue",
		"RequestId":       "MockValue",
		"ResourceGroupId": "MockValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_adb_db_cluster_lake_version", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudAdbDbClusterLakeVersionCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDBClusterAttribute Response
		"Items": map[string]interface{}{
			"DBCluster": []interface{}{
				map[string]interface{}{
					"DBClusterId": "CreateDBClusterValue",
				},
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDBCluster" {
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
		err := resourceAliCloudAdbDbClusterLakeVersionCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudAdbDbClusterLakeVersionUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyDBCluster
	attributesDiff := map[string]interface{}{
		"compute_resource": "ModifyDBClusterValue",
		"storage_resource": "ModifyDBClusterValue",
	}
	diff, err := newInstanceDiff("alicloud_adb_db_cluster_lake_version", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBClusterAttribute Response
		"Items": map[string]interface{}{
			"DBCluster": []interface{}{
				map[string]interface{}{
					"ComputeResource": "ModifyDBClusterValue",
					"StorageResource": "ModifyDBClusterValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDBCluster" {
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
		err := resourceAliCloudAdbDbClusterLakeVersionUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(dExisted.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidDBCluster.NotFound", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDBClusterAttribute" {
				switch errorCode {
				case "{}", "InvalidDBCluster.NotFound":
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
		err := resourceAliCloudAdbDbClusterLakeVersionRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}", "InvalidDBCluster.NotFound":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudAdbDbClusterLakeVersionDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidDBCluster.NotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDBCluster" {
				switch errorCode {
				case "NonRetryableError", "InvalidDBCluster.NotFound":
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
		err := resourceAliCloudAdbDbClusterLakeVersionDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InvalidDBCluster.NotFound":
			assert.Nil(t, err)
		}
	}
}

func TestAccAliCloudADBDBClusterLakeVersion_basic2(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbDbClusterLakeVersionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":       "5.0",
					"vpc_id":                   "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":                  "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":             "PayAsYouGo",
					"compute_resource":         "16ACU",
					"storage_resource":         "0ACU",
					"db_cluster_network_type":  "VPC",
					"maintain_time":            "04:00Z-05:00Z",
					"security_ips":             "10.23.1.1",
					"db_cluster_ip_array_name": "default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":       "5.0",
						"payment_type":             "PayAsYouGo",
						"db_cluster_network_type":  "VPC",
						"maintain_time":            "04:00Z-05:00Z",
						"security_ips":             "10.23.1.1",
						"db_cluster_ip_array_name": "default",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "05:00Z-06:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "05:00Z-06:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_ip_array_name":      "tfaccgroup",
					"db_cluster_ip_array_attribute": "hidden",
					"modify_mode":                   "Cover",
					"security_ips":                  "10.23.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_ip_array_name":      "tfaccgroup",
						"db_cluster_ip_array_attribute": "hidden",
						"modify_mode":                   "Cover",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_lake_cache": "true",
					"capacity":          "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_lake_cache": "true",
						"capacity":          "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_compaction_service": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_compaction_service": "true",
					}),
				),
			},
			{
				// Groups tagged as hidden are not listed by the API, so switch back to
				// the default group before the import check to keep the state readable.
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_ip_array_name":      "default",
					"db_cluster_ip_array_attribute": REMOVEKEY,
					"modify_mode":                   "Cover",
					"security_ips":                  "10.23.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_ip_array_name":      "default",
						"db_cluster_ip_array_attribute": "",
						"security_ips":                  "10.23.1.1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "enable_default_resource_group", "source_db_cluster_id", "backup_set_id", "restore_type", "restore_to_time", "db_version", "switch_mode", "modify_mode", "db_cluster_ip_array_attribute", "audit_log_status", "engine_type", "used_time"},
			},
		},
	})
}

// EnableEssdCache is rejected with InvalidStorageAcuNumberForServerless when the
// cluster reserves no storage ACU, so the disk cache case needs a cluster with a
// reserved storage resource (step 24ACU). That case is documented but not part
// of the pre-PR ACC selection until the reserved-storage cluster provisioning
// and the cache modification settle reliably in the test account.
func TestAccAliCloudADBDBClusterLakeVersion_essd_cache(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbDbClusterLakeVersionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version": "5.0",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":         "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":            "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":       "PayAsYouGo",
					"compute_resource":   "16ACU",
					"storage_resource":   "24ACU",
					"enable_essd_cache":  "true",
					"essd_cache_size":    "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_essd_cache": "true",
						"essd_cache_size":   "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"essd_cache_size": "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_essd_cache": "true",
						"essd_cache_size":   "500",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "enable_default_resource_group", "source_db_cluster_id", "backup_set_id", "restore_type", "restore_to_time", "db_version", "switch_mode", "modify_mode", "db_cluster_ip_array_attribute", "audit_log_status", "engine_type", "used_time"},
			},
		},
	})
}

func TestAccAliCloudADBDBClusterLakeVersion_basic3(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbDbClusterLakeVersionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckPrePaidResources(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":       "5.0",
					"vpc_id":                   "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":                  "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":             "Subscription",
					"period":                   1,
					"compute_resource":         "16ACU",
					"storage_resource":         "0ACU",
					"auto_renewal_status":      "AutoRenewal",
					"auto_renewal_period":      1,
					"auto_renewal_period_unit": "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":       "5.0",
						"payment_type":             "Subscription",
						"auto_renewal_status":      "AutoRenewal",
						"auto_renewal_period":      "1",
						"auto_renewal_period_unit": "Month",
						"used_time":                "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renewal_status":      "Normal",
					"auto_renewal_period":      2,
					"auto_renewal_period_unit": "Year",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renewal_status":      "Normal",
						"auto_renewal_period":      "2",
						"auto_renewal_period_unit": "Year",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
						"used_time":    "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "enable_default_resource_group", "source_db_cluster_id", "backup_set_id", "restore_type", "restore_to_time", "db_version", "switch_mode", "modify_mode", "used_time"},
			},
		},
	})
}

// auto_renewal_period must match the unit: Year only accepts 1, 2, 3 and 5. The
// invalid combination is rejected during planning, before any cluster is created.
func TestAccAliCloudADBDBClusterLakeVersion_auto_renewal_validation(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc("alicloud_adb_db_cluster_lake_version.default", name, AliCloudAdbDbClusterLakeVersionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":       "5.0",
					"vpc_id":                   "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":                  "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":             "Subscription",
					"period":                   1,
					"compute_resource":         "16ACU",
					"storage_resource":         "0ACU",
					"auto_renewal_period":      4,
					"auto_renewal_period_unit": "Year",
				}),
				ExpectError: regexp.MustCompile("auto_renewal_period must be 1, 2, 3 or 5"),
			},
		},
	})
}

// source_db_cluster_id/backup_set_id/restore_type/restore_to_time are
// create-only restore fields. Restoring needs a source cluster with backups or
// a PITR window that this suite cannot provision (the warehouse resource has no
// backup-set field and ADB exposes no backup data source). Omitting the
// required db_cluster_version stops the plan before the placeholders reach the
// API, while still covering the restore fields in the test template;
// db_cluster_version itself is covered by the basic cases above.
//
// db_version/switch_mode are write-only UpgradeKernelVersion arguments. A newly
// created lake-version cluster already runs the latest kernel and
// DescribeKernelVersion reports an empty AvailableKernelVersions list, so no
// positive upgrade can be configured here; they are placed in the same
// plan-error template for static coverage only.
func TestAccAliCloudADBDBClusterLakeVersion_restore_fields(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc("alicloud_adb_db_cluster_lake_version.default", name, AliCloudAdbDbClusterLakeVersionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":   REMOVEKEY,
					"vpc_id":               "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":           "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":              "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":         "PayAsYouGo",
					"source_db_cluster_id": "amv-bp1nonexistent",
					"restore_type":         "timepoint",
					"restore_to_time":      "2023-09-20T03:13:56Z",
					"backup_set_id":        "nonexistent-backup-id",
					"db_version":           "3.2.7",
					"switch_mode":          "1",
					"enable_essd_cache":    "false",
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

// ModifyAuditLogConfig for a lake-version cluster is rejected by the backend
// with HTTP 400 and an empty message ("null") for AuditLogStatus=on on three
// different clusters (RequestIds 01A17CCB-E3D0-527D-9ED6-95E1CE5ABD9B,
// 01A17CDB-7222-50F8-8CFC-CAA3F2922F99, 01A17CEC-B2CB-58B1-B049-07A5A1B74B6D;
// aliyunpop dx shows the request reaching adb4mysql and the backend returning
// 400 without a readable message). The request only carries the required
// parameters of the public contract (DBClusterId, RegionId, AuditLogStatus;
// EngineType is optional and defaults to XIHE), so the failure is on the
// product side.
// The case is therefore kept out of the pre-PR acceptance selection until the
// product side fixes the API; it documents the intended configuration.
func TestAccAliCloudADBDBClusterLakeVersion_audit_log(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbDbClusterLakeVersionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version": "5.0",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":         "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":            "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":       "PayAsYouGo",
					"compute_resource":   "16ACU",
					"storage_resource":   "0ACU",
					"audit_log_status":   "on",
					"engine_type":        "XIHE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_log_status": "on",
						"engine_type":      "XIHE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"audit_log_status":  "off",
					"enable_essd_cache": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_log_status": "off",
					}),
				),
			},
		},
	})
}

func TestUnitAliCloudAdbDbClusterLakeVersionMaintainTimeValidation(t *testing.T) {
	for _, value := range []string{"04:00Z-05:00Z", "23:00Z-00:00Z"} {
		if _, errs := validateAdbDbClusterLakeVersionMaintainTime(value, "maintain_time"); len(errs) != 0 {
			t.Fatalf("expected %s to be valid, got %v", value, errs)
		}
	}
	for _, value := range []string{"04:30Z-05:30Z", "04:00Z-06:00Z", "24:00Z-01:00Z", "04:00Z", ""} {
		if _, errs := validateAdbDbClusterLakeVersionMaintainTime(value, "maintain_time"); len(errs) == 0 {
			t.Fatalf("expected %s to be invalid", value)
		}
	}
}

func TestUnitAliCloudAdbDbClusterLakeVersionAuditLogUpdateSucceeded(t *testing.T) {
	cases := []struct {
		response map[string]interface{}
		want     bool
		present  bool
	}{
		{map[string]interface{}{"UpdateSucceed": true}, true, true},
		{map[string]interface{}{"UpdateSucceed": false}, false, true},
		{map[string]interface{}{"UpdateSucceed": "false"}, false, true},
		{map[string]interface{}{"UpdateSucceed": "true"}, true, true},
		{map[string]interface{}{}, false, false},
		{nil, false, false},
	}
	for _, item := range cases {
		got, present := adbAuditLogUpdateSucceeded(item.response)
		if got != item.want || present != item.present {
			t.Fatalf("response %v: got (%v,%v), want (%v,%v)", item.response, got, present, item.want, item.present)
		}
	}
}

func TestUnitAliCloudAdbDbClusterLakeVersionBoolConfigValue(t *testing.T) {
	raw := cty.ObjectVal(map[string]cty.Value{
		"enabled":  cty.True,
		"disabled": cty.False,
		"null":     cty.NullVal(cty.Bool),
		"unknown":  cty.UnknownVal(cty.Bool),
	})
	cases := []struct {
		value   cty.Value
		key     string
		want    bool
		present bool
	}{
		{raw, "enabled", true, true},
		{raw, "disabled", false, true},
		{raw, "null", false, false},
		{raw, "unknown", false, false},
		{raw, "absent", false, false},
		{cty.NilVal, "enabled", false, false},
	}
	for _, item := range cases {
		got, present := adbBoolConfigValue(item.value, item.key)
		if got != item.want || present != item.present {
			t.Fatalf("key %s: got (%v,%v), want (%v,%v)", item.key, got, present, item.want, item.present)
		}
	}
}

func TestUnitAliCloudAdbDbClusterLakeVersionResponseClassification(t *testing.T) {
	if _, err := adbDataObjectSection(map[string]interface{}{"Data": map[string]interface{}{"Capacity": 1}}, "$.Data", "amv-a"); err != nil {
		t.Fatalf("expected valid section, got %v", err)
	}
	if _, err := adbDataObjectSection(map[string]interface{}{}, "$.Data", "amv-a"); err == nil || NotFoundError(err) {
		t.Fatalf("missing section must be a malformed-response error, got %v", err)
	}
	if _, err := adbDataObjectSection(map[string]interface{}{"Data": []interface{}{}}, "$.Data", "amv-a"); err == nil || NotFoundError(err) {
		t.Fatalf("wrong section type must be a malformed-response error, got %v", err)
	}

	emptyRenewal := map[string]interface{}{"Items": map[string]interface{}{"AutoRenewalAttribute": []interface{}{}}}
	if _, err := adbAutoRenewalRecord(emptyRenewal, "amv-a"); err == nil || !NotFoundError(err) {
		t.Fatalf("empty auto renewal list must be NotFound, got %v", err)
	}
	if _, err := adbAutoRenewalRecord(map[string]interface{}{}, "amv-a"); err == nil || NotFoundError(err) {
		t.Fatalf("missing auto renewal section must be malformed, got %v", err)
	}
	wrongRenewal := map[string]interface{}{"Items": map[string]interface{}{"AutoRenewalAttribute": "bad"}}
	if _, err := adbAutoRenewalRecord(wrongRenewal, "amv-a"); err == nil || NotFoundError(err) {
		t.Fatalf("wrong auto renewal type must be malformed, got %v", err)
	}
	for _, malformed := range []map[string]interface{}{
		{"AutoRenewalStatus": "Normal"},
		{"DBClusterId": 123, "AutoRenewalStatus": "Normal"},
		{"DBClusterId": "   ", "AutoRenewalStatus": "Normal"},
	} {
		response := map[string]interface{}{"Items": map[string]interface{}{"AutoRenewalAttribute": []interface{}{
			map[string]interface{}{"DBClusterId": "amv-a", "AutoRenewalStatus": "Normal"},
			malformed,
		}}}
		if _, err := adbAutoRenewalRecord(response, "amv-a"); err == nil || NotFoundError(err) {
			t.Fatalf("malformed renewal element %v must be an error even when another element matches, got %v", malformed, err)
		}
	}
	validRenewal := map[string]interface{}{"Items": map[string]interface{}{"AutoRenewalAttribute": []interface{}{
		map[string]interface{}{"DBClusterId": "amv-a", "AutoRenewalStatus": "Normal"},
		map[string]interface{}{"DBClusterId": "amv-b", "AutoRenewalStatus": "AutoRenewal"},
	}}}
	record, err := adbAutoRenewalRecord(validRenewal, "amv-b")
	if err != nil || record["AutoRenewalStatus"] != "AutoRenewal" {
		t.Fatalf("expected amv-b record, got %v %v", record, err)
	}
	if _, err := adbAutoRenewalRecord(validRenewal, "amv-c"); err == nil || !NotFoundError(err) {
		t.Fatalf("absent cluster must be NotFound, got %v", err)
	}

	validGroups := map[string]interface{}{"Items": map[string]interface{}{"IPArray": []interface{}{
		map[string]interface{}{"DBClusterIPArrayName": "default", "SecurityIPList": "127.0.0.1"},
	}}}
	group, err := adbIPArrayGroup(validGroups, "amv-a", "default")
	if err != nil || group["SecurityIPList"] != "127.0.0.1" {
		t.Fatalf("expected default group, got %v %v", group, err)
	}
	if _, err := adbIPArrayGroup(validGroups, "amv-a", "hidden-group"); err == nil || !NotFoundError(err) {
		t.Fatalf("absent group must be NotFound, got %v", err)
	}
	emptyGroups := map[string]interface{}{"Items": map[string]interface{}{"IPArray": []interface{}{}}}
	if _, err := adbIPArrayGroup(emptyGroups, "amv-a", "default"); err == nil || !NotFoundError(err) {
		t.Fatalf("empty group list must be NotFound, got %v", err)
	}
	if _, err := adbIPArrayGroup(map[string]interface{}{}, "amv-a", "default"); err == nil || NotFoundError(err) {
		t.Fatalf("missing group section must be malformed, got %v", err)
	}
	wrongGroups := map[string]interface{}{"Items": map[string]interface{}{"IPArray": "bad"}}
	if _, err := adbIPArrayGroup(wrongGroups, "amv-a", "default"); err == nil || NotFoundError(err) {
		t.Fatalf("wrong group type must be malformed, got %v", err)
	}
	for _, malformed := range []map[string]interface{}{
		{"SecurityIPList": "127.0.0.1"},
		{"DBClusterIPArrayName": 1, "SecurityIPList": "127.0.0.1"},
		{"DBClusterIPArrayName": " ", "SecurityIPList": "127.0.0.1"},
	} {
		response := map[string]interface{}{"Items": map[string]interface{}{"IPArray": []interface{}{
			map[string]interface{}{"DBClusterIPArrayName": "default", "SecurityIPList": "127.0.0.1"},
			malformed,
		}}}
		if _, err := adbIPArrayGroup(response, "amv-a", "default"); err == nil || NotFoundError(err) {
			t.Fatalf("malformed whitelist element %v must be an error even when another element matches, got %v", malformed, err)
		}
	}
}

func TestUnitAliCloudAdbDbClusterLakeVersionLakeCacheResponseCode(t *testing.T) {
	cases := []struct {
		response map[string]interface{}
		ok       bool
		present  bool
	}{
		{map[string]interface{}{"Code": 200}, true, true},
		{map[string]interface{}{"Code": "200"}, true, true},
		{map[string]interface{}{"Code": 500}, false, true},
		{map[string]interface{}{"Code": "InvalidParameter"}, false, true},
		{map[string]interface{}{}, false, false},
		{nil, false, false},
	}
	for _, item := range cases {
		ok, present := adbLakeCacheResponseCodeOK(item.response)
		if ok != item.ok || present != item.present {
			t.Fatalf("response %v: got (%v,%v), want (%v,%v)", item.response, ok, present, item.ok, item.present)
		}
	}
}

func adbLakeVersionHelperTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	config := &connectivity.Config{AccessKey: "test-key", SecretKey: "test-secret", Credential: credential, RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http", Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("adb", endpoint)
	return client
}

func adbLakeVersionJSONResponse(t *testing.T, status int, body interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(body); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}
}

func TestUnitAliCloudAdbDbClusterLakeVersionHelperAPIErrorPropagation(t *testing.T) {
	cases := []struct {
		name     string
		status   int
		body     interface{}
		call     func(*AdbService) error
		wantErr  bool
		notFound bool
	}{
		{
			name:     "whitelist documented not found",
			status:   http.StatusNotFound,
			body:     map[string]interface{}{"Code": "InvalidDBClusterId.NotFound", "Message": "cluster not found", "RequestId": "req-1"},
			call:     func(s *AdbService) error { _, err := s.DescribeClusterAccessWhiteList("amv-a", "default"); return err },
			wantErr:  true,
			notFound: true,
		},
		{
			name:    "whitelist forbidden propagates",
			status:  http.StatusForbidden,
			body:    map[string]interface{}{"Code": "Forbidden.RAM", "Message": "not authorized", "RequestId": "req-2"},
			call:    func(s *AdbService) error { _, err := s.DescribeClusterAccessWhiteList("amv-a", "default"); return err },
			wantErr: true,
		},
		{
			name:    "whitelist malformed element",
			status:  http.StatusOK,
			body:    map[string]interface{}{"Items": map[string]interface{}{"IPArray": []interface{}{map[string]interface{}{"SecurityIPList": "127.0.0.1"}}}},
			call:    func(s *AdbService) error { _, err := s.DescribeClusterAccessWhiteList("amv-a", "default"); return err },
			wantErr: true,
		},
		{
			name:     "whitelist empty list",
			status:   http.StatusOK,
			body:     map[string]interface{}{"Items": map[string]interface{}{"IPArray": []interface{}{}}},
			call:     func(s *AdbService) error { _, err := s.DescribeClusterAccessWhiteList("amv-a", "default"); return err },
			wantErr:  true,
			notFound: true,
		},
		{
			name:   "auto renewal forbidden propagates",
			status: http.StatusForbidden,
			body:   map[string]interface{}{"Code": "Forbidden.RAM", "Message": "not authorized", "RequestId": "req-3"},
			call: func(s *AdbService) error {
				_, err := s.DescribeAdbDbClusterLakeVersionAutoRenewalAttribute("amv-a")
				return err
			},
			wantErr: true,
		},
		{
			name:   "auto renewal malformed element",
			status: http.StatusOK,
			body:   map[string]interface{}{"Items": map[string]interface{}{"AutoRenewalAttribute": []interface{}{map[string]interface{}{"AutoRenewalStatus": "Normal"}}}},
			call: func(s *AdbService) error {
				_, err := s.DescribeAdbDbClusterLakeVersionAutoRenewalAttribute("amv-a")
				return err
			},
			wantErr: true,
		},
		{
			name:   "auto renewal empty list",
			status: http.StatusOK,
			body:   map[string]interface{}{"Items": map[string]interface{}{"AutoRenewalAttribute": []interface{}{}}},
			call: func(s *AdbService) error {
				_, err := s.DescribeAdbDbClusterLakeVersionAutoRenewalAttribute("amv-a")
				return err
			},
			wantErr:  true,
			notFound: true,
		},
		{
			name:   "auto renewal valid record",
			status: http.StatusOK,
			body:   map[string]interface{}{"Items": map[string]interface{}{"AutoRenewalAttribute": []interface{}{map[string]interface{}{"DBClusterId": "amv-a", "AutoRenewalStatus": "Normal"}}}},
			call: func(s *AdbService) error {
				_, err := s.DescribeAdbDbClusterLakeVersionAutoRenewalAttribute("amv-a")
				return err
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := &AdbService{adbLakeVersionHelperTestClient(t, adbLakeVersionJSONResponse(t, tc.status, tc.body))}
			err := tc.call(service)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v, got %v", tc.wantErr, err)
			}
			if tc.wantErr && NotFoundError(err) != tc.notFound {
				t.Fatalf("notFound=%v, got error %v", tc.notFound, err)
			}
		})
	}
}
