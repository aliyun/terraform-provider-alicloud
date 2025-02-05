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
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
	}
	namespace, _ := resp.([]interface{})

	for _, v := range namespace {
		// item namespace
		item := v.(map[string]interface{})

		skip := true
		for _, prefix := range prefixes {
			app_name := ""
			if val, exist := item["NamespaceName"]; exist {
				app_name = val.(string)
			}
			if strings.Contains(strings.ToLower(app_name), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Sae Namespace: %s", item["NamespaceName"])
			continue
		}

		action := "/pop/v1/paas/namespace"
		request = map[string]*string{
			"NamespaceId": StringPointer(fmt.Sprint(item["NamespaceId"])),
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
		if err != nil && !IsExpectedErrors(err, []string{"Namespace.AppExists"}) {
			return WrapError(err)
		}
		log.Printf("[INFO] Delete Sae Namespace success: %v ", item["NamespaceName"])
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
