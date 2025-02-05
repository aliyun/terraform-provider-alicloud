package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

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
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
	}
	namespace, _ := resp.([]interface{})

	for _, v := range namespace {
		// item namespace
		item := v.(map[string]interface{})

		action := "/pop/v1/sam/configmap/listNamespacedConfigMaps"
		request["NamespaceId"] = StringPointer(item["NamespaceId"].(string))
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_sae_application", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.ConfigMaps", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.ConfigMaps", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				app_name := ""
				if val, exist := item["Name"]; exist {
					app_name = val.(string)
				}
				if strings.Contains(strings.ToLower(app_name), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Sae Config Map: %s (%s)", item["Name"], item["ConfigMapId"])
				continue
			}
			action := "/pop/v1/sam/configmap/configMap"
			request = map[string]*string{
				"ConfigMapId": StringPointer(fmt.Sprint(item["ConfigMapId"])),
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err = client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}

					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil && !IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) {
				return WrapError(err)
			}
			log.Printf("[INFO] Delete Sae Config Map success: %v ", item["Name"])
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
