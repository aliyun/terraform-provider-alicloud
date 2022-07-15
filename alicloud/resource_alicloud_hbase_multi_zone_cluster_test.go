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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudHBaseMultiZoneCluster_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbase_multi_zone_cluster.default"
	checkoutSupportedRegions(t, true, connectivity.HBaseMultiZoneClusterSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudHBaseMultiZoneClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HBaseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbaseMultiZoneCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbasemultizonecluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBaseMultiZoneClusterBasicDependence0)
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
					"arch_version":      "2.0",
					"engine_version":    "2.0",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"Created": "tfTestAcc0",
						"For":     "Tftestacc 0",
					},
					"engine":                 "hbaseue",
					"vpc_id":                 "${data.alicloud_vpcs.default.ids.0}",
					"core_disk_type":         "cloud_ssd",
					"master_instance_type":   "hbase.sn1.large",
					"log_disk_size":          "400",
					"core_instance_type":     "hbase.sn1.large",
					"cluster_name":           "${var.name}",
					"multi_zone_combination": "cn-hangzhou-bef-aliyun",
					"core_node_count":        "4",
					"log_instance_type":      "hbase.sn1.large",
					"log_node_count":         "4",
					"log_disk_type":          "cloud_ssd",
					"core_disk_size":         "400",
					"payment_type":           "PayAsYouGo",
					"primary_zone_id":        "cn-hangzhou-b",
					"primary_vswitch_id":     "${data.alicloud_vswitches.cn-hangzhou-b.ids.0}",
					"standby_zone_id":        "cn-hangzhou-e",
					"standby_vswitch_id":     "${data.alicloud_vswitches.cn-hangzhou-e.ids.0}",
					"arbiter_zone_id":        "cn-hangzhou-f",
					"arbiter_vswitch_id":     "${data.alicloud_vswitches.cn-hangzhou-f.ids.0}",
					"security_ip_list":       "127.0.0.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"arch_version":           "2.0",
						"engine_version":         "2.0",
						"resource_group_id":      CHECKSET,
						"tags.%":                 "2",
						"tags.Created":           "tfTestAcc0",
						"tags.For":               "Tftestacc 0",
						"engine":                 "hbaseue",
						"vpc_id":                 CHECKSET,
						"core_disk_type":         "cloud_ssd",
						"master_instance_type":   "hbase.sn1.large",
						"log_disk_size":          "400",
						"core_instance_type":     "hbase.sn1.large",
						"cluster_name":           name,
						"multi_zone_combination": "cn-hangzhou-bef-aliyun",
						"core_node_count":        "4",
						"log_instance_type":      "hbase.sn1.large",
						"log_node_count":         "4",
						"log_disk_type":          "cloud_ssd",
						"core_disk_size":         "400",
						"payment_type":           "PayAsYouGo",
						"primary_zone_id":        "cn-hangzhou-b",
						"standby_zone_id":        "cn-hangzhou-e",
						"arbiter_zone_id":        "cn-hangzhou-f",
						"standby_vswitch_id":     CHECKSET,
						"arbiter_vswitch_id":     CHECKSET,
						"primary_vswitch_id":     CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"immediate_delete_flag", "auto_renew_period", "security_ip_list", "period_unit", "primary_core_node_count", "period", "standby_core_node_count"},
			},
		},
	})
}

var AlicloudHBaseMultiZoneClusterMap0 = map[string]string{
	"status":                  CHECKSET,
	"immediate_delete_flag":   NOSET,
	"security_ip_list":        NOSET,
	"auto_renew_period":       NOSET,
	"period_unit":             NOSET,
	"primary_core_node_count": NOSET,
	"period":                  NOSET,
	"standby_core_node_count": NOSET,
}

func AlicloudHBaseMultiZoneClusterBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default"{
		status = "OK"
	}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "cn-hangzhou-b" {
  zone_id = "cn-hangzhou-b"
  vpc_id  = data.alicloud_vpcs.default.ids.0
}
data "alicloud_vswitches" "cn-hangzhou-e" {
  zone_id = "cn-hangzhou-e"
  vpc_id  = data.alicloud_vpcs.default.ids.0
}
data "alicloud_vswitches" "cn-hangzhou-f" {
  zone_id = "cn-hangzhou-f"
  vpc_id  = data.alicloud_vpcs.default.ids.0
}
`, name)
}

func TestAccAlicloudHBaseMultiZoneCluster_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbase_multi_zone_cluster.default"
	checkoutSupportedRegions(t, true, connectivity.HBaseMultiZoneClusterSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudHBaseMultiZoneClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HBaseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbaseMultiZoneCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbasemultizonecluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBaseMultiZoneClusterBasicDependence0)
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
					"arch_version":           "2.0",
					"engine_version":         "2.0",
					"cluster_name":           "${var.name}",
					"log_disk_type":          "cloud_ssd",
					"core_node_count":        "4",
					"vpc_id":                 "${data.alicloud_vpcs.default.ids.0}",
					"primary_zone_id":        "cn-hangzhou-b",
					"log_disk_size":          "400",
					"master_instance_type":   "hbase.sn1.large",
					"standby_vswitch_id":     "${data.alicloud_vswitches.cn-hangzhou-e.ids.0}",
					"core_instance_type":     "hbase.sn1.large",
					"multi_zone_combination": "cn-hangzhou-bef-aliyun",
					"core_disk_type":         "cloud_ssd",
					"arbiter_vswitch_id":     "${data.alicloud_vswitches.cn-hangzhou-f.ids.0}",
					"payment_type":           "PayAsYouGo",
					"arbiter_zone_id":        "cn-hangzhou-f",
					"core_disk_size":         "400",
					"primary_vswitch_id":     "${data.alicloud_vswitches.cn-hangzhou-b.ids.0}",
					"engine":                 "hbaseue",
					"log_node_count":         "4",
					"log_instance_type":      "hbase.sn1.large",
					"standby_zone_id":        "cn-hangzhou-e",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"arch_version":           "2.0",
						"engine_version":         "2.0",
						"cluster_name":           name,
						"log_disk_type":          "cloud_ssd",
						"core_node_count":        "4",
						"vpc_id":                 CHECKSET,
						"primary_zone_id":        "cn-hangzhou-b",
						"log_disk_size":          "400",
						"master_instance_type":   "hbase.sn1.large",
						"standby_vswitch_id":     CHECKSET,
						"core_instance_type":     "hbase.sn1.large",
						"multi_zone_combination": "cn-hangzhou-bef-aliyun",
						"core_disk_type":         "cloud_ssd",
						"arbiter_vswitch_id":     CHECKSET,
						"payment_type":           "PayAsYouGo",
						"arbiter_zone_id":        "cn-hangzhou-f",
						"core_disk_size":         "400",
						"primary_vswitch_id":     CHECKSET,
						"engine":                 "hbaseue",
						"log_node_count":         "4",
						"log_instance_type":      "hbase.sn1.large",
						"standby_zone_id":        "cn-hangzhou-e",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"master_instance_type": "hbase.sn1.2xlarge",
					"core_instance_type":   "hbase.sn1.2xlarge",
					"log_instance_type":    "hbase.sn1.2xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_instance_type": "hbase.sn1.2xlarge",
						"core_instance_type":   "hbase.sn1.2xlarge",
						"log_instance_type":    "hbase.sn1.2xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_disk_size":  "500",
					"core_disk_size": "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_disk_size":  "500",
						"core_disk_size": "500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"core_node_count": "5",
					"log_node_count":  "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"core_node_count": "5",
						"log_node_count":  "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"core_disk_size": "400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"core_disk_size": "400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"master_instance_type": "hbase.sn1.large",
					"log_disk_size":        "400",
					"core_instance_type":   "hbase.sn1.large",
					"core_node_count":      "4",
					"log_instance_type":    "hbase.sn1.large",
					"log_node_count":       "4",
					"core_disk_size":       "400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_instance_type": "hbase.sn1.large",
						"log_disk_size":        "400",
						"core_instance_type":   "hbase.sn1.large",
						"core_node_count":      "4",
						"log_instance_type":    "hbase.sn1.large",
						"log_node_count":       "4",
						"core_disk_size":       "400",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "standby_core_node_count", "immediate_delete_flag", "security_ip_list", "auto_renew_period", "period_unit", "primary_core_node_count"},
			},
		},
	})
}

func TestUnitAccAlicloudHbaseMultiZoneCluster(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_hbase_multi_zone_cluster"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_hbase_multi_zone_cluster"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"arch_version":      "CreateHbaseMultiZoneClusterValue",
		"engine_version":    "CreateHbaseMultiZoneClusterValue",
		"resource_group_id": "CreateHbaseMultiZoneClusterValue",
		"tags": map[string]string{
			"Created": "TF",
		},
		"engine":                  "CreateHbaseMultiZoneClusterValue",
		"vpc_id":                  "CreateHbaseMultiZoneClusterValue",
		"core_disk_type":          "CreateHbaseMultiZoneClusterValue",
		"master_instance_type":    "CreateHbaseMultiZoneClusterValue",
		"log_disk_size":           400,
		"auto_renew_period":       2,
		"core_instance_type":      "CreateHbaseMultiZoneClusterValue",
		"cluster_name":            "CreateHbaseMultiZoneClusterValue",
		"multi_zone_combination":  "CreateHbaseMultiZoneClusterValue",
		"core_node_count":         4,
		"log_instance_type":       "CreateHbaseMultiZoneClusterValue",
		"log_node_count":          4,
		"log_disk_type":           "CreateHbaseMultiZoneClusterValue",
		"core_disk_size":          400,
		"primary_core_node_count": 6,
		"standby_core_node_count": 6,
		"payment_type":            "Subscription",
		"period":                  1,
		"period_unit":             "month",
		"primary_zone_id":         "CreateHbaseMultiZoneClusterValue",
		"primary_vswitch_id":      "CreateHbaseMultiZoneClusterValue",
		"standby_zone_id":         "CreateHbaseMultiZoneClusterValue",
		"standby_vswitch_id":      "CreateHbaseMultiZoneClusterValue",
		"arbiter_zone_id":         "CreateHbaseMultiZoneClusterValue",
		"arbiter_vswitch_id":      "CreateHbaseMultiZoneClusterValue",
		"security_ip_list":        "CreateHbaseMultiZoneClusterValue",
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
		"Data": map[string]interface{}{
			"Name": "CreateHbaseMultiZoneClusterValue",
		},
		"Success": true,

		"StandbyZoneId":     "CreateHbaseMultiZoneClusterValue",
		"ResourceGroupId":   "CreateHbaseMultiZoneClusterValue",
		"ModuleId":          0,
		"StandbyVSwitchIds": "CreateHbaseMultiZoneClusterValue",
		"Engine":            "CreateHbaseMultiZoneClusterValue",
		"Tags": map[string]interface{}{
			"Tag": []interface{}{
				map[string]interface{}{
					"Created": "TF",
				},
			},
		},
		"Status":               "ACTIVATION",
		"EncryptionType":       "CreateHbaseMultiZoneClusterValue",
		"InstanceId":           "CreateHbaseMultiZoneClusterValue",
		"PayType":              "Prepaid",
		"InstanceName":         "CreateHbaseMultiZoneClusterValue",
		"VpcId":                "CreateHbaseMultiZoneClusterValue",
		"CoreDiskType":         "CreateHbaseMultiZoneClusterValue",
		"EncryptionKey":        "CreateHbaseMultiZoneClusterValue",
		"MasterInstanceType":   "CreateHbaseMultiZoneClusterValue",
		"PrimaryVSwitchIds":    "CreateHbaseMultiZoneClusterValue",
		"IsDeletionProtection": false,
		"LogDiskCount":         4,
		"LogDiskSize":          400,
		"ArbiterVSwitchIds":    "CreateHbaseMultiZoneClusterValue",
		"MaintainEndTime":      "CreateHbaseMultiZoneClusterValue",
		"NetworkType":          "CreateHbaseMultiZoneClusterValue",
		"CoreInstanceType":     "CreateHbaseMultiZoneClusterValue",
		"ClusterName":          "CreateHbaseMultiZoneClusterValue",
		"MaintainStartTime":    "CreateHbaseMultiZoneClusterValue",
		"ArbiterZoneId":        "CreateHbaseMultiZoneClusterValue",
		"MajorVersion":         "CreateHbaseMultiZoneClusterValue",
		"CoreDiskCount":        4,
		"PrimaryZoneId":        "CreateHbaseMultiZoneClusterValue",
		"MultiZoneCombination": "CreateHbaseMultiZoneClusterValue",
		"ClusterId":            "CreateHbaseMultiZoneClusterValue",
		"CoreNodeCount":        4,
		"CreatedTimeUTC":       "CreateHbaseMultiZoneClusterValue",
		"LogInstanceType":      "CreateHbaseMultiZoneClusterValue",
		"LogNodeCount":         4,
		"LogDiskType":          "CreateHbaseMultiZoneClusterValue",
		"RegionId":             "CreateHbaseMultiZoneClusterValue",
		"CoreDiskSize":         400,
	}
	CreateMockResponse := map[string]interface{}{
		"ClusterId": "CreateHbaseMultiZoneClusterValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_hbase_multi_zone_cluster", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHbaseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudHbaseMultiZoneClusterCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateMultiZoneCluster" {
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
		err := resourceAlicloudHbaseMultiZoneClusterCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbase_multi_zone_cluster"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHbaseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudHbaseMultiZoneClusterUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)

	// Tags
	attributesDiff := map[string]interface{}{
		"tags": map[string]string{
			"Created1": "TF",
		},
	}
	diff, err := newInstanceDiff("alicloud_hbase_multi_zone_cluster", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbase_multi_zone_cluster"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Tags": map[string]interface{}{
			"Tag": []interface{}{
				map[string]interface{}{
					"Created1": "TF",
				}},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "TagResources" {
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
			if *action == "UnTagResources" {
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
		err := resourceAlicloudHbaseMultiZoneClusterUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbase_multi_zone_cluster"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ResizeMultiZoneClusterDiskSize
	attributesDiff = map[string]interface{}{
		"core_disk_size": 500,
		"log_disk_size":  500,
	}
	diff, err = newInstanceDiff("alicloud_hbase_multi_zone_cluster", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbase_multi_zone_cluster"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"LogDiskSize":  500,
		"CoreDiskSize": 500,
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ResizeMultiZoneClusterDiskSize" {
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
		err := resourceAlicloudHbaseMultiZoneClusterUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbase_multi_zone_cluster"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifyMultiZoneClusterNodeType
	attributesDiff = map[string]interface{}{
		"core_instance_type":   "UpdateHbaseMultiZoneClusterValue",
		"log_instance_type":    "UpdateHbaseMultiZoneClusterValue",
		"master_instance_type": "UpdateHbaseMultiZoneClusterValue",
	}
	diff, err = newInstanceDiff("alicloud_hbase_multi_zone_cluster", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbase_multi_zone_cluster"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"CoreInstanceType":   "UpdateHbaseMultiZoneClusterValue",
		"LogInstanceType":    "UpdateHbaseMultiZoneClusterValue",
		"MasterInstanceType": "UpdateHbaseMultiZoneClusterValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyMultiZoneClusterNodeType" {
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
		err := resourceAlicloudHbaseMultiZoneClusterUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbase_multi_zone_cluster"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ResizeMultiZoneClusterNodeCount
	attributesDiff = map[string]interface{}{
		"core_node_count": 5,
		"log_node_count":  5,
	}
	diff, err = newInstanceDiff("alicloud_hbase_multi_zone_cluster", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbase_multi_zone_cluster"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"CoreNodeCount": 5,
		"LogNodeCount":  5,
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ResizeMultiZoneClusterNodeCount" {
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
		err := resourceAlicloudHbaseMultiZoneClusterUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbase_multi_zone_cluster"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_hbase_multi_zone_cluster", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbase_multi_zone_cluster"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeMultiZoneCluster" {
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
		err := resourceAlicloudHbaseMultiZoneClusterRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHbaseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudHbaseMultiZoneClusterDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_hbase_multi_zone_cluster", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbase_multi_zone_cluster"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteMultiZoneCluster" {
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
		err := resourceAlicloudHbaseMultiZoneClusterDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
