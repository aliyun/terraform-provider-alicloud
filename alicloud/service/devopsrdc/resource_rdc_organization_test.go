package devopsrdc_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/devopsrdc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// devopsrdcTestClientForRegion thin-delegates to the shared
// alicloud/sweep.SharedClientForRegion so the sweeper client reuses the
// platform's per-region endpoint cache (SharedEndpointCacheForRegion),
// shared signVersion and AccountId/AccountType/credential handling instead
// of a single global cache shared across all swept regions.
func devopsrdcTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	raw, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return nil, err
	}
	return raw.(*connectivity.AliyunClient), nil
}

func devopsrdcTestRegion() string {
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	return region
}

func init() {
	resource.AddTestSweepers("alicloud_rdc_organization", &resource.Sweeper{
		Name: "alicloud_rdc_organization",
		F:    testSweepRdcOrganization,
	})
}

func testSweepRdcOrganization(region string) error {
	client, err := devopsrdcTestClientForRegion(region)
	if err != nil {
		return errs.WrapError(err)
	}
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListUserOrganization"
	request := map[string]interface{}{}

	var response map[string]interface{}
	response, err = client.RpcPost("devops-rdc", "2020-03-03", action, nil, request, true)

	if err != nil {
		log.Printf("[ERROR] %s got an error: %s", action, err)
		return nil
	}
	resp, err := jsonpath.Get("$.Object", response)
	if err != nil {
		log.Printf("[ERROR] %s parsing response got an error: %s", action, err)
		return nil
	}

	result, _ := resp.([]interface{})

	for _, v := range result {
		item := v.(map[string]interface{})

		name := item["Name"]
		id := item["Id"]
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name.(string)), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping organization: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting organization: %s (%s)", name, id)
		action = "DeleteDevopsOrganization"
		request = map[string]interface{}{
			"OrgId": item["Id"],
		}

		response, err = client.RpcPost("devops-rdc", "2020-03-03", action, nil, request, false)

		if err != nil {
			log.Printf("[ERROR] Failed to delete organization(%s (%s)): %s", name, id, err)
		}
	}

	return nil
}

// resource alicloud_rdc_organization has been deprecated from version 1.238.0.
// TestAccAliCloudRdcOrganization_basic is a newly added representative
// acceptance case for the migrated alicloud_rdc_organization resource in the
// devopsrdc subpackage; it exercises real Create/Read/Destroy against the
// devops-rdc endpoint. The original disabled import case is preserved
// separately below in full; not being executed does not count as PASS.
func TestAccAliCloudRdcOrganization_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdcorganization%d", acctest.DefaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.RDCupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkRdcOrganizationDestroy("alicloud_rdc_organization.default"),
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudRdcOrganizationConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_rdc_organization.default", "organization_name", name),
					resource.TestCheckResourceAttr("alicloud_rdc_organization.default", "source", "tf-testaccsource"),
				),
			},
		},
	})
}

func checkRdcOrganizationDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		rdcService := devopsrdc.NewService(client)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_rdc_organization" {
				continue
			}
			if _, err := rdcService.DescribeRdcOrganization(rs.Primary.ID); err != nil {
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

func testAccAlicloudRdcOrganizationConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_rdc_organization" "default" {
  organization_name = var.name
  source            = "tf-testaccsource"
}
`, name)
}

// resource alicloud_rdc_organization has been deprecated from version 1.238.0
func SkipTestAccAlicloudRDCOrganization_basic1(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.RDCupportRegions)
	resourceId := "alicloud_rdc_organization.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRDCOrganizationMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return devopsrdc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRdcOrganization")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdcorganization%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRDCOrganizationBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"organization_name": name,
					"source":            "tf-testaccsource",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"organization_name": name,
						"source":            "tf-testaccsource",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"real_pk", "source", "desired_member_count"},
			},
		},
	})
}

var AlicloudRDCOrganizationMap0 = map[string]string{
	"real_pk":              acctest.NOSET,
	"source":               acctest.NOSET,
	"desired_member_count": acctest.NOSET,
}

func AlicloudRDCOrganizationBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudRDCOrganization(t *testing.T) {
	resourceSchema := devopsrdc.ResourceRdcOrganization()
	dInit, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"organization_name":    "CreateDevopsOrganizationValue",
		"source":               "CreateDevopsOrganizationValue",
		"desired_member_count": 1,
		"real_pk":              "CreateDevopsOrganizationValue",
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
	region := devopsrdcTestRegion()
	rawClient, err := devopsrdcTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
	}
	ReadMockResponse := map[string]interface{}{
		// ListUserOrganization
		"Object": []interface{}{
			map[string]interface{}{
				"Id":   "CreateDevopsOrganizationValue",
				"Name": "CreateDevopsOrganizationValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateDevopsOrganization
		"Object": "CreateDevopsOrganizationValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_rdc_organization", errorCode))
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
		// ListUserOrganization Response
		"Object": []interface{}{
			map[string]interface{}{
				"Id":   "CreateDevopsOrganizationValue",
				"Name": "CreateDevopsOrganizationValue",
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDevopsOrganization" {
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
	err = resourceSchema.Update(dExisted, rawClient)
	assert.Nil(t, err)

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ListUserOrganization" {
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidOrganization.NotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDevopsOrganization" {
				switch errorCode {
				case "NonRetryableError", "InvalidOrganization.NotFound":
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
		case "InvalidOrganization.NotFound":
			assert.Nil(t, err)
		}
	}

}
