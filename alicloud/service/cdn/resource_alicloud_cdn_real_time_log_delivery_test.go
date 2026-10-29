package cdn_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cdn"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCDNRealTimeLogDelivery_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cdn_real_time_log_delivery.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCDNRealTimeLogDeliveryMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cdn.NewCdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCdnRealTimeLogDelivery")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scdnrealtimelogdelivery%d", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCDNRealTimeLogDeliveryBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain":     "${alicloud_cdn_domain_new.default.domain_name}",
					"project":    "${alicloud_log_project.default.name}",
					"logstore":   "${alicloud_log_store.default.name}",
					"sls_region": "${data.alicloud_regions.default.regions.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain":     acctest.CHECKSET,
						"project":    name,
						"logstore":   name,
						"sls_region": acctest.CHECKSET,
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

var AlicloudCDNRealTimeLogDeliveryMap0 = map[string]string{}

func AlicloudCDNRealTimeLogDeliveryBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_regions" "default" {
  current = true
}
resource "alicloud_cdn_domain_new" "default" {
  domain_name = "%s"
  cdn_type = "web"
  scope = "overseas"
  sources {
	 content = "www.aliyuntest.com"
	 type = "domain"
	 priority = 20
	 port = 80
	 weight = 10
  }
}
resource "alicloud_log_project" "default" {
  name        = var.name
  description = "created by terraform"
}
resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.name
  name                  = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
`, name, fmt.Sprintf("%s.alicloud-provider.cn", name))
}

func TestAccAliCloudCdnRealTimeLogDelivery_basic9189_modify(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cdn_real_time_log_delivery.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCdnRealTimeLogDeliveryMap9189)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cdn.NewCdnServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCdnRealTimeLogDelivery")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scdnrealtimelogdelivery%d", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCdnRealTimeLogDeliveryBasicDependence9189)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project":    "${alicloud_log_project.slsPoject.name}",
					"logstore":   "${alicloud_log_store.logstore.name}",
					"domain":     "${alicloud_cdn_domain_new.domain.domain_name}",
					"sls_region": "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":    acctest.CHECKSET,
						"logstore":   acctest.CHECKSET,
						"domain":     acctest.CHECKSET,
						"sls_region": "cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"project":    "${alicloud_log_project.updateProject.name}",
					"logstore":   "${alicloud_log_store.updateLogstore.name}",
					"sls_region": "cn-beijing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":    acctest.CHECKSET,
						"logstore":   acctest.CHECKSET,
						"sls_region": "cn-beijing",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "online",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "online",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "online",
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

var AlicloudCdnRealTimeLogDeliveryMap9189 = map[string]string{}

func AlicloudCdnRealTimeLogDeliveryBasicDependence9189(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  alias  = "r1"
  region = "cn-hangzhou"
}

provider "alicloud" {
  alias  = "r2"
  region = "cn-beijing"
}

resource "alicloud_log_project" "slsPoject" {
  provider   = alicloud.r1
  description = "test"
  name        = var.name
}

resource "alicloud_log_store" "logstore" {
  provider         = alicloud.r1
  retention_period = "60"
  shard_count      = "2"
  project          = alicloud_log_project.slsPoject.name
  name             = format("%%s1", var.name)
}

resource "alicloud_cdn_domain_new" "domain" {
  domain_name = format("%%s.alicloud-provider.cn", var.name)
  cdn_type = "web"
  scope = "overseas"
  sources {
	 content = "www.aliyuntest.com"
	 type = "domain"
	 priority = 20
	 port = 80
	 weight = 10
  }
}

resource "alicloud_log_project" "updateProject" {
  provider   = alicloud.r2
  description = "test"
  name        = format("%%s3", var.name)
}

resource "alicloud_log_store" "updateLogstore" {
  provider         = alicloud.r2
  retention_period = "60"
  shard_count      = "2"
  project          = alicloud_log_project.updateProject.name
  name             = format("%%s4", var.name)
}


`, name)
}
