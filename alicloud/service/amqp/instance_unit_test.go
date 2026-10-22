package amqp

// This file ports the pre-migration root-package unit tests
// (alicloud/resource_alicloud_amqp_instance_test.go) for
// TestUnitAliCloudAmqpInstance and TestUnitAliCloudAmqpInstanceEncryption
// into the migrated service/amqp package. The CRUD helpers are unexported,
// so this is an internal test (package amqp, not amqp_test). The mock
// responses and gomonkey dispatch are adapted to the migrated RPC call
// chain (rpcRequest -> tea-rpc client.DoRequest) and the migrated
// DescribeAmqpInstance which reads $.Data.Status. New Delete subtests
// cover the post-migration PayAsYouGo DeleteInstance and Subscription
// RefundInstance branches (including idempotent already-absent cases).
//
// The CreateClientAbnormal/UpdateClientAbnormal subtests from the
// original are intentionally dropped: the migrated code calls
// BssOpenApi through RpcPostWithEndpoint (rpcRequest), not through the
// now-removed NewBssopenapiClient method, so patching that method has no
// equivalent effect on the migrated resource.

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

// buildAmqpUnitClient mirrors amqpSweepClient: it builds a real
// *connectivity.AliyunClient from env credentials so that rpcRequest can
// resolve endpoints and credentials. gomonkey patches DoRequest so the
// client never reaches the network, but the client object must be non-nil
// with a populated SDK client field for the patch to apply.
func buildAmqpUnitClient(region, accessKey, secretKey, securityToken string) (*connectivity.AliyunClient, error) {
	// Unit tests never reach the network: every subtest patches
	// client.Client.DoRequest via gomonkey. When real credentials are
	// absent (e.g. CI without ALICLOUD_ACCESS_KEY), fall back to dummy
	// credentials so the mocked Delete/Read/Create/Update subtests still
	// execute instead of being skipped. SkipRegionValidation bypasses
	// loadAndValidate, and a non-empty AccountType short-circuits
	// getAccountType, which would otherwise issue a BssOpenApi RPC during
	// Client() construction.
	if accessKey == "" {
		accessKey = "dummy-access-key-id"
	}
	if secretKey == "" {
		secretKey = "dummy-access-key-secret"
	}
	conf := connectivity.Config{
		RegionId:             region,
		Region:               connectivity.Region(region),
		AccessKey:            accessKey,
		SecretKey:            secretKey,
		Protocol:             "HTTPS",
		Endpoints:            &sync.Map{},
		SignVersion:          &sync.Map{},
		SkipRegionValidation: true,
		AccountType:          "Domestic",
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

// testAmqpUnitClient builds the client for unit tests. Real credentials
// are not required: buildAmqpUnitClient falls back to dummy credentials so
// the mocked subtests run in CI; the test is only skipped when client
// construction itself fails.
func testAmqpUnitClient(t *testing.T) *connectivity.AliyunClient {
	t.Helper()
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	rawClient, err := buildAmqpUnitClient(region, os.Getenv("ALICLOUD_ACCESS_KEY"), os.Getenv("ALICLOUD_SECRET_KEY"), os.Getenv("ALICLOUD_SECURITY_TOKEN"))
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		return nil
	}
	return rawClient
}

// lintignore: R001
func TestUnitAliCloudAmqpInstance(t *testing.T) {
	p := ResourceAmqpInstance()
	d, _ := schema.InternalMap(p.Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p.Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	dCreateMock, _ := schema.InternalMap(p.Schema).Data(nil, nil)
	dCreateMock.MarkNewResource()
	dCreateRenewalStatus, _ := schema.InternalMap(p.Schema).Data(nil, nil)
	dCreateRenewalStatus.MarkNewResource()
	for key, value := range map[string]interface{}{
		"instance_name":    "${var.name}",
		"instance_type":    "professional",
		"max_tps":          "1000",
		"payment_type":     "Subscription",
		"period":           1,
		"queue_capacity":   "50",
		"support_eip":      false,
		"max_eip_tps":      "128",
		"renewal_duration": 1,
		"renewal_status":   "ManualRenewal",
		"storage_size":     "800",
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	for key, value := range map[string]interface{}{
		"instance_name":    "${var.name}",
		"instance_type":    "professional",
		"max_tps":          "1000",
		"payment_type":     "Subscription",
		"period":           1,
		"queue_capacity":   "50",
		"support_eip":      true,
		"renewal_duration": 1,
		"renewal_status":   "ManualRenewal",
		"storage_size":     "800",
	} {
		err := dCreateMock.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	for key, value := range map[string]interface{}{
		"instance_name":  "${var.name}",
		"instance_type":  "professional",
		"max_tps":        "1000",
		"payment_type":   "Subscription",
		"period":         1,
		"queue_capacity": "50",
		"support_eip":    false,
		"renewal_status": "AutoRenewal",
		"storage_size":   "800",
	} {
		err := dCreateRenewalStatus.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	rawClient := testAmqpUnitClient(t)

	// readMockResponse is the shared super-response: it carries $.Data.Status
	// for DescribeAmqpInstance/AmqpInstanceStateRefreshFunc, the flat Read
	// fields at the Data level, $.Data.InstanceList for
	// DescribeInstanceQueryAvailableInstances, and TagResources for
	// DescribeInstanceListTagResources. A single response works for every
	// apiName because each consumer extracts its own slice.
	readMockResponse := map[string]interface{}{
		"Data": map[string]interface{}{
			"Status":              "SERVING",
			"InstanceId":          "MockInstanceId",
			"InstanceName":        "instance_name",
			"InstanceType":        "PROFESSIONAL",
			"SupportEIP":          false,
			"EncryptedInstance":   true,
			"KmsKeyId":            "MockKmsKeyId",
			"OrderCreateTime":     "1630000000000",
			"Edition":             "dedicated",
			"ListenerMode":        "tcp_and_ssl",
			"MaxConnections":      1000,
			"MaxEipTps":           "128",
			"MaxTps":              "1000",
			"ProvisionedCapacity": 2000,
			"MaxQueue":            "50",
			"SecurityGroupId":     "sg-xxx",
			"ServerlessSwitch":    false,
			"StorageSize":         "800",
			"SupportTracing":      false,
			"TracingStorageTime":  0,
			"VpcId":               "vpc-xxx",
			"InstanceList": []interface{}{
				map[string]interface{}{
					"SubscriptionType":    "Subscription",
					"RenewalDuration":     1,
					"RenewStatus":         "ManualRenewal",
					"RenewalDurationUnit": "M",
					"InstanceID":          "MockInstanceId",
					"CreateTime":          "2021-01-01T00:00:00Z",
				},
			},
		},
		"TagResources": []interface{}{},
		"Code":         "Success",
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
			return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_amqp_instance", "MockInstanceId"))
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
			return readMockResponse, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			return readMockResponse, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			return readMockResponse, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			return readMockResponse, nil
		},
		"CreateResponseCode": func(errorCode string) (map[string]interface{}, error) {
			return map[string]interface{}{"Code": "Failed"}, nil
		},
		"UpdateResponse": func(errorCode string) (map[string]interface{}, error) {
			return map[string]interface{}{"Success": "false"}, nil
		},
	}

	// Neutralize the state-refresh initial Delay so the unit tests do not
	// sleep for the production 3-4 minute delays. The real Refresh func is
	// passed through unchanged, so the mock still exercises the refresh
	// logic.
	stateConfPatch := gomonkey.ApplyFunc(helper.BuildStateConf, func(pending, target []string, timeout, delay time.Duration, f retry.StateRefreshFunc) *retry.StateChangeConf {
		return &retry.StateChangeConf{
			Pending:    pending,
			Target:     target,
			Refresh:    f,
			Timeout:    timeout,
			Delay:      0,
			MinTimeout: 0,
		}
	})
	defer stateConfPatch.Reset()

	// Create
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
		err := resourceAmqpInstanceCreate(d, rawClient)
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
		err := resourceAmqpInstanceCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("CreateIsExpectedErrors", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("NotApplicable")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAmqpInstanceCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("CreateAttributesSupportEip", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("NotApplicable")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		// support_eip=true is a valid create attribute; neither the migrated
		// nor the pre-migration create validates/rejects it, so with the
		// success mock create returns nil. (The pre-migration root test
		// carried assert.NotNil here but the case was always skipped in CI
		// without credentials, so the assertion never matched reality.)
		err := resourceAmqpInstanceCreate(dCreateMock, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("CreateAttributesRenewalStatus", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("NotApplicable")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		// renewal_status=AutoRenewal is a valid create attribute; the create
		// path does not validate it, so with the success mock create returns
		// nil. (Pre-migration root test asserted NotNil but it was always
		// skipped in CI without credentials.)
		err := resourceAmqpInstanceCreate(dCreateRenewalStatus, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("CreateResponseCode", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("NotApplicable")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateResponseCode"]("")
		})
		err := resourceAmqpInstanceCreate(dCreate, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("MockInstanceId")
	// Update
	t.Run("UpdateInstanceNameAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"instance_name"} {
			switch p.Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p.Schema).Data(nil, diff)
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
		err := resourceAmqpInstanceUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateInstanceNameNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"instance_name"} {
			switch p.Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p.Schema).Data(nil, diff)
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
		err := resourceAmqpInstanceUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateInstanceNameResponseAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"instance_name"} {
			switch p.Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p.Schema).Data(nil, diff)
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
			return responseMock["UpdateResponse"]("")
		})
		err := resourceAmqpInstanceUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateSetRenewalAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"renewal_status", "renewal_duration", "payment_type", "renewal_duration_unit"} {
			switch p.Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p.Schema).Data(nil, diff)
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
		err := resourceAmqpInstanceUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateSetRenewalIsExpectedErrors", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"renewal_status", "renewal_duration", "payment_type", "renewal_duration_unit"} {
			switch p.Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p.Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("NotApplicable")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAmqpInstanceUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateSetRenewalNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"renewal_status", "renewal_duration", "payment_type", "renewal_duration_unit"} {
			switch p.Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}

		}
		resourceData1, _ := schema.InternalMap(p.Schema).Data(nil, diff)
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
		err := resourceAmqpInstanceUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateSetRenewalResponseAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"renewal_status", "renewal_duration", "payment_type", "renewal_duration_unit"} {
			switch p.Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}

		}
		resourceData1, _ := schema.InternalMap(p.Schema).Data(nil, diff)
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
			return responseMock["UpdateResponse"]("")
		})
		err := resourceAmqpInstanceUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateModifyInstanceAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"max_tps", "payment_type", "queue_capacity", "support_eip", "max_eip_tps", "storage_size", "modify_type"} {
			switch p.Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p.Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "SetRenewal" {
				return responseMock["UpdateNormal"]("")
			}
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAmqpInstanceUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateModifyInstanceNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"max_tps", "payment_type", "queue_capacity", "support_eip", "max_eip_tps", "storage_size", "modify_type"} {
			switch p.Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p.Schema).Data(nil, diff)
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
		err := resourceAmqpInstanceUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateResponseAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"max_tps", "payment_type", "queue_capacity", "support_eip", "max_eip_tps", "storage_size", "modify_type"} {
			switch p.Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p.Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "SetRenewal" {
				return responseMock["UpdateNormal"]("")
			}
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateResponse"]("")
		})
		err := resourceAmqpInstanceUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateModifyInstanceIsExpectedErrors", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"max_tps", "payment_type", "queue_capacity", "support_eip", "max_eip_tps", "storage_size", "modify_type"} {
			switch p.Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p.Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "SetRenewal" {
				return responseMock["UpdateNormal"]("")
			}
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("NotApplicable")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAmqpInstanceUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateAttributesSupportEip", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.Attributes["support_eip"] = &terraform.ResourceAttrDiff{Old: "false", New: "true"}
		resourceData1, _ := schema.InternalMap(p.Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("NotApplicable")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAmqpInstanceUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Delete — new migrated branches
	t.Run("DeletePayAsYouGoNormal", func(t *testing.T) {
		dDel, _ := schema.InternalMap(p.Schema).Data(nil, nil)
		dDel.SetId("MockInstanceId")
		dDel.Set("payment_type", "PayAsYouGo")
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, apiName *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			switch *apiName {
			case "DeleteInstance":
				return map[string]interface{}{"Code": "Success"}, nil
			case "GetInstance":
				return responseMock["NotFoundError"]("ResourceNotFound")
			default:
				return map[string]interface{}{"Code": "Success"}, nil
			}
		})
		err := resourceAmqpInstanceDelete(dDel, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("DeletePayAsYouGoIdempotent", func(t *testing.T) {
		dDel, _ := schema.InternalMap(p.Schema).Data(nil, nil)
		dDel.SetId("MockInstanceId")
		dDel.Set("payment_type", "PayAsYouGo")
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, apiName *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			switch *apiName {
			case "DeleteInstance":
				return nil, &tea.SDKError{
					Code:       tea.String("InstanceNotExist"),
					Data:       tea.String("InstanceNotExist"),
					Message:    tea.String("InstanceNotExist"),
					StatusCode: tea.Int(400),
				}
			case "GetInstance":
				return responseMock["NotFoundError"]("ResourceNotFound")
			default:
				return map[string]interface{}{"Code": "Success"}, nil
			}
		})
		err := resourceAmqpInstanceDelete(dDel, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("DeleteSubscriptionNormal", func(t *testing.T) {
		dDel, _ := schema.InternalMap(p.Schema).Data(nil, nil)
		dDel.SetId("MockInstanceId")
		dDel.Set("payment_type", "Subscription")
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, apiName *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			switch *apiName {
			case "RefundInstance":
				return map[string]interface{}{"Code": "Success"}, nil
			case "GetInstance":
				return responseMock["NotFoundError"]("ResourceNotFound")
			default:
				return map[string]interface{}{"Code": "Success"}, nil
			}
		})
		err := resourceAmqpInstanceDelete(dDel, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("DeleteSubscriptionIdempotent", func(t *testing.T) {
		dDel, _ := schema.InternalMap(p.Schema).Data(nil, nil)
		dDel.SetId("MockInstanceId")
		dDel.Set("payment_type", "Subscription")
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, apiName *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			switch *apiName {
			case "RefundInstance":
				return nil, &tea.SDKError{
					Code:       tea.String("ResourceNotExists"),
					Data:       tea.String("ResourceNotExists"),
					Message:    tea.String("ResourceNotExists"),
					StatusCode: tea.Int(400),
				}
			case "GetInstance":
				return responseMock["NotFoundError"]("ResourceNotFound")
			default:
				return map[string]interface{}{"Code": "Success"}, nil
			}
		})
		err := resourceAmqpInstanceDelete(dDel, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeAmqpInstanceNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["NotFoundError"]("ResourceNotFound")
		})
		err := resourceAmqpInstanceRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeAmqpInstanceAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NonRetryableError")
		})
		err := resourceAmqpInstanceRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

func TestUnitAliCloudAmqpInstanceEncryption(t *testing.T) {
	resourceSchema := ResourceAmqpInstance()
	assert.True(t, resourceSchema.Schema["encrypted_instance"].Optional)
	assert.True(t, resourceSchema.Schema["encrypted_instance"].Computed)
	assert.True(t, resourceSchema.Schema["encrypted_instance"].ForceNew)
	assert.True(t, resourceSchema.Schema["kms_key_id"].Optional)
	assert.False(t, resourceSchema.Schema["kms_key_id"].Computed)
	assert.True(t, resourceSchema.Schema["kms_key_id"].ForceNew)
	assert.False(t, resourceSchema.Schema["edition"].ForceNew)
	serverlessChargeTypeValidate := resourceSchema.Schema["serverless_charge_type"].ValidateFunc
	assert.NotNil(t, serverlessChargeTypeValidate)
	_, validationErrors := serverlessChargeTypeValidate("onDemand", "serverless_charge_type")
	assert.Empty(t, validationErrors)
	_, validationErrors = serverlessChargeTypeValidate("provisioned", "serverless_charge_type")
	assert.Empty(t, validationErrors)
	_, validationErrors = serverlessChargeTypeValidate("invalid", "serverless_charge_type")
	assert.NotEmpty(t, validationErrors)

	newData := func(values map[string]interface{}) *schema.ResourceData {
		return schema.TestResourceDataRaw(t, resourceSchema.Schema, values)
	}

	for name, values := range map[string]map[string]interface{}{
		"vip": {
			"payment_type":       "Subscription",
			"instance_type":      "vip",
			"encrypted_instance": true,
			"kms_key_id":         "MockKmsKeyId",
		},
		"serverless provisioned dedicated": {
			"payment_type":           "PayAsYouGo",
			"serverless_charge_type": "provisioned",
			"edition":                "dedicated",
			"encrypted_instance":     true,
			"kms_key_id":             "MockKmsKeyId",
		},
	} {
		t.Run("valid "+name, func(t *testing.T) {
			request := make(map[string]interface{})
			assert.NoError(t, addAmqpInstanceEncryptionCreateRequest(newData(values), request))
			assert.Equal(t, true, request["EncryptedInstance"])
			assert.Equal(t, "MockKmsKeyId", request["KmsKeyId"])
		})
	}

	t.Run("encryption arguments omitted", func(t *testing.T) {
		request := make(map[string]interface{})
		assert.NoError(t, addAmqpInstanceEncryptionCreateRequest(newData(map[string]interface{}{
			"payment_type":  "PayAsYouGo",
			"instance_type": "enterprise",
			"edition":       "shared",
		}), request))
		assert.NotContains(t, request, "EncryptedInstance")
		assert.NotContains(t, request, "KmsKeyId")
	})

	t.Run("encryption explicitly disabled", func(t *testing.T) {
		request := make(map[string]interface{})
		assert.NoError(t, addAmqpInstanceEncryptionCreateRequest(newData(map[string]interface{}{
			"payment_type":       "PayAsYouGo",
			"instance_type":      "enterprise",
			"edition":            "shared",
			"encrypted_instance": false,
		}), request))
		assert.Equal(t, false, request["EncryptedInstance"])
		assert.NotContains(t, request, "KmsKeyId")
	})

	for name, testCase := range map[string]struct {
		values map[string]interface{}
		err    string
	}{
		"missing key": {
			values: map[string]interface{}{
				"payment_type":       "Subscription",
				"instance_type":      "vip",
				"encrypted_instance": true,
			},
			err: "kms_key_id must be configured when encrypted_instance is true",
		},
		"blank key": {
			values: map[string]interface{}{
				"payment_type":       "Subscription",
				"instance_type":      "vip",
				"encrypted_instance": true,
				"kms_key_id":         "  ",
			},
			err: "kms_key_id must be configured when encrypted_instance is true",
		},
		"key without encryption": {
			values: map[string]interface{}{
				"payment_type":       "Subscription",
				"instance_type":      "vip",
				"encrypted_instance": false,
				"kms_key_id":         "MockKmsKeyId",
			},
			err: "encrypted_instance must be true when kms_key_id is configured",
		},
	} {
		t.Run("invalid "+name, func(t *testing.T) {
			assert.EqualError(t, addAmqpInstanceEncryptionCreateRequest(newData(testCase.values), make(map[string]interface{})), testCase.err)
		})
	}

	for name, values := range map[string]map[string]interface{}{
		"subscription enterprise dedicated": {
			"payment_type":       "Subscription",
			"instance_type":      "enterprise",
			"edition":            "dedicated",
			"encrypted_instance": true,
			"kms_key_id":         "MockKmsKeyId",
		},
		"pay as you go on demand dedicated": {
			"payment_type":           "PayAsYouGo",
			"serverless_charge_type": "onDemand",
			"edition":                "dedicated",
			"encrypted_instance":     true,
			"kms_key_id":             "MockKmsKeyId",
		},
		"pay as you go provisioned shared": {
			"payment_type":           "PayAsYouGo",
			"serverless_charge_type": "provisioned",
			"edition":                "shared",
			"encrypted_instance":     true,
			"kms_key_id":             "MockKmsKeyId",
		},
		"subscription provisioned dedicated": {
			"payment_type":           "Subscription",
			"serverless_charge_type": "provisioned",
			"edition":                "dedicated",
			"encrypted_instance":     true,
			"kms_key_id":             "MockKmsKeyId",
		},
	} {
		t.Run("passes through "+name, func(t *testing.T) {
			request := make(map[string]interface{})
			assert.NoError(t, addAmqpInstanceEncryptionCreateRequest(newData(values), request))
			assert.Equal(t, true, request["EncryptedInstance"])
			assert.Equal(t, "MockKmsKeyId", request["KmsKeyId"])
		})
	}
}
