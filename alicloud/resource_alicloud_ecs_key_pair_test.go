package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECSKeyPair_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_key_pair.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsKeyPairMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsKeyPair")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsKeyPair%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsKeyPairBasicDependence)
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
					"key_name":          name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_name":          name,
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "Terraform",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "Terraform",
						"tags.For":     "Test",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudECSKeyPair_multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_key_pair.default.9"
	ra := resourceAttrInit(resourceId, AlicloudEcsKeyPairMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsKeyPair")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsKeyPair%d", rand)
	//testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsKeyPairBasicDependence)
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
					"count": "10",
				}),
				Check: nil,
			},
		},
	})
}

func TestAccAlicloudECSKeyPair_name_prefix(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_key_pair.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsKeyPairMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsKeyPair")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsKeyPair%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsKeyPairBasicDependence)
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
					"key_name_prefix": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_name_prefix": name,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudECSKeyPair_public_key(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_key_pair.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsKeyPairMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsKeyPair")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsKeyPair%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsKeyPairBasicDependence)
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
					"key_name":   name,
					"public_key": "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_name": name,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudECSKeyPair_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_key_pair.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsKeyPairMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsKeyPair")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsKeyPair%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsKeyPairBasicDependence)
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
					"key_name":          name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"Created": "Terraform",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_name":          name,
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "Terraform",
						"tags.For":          "Test",
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

var AlicloudEcsKeyPairMap = map[string]string{}

func AlicloudEcsKeyPairBasicDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_resource_manager_resource_groups" "default"{
		status = "OK"
	}
`)
}

func TestUnitAlicloudECSKeyPair(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_ecs_key_pair"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_ecs_key_pair"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	dCreateCompletion, _ := schema.InternalMap(p["alicloud_ecs_key_pair"].Schema).Data(nil, nil)
	dCreateCompletion.MarkNewResource()
	dCreateKeyName, _ := schema.InternalMap(p["alicloud_ecs_key_pair"].Schema).Data(nil, nil)
	dCreateKeyName.MarkNewResource()
	dCreateKeyNamePrefix, _ := schema.InternalMap(p["alicloud_ecs_key_pair"].Schema).Data(nil, nil)
	dCreateKeyNamePrefix.MarkNewResource()
	for key, value := range map[string]interface{}{
		"resource_group_id": "resource_group_id",
		"key_pair_name":     "key_pair_name",
		"public_key":        "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg",
		"tags": map[string]string{
			"Created": "Terraform",
			"For":     "Test",
		},
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	for keyCompletion, valueCompletion := range map[string]interface{}{
		"key_name": "key_name",
	} {
		err := dCreateCompletion.Set(keyCompletion, valueCompletion)
		assert.Nil(t, err)
		err = d.Set(keyCompletion, valueCompletion)
		assert.Nil(t, err)

	}
	for keyCompletion, valueCompletion := range map[string]interface{}{
		"key_name_prefix": "key_name_prefix",
	} {
		err := dCreateKeyNamePrefix.Set(keyCompletion, valueCompletion)
		assert.Nil(t, err)
		err = d.Set(keyCompletion, valueCompletion)
		assert.Nil(t, err)

	}
	for keyCompletion, valueCompletion := range map[string]interface{}{
		"resource_group_id": "resource_group_id",
		"key_file":          "key_file",
	} {
		err := dCreateKeyName.Set(keyCompletion, valueCompletion)
		assert.Nil(t, err)
		err = d.Set(keyCompletion, valueCompletion)
		assert.Nil(t, err)

	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"KeyPairs": map[string]interface{}{
			"KeyPair": []interface{}{
				map[string]interface{}{
					"KeyPairName":        "MockKeyPairName",
					"KeyPairFingerPrint": "finger_print",
					"ResourceGroupId":    "resource_group_id",
					"Tags": map[string]string{
						"Created": "Terraform",
						"For":     "Test",
					},
				},
			},
		},
		"TagResources": map[string]interface{}{
			"TagResource": []interface{}{
				map[string]interface{}{
					"Tags": map[string]string{
						"Created": "Terraform",
					},
				},
			},
		},
		"PrivateKeyBody": "private_key_body",
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["KeyPairName"] = "MockKeyPairName"
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadDescribeEcsKeyPairNotFound": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"KeyPairs": map[string]interface{}{
					"KeyPair": []interface{}{},
				},
			}
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudEcsKeyPairCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudEcsKeyPairCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudEcsKeyPairCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("CreateCompletionNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudEcsKeyPairCreate(dCreateCompletion, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("CreateKeyNamePrefixNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudEcsKeyPairCreate(dCreateKeyNamePrefix, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("CreateKeyNameNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudEcsKeyPairCreate(dCreateKeyName, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("MockKeyPairName")
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAliCloudEcsKeyPairUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateJoinResourceGroupAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"resource_group_id"} {
			switch p["alicloud_ecs_key_pair"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ecs_key_pair"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudEcsKeyPairUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateJoinResourceGroupNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"resource_group_id", "tags"} {
			switch p["alicloud_ecs_key_pair"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ecs_key_pair"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		patcheSetResourceTags := gomonkey.ApplyMethod(reflect.TypeOf(&VpcService{}), "SetResourceTags", func(*VpcService, *schema.ResourceData, string) error {
			_, err := responseMock["NoRetryError"]("NoRetryError")
			return err
		})
		err := resourceAliCloudEcsKeyPairUpdate(resourceData1, rawClient)
		patches.Reset()
		patcheSetResourceTags.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudEcsKeyPairDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudEcsKeyPairDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudEcsKeyPairDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("DeleteNonRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = true
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudEcsKeyPairDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeEcsKeyPairNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["ReadDescribeEcsKeyPairNotFound"]("")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudEcsKeyPairRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeEcsKeyPairAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudEcsKeyPairRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

// Test Ecs KeyPair. >>> Resource test cases, automatically generated.
// Case 3949
func TestAccAlicloudEcsKeyPair_basic3949(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_key_pair.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsKeyPairMap3949)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsKeyPair")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secskeypair%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsKeyPairBasicDependence3949)
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
					"key_pair_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${alicloud_resource_manager_resource_group.ResourceGroup.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_pair_name":     name + "_update",
					"resource_group_id": "${alicloud_resource_manager_resource_group.ResourceGroup.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name":     name + "_update",
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
				ImportStateVerifyIgnore: []string{"resource_type"},
			},
		},
	})
}

var AlicloudEcsKeyPairMap3949 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudEcsKeyPairBasicDependence3949(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_resource_manager_resource_group" "ResourceGroup" {
  display_name        = "test"
  resource_group_name = var.name
}


`, name)
}

// Case 3949  twin
func TestAccAlicloudEcsKeyPair_basic3949_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_key_pair.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsKeyPairMap3949)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsKeyPair")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secskeypair%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsKeyPairBasicDependence3949)
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
					"key_pair_name":     name,
					"resource_group_id": "${alicloud_resource_manager_resource_group.ResourceGroup.id}",
					"resource_type":     "keypair",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name":     name,
						"resource_group_id": CHECKSET,
						"resource_type":     "keypair",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_type"},
			},
		},
	})
}

// Test Ecs KeyPair. <<< Resource test cases, automatically generated.
