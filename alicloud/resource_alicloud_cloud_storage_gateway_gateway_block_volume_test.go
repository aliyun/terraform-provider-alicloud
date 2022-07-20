package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudStorageGatewayGatewayBlockVolume_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_storage_gateway_gateway_block_volume.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudStorageGatewayGatewayBlockVolumeMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SgwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudStorageGatewayGatewayBlockVolume")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftestacccsvolume%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudStorageGatewayGatewayBlockVolumeBasicDependence0)
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
					"gateway_id":                "${alicloud_cloud_storage_gateway_gateway.default.id}",
					"gateway_block_volume_name": name,
					"chunk_size":                "8192",
					"chap_enabled":              "false",
					"oss_endpoint":              "${alicloud_oss_bucket.default.extranet_endpoint}",
					"oss_bucket_name":           "${alicloud_oss_bucket.default.bucket}",
					"cache_mode":                "Cache",
					"local_path":                "${alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_file_path}",
					"protocol":                  "iSCSI",
					"oss_bucket_ssl":            "true",
					"size":                      "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_block_volume_name": name,
						"chunk_size":                "8192",
						"chap_enabled":              "false",
						"oss_endpoint":              CHECKSET,
						"oss_bucket_name":           CHECKSET,
						"cache_mode":                "Cache",
						"local_path":                CHECKSET,
						"protocol":                  "iSCSI",
						"oss_bucket_ssl":            "true",
						"gateway_id":                CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"chap_enabled":     "true",
					"chap_in_user":     "tftestAccnmSa123",
					"chap_in_password": "tftestAccnmSa456",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"chap_enabled": "true",
						"chap_in_user": "tftestAccnmSa123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"chap_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"chap_enabled": "false",
						"chap_in_user": "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_source_deletion", "recovery", "size", "chap_in_password"},
			},
		},
	})
}

var AlicloudCloudStorageGatewayGatewayBlockVolumeMap0 = map[string]string{
	"is_source_deletion": NOSET,
	"recovery":           NOSET,
}

func AlicloudCloudStorageGatewayGatewayBlockVolumeBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_cloud_storage_gateway_stocks" "default" {
  gateway_class = "Standard"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  description              = "tf-acctestDesalone"
  gateway_class            = "Standard"
  type                     = "Iscsi"
  payment_type             = "PayAsYouGo"
  vswitch_id               = local.vswitch_id
  release_after_expiration = true
  public_network_bandwidth = 10
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  location                 = "Cloud"
  gateway_name             = var.name
}


resource "alicloud_cloud_storage_gateway_gateway_cache_disk" "default" {
  cache_disk_category   = "cloud_efficiency"
  gateway_id            = alicloud_cloud_storage_gateway_gateway.default.id
  cache_disk_size_in_gb = 50
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
  acl    = "public-read-write"
}
`, name)
}

func TestUnitAlicloudCloudStorageGatewayGatewayBlockVolume(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cloud_storage_gateway_gateway_block_volume"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cloud_storage_gateway_gateway_block_volume"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"cache_mode":                "CreateGatewayBlockVolumeValue",
		"chap_enabled":              true,
		"chap_in_user":              "CreateGatewayBlockVolumeValue",
		"chap_in_password":          "CreateGatewayBlockVolumeValue",
		"chunk_size":                131072,
		"gateway_block_volume_name": "CreateGatewayBlockVolumeValue",
		"gateway_id":                "CreateGatewayBlockVolumeValue",
		"local_path":                "CreateGatewayBlockVolumeValue",
		"oss_bucket_name":           "CreateGatewayBlockVolumeValue",
		"oss_bucket_ssl":            true,
		"oss_endpoint":              "CreateGatewayBlockVolumeValue",
		"protocol":                  "CreateGatewayBlockVolumeValue",
		"recovery":                  true,
		"size":                      2,
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
		// DescribeGatewayBlockVolumes
		"BlockVolumes": map[string]interface{}{
			"BlockVolume": []interface{}{
				map[string]interface{}{
					"IndexId":       "CreateGatewayBlockVolumeValue",
					"GatewayId":     "CreateGatewayBlockVolumeValue",
					"CacheMode":     "CreateGatewayBlockVolumeValue",
					"ChapEnabled":   true,
					"ChapInUser":    "CreateGatewayBlockVolumeValue",
					"ChunkSize":     131072,
					"Name":          "CreateGatewayBlockVolumeValue",
					"LocalPath":     "CreateGatewayBlockVolumeValue",
					"OssBucketName": "CreateGatewayBlockVolumeValue",
					"OssBucketSsl":  true,
					"OssEndpoint":   "CreateGatewayBlockVolumeValue",
					"Protocol":      "CreateGatewayBlockVolumeValue",
					"VolumeState":   1,
				},
			},
		},
		"Tasks": map[string]interface{}{
			"SimpleTask": []interface{}{
				map[string]interface{}{
					"TaskId":            "CreateGatewayBlockVolumeValue",
					"StateCode":         "task.state.completed",
					"RelatedResourceId": "CreateGatewayBlockVolumeValue",
				},
			},
		},
		"TaskId":            "CreateGatewayBlockVolumeValue",
		"RelatedResourceId": "CreateGatewayBlockVolumeValue",
		"Success":           true,
	}
	CreateMockResponse := map[string]interface{}{
		// CreateGatewayBlockVolume
		"BlockVolumes": map[string]interface{}{
			"BlockVolume": []interface{}{
				map[string]interface{}{
					"IndexId": "CreateGatewayBlockVolumeValue",
				},
			},
		},
		"Tasks": map[string]interface{}{
			"SimpleTask": []interface{}{
				map[string]interface{}{
					"TaskId":            "CreateGatewayBlockVolumeValue",
					"StateCode":         "task.state.completed",
					"RelatedResourceId": "CreateGatewayBlockVolumeValue",
				},
			},
		},
		"TaskId":  "CreateGatewayBlockVolumeValue",
		"Success": true,
	}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cloud_storage_gateway_gateway_block_volume", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	t.Run("Create", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHcsSgwClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudCloudStorageGatewayGatewayBlockVolumeCreate(dInit, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		ReadMockResponseDiff = map[string]interface{}{
			// DescribeGatewayBlockVolumes Response
			"BlockVolumes": map[string]interface{}{
				"BlockVolume": []interface{}{
					map[string]interface{}{
						"IndexId": "CreateGatewayBlockVolumeValue",
					},
				},
			},
			"Success": true,
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1 // a counter used to cover retry scenario; the same below
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "CreateGatewayBlockVolume" {
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
			err := resourceAlicloudCloudStorageGatewayGatewayBlockVolumeCreate(dInit, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_cloud_storage_gateway_gateway_block_volume"].Schema).Data(dInit.State(), nil)
				for key, value := range attributes {
					_ = dCompare.Set(key, value)
				}
				assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
			}
			if retryIndex >= len(errorCodes)-1 {
				break
			}
		}
	})

	// Update
	t.Run("Update", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHcsSgwClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudCloudStorageGatewayGatewayBlockVolumeUpdate(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		//UpdateGatewayBlockVolume
		attributesDiff := map[string]interface{}{
			"chap_enabled":     false,
			"chap_in_password": "UpdateGatewayBlockVolumeValue",
			"chap_in_user":     "UpdateGatewayBlockVolumeValue",
			"size":             1,
		}
		diff, err := newInstanceDiff("alicloud_cloud_storage_gateway_gateway_block_volume", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_cloud_storage_gateway_gateway_block_volume"].Schema).Data(dInit.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			// DescribeGatewayBlockVolumes Response
			"BlockVolumes": map[string]interface{}{
				"BlockVolume": []interface{}{
					map[string]interface{}{
						"ChapEnabled":    false,
						"ChapInUser":     "UpdateGatewayBlockVolumeValue",
						"ChapInPassword": "UpdateGatewayBlockVolumeValue",
						"Size":           1,
					},
				},
			},
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "UpdateGatewayBlockVolume" {
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
			err := resourceAlicloudCloudStorageGatewayGatewayBlockVolumeUpdate(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_cloud_storage_gateway_gateway_block_volume"].Schema).Data(dExisted.State(), nil)
				for key, value := range attributes {
					_ = dCompare.Set(key, value)
				}
				assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
			}
			if retryIndex >= len(errorCodes)-1 {
				break
			}
		}
	})

	// Read
	t.Run("Read", func(t *testing.T) {
		attributesDiff := map[string]interface{}{}
		diff, err := newInstanceDiff("alicloud_cloud_storage_gateway_gateway_block_volume", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_cloud_storage_gateway_gateway_block_volume"].Schema).Data(dInit.State(), diff)
		errorCodes := []string{"NonRetryableError", "Throttling", "nil", "{}"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "DescribeGatewayBlockVolumes" {
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
			err := resourceAlicloudCloudStorageGatewayGatewayBlockVolumeRead(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "{}":
				assert.Nil(t, err)
			}
		}
	})

	// Delete
	t.Run("Delete", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHcsSgwClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudCloudStorageGatewayGatewayBlockVolumeDelete(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		attributesDiff := map[string]interface{}{}
		diff, err := newInstanceDiff("alicloud_cloud_storage_gateway_gateway_block_volume", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_cloud_storage_gateway_gateway_block_volume"].Schema).Data(dInit.State(), diff)
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "DeleteGatewayBlockVolumes" {
					switch errorCode {
					case "NonRetryableError":
						return failedResponseMock(errorCode)
					default:
						retryIndex++
						if errorCodes[retryIndex] == "nil" {
							ReadMockResponse = map[string]interface{}{
								"Tasks": map[string]interface{}{
									"SimpleTask": []interface{}{
										map[string]interface{}{
											"TaskId":    "CreateGatewayBlockVolumeValue",
											"StateCode": "task.state.completed",
										},
									},
								},
								"TaskId":  "CreateGatewayBlockVolumeValue",
								"Success": true,
							}
							return ReadMockResponse, nil
						}
						return failedResponseMock(errorCodes[retryIndex])
					}
				}
				return ReadMockResponse, nil
			})
			err := resourceAlicloudCloudStorageGatewayGatewayBlockVolumeDelete(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "nil":
				assert.Nil(t, err)
			}
		}
	})
}
