package vod_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vod"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccAlicloudVODEditingProject_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vod_editing_project.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.VODSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVODEditingProjectMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return vod.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEditingProject")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svodeditingproject%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVODEditingProjectBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"editing_project_name": "${var.name}",
					"timeline":             `{\"VideoTracks\":[{\"VideoTrackClips\":[{\"MediaId\":\"0c60e6f02dae71edbfaa472190a90102\",\"In\":2811}]}]}`,
					"title":                "${var.name}",
					"cover_url":            "https://demo.aliyundoc.com/6AB4D0E1E1C74468883516C2349D1FC2-6-2.png",
					"division":             "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"editing_project_name": name,
						"timeline":             "{\"VideoTracks\":[{\"VideoTrackClips\":[{\"MediaId\":\"0c60e6f02dae71edbfaa472190a90102\",\"In\":2811}]}]}",
						"title":                name,
						"cover_url":            "https://demo.aliyundoc.com/6AB4D0E1E1C74468883516C2349D1FC2-6-2.png",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"division"},
			},
		},
	})
}

var AlicloudVODEditingProjectMap0 = map[string]string{
	"editing_project_name": acctest.CHECKSET,
	"status":               acctest.CHECKSET,
	"title":                acctest.CHECKSET,
	"timeline":             acctest.CHECKSET,
}

func AlicloudVODEditingProjectBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
func TestAccAlicloudVODEditingProject_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vod_editing_project.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.VODSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVODEditingProjectMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return vod.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEditingProject")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svodeditingproject%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVODEditingProjectBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"title": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"title": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"editing_project_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"editing_project_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cover_url": "https://demo.aliyundoc.com/6AB4D0E1E1C74468883516C2349D1FC2-6-1.png",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cover_url": "https://demo.aliyundoc.com/6AB4D0E1E1C74468883516C2349D1FC2-6-1.png",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeline": `{\"VideoTracks\":[{\"VideoTrackClips\":[{\"MediaId\":\"0c60e6f02dae71edbfaa472190a90102\",\"In\":2811}]}]}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeline": "{\"VideoTracks\":[{\"VideoTrackClips\":[{\"MediaId\":\"0c60e6f02dae71edbfaa472190a90102\",\"In\":2811}]}]}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"editing_project_name": "${var.name}_update",
					"timeline":             `{\"VideoTracks\":[{\"VideoTrackClips\":[{\"MediaId\":\"0c60e6f02dae71edbfaa472190a90102\",\"In\":2812}]}]}`,
					"title":                "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"editing_project_name": name + "_update",
						"timeline":             "{\"VideoTracks\":[{\"VideoTrackClips\":[{\"MediaId\":\"0c60e6f02dae71edbfaa472190a90102\",\"In\":2812}]}]}",
						"title":                name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"division"},
			},
		},
	})
}

var AlicloudVODEditingProjectMap1 = map[string]string{
	"status": acctest.CHECKSET,
	"title":  acctest.CHECKSET,
}

func AlicloudVODEditingProjectBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

// lintignore: R001
func TestUnitAccAlicloudVodEditingProject(t *testing.T) {
	resourceSchema := vod.ResourceEditingProject()
	dInit, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"editing_project_name": "CreateVodEditingProjectValue",
		"timeline":             "CreateVodEditingProjectValue",
		"title":                "CreateVodEditingProjectValue",
		"cover_url":            "CreateVodEditingProjectValue",
		"division":             "CreateVodEditingProjectValue",
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
	rawClient, err := vodTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	ReadMockResponse := map[string]interface{}{
		"Project": map[string]interface{}{
			"Status":      "Normal",
			"Timeline":    "CreateVodEditingProjectValue",
			"Description": "CreateVodEditingProjectValue",
			"ProjectId":   "CreateVodEditingProjectValue",
			"Title":       "CreateVodEditingProjectValue",
			"CoverURL":    "CreateVodEditingProjectValue",
			"RegionId":    "CreateVodEditingProjectValue",
		},
	}
	CreateMockResponse := map[string]interface{}{
		"Project": map[string]interface{}{
			"ProjectId": "CreateVodEditingProjectValue",
		},
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String(errorCode),
			Data:       tea.String(errorCode),
			Message:    tea.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_vod_editing_project", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceSchema.Create(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "AddEditingProject" {
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
		err := resourceSchema.Create(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceSchema.Update(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"editing_project_name": "UpdateVodEditingProjectValue",
		"timeline":             "UpdateVodEditingProjectValue",
		"title":                "UpdateVodEditingProjectValue",
		"cover_url":            "UpdateVodEditingProjectValue",
	}
	diff := buildVodEditingProjectDiff(attributes, attributesDiff)
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Project": map[string]interface{}{
			"Timeline":    "UpdateVodEditingProjectValue",
			"Description": "UpdateVodEditingProjectValue",
			"Title":       "UpdateVodEditingProjectValue",
			"CoverURL":    "UpdateVodEditingProjectValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateEditingProject" {
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
		err := resourceSchema.Update(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(resourceSchema.Schema).Data(dExisted.State(), nil)
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
	diff = buildVodEditingProjectDiff(attributes, attributesDiff)
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetEditingProject" {
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
		err := resourceSchema.Read(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceSchema.Delete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff = buildVodEditingProjectDiff(attributes, attributesDiff)
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteEditingProject" {
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
		err := resourceSchema.Delete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

func buildVodEditingProjectDiff(attributes, attributesDiff map[string]interface{}) *sdkterraform.InstanceDiff {
	diff := sdkterraform.NewInstanceDiff()
	for key, newVal := range attributesDiff {
		diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{
			Old: attributes[key].(string),
			New: newVal.(string),
		}
	}
	return diff
}

func mapMerge(dst map[string]interface{}, src map[string]interface{}) {
	for k, v := range src {
		dst[k] = v
	}
}

var vodTestEndpoints sync.Map
var vodTestSignVersion sync.Map

func vodTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	var accessKey, secretKey string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}
	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}
	securityToken := os.Getenv("ALICLOUD_SECURITY_TOKEN")
	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &vodTestEndpoints,
		SignVersion: &vodTestSignVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential
	return conf.Client()
}
