package dts_test

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/bssopenapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

const synchronizationInstanceResourceId = "alicloud_dts_synchronization_instance.default"

func TestAccAliCloudDTSSynchronizationInstance_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationinstance%d", dtsTestRegion(), rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-beijing")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkSynchronizationInstanceDestroy(synchronizationInstanceResourceId),
		Steps: []resource.TestStep{
			{
				Config: synchronizationInstanceBasicDependence0(name) + `
resource "alicloud_dts_synchronization_instance" "default" {
  auto_pay                         = "false"
  auto_start                       = "false"
  compute_unit                     = "2"
  database_count                   = "1"
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = "${data.alicloud_regions.default.regions.0.id}"
  instance_class                   = "small"
  payment_duration                 = "1"
  payment_duration_unit            = "Month"
  payment_type                     = "PayAsYouGo"
  quantity                          = "1"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region            = "${data.alicloud_regions.default.regions.0.id}"
  sync_architecture                 = "oneway"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "payment_type", "PayAsYouGo"),
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttrSet(synchronizationInstanceResourceId, "source_endpoint_region"),
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "destination_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttrSet(synchronizationInstanceResourceId, "destination_endpoint_region"),
					resource.TestCheckNoResourceAttr(synchronizationInstanceResourceId, "period"),
					resource.TestCheckNoResourceAttr(synchronizationInstanceResourceId, "used_time"),
					resource.TestCheckNoResourceAttr(synchronizationInstanceResourceId, "order_type"),
					resource.TestCheckNoResourceAttr(synchronizationInstanceResourceId, "synchronization_direction"),
					resource.TestCheckNoResourceAttr(synchronizationInstanceResourceId, "status"),
				),
			},
			{
				// auto_pay is ignored for PayAsYouGo instances: the provider must not
				// touch BSS, so the configured value stays in state unchanged.
				Config: synchronizationInstanceBasicDependence0(name) + fmt.Sprintf(`
resource "alicloud_dts_synchronization_instance" "default" {
%s}
`, dtsRenderAttrs(map[string]interface{}{
					"auto_pay":                         "true",
					"auto_start":                       "false",
					"compute_unit":                     "2",
					"database_count":                   "1",
					"destination_endpoint_engine_name": "MySQL",
					"destination_endpoint_region":      "${data.alicloud_regions.default.regions.0.id}",
					"instance_class":                   "small",
					"payment_duration":                 "1",
					"payment_duration_unit":            "Month",
					"payment_type":                     "PayAsYouGo",
					"quantity":                         "1",
					"source_endpoint_engine_name":      "MySQL",
					"source_endpoint_region":           "${data.alicloud_regions.default.regions.0.id}",
					"sync_architecture":                "oneway",
				})),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "auto_pay", "true"),
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "payment_type", "PayAsYouGo"),
				),
			},
			{
				// Removing auto_pay from a pay-as-you-go instance must not fail and
				// must converge: the optional-computed attribute keeps its value in
				// state and the provider never reaches BSS for this payment type.
				Config: synchronizationInstanceBasicDependence0(name) + fmt.Sprintf(`
resource "alicloud_dts_synchronization_instance" "default" {
%s}
`, dtsRenderAttrs(map[string]interface{}{
					"auto_start":                       "false",
					"compute_unit":                     "2",
					"database_count":                   "1",
					"destination_endpoint_engine_name": "MySQL",
					"destination_endpoint_region":      "${data.alicloud_regions.default.regions.0.id}",
					"instance_class":                   "small",
					"payment_duration":                 "1",
					"payment_duration_unit":            "Month",
					"payment_type":                     "PayAsYouGo",
					"quantity":                         "1",
					"source_endpoint_engine_name":      "MySQL",
					"source_endpoint_region":           "${data.alicloud_regions.default.regions.0.id}",
					"sync_architecture":                "oneway",
				})),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "auto_pay", "true"),
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "payment_type", "PayAsYouGo"),
				),
			},
			{
				// The removal above must leave an empty plan.
				Config: synchronizationInstanceBasicDependence0(name) + fmt.Sprintf(`
resource "alicloud_dts_synchronization_instance" "default" {
%s}
`, dtsRenderAttrs(map[string]interface{}{
					"auto_start":                       "false",
					"compute_unit":                     "2",
					"database_count":                   "1",
					"destination_endpoint_engine_name": "MySQL",
					"destination_endpoint_region":      "${data.alicloud_regions.default.regions.0.id}",
					"instance_class":                   "small",
					"payment_duration":                 "1",
					"payment_duration_unit":            "Month",
					"payment_type":                     "PayAsYouGo",
					"quantity":                         "1",
					"source_endpoint_engine_name":      "MySQL",
					"source_endpoint_region":           "${data.alicloud_regions.default.regions.0.id}",
					"sync_architecture":                "oneway",
				})),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:            synchronizationInstanceResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_count", "status", "quantity", "sync_architecture", "auto_pay", "auto_start", "compute_unit", "payment_duration", "payment_duration_unit"},
			},
		},
	})
}

func checkSynchronizationInstanceDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_dts_synchronization_instance" {
				continue
			}
			client := acctest.Provider.Meta().(*connectivity.AliyunClient)
			if _, err := dts.NewService(client).DescribeDtsSynchronizationInstance(rs.Primary.ID); err != nil {
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

func synchronizationInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

// The region has to come from the provider rather than from ALICLOUD_REGION: that
// variable is not always set, and an empty region_id makes CreateDtsInstance fail
// with MissingDestinationRegion before the instance is created.
data "alicloud_regions" "default" {
  current = true
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudDTSSynchronizationInstance(t *testing.T) {
	resourceSchema := dts.ResourceSynchronizationInstance()
	dInit, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"auto_pay":                         "false",
		"auto_start":                       "false",
		"payment_type":                     "PayAsYouGo",
		"source_endpoint_engine_name":      "CreateDtsInstanceValue",
		"source_endpoint_region":           "CreateDtsInstanceValue",
		"destination_endpoint_engine_name": "CreateDtsInstanceValue",
		"destination_endpoint_region":      "CreateDtsInstanceValue",
		"instance_class":                   "CreateDtsInstanceValue",
		"compute_unit":                     1,
		"database_count":                   1,
		"payment_duration_unit":            "CreateDtsInstanceValue",
		"quantity":                         1,
		"sync_architecture":                "CreateDtsInstanceValue",
		"payment_duration":                 1,
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
	}
	region := "cn-beijing"
	rawClient, err := dtsTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	ReadMockResponse := map[string]interface{}{
		// DescribeDtsJobDetail
		"SourceEndpoint": map[string]interface{}{
			"EngineName": "CreateDtsInstanceValue",
			"Region":     "CreateDtsInstanceValue",
		},
		"DestinationEndpoint": map[string]interface{}{
			"EngineName": "CreateDtsInstanceValue",
			"Region":     "CreateDtsInstanceValue",
		},
		"PayType":     "PostPaid",
		"DtsJobClass": "CreateDtsInstanceValue",
		"InstanceId":  "CreateDtsInstanceValue",
	}
	CreateMockResponse := map[string]interface{}{
		//CreateDtsInstance
		"InstanceId": "CreateDtsInstanceValue",
	}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String(errorCode),
			Data:       tea.String(errorCode),
			Message:    tea.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_dts_synchronization_instance", errorCode))
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
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDtsInstance" {
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
	err = resourceSchema.Update(dExisted, rawClient)
	assert.Nil(t, err)

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDtsJobDetail" {
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidJobId"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteSynchronizationJob" {
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
		case "InvalidJobId":
			assert.Nil(t, err)
		}
	}
}

// TestAccAliCloudDTSSynchronizationInstance_renewal covers the auto-renew
// update path for a prepaid instance. It purchases a Subscription instance, so
// it is skipped unless the exact function name is selected or
// TF_ACC_DTS_AUTO_PAY is set.
//
// Cleanup: Terraform destroy only removes a Subscription instance from state
// (the resource Delete is a no-op), so the instance keeps running. DTS does not
// support releasing a prepaid instance through an API, and a BSS refund inquiry
// for ProductCode=dts returns CommodityNotSupported. Clean up the instance
// through the console unsubscribe flow or a support ticket, and record the
// instance ID and the final cleanup result in the test evidence.
func TestAccAliCloudDTSSynchronizationInstance_renewal(t *testing.T) {
	skipUnlessDtsAutoPaySelected(t)
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationinstance%d", dtsTestRegion(), rand)
	config := func(attrs map[string]interface{}) string {
		return synchronizationInstanceBasicDependence0(name) + fmt.Sprintf(`
resource "alicloud_dts_synchronization_instance" "default" {
%s}
`, dtsRenderAttrs(attrs))
	}
	subscriptionAttrs := func(autoPay string) map[string]interface{} {
		attrs := map[string]interface{}{
			"auto_start":                       "false",
			"destination_endpoint_engine_name": "MySQL",
			"destination_endpoint_region":      "${data.alicloud_regions.default.regions.0.id}",
			"instance_class":                   "small",
			"payment_duration":                 "1",
			"payment_duration_unit":            "Month",
			"payment_type":                     "Subscription",
			"source_endpoint_engine_name":      "MySQL",
			"source_endpoint_region":           "${data.alicloud_regions.default.regions.0.id}",
			"sync_architecture":                "oneway",
		}
		if autoPay != "" {
			attrs["auto_pay"] = autoPay
		}
		return attrs
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-beijing")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// CheckDestroy stays nil: Delete only drops a Subscription instance from
		// state, because DTS has no release API and a BSS refund is unsupported for
		// dts. The generic destroy check can never pass for this case; unsubscribe
		// the instance through the console or a support ticket after the run.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config(subscriptionAttrs("false")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "auto_pay", "false"),
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "payment_type", "Subscription"),
					testAccCheckDtsSynchronizationInstanceRenewStatus(synchronizationInstanceResourceId, "ManualRenewal"),
				),
			},
			{
				Config: config(subscriptionAttrs("true")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "auto_pay", "true"),
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "payment_type", "Subscription"),
					testAccCheckDtsSynchronizationInstanceRenewStatus(synchronizationInstanceResourceId, "AutoRenewal"),
				),
			},
			{
				// Removing the argument leaves the computed value in state and must not
				// write BSS: the instance keeps the auto renewal setting from step 2.
				Config: config(subscriptionAttrs("")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "auto_pay", "true"),
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "payment_type", "Subscription"),
					testAccCheckDtsSynchronizationInstanceRenewStatus(synchronizationInstanceResourceId, "AutoRenewal"),
				),
			},
			{
				Config: config(subscriptionAttrs("false")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "auto_pay", "false"),
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "payment_type", "Subscription"),
					testAccCheckDtsSynchronizationInstanceRenewStatus(synchronizationInstanceResourceId, "ManualRenewal"),
				),
			},
			{
				ResourceName:            synchronizationInstanceResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_count", "status", "quantity", "sync_architecture", "auto_start", "compute_unit", "payment_duration", "payment_duration_unit"},
			},
		},
	})
}

// skipUnlessDtsAutoPaySelected keeps the prepaid renewal case out of broad or
// accidental runs: it purchases a Subscription instance and Terraform destroy
// only removes it from state. Run it by selecting this exact function name, or
// set TF_ACC_DTS_AUTO_PAY=1 to allow any invocation.
func skipUnlessDtsAutoPaySelected(t *testing.T) {
	t.Helper()
	if os.Getenv("TF_ACC_DTS_AUTO_PAY") != "" {
		return
	}
	runFlag := flag.Lookup("test.run")
	if runFlag == nil || runFlag.Value == nil {
		t.Skip("Skipping prepaid auto_pay test: select TestAccAliCloudDTSSynchronizationInstance_renewal by exact name or set TF_ACC_DTS_AUTO_PAY=1")
	}
	pattern := strings.TrimSpace(runFlag.Value.String())
	pattern = strings.TrimPrefix(pattern, "^")
	pattern = strings.TrimSuffix(pattern, "$")
	if strings.HasPrefix(pattern, "(") && strings.HasSuffix(pattern, ")") {
		pattern = pattern[1 : len(pattern)-1]
	}
	for _, part := range strings.Split(pattern, "|") {
		if part == t.Name() {
			return
		}
	}
	t.Skip("Skipping prepaid auto_pay test: select TestAccAliCloudDTSSynchronizationInstance_renewal by exact name or set TF_ACC_DTS_AUTO_PAY=1")
}

// testAccCheckDtsSynchronizationInstanceRenewStatus reads the renewal setting
// from BSS, which is the only place that records it for DTS instances.
func testAccCheckDtsSynchronizationInstanceRenewStatus(resourceId, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		deadline := time.Now().Add(60 * time.Second)
		var last string
		for {
			renewObject, err := bssopenapi.NewService(client).QueryAvailableInstances(resourceState.Primary.ID, "", "dts", "dts")
			if err == nil {
				last = fmt.Sprint(renewObject["RenewStatus"])
				if last == expected {
					return nil
				}
			} else {
				last = err.Error()
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("expected RenewStatus %s for %s, got %s", expected, resourceState.Primary.ID, last)
			}
			time.Sleep(5 * time.Second)
		}
	}
}

// buildAutoPayUnitResourceData builds the planned diff of the auto_pay update path.
func buildAutoPayUnitResourceData(t *testing.T, paymentType, oldAutoPay, newAutoPay string, extra map[string]string) *schema.ResourceData {
	t.Helper()
	diff := sdkterraform.NewInstanceDiff()
	diff.Attributes["auto_pay"] = &sdkterraform.ResourceAttrDiff{Old: oldAutoPay, New: newAutoPay}
	diff.Attributes["payment_type"] = &sdkterraform.ResourceAttrDiff{Old: paymentType, New: paymentType}
	for key, value := range extra {
		diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: "", New: value}
	}
	d, err := schema.InternalMap(dts.ResourceSynchronizationInstance().Schema).Data(nil, diff)
	assert.Nil(t, err)
	d.SetId("dts-test123")
	return d
}

// runAutoPayUnitUpdate runs the resource update with a mocked RPC client and
// reports the SetRenewal/QueryAvailableInstances calls plus the SetRenewal body.
func runAutoPayUnitUpdate(t *testing.T, d *schema.ResourceData, payType, renewStatus string) (int, int, map[string]interface{}, error) {
	t.Helper()
	endpoints := new(sync.Map)
	endpoints.Store("dts", "dts.cn-hangzhou.aliyuncs.com")
	endpoints.Store("bssopenapi", "bssopenapi.aliyuncs.com")
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId("test-access-key").SetAccessKeySecret("test-secret-key")
	testCredential, credErr := credentials.NewCredential(credentialConfig)
	assert.Nil(t, credErr)
	providerConfig := connectivity.Config{
		Region:               connectivity.Region("cn-hangzhou"),
		RegionId:             "cn-hangzhou",
		Protocol:             "HTTPS",
		AccessKey:            "test-access-key",
		SecretKey:            "test-secret-key",
		AccountType:          "Domestic",
		SkipRegionValidation: true,
		Endpoints:            endpoints,
		Credential:           testCredential,
	}
	aliyunClient, err := providerConfig.Client()
	assert.Nil(t, err)

	setRenewalCalls := 0
	queryCalls := 0
	var setRenewalRequest map[string]interface{}

	patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		switch *action {
		case "SetRenewal":
			setRenewalCalls++
			setRenewalRequest = body
			return map[string]interface{}{"Success": true}, nil
		case "QueryAvailableInstances":
			queryCalls++
			return map[string]interface{}{
				"Data": map[string]interface{}{
					"InstanceList": []interface{}{
						map[string]interface{}{"InstanceID": "dts-test123", "RenewStatus": renewStatus},
					},
				},
			}, nil
		default:
			return map[string]interface{}{
				"PayType":             payType,
				"DtsJobClass":         "small",
				"SourceEndpoint":      map[string]interface{}{"EngineName": "MySQL", "Region": "cn-hangzhou"},
				"DestinationEndpoint": map[string]interface{}{"EngineName": "MySQL", "Region": "cn-hangzhou"},
			}, nil
		}
	})
	defer patches.Reset()

	err = dts.ResourceSynchronizationInstance().Update(d, aliyunClient)
	return setRenewalCalls, queryCalls, setRenewalRequest, err
}

// TestUnitDtsSynchronizationInstanceAutoPayUpdate covers the update path: the
// BSS SetRenewal request, the read-back mapping of the renewal status, the
// pay-as-you-go skip and the rejection of an incomplete auto renewal request.
func TestUnitDtsSynchronizationInstanceAutoPayUpdate(t *testing.T) {
	t.Run("enable auto renewal with monthly period", func(t *testing.T) {
		d := buildAutoPayUnitResourceData(t, "Subscription", "false", "true", map[string]string{
			"payment_duration":      "3",
			"payment_duration_unit": "Month",
		})
		calls, queryCalls, request, err := runAutoPayUnitUpdate(t, d, "PrePaid", "AutoRenewal")
		assert.Nil(t, err)
		assert.Equal(t, 1, calls)
		assert.GreaterOrEqual(t, queryCalls, 1)
		assert.Equal(t, "dts-test123", request["InstanceIDs"])
		assert.Equal(t, "dts", request["ProductCode"])
		assert.Equal(t, "Subscription", request["SubscriptionType"])
		assert.Equal(t, "AutoRenewal", request["RenewalStatus"])
		assert.Equal(t, 3, request["RenewalPeriod"])
		assert.Equal(t, "M", request["RenewalPeriodUnit"])
		assert.Equal(t, "true", d.Get("auto_pay"))
	})

	t.Run("enable auto renewal with yearly period", func(t *testing.T) {
		d := buildAutoPayUnitResourceData(t, "Subscription", "false", "true", map[string]string{
			"payment_duration":      "1",
			"payment_duration_unit": "Year",
		})
		calls, _, request, err := runAutoPayUnitUpdate(t, d, "PrePaid", "AutoRenewal")
		assert.Nil(t, err)
		assert.Equal(t, 1, calls)
		assert.Equal(t, "AutoRenewal", request["RenewalStatus"])
		assert.Equal(t, 1, request["RenewalPeriod"])
		assert.Equal(t, "Y", request["RenewalPeriodUnit"])
	})

	t.Run("disable auto renewal", func(t *testing.T) {
		d := buildAutoPayUnitResourceData(t, "Subscription", "true", "false", map[string]string{
			"payment_duration":      "1",
			"payment_duration_unit": "Month",
		})
		calls, queryCalls, request, err := runAutoPayUnitUpdate(t, d, "PrePaid", "ManualRenewal")
		assert.Nil(t, err)
		assert.Equal(t, 1, calls)
		assert.GreaterOrEqual(t, queryCalls, 1)
		assert.Equal(t, "ManualRenewal", request["RenewalStatus"])
		_, ok := request["RenewalPeriod"]
		assert.False(t, ok)
		_, ok = request["RenewalPeriodUnit"]
		assert.False(t, ok)
		assert.Equal(t, "false", d.Get("auto_pay"))
	})

	t.Run("skip auto renewal update for pay as you go instance", func(t *testing.T) {
		d := buildAutoPayUnitResourceData(t, "PayAsYouGo", "false", "true", nil)
		calls, queryCalls, _, err := runAutoPayUnitUpdate(t, d, "PostPaid", "ManualRenewal")
		assert.Nil(t, err)
		assert.Equal(t, 0, calls)
		assert.Equal(t, 0, queryCalls)
	})

	t.Run("skip auto renewal update for pay as you go instance when auto_pay is removed", func(t *testing.T) {
		// Removing an optional-computed attribute reaches the update path with an
		// explicit removal diff; a pay-as-you-go instance must still stay out of
		// BSS and must not have a fabricated auto_pay value written back.
		d := buildAutoPayUnitResourceData(t, "PayAsYouGo", "true", "", nil)
		calls, queryCalls, _, err := runAutoPayUnitUpdate(t, d, "PostPaid", "ManualRenewal")
		assert.Nil(t, err)
		assert.Equal(t, 0, calls)
		assert.Equal(t, 0, queryCalls)
		assert.Equal(t, "", d.Get("auto_pay"))
	})

	t.Run("reject auto renewal without duration unit", func(t *testing.T) {
		d := buildAutoPayUnitResourceData(t, "Subscription", "false", "true", map[string]string{
			"payment_duration": "1",
		})
		calls, queryCalls, _, err := runAutoPayUnitUpdate(t, d, "PrePaid", "ManualRenewal")
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "payment_duration_unit must be set when auto_pay is true")
		}
		assert.Equal(t, 0, calls)
		assert.Equal(t, 0, queryCalls)
	})

	t.Run("reject auto renewal without duration", func(t *testing.T) {
		d := buildAutoPayUnitResourceData(t, "Subscription", "false", "true", map[string]string{
			"payment_duration_unit": "Month",
		})
		calls, queryCalls, _, err := runAutoPayUnitUpdate(t, d, "PrePaid", "ManualRenewal")
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "payment_duration must be set when auto_pay is true")
		}
		assert.Equal(t, 0, calls)
		assert.Equal(t, 0, queryCalls)
	})

	t.Run("reject unsupported renewal period", func(t *testing.T) {
		d := buildAutoPayUnitResourceData(t, "Subscription", "false", "true", map[string]string{
			"payment_duration":      "24",
			"payment_duration_unit": "Month",
		})
		calls, queryCalls, _, err := runAutoPayUnitUpdate(t, d, "PrePaid", "ManualRenewal")
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "payment_duration must be one of 1, 2, 3, 6, 12 when auto_pay is true, got 24")
		}
		assert.Equal(t, 0, calls)
		assert.Equal(t, 0, queryCalls)
	})
}

// TestUnitDtsSynchronizationInstanceAutoPayResponse maps the BSS RenewStatus
// value onto the auto_pay attribute through the read path.
func TestUnitDtsSynchronizationInstanceAutoPayResponse(t *testing.T) {
	cases := []struct {
		renewStatus string
		expected    string
	}{
		{renewStatus: "AutoRenewal", expected: "true"},
		{renewStatus: "ManualRenewal", expected: "false"},
		{renewStatus: "NotRenewal", expected: "false"},
	}
	for _, c := range cases {
		t.Run(c.renewStatus, func(t *testing.T) {
			d := buildAutoPayUnitResourceData(t, "Subscription", "true", "false", nil)
			_, _, _, err := runAutoPayUnitUpdate(t, d, "PrePaid", c.renewStatus)
			assert.Nil(t, err)
			assert.Equal(t, c.expected, d.Get("auto_pay"))
		})
	}
}
