package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_dcdn_kv_namespace",
		&resource.Sweeper{
			Name: "alicloud_dcdn_kv_namespace",
			F:    testSweepDcdnKvNamespace,
		})
}

func testSweepDcdnKvNamespace(region string) error {
	// DCDN KV is a global/account-scoped service — no region gate (cn-beijing must run).
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	// RED LINE #5: prefix guard always applied (even under SWEEP_ALL=1).
	prefixes := []string{
		"tf-",
		"tf-testacc",
		"tf_testacc",
		"test-",
		"testacc",
		"test_acc",
	}

	action := "DescribeDcdnKvAccount"
	request := map[string]interface{}{}

	var response map[string]interface{}
	listWait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				listWait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		// KV account may not be activated in the sweep account — non-fatal.
		log.Printf("[ERROR] %s got an error: %v", action, err)
		return nil
	}

	resp, err := jsonpath.Get("$.NamespaceList", response)
	if err != nil {
		log.Printf("[ERROR] Getting %s attribute by path %s failed: %v", action, "$.NamespaceList", err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := fmt.Sprint(item["Namespace"])

		// RED LINE #2: NODELETING / do-not-delete marker protection.
		if protectedSweepVpcName(name) {
			log.Printf("[INFO] Skipping Dcdn Kv Namespace (NODELETING): %s", name)
			continue
		}
		// RED LINE #5: prefix guard always applied.
		matches := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), prefix) {
				matches = true
				break
			}
		}
		if !matches {
			log.Printf("[INFO] Skipping Dcdn Kv Namespace (no test prefix): %s", name)
			continue
		}

		deleteAction := "DeleteDcdnKvNamespace"
		deleteRequest := map[string]interface{}{
			"Namespace": name,
		}
		log.Printf("[INFO] Deleting Dcdn Kv Namespace: %s", name)
		deleteWait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("dcdn", "2018-01-15", deleteAction, nil, deleteRequest, false)
			if err != nil {
				if NeedRetry(err) {
					deleteWait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(deleteAction, response, deleteRequest)
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Dcdn Kv Namespace (%s): %v", name, err)
			continue
		}
		log.Printf("[INFO] Deleted Dcdn Kv Namespace: %s", name)
	}

	return nil
}

// Case 1
func TestAccAliCloudDcdnKvNamespace_basic2278(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_kv_namespace.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnKvNamespaceMap2278)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnKvNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccDcdnKvNamespace%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnKvNamespaceBasicDependence2278)
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
					"description": "wkmtest",
					"namespace":   "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "wkmtest",
						"namespace":   name,
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudDcdnKvNamespaceMap2278 = map[string]string{}

func AlicloudDcdnKvNamespaceBasicDependence2278(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_dcdn_kv_account" "current" {
  status = "online"
}

`, name)
}
