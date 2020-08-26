package alicloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_eci_image_cache",
		&resource.Sweeper{
			Name: "alicloud_eci_image_cache",
			F:    testSweepEciImageCache,
		})
}

func testSweepEciImageCache(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	queryRequest := eci.CreateDescribeImageCachesRequest()
	var allCaches []eci.DescribeImageCachesImageCache0

	raw, err := client.WithEciClient(func(eciClient *eci.Client) (interface{}, error) {
		return eciClient.DescribeImageCaches(queryRequest)
	})
	if err != nil {
		log.Printf("[ERROR] %s get an error %#v", queryRequest.GetActionName(), err)
	}
	addDebug(queryRequest.GetActionName(), raw)
	response, _ := raw.(*eci.DescribeImageCachesResponse)
	for _, cache := range response.ImageCaches {
		if strings.HasPrefix(cache.ImageCacheName, "tf-testacc") {
			allCaches = append(allCaches, cache)
		} else {
			log.Printf("Skip %#v", cache)
		}
	}

	removeRequest := eci.CreateDeleteImageCacheRequest()
	removeRequest.ImageCacheId = ""
	for _, cache := range allCaches {
		removeRequest.ImageCacheId = cache.ImageCacheId
		raw, err := client.WithEciClient(func(eciClient *eci.Client) (interface{}, error) {
			return eciClient.DeleteImageCache(removeRequest)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error %s", removeRequest.GetActionName(), err)
		}
		addDebug(removeRequest.GetActionName(), raw)
	}
	return nil
}

func TestAccAlicloudECIOpenAPIImageCache_basic(t *testing.T) {
	var v eci.DescribeImageCachesImageCache0
	resourceId := "alicloud_eci_image_cache.default"
	ra := resourceAttrInit(resourceId, EciOpenapiImageCacheMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciImageCache")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccEciImageCache%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, EciOpenapiImageCacheBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_cache_name":  strings.ToLower(name),
					"images":            []string{"registry.cn-beijing.aliyuncs.com/sceneplatform/sae-image-demo:latest"},
					"security_group_id": "${alicloud_security_group.group.id}",
					"vswitch_id":        "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}",
					"eip_instance_id":   "${alicloud_eip.default.id}",
					"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
					"depends_on ":       []string{"alicloud_eip.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_cache_name":  strings.ToLower(name),
						"images.#":          "1",
						"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_group_id", "vswitch_id", "resource_group_id", "eip_instance_id"},
			},
		},
	})
}

var EciOpenapiImageCacheMap = map[string]string{
	"container_group_id": CHECKSET,
	"status":             CHECKSET,
	"eip_instance_id":    CHECKSET,
}

func EciOpenapiImageCacheBasicdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	  is_default = true
	}
	data "alicloud_vswitches" "default" {
	  ids = [data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0]
	}
	resource "alicloud_security_group" "group" {
	  name        = "%[1]s"
	  description = "tf-eci-image-test"
	  vpc_id      = data.alicloud_vpcs.default.vpcs.0.id
	}
	resource "alicloud_eip" "default" {
	  name = "%[1]s"
	}
`, name)
}
