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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudResourceManagerResourceShare_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := resourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
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
		CheckDestroy:  rac.checkResourceDestroy(),
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
				Config: testAccConfig(map[string]interface{}{
					"allow_external_targets": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_external_targets": "true",
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
				ImportStateVerifyIgnore: []string{"permission_names", "resource_arns", "resource_properties", "resources", "targets"},
			},
		},
	})
}

func TestAccAliCloudResourceManagerResourceShare_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := resourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccResourceManagerResourceShare%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudResourceManagerResourceShareBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name":    name,
					"allow_external_targets": "true",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"permission_names":       []string{"AliyunRSDefaultPermissionVSwitch"},
					"resource_arns":          []string{"acs:vpc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:vswitch/${data.alicloud_vswitches.default.ids.0}"},
					"targets":                []string{"${data.alicloud_resource_manager_accounts.default.ids.1}"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name":    name,
						"allow_external_targets": "true",
						"resource_group_id":      CHECKSET,
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"permission_names", "resource_arns", "resource_properties", "resources", "targets"},
			},
		},
	})
}

func TestAccAliCloudResourceManagerResourceShare_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := resourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccResourceManagerResourceShare%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudResourceManagerResourceShareBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name":    name,
					"allow_external_targets": "true",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"permission_names":       []string{"AliyunRSDefaultPermissionServiceCatalogPortfolio"},
					"targets":                []string{"${data.alicloud_resource_manager_accounts.default.ids.1}"},
					"resource_properties": []map[string]interface{}{
						{
							"property":     `{\"sharePrincipals\":false,\"shareTagOptions\":false}`,
							"resource_arn": "${alicloud_service_catalog_portfolio.default.portfolio_arn}",
						},
					},
					"resources": []map[string]interface{}{
						{
							"resource_type": "ServiceCatalogPortfolio",
							"resource_id":   "${alicloud_service_catalog_portfolio.default.id}",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name":    name,
						"allow_external_targets": "true",
						"resource_group_id":      CHECKSET,
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"permission_names", "resource_arns", "resource_properties", "resources", "targets"},
			},
		},
	})
}

var AliCloudResourceManagerResourceShareMap0 = map[string]string{
	"resource_share_owner": CHECKSET,
	"resource_group_id":    CHECKSET,
	"create_time":          CHECKSET,
	"status":               CHECKSET,
}

func AliCloudResourceManagerResourceShareBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_account" "default" {
}

data "alicloud_vswitches" "default" {
}

data "alicloud_resource_manager_accounts" "default" {
  status  = "CreateSuccess"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_service_catalog_portfolio" "default" {
  portfolio_name = var.name
  provider_name  = var.name
}
`, name)
}

// lintignore: R001
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

// Case 全生命周期带资源组-依赖资源版本_20260311_resourceId 12595
func TestAccAliCloudResourceManagerResourceShare_basic12595(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := resourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap12595)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudResourceManagerResourceShareBasicDependence12595)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
				Config: testAccConfig(map[string]interface{}{
					"allow_external_targets": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_external_targets": "true",
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
				ImportStateVerifyIgnore: []string{"permission_names", "resource_arns", "resource_properties", "resources", "targets"},
			},
		},
	})
}

func TestAccAliCloudResourceManagerResourceShare_basic12595_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := resourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap12595)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudResourceManagerResourceShareBasicDependence12595)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name":    name,
					"allow_external_targets": "true",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"permission_names":       []string{"AliyunRSDefaultPermissionVSwitch"},
					"resource_arns":          []string{"acs:vpc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:vswitch/${data.alicloud_vswitches.default.ids.0}"},
					"targets":                []string{"${data.alicloud_resource_manager_accounts.default.ids.1}"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name":    name,
						"allow_external_targets": "true",
						"resource_group_id":      CHECKSET,
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"permission_names", "resource_arns", "resource_properties", "resources", "targets"},
			},
		},
	})
}

var AliCloudResourceManagerResourceShareMap12595 = map[string]string{
	"resource_share_owner": CHECKSET,
	"resource_group_id":    CHECKSET,
	"create_time":          CHECKSET,
	"status":               CHECKSET,
}

func AliCloudResourceManagerResourceShareBasicDependence12595(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_account" "default" {
}

data "alicloud_vswitches" "default" {
}

data "alicloud_resource_manager_accounts" "default" {
  status  = "CreateSuccess"
}

data "alicloud_resource_manager_resource_groups" "default" {
}
`, name)
}

// Case 全生命周期带资源组-依赖资源版本_20260311_resourceArn 12596
func TestAccAliCloudResourceManagerResourceShare_basic12596(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := resourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap12596)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudResourceManagerResourceShareBasicDependence12596)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
				Config: testAccConfig(map[string]interface{}{
					"allow_external_targets": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_external_targets": "true",
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
				ImportStateVerifyIgnore: []string{"permission_names", "resource_arns", "resource_properties", "resources", "targets"},
			},
		},
	})
}

func TestAccAliCloudResourceManagerResourceShare_basic12596_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := resourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap12596)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceShare")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudResourceManagerResourceShareBasicDependence12596)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_name":    name,
					"allow_external_targets": "true",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"permission_names":       []string{"AliyunRSDefaultPermissionServiceCatalogPortfolio"},
					"targets":                []string{"${data.alicloud_resource_manager_accounts.default.ids.1}"},
					"resource_properties": []map[string]interface{}{
						{
							"property":     `{\"sharePrincipals\":false,\"shareTagOptions\":false}`,
							"resource_arn": "${alicloud_service_catalog_portfolio.default.portfolio_arn}",
						},
					},
					"resources": []map[string]interface{}{
						{
							"resource_type": "ServiceCatalogPortfolio",
							"resource_id":   "${alicloud_service_catalog_portfolio.default.id}",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_name":    name,
						"allow_external_targets": "true",
						"resource_group_id":      CHECKSET,
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"permission_names", "resource_arns", "resource_properties", "resources", "targets"},
			},
		},
	})
}

var AliCloudResourceManagerResourceShareMap12596 = map[string]string{
	"resource_share_owner": CHECKSET,
	"resource_group_id":    CHECKSET,
	"create_time":          CHECKSET,
	"status":               CHECKSET,
}

func AliCloudResourceManagerResourceShareBasicDependence12596(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

data "alicloud_resource_manager_accounts" "default" {
  status  = "CreateSuccess"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_service_catalog_portfolio" "default" {
  portfolio_name = var.name
  provider_name  = var.name
}
`, name)
}

// Test ResourceManager ResourceShare. <<< Resource test cases, automatically generated.
