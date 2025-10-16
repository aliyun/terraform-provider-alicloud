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

func TestAccAliCloudResourceManagerResourceShare_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := resourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcesharingService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccResourceManagerResourceShare%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudResourceManagerResourceShareBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name": name + "update",
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

var AliCloudResourceManagerResourceShareMap0 = map[string]string{
	"resource_share_owner": CHECKSET,
	"status":               CHECKSET,
}

func AliCloudResourceManagerResourceShareBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
`, name)
}

func TestUnitAliCloudResourceManagerResourceShare(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_resource_manager_resource_share"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_resource_manager_resource_share"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"resource_share_name": "CreateResourceShareValue",
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
		// ListResourceShares
		"ResourceShares": []interface{}{
			map[string]interface{}{
				"ResourceShareName":   "CreateResourceShareValue",
				"ResourceShareOwner":  "CreateResourceShareValue",
				"ResourceShareStatus": "CreateResourceShareValue",
				"ResourceShareId":     "CreateResourceShareValue",
			},
		},
		"ResourceShare": map[string]interface{}{
			"ResourceShareId": "CreateResourceShareValue",
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateResourceShare
		"ResourceShare": map[string]interface{}{
			"ResourceShareId": "CreateResourceShareValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_resource_manager_resource_share", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewRessharingClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudResourceManagerResourceShareCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// ListResourceShares Response
		"ResourceShare": map[string]interface{}{
			"ResourceShareId": "CreateResourceShareValue",
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateResourceShare" {
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
		err := resourceAliCloudResourceManagerResourceShareCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_resource_manager_resource_share"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewRessharingClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudResourceManagerResourceShareUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateResourceShare
	attributesDiff := map[string]interface{}{
		"resource_share_name": "UpdateResourceShareValue",
	}
	diff, err := newInstanceDiff("alicloud_resource_manager_resource_share", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_resource_manager_resource_share"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// ListResourceShares Response
		"ResourceShares": []interface{}{
			map[string]interface{}{
				"ResourceShareName": "UpdateResourceShareValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateResourceShare" {
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
		err := resourceAliCloudResourceManagerResourceShareUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_resource_manager_resource_share"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
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
			if *action == "ListResourceShares" {
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
		err := resourceAliCloudResourceManagerResourceShareRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewRessharingClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudResourceManagerResourceShareDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_resource_manager_resource_share", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_resource_manager_resource_share"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteResourceShare" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"ResourceShares": []interface{}{
								map[string]interface{}{
									"ResourceShareId":     "CreateResourceShareValue",
									"ResourceShareStatus": "Deleted",
								},
							},
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudResourceManagerResourceShareDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}

// Test ResourceManager ResourceShare. >>> Resource test cases, automatically generated.
// Case 全生命周期带资源组-依赖资源版本_20250926 11575
func TestAccAliCloudResourceManagerResourceShare_basic11575(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := resourceAttrInit(resourceId, AlicloudResourceManagerResourceShareMap11575)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerResourceShareBasicDependence11575)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name": name,
					"resources": []map[string]interface{}{
						{
							"resource_type": "VSwitch",
							"resource_id":   "${alicloud_vswitch.vs.id}",
						},
					},
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"allow_external_targets": "false",
					"targets": []string{
						"${alicloud_resource_manager_account.defaultEYqcDi.id}"},
					"permission_names": []string{
						"AliyunRSDefaultPermissionVSwitch"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name":    name,
						"resources.#":            "1",
						"resource_group_id":      CHECKSET,
						"allow_external_targets": "false",
						"targets.#":              "1",
						"permission_names.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name":    name + "_update",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"allow_external_targets": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name":    name + "_update",
						"resource_group_id":      CHECKSET,
						"allow_external_targets": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"permission_names", "resources", "targets"},
			},
		},
	})
}

var AlicloudResourceManagerResourceShareMap11575 = map[string]string{
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"resource_share_owner": CHECKSET,
}

func AlicloudResourceManagerResourceShareBasicDependence11575(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "vpc" {
  vpc_name   = "tf01"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "vs" {
  vpc_id     = alicloud_vpc.vpc.id
  zone_id    = "cn-shanghai-f"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_resource_manager_account" "defaultEYqcDi" {
  display_name = var.name
  force_delete = true
}


`, name)
}

// Case 全生命周期带资源组-打包 11306
func TestAccAliCloudResourceManagerResourceShare_basic11306(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := resourceAttrInit(resourceId, AlicloudResourceManagerResourceShareMap11306)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerResourceShareBasicDependence11306)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name": name,
					"resources": []map[string]interface{}{
						{
							"resource_type": "VSwitch",
							"resource_id":   "${alicloud_vswitch.vs.id}",
						},
					},
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"allow_external_targets": "false",
					"targets": []string{
						"${alicloud_resource_manager_account.defaultEYqcDi.id}"},
					"permission_names": []string{
						"AliyunRSDefaultPermissionVSwitch"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name":    name,
						"resources.#":            "1",
						"resource_group_id":      CHECKSET,
						"allow_external_targets": "false",
						"targets.#":              "1",
						"permission_names.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name":    name + "_update",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"allow_external_targets": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name":    name + "_update",
						"resource_group_id":      CHECKSET,
						"allow_external_targets": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"permission_names", "resources", "targets"},
			},
		},
	})
}

var AlicloudResourceManagerResourceShareMap11306 = map[string]string{
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"resource_share_owner": CHECKSET,
}

func AlicloudResourceManagerResourceShareBasicDependence11306(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "vpc" {
  vpc_name   = "tf01"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "vs" {
  vpc_id     = alicloud_vpc.vpc.id
  zone_id    = "cn-shanghai-f"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_resource_manager_account" "defaultEYqcDi" {
  display_name = var.name
  force_delete = true
}


`, name)
}

// Case 全生命周期带资源组-依赖资源版本_20250821 11273
func TestAccAliCloudResourceManagerResourceShare_basic11273(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := resourceAttrInit(resourceId, AlicloudResourceManagerResourceShareMap11273)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerResourceShareBasicDependence11273)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name": name,
					"resources": []map[string]interface{}{
						{
							"resource_type": "VSwitch",
							"resource_id":   "${alicloud_vswitch.vs.id}",
						},
					},
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"allow_external_targets": "false",
					"targets": []string{
						"${alicloud_resource_manager_account.defaultEYqcDi.id}"},
					"permission_names": []string{
						"AliyunRSDefaultPermissionVSwitch"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name":    name,
						"resources.#":            "1",
						"resource_group_id":      CHECKSET,
						"allow_external_targets": "false",
						"targets.#":              "1",
						"permission_names.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name":    name + "_update",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"allow_external_targets": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name":    name + "_update",
						"resource_group_id":      CHECKSET,
						"allow_external_targets": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"permission_names", "resources", "targets"},
			},
		},
	})
}

var AlicloudResourceManagerResourceShareMap11273 = map[string]string{
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"resource_share_owner": CHECKSET,
}

func AlicloudResourceManagerResourceShareBasicDependence11273(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "vpc" {
  vpc_name   = "tf01"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "vs" {
  vpc_id     = alicloud_vpc.vpc.id
  zone_id    = "cn-shanghai-f"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_resource_manager_account" "defaultEYqcDi" {
  display_name = var.name
  force_delete = true
}


`, name)
}

// Case 全生命周期带资源组-依赖资源版本_20250320 10573
func TestAccAliCloudResourceManagerResourceShare_basic10573(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := resourceAttrInit(resourceId, AlicloudResourceManagerResourceShareMap10573)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerResourceShareBasicDependence10573)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name": name,
					"resources": []map[string]interface{}{
						{
							"resource_type": "VSwitch",
							"resource_id":   "${alicloud_vswitch.vs.id}",
						},
					},
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"allow_external_targets": "false",
					"targets": []string{
						"${alicloud_resource_manager_account.defaultEYqcDi.id}"},
					"permission_names": []string{
						"AliyunRSDefaultPermissionVSwitch"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name":    name,
						"resources.#":            "1",
						"resource_group_id":      CHECKSET,
						"allow_external_targets": "false",
						"targets.#":              "1",
						"permission_names.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name":    name + "_update",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"allow_external_targets": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name":    name + "_update",
						"resource_group_id":      CHECKSET,
						"allow_external_targets": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"permission_names", "resources", "targets"},
			},
		},
	})
}

var AlicloudResourceManagerResourceShareMap10573 = map[string]string{
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"resource_share_owner": CHECKSET,
}

func AlicloudResourceManagerResourceShareBasicDependence10573(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "vpc" {
  vpc_name   = "tf01"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "vs" {
  vpc_id     = alicloud_vpc.vpc.id
  zone_id    = "cn-shanghai-f"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_resource_manager_account" "defaultEYqcDi" {
  display_name = var.name
  force_delete = true
}


`, name)
}

// Case 全生命周期带资源组-依赖资源版本_20241018 8293
func TestAccAliCloudResourceManagerResourceShare_basic8293(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := resourceAttrInit(resourceId, AlicloudResourceManagerResourceShareMap8293)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerResourceShareBasicDependence8293)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name": name,
					"resources": []map[string]interface{}{
						{
							"resource_type": "VSwitch",
							"resource_id":   "${alicloud_vswitch.vs.id}",
						},
					},
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"allow_external_targets": "true",
					"targets": []string{
						"1524526934643077"},
					"permission_names": []string{
						"AliyunRSDefaultPermissionVSwitch"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name":    name,
						"resources.#":            "1",
						"resource_group_id":      CHECKSET,
						"allow_external_targets": "true",
						"targets.#":              "1",
						"permission_names.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name": name + "_update",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name": name + "_update",
						"resource_group_id":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"permission_names", "resources", "targets"},
			},
		},
	})
}

var AlicloudResourceManagerResourceShareMap8293 = map[string]string{
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"resource_share_owner": CHECKSET,
}

func AlicloudResourceManagerResourceShareBasicDependence8293(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "vpc" {
  vpc_name   = "tf01"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "vs" {
  vpc_id     = alicloud_vpc.vpc.id
  zone_id    = "cn-shanghai-f"
  cidr_block = "192.168.1.0/24"
}


`, name)
}

// Test ResourceManager ResourceShare. <<< Resource test cases, automatically generated.
