package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// A GPDB ApiKey requires a pre-existing GPDB workspace, but there is no
// Terraform resource for a GPDB workspace, so it must be provisioned out of
// band before the test case is built. It cannot be created inside a step's
// PreConfig: the SDK configures the provider lazily on the first apply, so
// testAccProvider.Meta() is still nil when the first step's PreConfig runs,
// and the step Config string is rendered eagerly (before any PreConfig), so a
// workspace id captured in PreConfig would never reach the config. We therefore
// build a client from the acceptance credentials (sharedClientForRegion) and
// create the workspace up front, guarded by TF_ACC so non-acceptance runs
// (plain `go test`, vet, CI compile) never touch the API.
func gpdbApiKeyTestCreateWorkspace(t *testing.T, rand int) (*connectivity.AliyunClient, string) {
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Fatalf("failed to get AliCloud client for gpdb workspace: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"WorkspaceName": fmt.Sprintf("tfacc-gpdb-ws-%d", rand),
		"RegionId":      client.RegionId,
	}
	response, err := client.RpcPost("gpdb", "2016-05-03", "CreateWorkspace", nil, request, true)
	if err != nil {
		t.Fatalf("failed to create gpdb workspace: %s", err)
	}
	return client, fmt.Sprint(response["WorkspaceId"])
}

func gpdbApiKeyTestDeleteWorkspace(client *connectivity.AliyunClient, workspaceId string) {
	if client == nil || workspaceId == "" {
		return
	}
	request := map[string]interface{}{
		"WorkspaceId": workspaceId,
		"RegionId":    client.RegionId,
	}
	_, _ = client.RpcPost("gpdb", "2016-05-03", "DeleteWorkspace", nil, request, true)
}

// Test Gpdb ApiKey. >>> Resource test cases, automatically generated.
// Case: basic test
func TestAccAliCloudGpdbApiKey_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_api_key.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbApiKeyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbApiKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgpdbapikey%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbApiKeyBasicDependence)

	var workspaceId string
	var wsClient *connectivity.AliyunClient
	if os.Getenv("TF_ACC") != "" {
		wsClient, workspaceId = gpdbApiKeyTestCreateWorkspace(t, rand)
		defer gpdbApiKeyTestDeleteWorkspace(wsClient, workspaceId)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace_id": workspaceId,
					"key_name":     name,
					"description":  "terraform test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace_id": workspaceId,
						"key_name":     name,
						"description":  "terraform test",
						"key_id":       CHECKSET,
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

// Case: twin test with service_ids
func TestAccAliCloudGpdbApiKey_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_api_key.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbApiKeyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbApiKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgpdbapikey%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbApiKeyBasicDependence)

	var workspaceId string
	var wsClient *connectivity.AliyunClient
	if os.Getenv("TF_ACC") != "" {
		wsClient, workspaceId = gpdbApiKeyTestCreateWorkspace(t, rand)
		defer gpdbApiKeyTestDeleteWorkspace(wsClient, workspaceId)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace_id": workspaceId,
					"key_name":     name,
					"description":  "terraform test twin",
					"service_ids":  []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace_id":  workspaceId,
						"key_name":      name,
						"description":   "terraform test twin",
						"key_id":        CHECKSET,
						"service_ids.#": "0",
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

var AliCloudGpdbApiKeyMap = map[string]string{
	"key_id":       CHECKSET,
	"key_name":     CHECKSET,
	"description":  CHECKSET,
	"workspace_id": CHECKSET,
}

func AliCloudGpdbApiKeyBasicDependence(name string) string {
	return ""
}

// Test Gpdb ApiKey. <<< Resource test cases, automatically generated.
