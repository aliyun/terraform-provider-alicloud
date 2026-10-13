package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCdnRefreshObjectCache_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cdn_refresh_object_cache.default"
	ra := resourceAttrInit(resourceId, AlicloudCdnRefreshObjectCacheMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CdnServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCdnRefreshObjectCache")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCdnRefreshObjectCacheBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"object_path": "http://${alicloud_cdn_domain_new.default.domain_name}/test.html",
					"object_type": "File",
					"force":       true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"object_path":   CHECKSET,
						"object_type":   "File",
						"force":         "true",
						"status":        CHECKSET,
						"process":       CHECKSET,
						"creation_time": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

// AlicloudCdnRefreshObjectCacheBasicDependence0 builds a self-owned CDN domain so
// that RefreshObjectCaches accepts the object_path. The domain_new Create waits
// until the domain reaches the "online" state, and object_path references its
// domain_name, so the refresh task is submitted only after the domain is owned
// and online by the test account.
func AlicloudCdnRefreshObjectCacheBasicDependence0(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_cdn_domain_new" "default" {
		domain_name = "%s"
		cdn_type    = "web"
		scope       = "overseas"
		sources {
			content  = "www.aliyuntest.com"
			type     = "domain"
			priority = 20
			port     = 80
			weight   = 10
		}
	}
`, name)
}

var AlicloudCdnRefreshObjectCacheMap0 = map[string]string{
	"object_path":   CHECKSET,
	"object_type":   CHECKSET,
	"force":         CHECKSET,
	"status":        CHECKSET,
	"process":       CHECKSET,
	"creation_time": CHECKSET,
}
