package dts_test

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

const migrationInstanceResourceId = "alicloud_dts_migration_instance.default"

func TestAccAliCloudDTSMigrationInstance_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsmigrationinstance%d", dtsTestRegion(), rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// The flat-suite gates (checkoutSupportedRegions + testAccPreCheckWithRegions)
			// converged to cn-beijing for every input region, so pin it here.
			os.Setenv("ALICLOUD_REGION", "cn-beijing")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkMigrationInstanceDestroy(migrationInstanceResourceId),
		Steps: []resource.TestStep{
			{
				Config: migrationInstanceBasicConfig(name, ""),
				Check:  migrationInstanceBasicChecks(),
			},
			{
				Config: migrationInstanceBasicConfig(name, `
  tags = {
    Created = "TF"
    For     = "Test"
    From    = "acceptance"
  }
`),
				Check: migrationInstanceBasicChecks(
					resource.TestCheckResourceAttr(migrationInstanceResourceId, "tags.%", "3"),
				),
			},
			{
				Config: migrationInstanceBasicConfig(name, `
  tags = {
    Created = "TF"
    For     = "acceptance-Test"
  }
`),
				Check: migrationInstanceBasicChecks(
					resource.TestCheckResourceAttr(migrationInstanceResourceId, "tags.%", "2"),
				),
			},
			{
				ResourceName:            migrationInstanceResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_count", "sync_architecture", "compute_unit"},
			},
		},
	})
}

func TestAccAliCloudDTSMigrationInstance_mariadb(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsmigrationinstancemariadb%d", dtsTestRegion(), rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-beijing")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkMigrationInstanceDestroy(migrationInstanceResourceId),
		Steps: []resource.TestStep{
			{
				Config: migrationInstanceDependence(name) + `
resource "alicloud_dts_migration_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "mariadb"
  source_endpoint_region           = data.alicloud_regions.default.regions.0.id
  destination_endpoint_engine_name = "mariadb"
  destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
  instance_class                   = "small"
  sync_architecture                = "oneway"
  compute_unit                     = 2
  database_count                   = 1
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "status"),
					resource.TestCheckResourceAttr(migrationInstanceResourceId, "payment_type", "PayAsYouGo"),
					resource.TestCheckResourceAttr(migrationInstanceResourceId, "instance_class", "small"),
					resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "source_endpoint_engine_name"),
					resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "source_endpoint_region"),
					resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "destination_endpoint_engine_name"),
					resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "destination_endpoint_region"),
				),
			},
			{
				ResourceName:            migrationInstanceResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_count", "sync_architecture", "compute_unit"},
			},
		},
	})
}

func migrationInstanceBasicChecks(extras ...resource.TestCheckFunc) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "status"),
		resource.TestCheckResourceAttr(migrationInstanceResourceId, "payment_type", "PayAsYouGo"),
		resource.TestCheckResourceAttr(migrationInstanceResourceId, "instance_class", "small"),
		resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "destination_endpoint_region"),
		resource.TestCheckResourceAttr(migrationInstanceResourceId, "source_endpoint_engine_name", "MySQL"),
		resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "source_endpoint_region"),
		resource.TestCheckResourceAttr(migrationInstanceResourceId, "destination_endpoint_engine_name", "MySQL"),
	}
	return resource.ComposeTestCheckFunc(append(checks, extras...)...)
}

func checkMigrationInstanceDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_dts_migration_instance" {
				continue
			}
			client := acctest.Provider.Meta().(*connectivity.AliyunClient)
			if _, err := dts.NewService(client).DescribeMigrationInstance(rs.Primary.ID); err != nil {
				if errs.NotFoundError(err) {
					continue
				}
				return err
			}
			return fmt.Errorf("the resource %s %s was not destroyed", resourceId, rs.Primary.ID)
		}
		return nil
	}
}

func migrationInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_regions" "default" {
  current = true
}
`, name)
}

func migrationInstanceBasicConfig(name string, extra string) string {
	return migrationInstanceDependence(name) + fmt.Sprintf(`
resource "alicloud_dts_migration_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region           = data.alicloud_regions.default.regions.0.id
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
  instance_class                   = "small"
  sync_architecture                = "oneway"
%s}
`, extra)
}

func dtsTestRegion() string {
	if region := os.Getenv("ALICLOUD_REGION"); region != "" {
		return region
	}
	return "cn-beijing"
}

// lintignore: R001
func TestUnitAlicloudDTSMigrationInstance(t *testing.T) {
	resourceSchema := dts.ResourceMigrationInstance()
	d, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"payment_type":                     "PayAsYouGo",
		"source_endpoint_engine_name":      "MySQL",
		"source_endpoint_region":           "source_endpoint_region",
		"destination_endpoint_engine_name": "MySQL",
		"destination_endpoint_region":      "destination_endpoint_region",
		"instance_class":                   "small",
		"sync_architecture":                "oneway",
		"compute_unit":                     1,
		"database_count":                   1,
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := "cn-beijing"
	rawClient, err := dtsTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
	}
	ReadMockResponse := map[string]interface{}{
		"DtsInstanceID": "dts_instance_id",
		"DtsJobClass":   "xlarge",
		"PayType":       "PostPaid",
		"DestinationEndpoint": map[string]interface{}{
			"EngineName": "MySQL",
			"Region":     "destination_endpoint_region",
		},
		"SourceEndpoint": map[string]interface{}{
			"EngineName": "MySQL",
			"Region":     "source_endpoint_region",
		},
		"Status": "Migrating",
		"TagResources": map[string]interface{}{
			"TagResource": []interface{}{
				map[string]interface{}{
					"TagKey":   "tag_key",
					"TagValue": "tag_value",
				},
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String(errorCode),
				Data:       tea.String(errorCode),
				Message:    tea.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_dts_migration_instance", "dts_instance_id"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String(errorCode),
				Data:       tea.String(errorCode),
				Message:    tea.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateStatusNormal": func(errorCode string) (map[string]interface{}, error) {
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
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String("loadEndpoint error"),
				Data:       tea.String("loadEndpoint error"),
				Message:    tea.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceSchema.Create(d, rawClient)
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
		err := resourceSchema.Create(d, rawClient)
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
		err := resourceSchema.Create(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("CreateRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceSchema.Create(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String("loadEndpoint error"),
				Data:       tea.String("loadEndpoint error"),
				Message:    tea.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceSchema.Update(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("dts_instance_id")

	t.Run("UpdateModifyDtsMigrationInstanceTagsAbnormal", func(t *testing.T) {
		diff := sdkterraform.NewInstanceDiff()
		for _, key := range []string{"tags"} {
			switch resourceSchema.Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &sdkterraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &sdkterraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &sdkterraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, diff)
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
		err := resourceSchema.Update(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyDTSMigrationInstanceTagsNoRetryErrorAbnormal", func(t *testing.T) {
		diff := sdkterraform.NewInstanceDiff()
		for _, key := range []string{"tags"} {
			switch resourceSchema.Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &sdkterraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &sdkterraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &sdkterraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		patcheSetResourceTags := gomonkey.ApplyMethod(reflect.TypeOf(&dts.Service{}), "SetResourceTags", func(*dts.Service, *schema.ResourceData, string) error {
			_, err := responseMock["NoRetryError"]("NoRetryError")
			return err
		})
		err := resourceSchema.Update(resourceData, rawClient)
		patches.Reset()
		patcheSetResourceTags.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyDTSMigrationInstanceTagsNoRetryErrorNormal", func(t *testing.T) {
		diff := sdkterraform.NewInstanceDiff()
		for _, key := range []string{"tags"} {
			switch resourceSchema.Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &sdkterraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &sdkterraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &sdkterraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		patcheSetResourceTags := gomonkey.ApplyMethod(reflect.TypeOf(&dts.Service{}), "SetResourceTags", func(*dts.Service, *schema.ResourceData, string) error {
			return nil
		})
		err := resourceSchema.Update(resourceData, rawClient)
		patches.Reset()
		patcheSetResourceTags.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String("loadEndpoint error"),
				Data:       tea.String("loadEndpoint error"),
				Message:    tea.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceSchema.Update(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				// retry until the timeout comes
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceSchema.Delete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceSchema.Delete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("DeleteNonRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceSchema.Delete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteDescribeVpcDTSMigrationInstanceExpectedError", func(t *testing.T) {
		patchRequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("Forbidden.InstanceNotFound")
		})
		err := resourceSchema.Delete(d, rawClient)
		patchRequest.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeVpcDTSMigrationInstanceNotFound", func(t *testing.T) {
		diff := sdkterraform.NewInstanceDiff()
		resourceData, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&dts.Service{}), "DescribeMigrationInstance", func(*dts.Service, string) (map[string]interface{}, error) {
			return responseMock["NotFoundError"]("NotFoundError")
		})
		err := resourceSchema.Read(resourceData, rawClient)
		patcheDescribe.Reset()
		assert.Nil(t, err)
	})
	t.Run("ReadDescribeDTSMigrationInstanceAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["ReadNormal"]("")
		})
		err := resourceSchema.Delete(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})
}
