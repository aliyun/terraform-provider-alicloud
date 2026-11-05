package eci_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	eciclient "github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eci"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	rawClient, err := eciTestClientForRegion(region)
	if err != nil {
		return errs.WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	queryRequest := eciclient.CreateDescribeImageCachesRequest()
	var allCaches []eciclient.DescribeImageCachesImageCache0

	raw, err := client.WithEciClient(func(eciClient *eciclient.Client) (interface{}, error) {
		return eciClient.DescribeImageCaches(queryRequest)
	})
	if err != nil {
		log.Printf("[ERROR] %s get an error %#v", queryRequest.GetActionName(), err)
	}
	logs.AddDebug(queryRequest.GetActionName(), raw)
	response, _ := raw.(*eciclient.DescribeImageCachesResponse)
	for _, cache := range response.ImageCaches {
		if strings.HasPrefix(cache.ImageCacheName, "tf-testacc") {
			allCaches = append(allCaches, cache)
		} else {
			log.Printf("Skip %#v", cache)
		}
	}

	removeRequest := eciclient.CreateDeleteImageCacheRequest()
	removeRequest.ImageCacheId = ""
	for _, cache := range allCaches {
		removeRequest.ImageCacheId = cache.ImageCacheId
		raw, err := client.WithEciClient(func(eciClient *eciclient.Client) (interface{}, error) {
			return eciClient.DeleteImageCache(removeRequest)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error %s", removeRequest.GetActionName(), err)
		}
		logs.AddDebug(removeRequest.GetActionName(), raw)
	}
	return nil
}

func TestAccAliCloudEciImageCache_basic(t *testing.T) {
	var v eciclient.DescribeImageCachesImageCache0
	acctest.CheckoutSupportedRegions(t, true, connectivity.EciContainerGroupRegions)
	resourceId := "alicloud_eci_image_cache.default"
	ra := acctest.ResourceAttrInit(resourceId, EciOpenapiImageCacheMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eci.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEciImageCache")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccEciImageCache%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, EciOpenapiImageCacheBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_cache_name":  strings.ToLower(name),
					"images":            []string{fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest)},
					"security_group_id": "${alicloud_security_group.group.id}",
					"vswitch_id":        "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}",
					"eip_instance_id":   "${alicloud_eip_address.default.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"image_cache_size":  20,
					"retention_days":    7,
					"zone_id":           "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"image_registry_credential": []map[string]interface{}{{
						"password":  "tftestacc123",
						"server":    fmt.Sprintf("registry-vpc.%s.aliyuncs.com", acctest.DefaultRegionToTest),
						"user_name": name,
					}},
					"depends_on ": []string{"alicloud_eip_address.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_cache_name":  strings.ToLower(name),
						"images.#":          "1",
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_cache_name":  strings.ToLower(name),
					"images":            []string{fmt.Sprintf("registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine", acctest.DefaultRegionToTest)},
					"security_group_id": "${alicloud_security_group.group.id}",
					"vswitch_id":        "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}",
					"eip_instance_id":   "${alicloud_eip_address.default.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"image_cache_size":  20,
					"retention_days":    7,
					"zone_id":           "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"image_registry_credential": []map[string]interface{}{{
						"password":  "tftestacc456",
						"server":    fmt.Sprintf("registry.%s.aliyuncs.com", acctest.DefaultRegionToTest),
						"user_name": name + "2",
					}},
					"depends_on ": []string{"alicloud_eip_address.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_cache_name":  strings.ToLower(name),
						"images.#":          "1",
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_group_id", "vswitch_id", "resource_group_id", "eip_instance_id", "image_cache_size", "image_registry_credential", "retention_days", "zone_id"},
			},
		},
	})
}

var EciOpenapiImageCacheMap = map[string]string{
	"container_group_id": acctest.CHECKSET,
	"status":             acctest.CHECKSET,
	"eip_instance_id":    acctest.CHECKSET,
}

func EciOpenapiImageCacheBasicdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	  name_regex = "^default-NODELETING$"
	}
	data "alicloud_vswitches" "default" {
	  ids = [data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0]
	}
	resource "alicloud_security_group" "group" {
	  name        = "%[1]s"
	  description = "tf-eci-image-test"
	  vpc_id      = data.alicloud_vpcs.default.vpcs.0.id
	}
	resource "alicloud_eip_address" "default" {
	  address_name = "%[1]s"
	}
	data "alicloud_resource_manager_resource_groups" "default" {
	  name_regex = "default"
	}
`, name)
}
