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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_sae_config_map",
		&resource.Sweeper{
			Name: "alicloud_sae_config_map",
			F:    testSweepSaeConfigMap,
		})
}

func testSweepSaeConfigMap(region string) error {
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
		// empty-list guard: treat missing/malformed namespace list as empty
		log.Printf("[WARN] %s: $.Data not found or malformed, treating as empty list", action)
		return nil
	}
	namespace, _ := resp.([]interface{})

	for _, v := range namespace {
		// item namespace
		item := v.(map[string]interface{})

		action := "/pop/v1/sam/configmap/listNamespacedConfigMaps"
		request["NamespaceId"] = StringPointer(item["NamespaceId"].(string))
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_sae_config_map", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.ConfigMaps", response)
		if err != nil {
			// empty-list guard: namespace has no config maps → continue to next namespace
			continue
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			cmName := ""
			if val, exist := item["Name"]; exist {
				cmName = val.(string)
			}
			// red-line #2: skip NODELETING + do-not-delete markers
			if protectedSweepVpcName(cmName) {
				log.Printf("[INFO] Skipping Sae Config Map (NODELETING): %s (%s)", cmName, item["ConfigMapId"])
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.Contains(strings.ToLower(cmName), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Sae Config Map: %s (%s)", cmName, item["ConfigMapId"])
				continue
			}
			action := "/pop/v1/sam/configmap/configMap"
			request = map[string]*string{
				"ConfigMapId": StringPointer(fmt.Sprint(item["ConfigMapId"])),
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
			if err != nil && !IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) {
				return WrapError(err)
			}
			log.Printf("[INFO] Delete Sae Config Map success: %v", cmName)
		}
	}
	return nil
}

func TestAccAlicloudSAEConfigMap_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_config_map.default"
	ra := resourceAttrInit(resourceId, AlicloudSAEConfigMapMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeConfigMap")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%ssaeconfigmap%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAEConfigMapBasicDependence0)
	resource.Test(t, resource.TestCase{
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_id": "${alicloud_sae_namespace.default.namespace_id}",
					"name":         "tftestaccname",
					"data":         `{\"env.home\":\"/root\",\"envtest.shell\":\"/bin/sh\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_id": CHECKSET,
						"name":         "tftestaccname",
						"data":         "{\"env.home\":\"/root\",\"envtest.shell\":\"/bin/sh\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testaccdescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testaccdescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data": `{\"env.home\":\"/root\",\"env.shell\":\"/bin/sh\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data": "{\"env.home\":\"/root\",\"env.shell\":\"/bin/sh\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testAccDesc",
					"data":        `{\"env.home\":\"/root\",\"envtest.shell\":\"/bin/sh\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAccDesc",
						"data":        "{\"env.home\":\"/root\",\"envtest.shell\":\"/bin/sh\"}",
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

var AlicloudSAEConfigMapMap0 = map[string]string{
	"namespace_id": CHECKSET,
	"name":         CHECKSET,
}

func AlicloudSAEConfigMapBasicDependence0(name string) string {
	rand := acctest.RandIntRange(1, 100)
	return fmt.Sprintf(` 
resource "alicloud_sae_namespace" "default" {
  namespace_description = "namespace_desc"
  namespace_id = "%s:configmaptest%d"
  namespace_name = "namespace_name"
}

variable "name" {
  default = "%s"
}
`, defaultRegionToTest, rand, name)
}
