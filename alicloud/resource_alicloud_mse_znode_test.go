package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
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
	resource.AddTestSweepers("alicloud_mse_znode", &resource.Sweeper{
		Name: "alicloud_mse_znode",
		F:    testSweepMseZnode,
	})
}

// testSweepMseZnode sweeps leftover ZooKeeper znodes across every MSE ZooKeeper
// cluster in the region. A znode is identified by (ClusterId, Path). The sweeper
// walks every cluster, traverses the znode tree from "/" via ListZnodeChildren
// (the data source API, bounded by depth and node count), and deletes each
// matching leaf via DeleteZnode mirroring resourceAlicloudMseZnodeDelete.
// Directory znodes are traversed but deletion of a non-empty directory is
// expected to fail server-side and is logged (non-fatal). The sweeper always
// returns nil.
func testSweepMseZnode(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	clusters, err := listMseSweepClusters(client)
	if err != nil {
		log.Printf("[ERROR] ListClusters for mse_znode sweep failed: %s", err)
		return nil
	}

	for _, cluster := range clusters {
		clusterName := fmt.Sprint(cluster["ClusterAliasName"])
		clusterId := fmt.Sprint(cluster["ClusterId"])
		if clusterId == "" {
			clusterId = fmt.Sprint(cluster["InstanceId"])
		}
		if !mseSweepNameAllowed(clusterName) {
			log.Printf("[INFO] Skipping znodes of cluster %s (%s): name not allowed", clusterName, clusterId)
			continue
		}
		if err := sweepMseZnodesForCluster(client, clusterId); err != nil {
			log.Printf("[ERROR] Znode traversal for cluster %s failed: %s", clusterId, err)
			continue
		}
	}
	return nil
}

const (
	mseZnodeSweepMaxDepth = 6
	mseZnodeSweepMaxNodes = 5000
)

func sweepMseZnodesForCluster(client *connectivity.AliyunClient, clusterId string) error {
	type queueItem struct {
		path  string
		depth int
	}
	queue := []queueItem{{path: "/", depth: 0}}
	visited := 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		visited++
		if visited > mseZnodeSweepMaxNodes {
			log.Printf("[WARN] mse znode sweep capped at %d nodes for cluster %s; stopping traversal", mseZnodeSweepMaxNodes, clusterId)
			break
		}
		children, err := listMseSweepZnodeChildren(client, clusterId, cur.path)
		if err != nil {
			log.Printf("[ERROR] ListZnodeChildren for cluster %s path %s failed: %s", clusterId, cur.path, err)
			continue
		}
		for _, child := range children {
			childPath := fmt.Sprint(child["Path"])
			isDir := false
			if d, ok := child["Dir"]; ok {
				isDir = fmt.Sprint(d) == "true"
			}
			name := strings.TrimPrefix(childPath, "/")
			if !mseSweepNameAllowed(name) {
				log.Printf("[INFO] Skipping mse znode: %s:%s", clusterId, childPath)
			} else {
				deleteAction := "DeleteZnode"
				deleteRequest := map[string]interface{}{
					"ClusterId": clusterId,
					"Path":      childPath,
				}
				if _, derr := client.RpcPost("mse", "2019-05-31", deleteAction, nil, deleteRequest, false); derr != nil {
					log.Printf("[ERROR] Failed to delete mse znode %s:%s: %s", clusterId, childPath, derr)
				} else {
					log.Printf("[INFO] Delete mse znode success: %s:%s", clusterId, childPath)
				}
			}
			if isDir && cur.depth+1 < mseZnodeSweepMaxDepth && childPath != "" && childPath != "/" {
				queue = append(queue, queueItem{path: childPath, depth: cur.depth + 1})
			}
		}
	}
	return nil
}

func listMseSweepZnodeChildren(client *connectivity.AliyunClient, clusterId, path string) ([]map[string]interface{}, error) {
	action := "ListZnodeChildren"
	request := map[string]interface{}{
		"ClusterId": clusterId,
		"Path":      path,
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
		return nil, err
	}
	addDebug(action, response, request)
	if fmt.Sprint(response["Success"]) == "false" {
		return nil, nil
	}
	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return nil, err
	}
	result, _ := resp.([]interface{})
	out := make([]map[string]interface{}, 0, len(result))
	for _, v := range result {
		out = append(out, v.(map[string]interface{}))
	}
	return out, nil
}

func TestAccAlicloudMSEZnode_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_znode.default"
	checkoutSupportedRegions(t, true, connectivity.MSESupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudMSEZnodeMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMseZnode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("/tf-testaccmseznode%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMSEZnodeBasicDependence0)
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
					"path":       "${var.name}",
					"data":       "${var.name}",
					"cluster_id": "${alicloud_mse_cluster.default.cluster_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"path":       name,
						"data":       name,
						"cluster_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data": name,
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

var AlicloudMSEZnodeMap0 = map[string]string{
	"path":       CHECKSET,
	"cluster_id": CHECKSET,
}

func AlicloudMSEZnodeBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = "${data.alicloud_mongodb_zones.default.zones.0.id}"
}
resource "alicloud_mse_cluster" "default" {
  cluster_specification = "MSE_SC_1_2_60_c"
  cluster_type = "ZooKeeper"
  cluster_version = "ZooKeeper_3_8_0"
  instance_count = 1
  net_type = "privatenet"
  vswitch_id = data.alicloud_vswitches.default.ids.0
  pub_network_flow = "1"
  cluster_alias_name= "tf-testaccmseznode"
  mse_version = "mse_dev"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudMSEZnode(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_mse_znode"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_mse_znode"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"path":            "CreateMseZnodeValue",
		"data":            "CreateMseZnodeValue",
		"cluster_id":      "CreateMseZnodeValue",
		"accept_language": "CreateMseZnodeValue",
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
		"Data": map[string]interface{}{
			"Path": "CreateMseZnodeValue",
			"Data": "CreateMseZnodeValue",
			"Dir":  true,
			"Name": "CreateMseZnodeValue",
		},
		"Success": true,
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_mse_znode", errorCode))
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
	err = resourceAlicloudMseZnodeCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateZnode" {
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
		err := resourceAlicloudMseZnodeCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mse_znode"].Schema).Data(dInit.State(), nil)
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
	err = resourceAlicloudMseZnodeUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"data": "UpdateZnodeValue",
	}
	diff, err := newInstanceDiff("alicloud_mse_znode", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mse_znode"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Data": map[string]interface{}{
			"Data": "UpdateZnodeValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateZnode" {
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
		err := resourceAlicloudMseZnodeUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mse_znode"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_mse_znode", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mse_znode"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "QueryZnodeDetail" {
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
		err := resourceAlicloudMseZnodeRead(dExisted, rawClient)
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
	err = resourceAlicloudMseZnodeDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_mse_znode", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mse_znode"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteZnode" {
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
		err := resourceAlicloudMseZnodeDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
