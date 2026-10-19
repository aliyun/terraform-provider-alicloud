package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_adb_db_cluster_lake_version", &resource.Sweeper{
		Name: "alicloud_adb_db_cluster_lake_version",
		F:    testSweepAdbDbClusterLakeVersions,
	})
}

func testSweepAdbDbClusterLakeVersions(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeDBClusters"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	for {
		response, err = client.RpcPost("adb", "2021-12-01", action, nil, request, true)
		if err != nil {
			log.Println(WrapErrorf(err, DataDefaultErrorMsg, "alicloud_adb_db_cluster_lake_versions", action, AlibabaCloudSdkGoERROR))
			break
		}

		resp, err := jsonpath.Get("$.Items.DBCluster", response)
		if err != nil {
			log.Println(WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items.DBCluster", response))
			break
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["DBClusterDescription"])
			id := fmt.Sprint(item["DBClusterId"])
			if protectedSweepVpcName(name) {
				log.Printf("[INFO] Skipping ADB LakeVersion (NODELETING): %s (%s)", name, id)
				continue
			}
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping ADB LakeVersion: %s (%s)", name, id)
					continue
				}
			}
			log.Printf("[INFO] Deleting adb LakeVersion: %s (%s)", name, id)
			deleteAction := "DeleteDBCluster"
			deleteRequest := map[string]interface{}{
				"DBClusterId": id,
				"RegionId":    client.RegionId,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(5*time.Minute, func() *retry.RetryError {
				_, err = client.RpcPost("adb", "2021-12-01", deleteAction, nil, deleteRequest, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] Deleting ADB LakeVersion %s (%s) failed: %#v", name, id, err)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

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
		IDRefreshName:     resourceId,
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
				ImportStateVerifyIgnore: []string{"period", "enable_default_resource_group", "source_db_cluster_id", "backup_set_id", "restore_type", "restore_to_time"},
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
		IDRefreshName:     resourceId,
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
				ImportStateVerifyIgnore: []string{"period", "enable_default_resource_group", "source_db_cluster_id", "backup_set_id", "restore_type", "restore_to_time"},
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
		IDRefreshName:     resourceId,
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
		IDRefreshName:     resourceId,
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
