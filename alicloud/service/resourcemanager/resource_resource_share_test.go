package resourcemanager_test

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

// collection-order stock-debt declarations (pure rename migration; see scripts/collection-order/README.md)
// collection-order: bypass permission_names -- pre-existing stock TypeList; this is a pure package-rename migration (alicloud/ -> alicloud/service/resourcemanager/), schema structurally unchanged vs merge-base 2f8f887c52e916d4e87f2e7ed34bdf2aba355259
// collection-order: bypass resource_arns -- pre-existing stock TypeList; this is a pure package-rename migration (alicloud/ -> alicloud/service/resourcemanager/), schema structurally unchanged vs merge-base 2f8f887c52e916d4e87f2e7ed34bdf2aba355259
// collection-order: bypass resource_properties -- pre-existing stock TypeList; this is a pure package-rename migration (alicloud/ -> alicloud/service/resourcemanager/), schema structurally unchanged vs merge-base 2f8f887c52e916d4e87f2e7ed34bdf2aba355259
// collection-order: bypass resources -- pre-existing stock TypeList; this is a pure package-rename migration (alicloud/ -> alicloud/service/resourcemanager/), schema structurally unchanged vs merge-base 2f8f887c52e916d4e87f2e7ed34bdf2aba355259
// collection-order: bypass targets -- pre-existing stock TypeList; this is a pure package-rename migration (alicloud/ -> alicloud/service/resourcemanager/), schema structurally unchanged vs merge-base 2f8f887c52e916d4e87f2e7ed34bdf2aba355259

func TestAccAliCloudResourceManagerResourceShare_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerResourceShare")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccResourceManagerResourceShare%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudResourceManagerResourceShareBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id": acctest.CHECKSET,
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
				ImportStateVerifyIgnore: []string{"permission_names", "resource_arns", "resource_properties", "resources", "targets"},
			},
		},
	})
}

func TestAccAliCloudResourceManagerResourceShare_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerResourceShare")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccResourceManagerResourceShare%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudResourceManagerResourceShareBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id":      acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerResourceShare")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccResourceManagerResourceShare%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudResourceManagerResourceShareBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id":      acctest.CHECKSET,
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
	"resource_share_owner": acctest.CHECKSET,
	"resource_group_id":    acctest.CHECKSET,
	"create_time":          acctest.CHECKSET,
	"status":               acctest.CHECKSET,
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
	p := acctest.Provider.ResourcesMap
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
	rawClient, err := sweep.SharedClientForRegion(region)
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
			Code:       helper.String(errorCode),
			Data:       helper.String(errorCode),
			Message:    helper.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_resource_manager_resource_share", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewRessharingClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourcemanager.ResourceResourceShareCreate(dInit, rawClient)
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
		err := resourcemanager.ResourceResourceShareCreate(dInit, rawClient)
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
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourcemanager.ResourceResourceShareUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateResourceShare
	attributesDiff := map[string]interface{}{
		"resource_share_name": "UpdateResourceShareValue",
	}
	diff, err := acctest.NewInstanceDiff("alicloud_resource_manager_resource_share", attributes, attributesDiff, dInit.State())
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
		err := resourcemanager.ResourceResourceShareUpdate(dExisted, rawClient)
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
		err := resourcemanager.ResourceResourceShareRead(dExisted, rawClient)
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
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourcemanager.ResourceResourceShareDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = acctest.NewInstanceDiff("alicloud_resource_manager_resource_share", attributes, attributesDiff, dInit.State())
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
		err := resourcemanager.ResourceResourceShareDelete(dExisted, rawClient)
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap12595)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerResourceShare")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudResourceManagerResourceShareBasicDependence12595)
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
						"resource_group_id": acctest.CHECKSET,
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
				ImportStateVerifyIgnore: []string{"permission_names", "resource_arns", "resource_properties", "resources", "targets"},
			},
		},
	})
}

func TestAccAliCloudResourceManagerResourceShare_basic12595_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap12595)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerResourceShare")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudResourceManagerResourceShareBasicDependence12595)
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
						"resource_group_id":      acctest.CHECKSET,
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
	"resource_share_owner": acctest.CHECKSET,
	"resource_group_id":    acctest.CHECKSET,
	"create_time":          acctest.CHECKSET,
	"status":               acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap12596)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerResourceShare")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudResourceManagerResourceShareBasicDependence12596)
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
						"resource_group_id": acctest.CHECKSET,
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
				ImportStateVerifyIgnore: []string{"permission_names", "resource_arns", "resource_properties", "resources", "targets"},
			},
		},
	})
}

func TestAccAliCloudResourceManagerResourceShare_basic12596_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_share.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudResourceManagerResourceShareMap12596)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerResourceShare")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudResourceManagerResourceShareBasicDependence12596)
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
						"resource_group_id":      acctest.CHECKSET,
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
	"resource_share_owner": acctest.CHECKSET,
	"resource_group_id":    acctest.CHECKSET,
	"create_time":          acctest.CHECKSET,
	"status":               acctest.CHECKSET,
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
