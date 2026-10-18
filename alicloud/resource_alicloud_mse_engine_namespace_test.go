package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_mse_engine_namespace", &resource.Sweeper{
		Name: "alicloud_mse_engine_namespace",
		F:    testSweepMseEngineNamespace,
	})
}

// testSweepMseEngineNamespace sweeps leftover engine namespaces from every MSE
// cluster in the region. Engine namespaces are children of a cluster, so the
// sweeper walks every cluster, lists its namespaces via ListEngineNamespaces
// (the same API the provider Read uses) and deletes each one whose name passes
// the NODELETING / prefix guard via DeleteEngineNamespace (mirroring
// resourceAlicloudMseEngineNamespaceDelete). Errors per namespace are logged
// and skipped; the sweeper always returns nil.
func testSweepMseEngineNamespace(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	clusters, err := listMseSweepClusters(client)
	if err != nil {
		log.Printf("[ERROR] ListClusters for mse_engine_namespace sweep failed: %s", err)
		return nil
	}

	for _, cluster := range clusters {
		clusterName := fmt.Sprint(cluster["ClusterAliasName"])
		instanceId := fmt.Sprint(cluster["InstanceId"])
		if !mseSweepNameAllowed(clusterName) {
			log.Printf("[INFO] Skipping engine namespaces of cluster %s (%s): name not allowed", clusterName, instanceId)
			continue
		}
		if err := sweepMseEngineNamespacesForCluster(client, instanceId); err != nil {
			log.Printf("[ERROR] Listing engine namespaces for cluster %s failed: %s", instanceId, err)
			continue
		}
	}
	return nil
}

func sweepMseEngineNamespacesForCluster(client *connectivity.AliyunClient, instanceId string) error {
	action := "ListEngineNamespaces"
	request := map[string]interface{}{
		"InstanceId": instanceId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var response map[string]interface{}
	var err error
	if err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return err
	}
	addDebug(action, response, request)
	// The cluster may already be deleted; the API surfaces that as success=false.
	if fmt.Sprint(response["Success"]) == "false" {
		return nil
	}
	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return err
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := fmt.Sprint(item["NamespaceShowName"])
		namespaceId := fmt.Sprint(item["Namespace"])
		if !mseSweepNameAllowed(name) {
			log.Printf("[INFO] Skipping mse engine namespace: %s (%s)", name, namespaceId)
			continue
		}
		deleteAction := "DeleteEngineNamespace"
		deleteRequest := map[string]interface{}{
			"Id":         namespaceId,
			"InstanceId": instanceId,
		}
		if _, derr := client.RpcPost("mse", "2019-05-31", deleteAction, nil, deleteRequest, false); derr != nil {
			log.Printf("[ERROR] Failed to delete mse engine namespace %s:%s: %s", instanceId, namespaceId, derr)
			continue
		}
		log.Printf("[INFO] Delete mse engine namespace success: %s:%s", instanceId, namespaceId)
	}
	return nil
}

func TestAccAliCloudMSEEngineNamespace_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_engine_namespace.default"
	checkoutSupportedRegions(t, true, connectivity.MSESupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudMSEEngineNamespaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMseEngineNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-mseenginenamespace%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMSEEngineNamespaceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":         "${alicloud_mse_cluster.default.id}",
					"cluster_id":          "${alicloud_mse_cluster.default.cluster_id}",
					"namespace_show_name": "${var.name}",
					"namespace_id":        "${var.name}",
					"accept_language":     "zh",
					"namespace_desc":      "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":          CHECKSET,
						"namespace_show_name": name,
						"namespace_id":        name,
						"namespace_desc":      "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_show_name": "${var.name}_update",
					"namespace_desc":      "test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_show_name": name + "_update",
						"namespace_desc":      "test_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accept_language"},
			},
		},
	})
}

var AlicloudMSEEngineNamespaceMap0 = map[string]string{}

func AlicloudMSEEngineNamespaceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
 available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
 vpc_name   = "default"
 cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
 vswitch_name = "default"
 cidr_block   = "172.17.3.0/24"
 vpc_id       = alicloud_vpc.default.id
 zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_mse_cluster" "default" {
 cluster_specification = "MSE_SC_1_2_60_c"
 cluster_type          = "Nacos-Ans"
 cluster_version       = "NACOS_2_0_0"
 instance_count        = 3
 net_type              = "privatenet"
 vswitch_id            = alicloud_vswitch.default.id
 connection_type       = "slb"
 pub_network_flow      = "1"
 mse_version           = "mse_pro"
 vpc_id                = alicloud_vpc.default.id
}

data "alicloud_mse_clusters" "default" {
	name_regex = "tftest"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudMSEEngineNamespace(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_mse_engine_namespace"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_mse_engine_namespace"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"instance_id":         "CreateMseEngineNamespaceValue",
		"namespace_show_name": "CreateMseEngineNamespaceValue",
		"namespace_id":        "CreateMseEngineNamespaceValue",
		"accept_language":     "CreateMseEngineNamespaceValue",
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
		"Data": []interface{}{
			map[string]interface{}{
				"Type":              0,
				"Quota":             10000,
				"ConfigCount":       0,
				"NamespaceShowName": "CreateMseEngineNamespaceValue",
				"Namespace":         "CreateMseEngineNamespaceValue",
				"ServiceCount":      0,
			},
		},
	}
	CreateMockResponse := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_mse_engine_namespace", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewMseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudMseEngineNamespaceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateEngineNamespace" {
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
		err := resourceAlicloudMseEngineNamespaceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mse_engine_namespace"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewMseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudMseEngineNamespaceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"namespace_show_name": "UpdateMseEngineNamespaceValue",
		"accept_language":     "UpdateMseEngineNamespaceValue",
	}
	diff, err := newInstanceDiff("alicloud_mse_engine_namespace", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mse_engine_namespace"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Data": []interface{}{
			map[string]interface{}{
				"Type":              0,
				"Quota":             10000,
				"ConfigCount":       0,
				"NamespaceShowName": "UpdateMseEngineNamespaceValue",
				"Namespace":         "CreateMseEngineNamespaceValue",
				"ServiceCount":      0,
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateEngineNamespace" {
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
		err := resourceAlicloudMseEngineNamespaceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mse_engine_namespace"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_mse_engine_namespace", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mse_engine_namespace"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ListEngineNamespaces" {
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
		err := resourceAlicloudMseEngineNamespaceRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewMseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudMseEngineNamespaceDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_mse_engine_namespace", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mse_engine_namespace"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteEngineNamespace" {
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
		err := resourceAlicloudMseEngineNamespaceDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
