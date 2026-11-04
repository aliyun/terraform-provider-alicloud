package fnf_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/fnf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

var fnfTestEndpoints sync.Map
var fnfTestSignVersion sync.Map

// fnfTestClientForRegion rebuilds the package-level sweeper client on the
// exported connectivity API: the flat-suite sharedClientForRegion is not
// reachable from this package.
func fnfTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
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
		Endpoints:   &fnfTestEndpoints,
		SignVersion: &fnfTestSignVersion,
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

func fnfTestRegion() string {
	if region := os.Getenv("ALICLOUD_REGION"); region != "" {
		return region
	}
	return "cn-beijing"
}

func init() {
	resource.AddTestSweepers("alicloud_fnf_flow", &resource.Sweeper{
		Name:         "alicloud_fnf_flow",
		F:            testSweepFnfFlow,
		Dependencies: []string{"alicloud_fnf_schedule"},
	})
}

func testSweepFnfFlow(region string) error {
	client, err := fnfTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	support := false
	for _, v := range connectivity.FnfSupportRegions {
		if v == connectivity.Region(region) {
			support = true
			break
		}
	}
	if !support {
		return nil
	}

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListFlows"
	request := make(map[string]interface{})
	var response map[string]interface{}
	response, err = client.RpcGet("fnf", "2019-03-15", action, request, nil)
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_fnf_flows", action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	resp, err := jsonpath.Get("$.Flows", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Flows", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := item["Name"].(string)
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(name, prefix) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Fnf Flow: %s ", name)
			continue
		}
		log.Printf("[Info] Delete Fnf Flow: %s", name)

		action := "DeleteFlow"
		request := map[string]interface{}{
			"Name": name,
		}
		_, err = client.RpcGet("fnf", "2019-03-15", action, request, nil)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Fnf Flow (%s): %s", name, err)
		}
	}
	return nil
}

const fnfFlowResourceId = "alicloud_fnf_flow.default"

func TestAccAlicloudFnfFlow_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudFnfFlow%d", fnfTestRegion(), rand)
	definition1 := `version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworld`
	definition2 := `version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworldchange`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.FnfSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkFnfFlowDestroy(fnfFlowResourceId),
		Steps: []resource.TestStep{
			{
				Config: fnfFlowTestAccDependence(name) + fmt.Sprintf(`
resource "alicloud_fnf_flow" "default" {
  definition  = "%s"
  description = "tf-testaccFnFFlow983041"
  name        = var.name
  type        = "FDL"
}
`, definition1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fnfFlowResourceId, "definition", strings.Replace(definition1, `\n`, "\n", -1)),
					resource.TestCheckResourceAttr(fnfFlowResourceId, "description", "tf-testaccFnFFlow983041"),
					resource.TestCheckResourceAttr(fnfFlowResourceId, "name", name),
					resource.TestCheckResourceAttr(fnfFlowResourceId, "type", "FDL"),
					resource.TestCheckResourceAttrSet(fnfFlowResourceId, "flow_id"),
					resource.TestCheckResourceAttrSet(fnfFlowResourceId, "last_modified_time"),
				),
			},
			{
				ResourceName:      fnfFlowResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fnfFlowTestAccDependence(name) + fmt.Sprintf(`
resource "alicloud_fnf_flow" "default" {
  definition  = "%s"
  description = "tf-testaccFnFFlow983041"
  name        = var.name
  type        = "FDL"
}
`, definition2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fnfFlowResourceId, "definition", strings.Replace(definition2, `\n`, "\n", -1)),
				),
			},
			{
				Config: fnfFlowTestAccDependence(name) + fmt.Sprintf(`
resource "alicloud_fnf_flow" "default" {
  definition  = "%s"
  description = "tf-testaccFnFFlow813242"
  name        = var.name
  type        = "FDL"
}
`, definition2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fnfFlowResourceId, "description", "tf-testaccFnFFlow813242"),
				),
			},
			{
				Config: fnfFlowTestAccDependence(name) + fmt.Sprintf(`
resource "alicloud_fnf_flow" "default" {
  definition  = "%s"
  description = "tf-testaccFnFFlow813242"
  name        = var.name
  role_arn    = alicloud_ram_role.default.arn
  type        = "FDL"
}
`, definition2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(fnfFlowResourceId, "role_arn"),
				),
			},
			{
				Config: fnfFlowTestAccDependence(name) + fmt.Sprintf(`
resource "alicloud_fnf_flow" "default" {
  definition  = "%s"
  description = "tf-testaccFnFFlow813242"
  name        = var.name
  role_arn    = alicloud_ram_role.default.arn
  type        = "FDL"
}
`, definition2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fnfFlowResourceId, "type", "FDL"),
				),
			},
			{
				Config: fnfFlowTestAccDependence(name) + fmt.Sprintf(`
resource "alicloud_fnf_flow" "default" {
  definition  = "%s"
  description = "tf-testaccFnFFlow983041"
  name        = var.name
  role_arn    = alicloud_ram_role.default.arn
  type        = "FDL"
}
`, definition1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fnfFlowResourceId, "definition", strings.Replace(definition1, `\n`, "\n", -1)),
					resource.TestCheckResourceAttr(fnfFlowResourceId, "description", "tf-testaccFnFFlow983041"),
					resource.TestCheckResourceAttrSet(fnfFlowResourceId, "role_arn"),
					resource.TestCheckResourceAttr(fnfFlowResourceId, "type", "FDL"),
				),
			},
		},
	})
}

func checkFnfFlowDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		fnfService := fnf.NewService(client)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_fnf_flow" || rs.Primary.ID != resourceId {
				continue
			}
			if _, err := fnfService.DescribeFnfFlow(rs.Primary.ID); err != nil {
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

func fnfFlowTestAccDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "alicloud_ram_role" "default" {
  name = var.name
  document = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "fnf.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudFnfFlow(t *testing.T) {
	resourceSchema := fnf.ResourceFnfFlow()
	dInit, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"definition":  "CreateFlowValue",
		"description": "CreateFlowValue",
		"name":        "CreateFlowValue",
		"type":        "CreateFlowValue",
		"role_arn":    "CreateFlowValue",
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
	region := fnfTestRegion()
	rawClient, err := fnfTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
	}
	ReadMockResponse := map[string]interface{}{
		// DescribeFlow
		"Name":             "CreateFlowValue",
		"Definition":       "CreateFlowValue",
		"Description":      "CreateFlowValue",
		"Id":               "CreateFlowValue",
		"LastModifiedTime": "CreateFlowValue",
		"RoleArn":          "CreateFlowValue",
		"Type":             "CreateFlowValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateFlow
		"Name": "CreateFlowValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_fnf_flow", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			for k, v := range operationMockResponse {
				ReadMockResponse[k] = v
			}
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
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeFlow Response
		"Name": "CreateFlowValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateFlow" {
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
				dCompare.Set(key, value)
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
	// UpdateFlow
	attributesDiff := map[string]interface{}{
		"definition":  "UpdateFlowValue",
		"description": "UpdateFlowValue",
		"role_arn":    "UpdateFlowValue",
		"type":        "UpdateFlowValue",
	}
	diff := fnfStringAttrDiff(resourceSchema, attributes, attributesDiff)
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeFlow Response
		"Definition":  "UpdateFlowValue",
		"Description": "UpdateFlowValue",
		"RoleArn":     "UpdateFlowValue",
		"Type":        "UpdateFlowValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "ConcurrentUpdateError", "InternalServerError", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateFlow" {
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
			if *action == "DescribeFlow" {
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcGet", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}) (map[string]interface{}, error) {
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
	errorCodes = []string{"NonRetryableError", "Throttling", "InternalServerError", "nil", "FlowNotExists"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteFlow" {
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
		case "FlowNotExists":
			assert.Nil(t, err)
		}
	}
}

// fnfStringAttrDiff builds an InstanceDiff for string attributes whose values
// change from oldAttrs to newAttrs; it mirrors the flat-suite newInstanceDiff
// helper that is not reachable from this external test package.
func fnfStringAttrDiff(resourceSchema *schema.Resource, oldAttrs, newAttrs map[string]interface{}) *sdkterraform.InstanceDiff {
	diff := sdkterraform.NewInstanceDiff()
	for key, newVal := range newAttrs {
		old := ""
		if v, ok := oldAttrs[key]; ok {
			old = fmt.Sprintf("%v", v)
		}
		switch resourceSchema.Schema[key].Type {
		case schema.TypeBool:
			diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: old, New: fmt.Sprintf("%v", newVal)}
		default:
			diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: old, New: fmt.Sprintf("%v", newVal)}
		}
	}
	return diff
}
