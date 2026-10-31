package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_sae_namespace",
		&resource.Sweeper{
			Name:         "alicloud_sae_namespace",
			F:            testSweepSaeNamespace,
			// namespace delete fails if it still has apps or config maps;
			// declare dependencies so the app + config_map sweepers run first (children before parent)
			Dependencies: []string{"alicloud_sae_application", "alicloud_sae_config_map"},
		})
}

func testSweepSaeNamespace(region string) error {
	prefixes := []string{
		"tftestacc",
	}
	rawClient, err := sharedClientForRegionWithBackendRegions(region, true, connectivity.SaeSupportRegions)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	var response map[string]interface{}
	request := make(map[string]*string)

	request["ContainCustom"] = StringPointer(strconv.FormatBool(true))
	action := "/pop/v1/sam/namespace/describeNamespaceList"
	response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
	if err != nil {
		log.Printf("[ERROR] %s got an error: %s", action, err)
		return nil
	}
	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		// empty-list guard: treat missing/malformed list key as empty (no namespaces to sweep)
		log.Printf("[WARN] %s: $.Data not found or malformed, treating as empty list", action)
		return nil
	}
	namespace, _ := resp.([]interface{})

	for _, v := range namespace {
		// item namespace
		item := v.(map[string]interface{})

		nsName := ""
		if val, exist := item["NamespaceName"]; exist {
			nsName = val.(string)
		}
		// red-line #2: skip NODELETING + do-not-delete markers (both sweepAll + normal mode)
		if protectedSweepVpcName(nsName) {
			log.Printf("[INFO] Skipping Sae Namespace (NODELETING): %s", nsName)
			continue
		}
		skip := true
		for _, prefix := range prefixes {
			if strings.Contains(strings.ToLower(nsName), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Sae Namespace: %s", nsName)
			continue
		}

		log.Printf("[INFO] Deleting Sae Namespace: %s (NamespaceId: %v)", nsName, item["NamespaceId"])
		action := "/pop/v1/paas/namespace"
		request = map[string]*string{
			"NamespaceId": StringPointer(fmt.Sprint(item["NamespaceId"])),
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(10*time.Minute, func() *retry.RetryError {
			response, err = client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}

				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil && !IsExpectedErrors(err, []string{
			"Namespace.AppExists",
			// namespace still has config maps (config_map sweeper should have run first via Dependencies,
			// but non-test config maps may remain) — record as leftover, do not fail the sweep
			"sae.errorcode.namespace.has.configmap.when.delete.message",
		}) {
			return WrapError(err)
		}
		if err != nil {
			log.Printf("[WARN] Sae Namespace not empty (apps/configmaps exist), skipping: %v", nsName)
		} else {
			log.Printf("[INFO] Delete Sae Namespace success: %v", nsName)
		}
	}
	return nil
}

func TestAccAlicloudSAENamespace_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_namespace.default"
	ra := resourceAttrInit(resourceId, AlicloudSAENamespaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%ssaenamespace%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAENamespaceBasicDependence0)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name":            name,
					"namespace_id":              fmt.Sprintf("%s:tftest%d", defaultRegionToTest, rand),
					"namespace_description":     "tftestaccdescription",
					"enable_micro_registration": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":            name,
						"namespace_id":              fmt.Sprintf("%s:tftest%d", defaultRegionToTest, rand),
						"namespace_description":     "tftestaccdescription",
						"enable_micro_registration": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":        name + "update",
						"namespace_description": "tftestaccdescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_description": "tftestaccdescriptionupdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":        name + "update",
						"namespace_description": "tftestaccdescriptionupdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name":        name + "updateall",
					"namespace_description": "tftestaccdescriptionupdateall",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":        name + "updateall",
						"namespace_description": "tftestaccdescriptionupdateall",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_micro_registration": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_micro_registration": "true",
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

func TestAccAlicloudSAENamespace_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_namespace.default"
	ra := resourceAttrInit(resourceId, AlicloudSAENamespaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%ssaenamespace%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAENamespaceBasicDependence0)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name":            name,
					"namespace_short_id":        fmt.Sprintf("tftest%d", rand),
					"namespace_description":     "tftestaccdescription",
					"enable_micro_registration": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":            name,
						"namespace_short_id":        CHECKSET,
						"namespace_description":     "tftestaccdescription",
						"enable_micro_registration": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":        name + "update",
						"namespace_description": "tftestaccdescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_description": "tftestaccdescriptionupdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":        name + "update",
						"namespace_description": "tftestaccdescriptionupdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name":        name + "updateall",
					"namespace_description": "tftestaccdescriptionupdateall",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":        name + "updateall",
						"namespace_description": "tftestaccdescriptionupdateall",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_micro_registration": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_micro_registration": "true",
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

var AlicloudSAENamespaceMap0 = map[string]string{
	"namespace_id":              CHECKSET,
	"namespace_short_id":        CHECKSET,
	"enable_micro_registration": CHECKSET,
}

func AlicloudSAENamespaceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}
`, name)
}
