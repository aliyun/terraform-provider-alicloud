package alicloud

import (
	"fmt"
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
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAccAlicloudServiceMeshUserPermission_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_mesh_user_permission.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceMeshServiceMeshMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServicemeshService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeUserPermissions")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sservicemesh%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceMeshUserPermissionBasicDependence0)
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
					"sub_account_user_id": "${alicloud_ram_user.default.id}",
					"permissions": []map[string]interface{}{
						{
							"role_name":       "istio-ops",
							"service_mesh_id": "${alicloud_service_mesh_service_mesh.default1.id}",
							"role_type":       "custom",
							"is_custom":       "true",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permissions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"permissions": []map[string]interface{}{
						{
							"role_name":       "istio-admin",
							"service_mesh_id": "${alicloud_service_mesh_service_mesh.default1.id}",
							"role_type":       "custom",
							"is_custom":       "true",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permissions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"permissions": []map[string]interface{}{
						{
							"role_name":       "istio-readonly",
							"service_mesh_id": "${alicloud_service_mesh_service_mesh.default1.id}",
							"role_type":       "custom",
							"is_custom":       "true",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permissions.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudServiceMeshUserPermissionMap0 = map[string]string{}

func AlicloudServiceMeshUserPermissionBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
variable "service_mesh_name"{
  default = "%s"
}
data "alicloud_service_mesh_versions" "default" {
	edition = "Default"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id     	= data.alicloud_zones.default.zones.0.id
}

resource "alicloud_ram_user" "default" {
  name = var.name
}

resource "alicloud_service_mesh_service_mesh" "default1" {
	service_mesh_name = var.name
    edition = "Default"
    cluster_spec = "standard"
    version = data.alicloud_service_mesh_versions.default.versions.0.version
	network {
		vpc_id = data.alicloud_vpcs.default.ids.0
		vswitche_list = [data.alicloud_vswitches.default.ids.0]
	}
	load_balancer {
		pilot_public_eip = false
		api_server_public_eip = false
	}
}



`, name, fmt.Sprintf("%s%d", name, acctest.RandIntRange(1, 999)))
}

func TestUnitAlicloudServiceMeshUserPermission(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_service_mesh_user_permission"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_service_mesh_user_permission"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"sub_account_user_id": "sub_account_user_id",
		"permissions": []map[string]interface{}{
			{
				"role_name":       "istio-ops",
				"service_mesh_id": "service_mesh_id",
				"role_type":       "custom",
				"is_custom":       true,
				"is_ram_role":     false,
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
		"Permissions": []interface{}{
			map[string]interface{}{
				"ParentId":     "ParentId",
				"ResourceType": "cluster",
				"RoleType":     "custom",
				"RoleName":     "istio-ops",
				"ResourceId":   "service_mesh_id",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	t.Run("Create", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewServicemeshClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudServiceMeshUserPermissionCreate(dInit, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		ReadMockResponseDiff = map[string]interface{}{}
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1 // a counter used to cover retry scenario; the same below
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "GrantUserPermissions" {
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
			err := resourceAlicloudServiceMeshUserPermissionCreate(dInit, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_service_mesh_user_permission"].Schema).Data(dInit.State(), nil)
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewServicemeshClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudServiceMeshUserPermissionUpdate(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		attributesDiff := map[string]interface{}{
			"permissions": []map[string]interface{}{
				{
					"role_name":       "istio-admin",
					"service_mesh_id": "service_mesh_id",
					"role_type":       "custom",
					"is_custom":       true,
					"is_ram_role":     false,
				},
			},
		}
		diff, err := newInstanceDiff("alicloud_service_mesh_user_permission", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_service_mesh_user_permission"].Schema).Data(dInit.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			"Permissions": []interface{}{
				map[string]interface{}{
					"ParentId":     "ParentId",
					"ResourceType": "cluster",
					"RoleType":     "custom",
					"RoleName":     "istio-admin",
					"ResourceId":   "service_mesh_id",
				},
			},
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "GrantUserPermissions" {
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
			err := resourceAlicloudServiceMeshUserPermissionUpdate(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_service_mesh_user_permission"].Schema).Data(dExisted.State(), nil)
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
}
