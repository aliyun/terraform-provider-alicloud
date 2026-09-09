// Package alicloud. Hand-written acceptance tests for the CMS MemoryStoreAPIKey resource.
package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// TestAccAliCloudCmsMemoryStoreAPIKey_basic verifies the full lifecycle of the
// alicloud_cms_memory_store_api_key resource: create -> read -> import -> destroy.
// The API key resource is a child of a CMS Workspace + MemoryStore. Since
// alicloud_cms_memory_store is not yet a Terraform resource, the Workspace and
// MemoryStore prerequisites are created via the CMS RESTful API in PreConfig
// using a client constructed from the standard ALICLOUD_* environment variables.
func TestAccAliCloudCmsMemoryStoreAPIKey_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_memory_store_api_key.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsMemoryStoreAPIKeyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsMemoryStoreAPIKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsMemoryStoreAPIKeyBasicDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rac.checkResourceDestroy(),
			func(s *terraform.State) error {
				client, err := sharedClientForRegion(defaultRegionToTest)
				if err != nil {
					log.Printf("[WARN] Failed to create client for cleanup: %s", err)
					return nil
				}
				ac := client.(*connectivity.AliyunClient)
				if err := deleteCmsMemoryStoreForTest(ac, name, name); err != nil {
					log.Printf("[WARN] Failed to clean up memory store %s during destroy: %s", name, err)
				}
				if err := deleteCmsWorkspaceForTest(ac, name); err != nil {
					log.Printf("[WARN] Failed to clean up workspace %s during destroy: %s", name, err)
				}
				return nil
			},
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					client, err := sharedClientForRegion(defaultRegionToTest)
					if err != nil {
						t.Fatalf("failed to create test client: %v", err)
					}
					ac := client.(*connectivity.AliyunClient)
					if err := createCmsWorkspaceForTest(ac, name, name); err != nil {
						t.Fatalf("failed to create workspace for test: %v", err)
					}
					if err := createCmsMemoryStoreForTest(ac, name, name); err != nil {
						t.Fatalf("failed to create memory store for test: %v", err)
					}
				},
				Config: testAccConfig(map[string]interface{}{
					"workspace":         name,
					"memory_store_name": name,
					"name":              name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":         name,
						"memory_store_name": name,
						"name":              name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// createCmsWorkspaceForTest creates a CMS Workspace via the RESTful API for test setup.
func createCmsWorkspaceForTest(client *connectivity.AliyunClient, workspaceName, slsProject string) error {
	action := fmt.Sprintf("/workspace/%s", workspaceName)
	query := make(map[string]*string)
	body := map[string]interface{}{
		"slsProject": slsProject,
	}
	_, err := client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
	if err != nil {
		return fmt.Errorf("failed to create workspace %s for test: %w", workspaceName, err)
	}
	return nil
}

// deleteCmsWorkspaceForTest cleans up a CMS Workspace created via the API.
func deleteCmsWorkspaceForTest(client *connectivity.AliyunClient, workspaceName string) error {
	action := fmt.Sprintf("/workspace/%s", workspaceName)
	query := make(map[string]*string)
	_, err := client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"WorkspaceNotExist"}) {
			return nil
		}
		return fmt.Errorf("failed to delete workspace %s for test: %w", workspaceName, err)
	}
	return nil
}

// createCmsMemoryStoreForTest creates a CMS MemoryStore via the RESTful API for test setup.
func createCmsMemoryStoreForTest(client *connectivity.AliyunClient, workspace, memoryStoreName string) error {
	action := fmt.Sprintf("/workspace/%s/memorystore", workspace)
	query := make(map[string]*string)
	body := map[string]interface{}{
		"memoryStoreName": memoryStoreName,
		"shortTermTtl":    3600,
		"sourceType":      "None",
	}
	_, err := client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
	if err != nil {
		return fmt.Errorf("failed to create memory store %s for test: %w", memoryStoreName, err)
	}
	return nil
}

// deleteCmsMemoryStoreForTest cleans up a CMS MemoryStore created via the API.
func deleteCmsMemoryStoreForTest(client *connectivity.AliyunClient, workspace, memoryStoreName string) error {
	action := fmt.Sprintf("/workspace/%s/memorystore/%s", workspace, memoryStoreName)
	query := make(map[string]*string)
	_, err := client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"MemoryStoreNotExist", "WorkspaceNotExist"}) {
			return nil
		}
		return fmt.Errorf("failed to delete memory store %s for test: %w", memoryStoreName, err)
	}
	return nil
}

func AliCloudCmsMemoryStoreAPIKeyBasicDependence(name string) string {
	return fmt.Sprintf(`variable "name" {
  default = "%s"
}
`, name)
}

var AliCloudCmsMemoryStoreAPIKeyMap = map[string]string{
	"api_key":     CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}
