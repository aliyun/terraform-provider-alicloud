package alicloud

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms ContextStoreApiKey. >>> Resource test cases.
func TestAccAliCloudCmsContextStoreApiKey_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_context_store_api_key.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsContextStoreApiKeyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsContextStoreApiKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccsak%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsContextStoreApiKeyBasicDependence0)
	if os.Getenv("TF_ACC") != "" {
		testAccCmsContextStorePrepare(t, name)
	}
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
					"workspace":          name,
					"context_store_name": name,
					"name":               name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":          name,
						"context_store_name": name,
						"name":               name,
						"api_key":            CHECKSET,
						"create_time":        CHECKSET,
						"region_id":          CHECKSET,
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

func TestAccAliCloudCmsContextStoreApiKey_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_context_store_api_key.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsContextStoreApiKeyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsContextStoreApiKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccsak%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsContextStoreApiKeyBasicDependence0)
	if os.Getenv("TF_ACC") != "" {
		testAccCmsContextStorePrepare(t, name)
	}
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
					"workspace":          name,
					"context_store_name": name,
					"name":               fmt.Sprintf("%s-key", name),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":          name,
						"context_store_name": name,
						"name":               fmt.Sprintf("%s-key", name),
						"api_key":            CHECKSET,
						"create_time":        CHECKSET,
						"region_id":          CHECKSET,
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

var AliCloudCmsContextStoreApiKeyMap0 = map[string]string{
	"api_key":     CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudCmsContextStoreApiKeyBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

// testAccCmsContextStorePrepare pre-creates the parent resources (an SLS project, a Cms
// workspace binding it and a context store) through raw API calls, because the parent
// context store is not managed by the provider yet. All of them are removed on cleanup.
func testAccCmsContextStorePrepare(t *testing.T, name string) {
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Fatalf("Error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	retryAll := func(timeout time.Duration, call func() error) error {
		var lastErr error
		e := resource.Retry(timeout, func() *resource.RetryError {
			lastErr = call()
			if lastErr != nil {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(lastErr)
			}
			return nil
		})
		if e != nil && lastErr != nil {
			return lastErr
		}
		return e
	}

	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	slsBody := map[string]interface{}{
		"projectName": name,
		"description": name,
	}
	err = retryAll(3*time.Minute, func() error {
		_, e := client.Do("Sls", roaParam("POST", "2020-12-30", "CreateProject", "/"), query, slsBody, nil, hostMap, false)
		return e
	})
	if err != nil {
		t.Fatalf("Error creating SLS project %s: %s", name, err)
	}

	wsAction := fmt.Sprintf("/workspace/%s", name)
	wsBody := map[string]interface{}{
		"slsProject":  name,
		"description": name,
	}
	err = retryAll(3*time.Minute, func() error {
		_, e := client.RoaPost("Cms", "2024-03-30", wsAction, make(map[string]*string), nil, wsBody, true)
		return e
	})
	if err != nil {
		t.Fatalf("Error creating Cms workspace %s: %s", name, err)
	}

	csAction := fmt.Sprintf("/workspace/%s/contextstore", name)
	csBody := map[string]interface{}{
		"contextStoreName": name,
		"contextType":      "memory",
		"description":      name,
	}
	err = retryAll(3*time.Minute, func() error {
		_, e := client.RoaPost("Cms", "2024-03-30", csAction, make(map[string]*string), nil, csBody, true)
		return e
	})
	if err != nil {
		t.Fatalf("Error creating Cms context store %s: %s", name, err)
	}

	t.Cleanup(func() {
		csDelAction := fmt.Sprintf("/workspace/%s/contextstore/%s", name, name)
		_, _ = client.RoaDelete("Cms", "2024-03-30", csDelAction, make(map[string]*string), nil, nil, true)
		_, _ = client.RoaDelete("Cms", "2024-03-30", wsAction, make(map[string]*string), nil, nil, true)
		delHostMap := make(map[string]*string)
		delHostMap["project"] = StringPointer(name)
		_, _ = client.Do("Sls", roaParam("DELETE", "2020-12-30", "DeleteProject", "/"), query, nil, nil, delHostMap, false)
	})
}

// Test Cms ContextStoreApiKey. <<< Resource test cases.
